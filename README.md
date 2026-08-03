# serv-kun

A CLI toolbox of guided setup wizards for provisioning a Linux server — no
config file editing required.

- **nginx** — reverse proxy config + SSL via Let's Encrypt
- **user** — create a Linux user with sudo privileges
- **postgres** — install PostgreSQL, create a deployment role/database, and back up/restore
- **node** — install Node.js via nvm for a user and pick a default version
- **springboot** — set up and deploy a Dockerized Spring Boot app shaped like [Kuniel](https://github.com/lahat-m/kuniel)'s `docker-compose.prod.yml`

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/lahat-m/serv-kun/main/install.sh | sudo bash
```

Or download a binary directly from [Releases](https://github.com/lahat-m/serv-kun/releases):

```bash
sudo curl -fsSL https://github.com/lahat-m/serv-kun/releases/latest/download/serv-kun-linux-amd64 \
  -o /usr/local/bin/serv-kun && sudo chmod +x /usr/local/bin/serv-kun
```

> Use `serv-kun-linux-arm64` for ARM-based servers (AWS Graviton, etc.)

## Usage

All commands require root:

### nginx — reverse proxy + SSL

Full setup wizard. Installs nginx if needed, writes a reverse-proxy config, and optionally sets up SSL.

```bash
sudo serv-kun nginx setup
```

Prompts for:
- Upstream URL (e.g. `http://localhost:3000`)
- Domain(s) (e.g. `example.com www.example.com`)
- Max upload size (default: `100M`)

Shows the generated config for review before writing anything. Creates a timestamped backup of any existing config before overwriting.

Add SSL to an existing config:

```bash
sudo serv-kun nginx ssl
```

Prompts for the path to your existing nginx config and your domain(s), verifies DNS A records point to this server, then runs certbot.

#### What gets configured

The generated nginx config includes:

- WebSocket support (`Upgrade`, `Connection` headers)
- Standard forwarding headers (`X-Real-IP`, `X-Forwarded-For`, `X-Forwarded-Proto`)
- Configurable `client_max_body_size`
- Long-lived connection timeout (`proxy_read_timeout 86400s`)
- `proxy_buffering off`
- HTTP → HTTPS redirect (after SSL is set up)

#### DNS verification

Before running certbot, the tool resolves each domain's A record and compares it to the server's public IP. If any domain doesn't match, you'll see which ones failed and can choose to abort or proceed anyway.

### user — create a sudo user

```bash
sudo serv-kun user
```

Prompts for a username, creates the account, adds it to the `sudo` (Debian) or `wheel` (RHEL) group, then asks how it should authenticate over SSH — an SSH public key, a password, or both.

### postgres — install PostgreSQL + a deployment role/database

```bash
sudo serv-kun postgres setup
```

Installs PostgreSQL (via the OS package manager) if it isn't already present, then prompts for a role name, password, and database name. The role is granted ownership of the database, and `pg_hba.conf` is patched so that TCP connections from `127.0.0.1`/`::1` authenticate with a password — the fix is applied by prepending rules (pg_hba.conf matches top-down), so it takes effect without touching any existing entries.

#### Backup

```bash
sudo serv-kun postgres backup
```

Lets you pick a database and a destination path, then writes a gzip-compressed plain-SQL `pg_dump` there (`--format=plain --clean --if-exists --no-owner`, piped through gzip). Plain SQL has no archive-format version header, so unlike `pg_dump`'s custom/directory/tar formats it can always be replayed by any `psql`, regardless of version skew between the machine that took the backup and the one restoring it. It's not encrypted, though — gzip is compression only, so treat the backup file as sensitive as the database itself.

#### Restore

```bash
sudo serv-kun postgres restore
```

Prompts for a backup file and a target database — either an existing one or a new one to create on the spot. **Restoring into an existing database is destructive**: the backup's own `--clean --if-exists` statements drop conflicting objects as they're replayed, so you're asked to confirm before anything happens.

### node — install Node.js via nvm

```bash
sudo serv-kun node
```

Prompts for which Linux user to install Node.js for (nvm and Node live in that user's home directory), installs [nvm](https://github.com/nvm-sh/nvm) if needed, then lets you install one or more Node versions (`--lts`, `20`, `22.4.0`, ...) and pick which one becomes that user's default.

The default version's `node`/`npm`/`npx` are also symlinked into `/usr/local/bin`, which is on every shell's `PATH` already — unlike `~/.nvm`, which only gets wired up when a shell's rc file is sourced at login. That means they're usable immediately, including in the session you just ran this from, with no need to log out and back in. Re-running `serv-kun node` and picking a different default updates those symlinks to match.

### springboot — set up and deploy a Dockerized Spring Boot app

Unlike the other wizards, this one is opinionated about project shape rather
than generic: it targets a `docker-compose.prod.yml` with `app`, `postgres`,
`nginx`, and `certbot` services, driven by a `.env.prod` file — the shape
used by [Kuniel](https://github.com/lahat-m/kuniel) and any other Spring
Boot + Postgres project set up the same way. Run these from the project
directory (the one containing `docker-compose.prod.yml`), or pass `--dir`.

```bash
serv-kun springboot init
```

Scaffolds that whole shape for a brand-new project: `docker-compose.prod.yml`,
`nginx/templates/app.conf.template`, `nginx/conf.d/hardening.conf`,
`scripts/backup-postgres.sh`, `.github/workflows/deploy.yml`, and a
`DEPLOY.md` — parameterized by app name, the port the app listens on, max
upload size, and whether to mount `certs/` for JWT keys. Doesn't require
root (it only writes files); does require confirmation before overwriting
anything that already exists, keeping a timestamped backup either way. You
still need to add your own `Dockerfile` — that's build-specific enough per
project that it isn't scaffolded.

```bash
sudo serv-kun springboot setup
```

Installs Docker if it isn't already present, then prompts for the fields
every such project needs (domain, image, Postgres credentials, CORS origin).
JWT keys, a seeded admin account, Resend email, and S3-compatible object
storage are each asked about as a yes/no first — not every project needs all
of them, and declining one skips its prompts and omits it from `.env.prod`
entirely rather than writing it out blank. If JWT is enabled, it generates
an RSA keypair into `certs/`. It then verifies the domain's DNS A record
points at this server, bootstraps the first Let's Encrypt certificate
(writes a throwaway self-signed cert so nginx can start, swaps in the real
one via the HTTP-01 challenge, reloads nginx), brings the whole stack up,
and polls `/actuator/health` until it responds before declaring success.

```bash
sudo serv-kun springboot deploy
```

Pulls the latest `app` image and recreates that one container — the
existing `postgres`, `nginx`, and `certbot` services are left alone. Before
pulling, it records the currently-running image's ID; if the new image fails
its `/actuator/health` check, it automatically re-tags that ID back onto the
image reference and recreates the container again, so a bad deploy doesn't
leave the site down. A rollback still exits non-zero — the deploy failed and
needs investigating even though the site recovered.

```bash
sudo serv-kun springboot status   # docker compose ps + a health-endpoint check
sudo serv-kun springboot logs     # tail -f the app container's logs
```

## Supported systems

| Distro | Package manager |
|---|---|
| Ubuntu / Debian | `apt` |
| CentOS / RHEL / Rocky / AlmaLinux / Fedora | `dnf` |

## Building from source

Requires Go 1.25+.

```bash
git clone https://github.com/lahat-m/serv-kun.git
cd serv-kun
go build -o serv-kun .
sudo mv serv-kun /usr/local/bin/
```

## Releasing

Push a `v*` tag to trigger a GitHub Actions build that cross-compiles for `linux/amd64` and `linux/arm64` and publishes the binaries to GitHub Releases.

```bash
git tag v1.0.0
git push origin v1.0.0
```
