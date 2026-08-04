package postgres

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lahat-m/serv-kun/internal/system"
)

func IsInstalled() bool {
	_, err := exec.LookPath("psql")
	return err == nil
}

func Install(osInfo *system.Info) error {
	switch osInfo.Family {
	case system.Debian:
		if err := system.RunCmd("apt-get", "update", "-qq"); err != nil {
			return fmt.Errorf("apt-get update: %w", err)
		}
		return system.RunCmd("apt-get", "install", "-y", "postgresql", "postgresql-contrib")

	case system.RHEL:
		if err := system.RunCmd("dnf", "install", "-y", "postgresql-server", "postgresql-contrib"); err != nil {
			return err
		}
		// RHEL-family packages don't auto-initialize the data directory.
		if _, err := os.Stat("/var/lib/pgsql/data/PG_VERSION"); os.IsNotExist(err) {
			if err := system.RunCmd("postgresql-setup", "--initdb"); err != nil {
				return fmt.Errorf("initdb: %w", err)
			}
		}
		return nil
	}

	return fmt.Errorf("unsupported OS family")
}

// Enable starts PostgreSQL and enables it on boot.
func Enable() error {
	return system.RunCmd("systemctl", "enable", "--now", "postgresql")
}

// Restart restarts PostgreSQL so pg_hba.conf changes take effect.
func Restart() error {
	return system.RunCmd("systemctl", "restart", "postgresql")
}

// RoleExists reports whether a login role of this name is already present.
func RoleExists(name string) (bool, error) {
	out, err := runToolAsPostgres("psql", "-tAc",
		fmt.Sprintf("SELECT 1 FROM pg_roles WHERE rolname = %s", quoteLiteral(name)))
	if err != nil {
		return false, fmt.Errorf("%w\n%s", err, out)
	}
	return strings.TrimSpace(out) == "1", nil
}

// CreateRole creates a LOGIN role with the given password.
func CreateRole(name, password string) error {
	sql := fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD %s;",
		quoteIdent(name), quoteLiteral(password))
	out, err := runToolAsPostgres("psql", "-v", "ON_ERROR_STOP=1", "-c", sql)
	if err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}

// SetRolePassword resets an existing role's password. Callers need this
// because CreateRole can't: re-running setup against a role that already
// exists used to discard the password that was just typed, while still
// reporting success and printing a connection string that wouldn't work.
func SetRolePassword(name, password string) error {
	sql := fmt.Sprintf("ALTER ROLE %s WITH LOGIN PASSWORD %s;",
		quoteIdent(name), quoteLiteral(password))
	out, err := runToolAsPostgres("psql", "-v", "ON_ERROR_STOP=1", "-c", sql)
	if err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}

// CreateDatabase creates a database owned by owner if it doesn't already exist.
func CreateDatabase(name, owner string) error {
	exists, err := databaseExists(name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	sql := fmt.Sprintf("CREATE DATABASE %s OWNER %s;", quoteIdent(name), quoteIdent(owner))
	out, err := runToolAsPostgres("psql", "-v", "ON_ERROR_STOP=1", "-c", sql)
	if err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}

func databaseExists(name string) (bool, error) {
	out, err := runToolAsPostgres("psql", "-tAc", fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = %s", quoteLiteral(name)))
	if err != nil {
		return false, fmt.Errorf("%w\n%s", err, out)
	}
	return strings.TrimSpace(out) == "1", nil
}

// ListDatabases returns the names of all non-template databases.
func ListDatabases() ([]string, error) {
	out, err := runToolAsPostgres("psql", "-tAc", "SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname;")
	if err != nil {
		return nil, fmt.Errorf("%w\n%s", err, out)
	}

	var dbs []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			dbs = append(dbs, line)
		}
	}
	return dbs, nil
}

// Backup dumps database as a gzip-compressed plain-SQL pg_dump to destPath.
// Plain SQL has no archive-format version header, so — unlike pg_dump's
// custom/directory/tar formats — restoring it never runs into pg_dump/pg_restore
// version-skew errors ("unsupported version in file header"); any psql can
// replay it. Note: gzip is compression, not encryption — the backup file is
// exactly as sensitive as the database itself and should be handled that way.
func Backup(database, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create backup file: %w", err)
	}
	defer f.Close()

	gz := gzip.NewWriter(f)

	cmd, err := postgresCommand("pg_dump", "--format=plain", "--clean", "--if-exists", "--no-owner", "-d", database)
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stdout = gz
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump: %w\n%s", err, stderr.String())
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("finalize backup file: %w", err)
	}
	return nil
}

// Restore replays a gzip-compressed plain-SQL backup produced by Backup into
// database. The backup's own --clean --if-exists statements drop conflicting
// objects as they're replayed, so any existing data in database that overlaps
// with the backup is overwritten.
func Restore(database, srcPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", srcPath, err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("read %s as gzip: %w", srcPath, err)
	}
	defer gz.Close()

	cmd, err := postgresCommand("psql", "-v", "ON_ERROR_STOP=1", "-d", database)
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stdin = gz
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore: %w\n%s", err, stderr.String())
	}
	return nil
}

// EnableLocalPasswordAuth finds pg_hba.conf for the OS and prepends rules so
// TCP connections from localhost use password auth, regardless of any other
// rules further down the file (pg_hba.conf matches top-down, first match wins).
// If the file can't be located, this is a no-op — auth config is left alone.
func EnableLocalPasswordAuth(osInfo *system.Info) error {
	path, err := hbaPath(osInfo)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read pg_hba.conf: %w", err)
	}
	if strings.Contains(string(data), "# added by serv-kun") {
		return nil
	}

	backup := path + ".bak." + time.Now().Format("20060102-150405")
	if err := os.WriteFile(backup, data, 0640); err != nil {
		return fmt.Errorf("backup pg_hba.conf: %w", err)
	}

	method := passwordAuthMethod()
	header := fmt.Sprintf(
		"# added by serv-kun — force password auth for local TCP connections\n"+
			"host    all             all             127.0.0.1/32            %s\n"+
			"host    all             all             ::1/128                 %s\n\n",
		method, method)

	if err := os.WriteFile(path, append([]byte(header), data...), 0640); err != nil {
		return fmt.Errorf("write pg_hba.conf: %w", err)
	}
	return nil
}

// passwordAuthMethod picks the pg_hba method matching how this server
// actually stores passwords. Writing scram-sha-256 against a server still on
// md5 storage (PostgreSQL 13 and older, as shipped by RHEL 9) locks the role
// out, so the server's own setting decides rather than a hardcoded guess.
func passwordAuthMethod() string {
	out, err := runToolAsPostgres("psql", "-tAc", "SHOW password_encryption;")
	if err == nil && strings.TrimSpace(out) == "scram-sha-256" {
		return "scram-sha-256"
	}
	return "md5"
}

func hbaPath(osInfo *system.Info) (string, error) {
	if osInfo.Family == system.RHEL {
		p := "/var/lib/pgsql/data/pg_hba.conf"
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		return "", nil
	}

	matches, err := filepath.Glob("/etc/postgresql/*/main/pg_hba.conf")
	if err != nil || len(matches) == 0 {
		return "", nil
	}
	return matches[0], nil
}

// runToolAsPostgres runs a Postgres client tool (psql, pg_dump, ...) as the
// postgres OS user and returns its combined output.
func runToolAsPostgres(tool string, args ...string) (string, error) {
	cmd, err := postgresCommand(tool, args...)
	if err != nil {
		return "", err
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// postgresCommand builds a command for a Postgres client tool (psql, pg_dump, ...)
// that runs as the postgres OS user, dropping privileges directly (no shell
// involved) so arguments never need shell escaping — only the SQL-level quoting
// in quoteIdent/quoteLiteral applies to values embedded in SQL text. Stdin/Stdout/
// Stderr are left for the caller to wire up.
func postgresCommand(tool string, args ...string) (*exec.Cmd, error) {
	u, uid, gid, err := postgresUser()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(tool, args...)
	cmd.Dir = "/tmp"
	cmd.Env = append(os.Environ(), "HOME="+u.HomeDir, "USER=postgres")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
	}
	return cmd, nil
}

func postgresUser() (*user.User, int, int, error) {
	u, err := user.Lookup("postgres")
	if err != nil {
		return nil, 0, 0, fmt.Errorf("lookup postgres user: %w", err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("parse postgres uid: %w", err)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("parse postgres gid: %w", err)
	}
	return u, uid, gid, nil
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func quoteLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
