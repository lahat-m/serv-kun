
# serv-kun

A CLI toolbox of guided setup wizards for provisioning a Linux server — no
config file editing required.

- **nginx** — reverse proxy config + SSL via Let's Encrypt
- **user** — create a Linux user with sudo privileges
- **ssh** — disable root login and password auth, with a lockout safety check
- **firewall** — default-deny inbound via ufw/firewalld, allowing SSH, HTTP and HTTPS
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

## From your workstation

serv-kun runs *on the server*, so these steps stay manual — they happen on
the machine you're connecting **from**, before and after the wizards run.

### 1. Generate an SSH key

```bash
ssh-keygen -t ed25519 -C "your-comment"
```

Ed25519 keys are a fixed size, so there's no `-b` to pass — `ssh-keygen`
silently ignores it for this key type (it only applies to `-t rsa`, where
you'd want `-b 4096`). This writes `~/.ssh/id_ed25519` and the public half
`~/.ssh/id_ed25519.pub`.

### 2. Get the public key onto the server

If you can still log in with a password:

```bash
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@server
```

Otherwise print it and paste it when `sudo serv-kun user` asks for it —
that wizard installs it with the right ownership and permissions
(`~/.ssh` `0700`, `authorized_keys` `0600`), which is the part that most
often gets silently wrong by hand:

```bash
cat ~/.ssh/id_ed25519.pub
```

### 3. Log in

```bash
ssh user@server
```

### 4. Set up an SSH alias

Saves retyping the host, user, key and port every time. Add to `~/.ssh/config`
on your workstation:

```
Host myserver
    HostName 203.0.113.10
    User lahat
    IdentityFile ~/.ssh/id_ed25519
    IdentitiesOnly yes
```

Then `ssh myserver` is enough. `scp`, `rsync`, and `git` understand the same
alias. `IdentitiesOnly yes` makes ssh offer only that key rather than every
key in your agent, which matters once the server stops accepting passwords
and starts counting failed attempts against `MaxAuthTries`.

## First deployment: a worked example

A complete first run on a fresh Ubuntu VPS, in order. Every answer below is
sample data — substitute your own, but keep to the same shape.

| | Sample value |
|---|---|
| Server IP | `203.0.113.10` |
| Domain | `kuniel.com` |
| Sudo user | `lahat` |
| App name | `kuniel` |
| App image | `ghcr.io/lahat-m/kuniel:latest` |

### Step 1 — Point DNS at the server

Create an `A` record for `kuniel.com` → `203.0.113.10` before anything else.
Let's Encrypt validates over HTTP, so certificate issuance in step 9 fails if
DNS hasn't propagated. Behind Cloudflare, set SSL/TLS mode to **Full (strict)**.

### Step 2 — Generate a key and install serv-kun

Generate an SSH key on your workstation ([above](#from-your-workstation)), then
on the server, as root:

```bash
curl -fsSL https://raw.githubusercontent.com/lahat-m/serv-kun/main/install.sh | sudo bash
```

### Step 3 — Create the sudo user

```bash
sudo serv-kun user
```

| Prompt | Sample answer | What to use |
|---|---|---|
| Username for the new sudo user | `lahat` | Lowercase, no spaces. This is who you'll SSH as from now on. |
| How should this user authenticate over SSH? | **both** `SSH public key` **and** `Password` | See the warning below — pick both. |
| Paste the SSH public key | contents of `~/.ssh/id_ed25519.pub` | The `.pub` file, one line, starting `ssh-ed25519 `. Never the private key. |
| Set a password for this user | a strong password | Used for `sudo`, not for SSH login. Store it in your password manager. |

> **Select both auth methods.** A key-only account is created with a *locked*
> password, and Ubuntu's default sudoers requires a password to run `sudo` —
> so a key-only user can log in but gets `sudo: a password is required` and
> can administer nothing. The key is for SSH; the password is for `sudo`. You
> still get key-only *login* once step 5 disables password authentication,
> because that setting governs SSH auth, not `sudo`. (`ssh harden` checks for
> this in step 5 and refuses to run, but it's easier to get right here.)

### Step 4 — Verify access before hardening anything

Open a **new terminal** and confirm both halves work:

```bash
ssh lahat@203.0.113.10
sudo whoami          # must print: root
```

Do not continue until both succeed. Steps 5 and 6 each remove a way into the
machine; this is the check that makes them safe.

### Step 5 — Harden SSH

```bash
sudo serv-kun ssh harden
```

| Prompt | Sample answer | What to use |
|---|---|---|
| Admin user that must keep access | `lahat` | Defaults to `$SUDO_USER`. Must be the user you just verified in step 4. |
| Disable root login over SSH? | `Yes` | |
| Disable password authentication (keys only)? | `Yes` | Say `No` only if you need password SSH from machines without your key. |
| Apply this configuration? | `Yes` | Shown after a preview of the exact directives. |

Keep the session open afterwards and re-run the step 4 checks in a second
terminal. If something is wrong, the still-open session can undo it.

### Step 6 — Turn on the firewall

```bash
sudo serv-kun firewall
```

| Prompt | Sample answer | What to use |
|---|---|---|
| ufw is not installed. Install it now? | `Yes` | |
| Additional TCP ports to open | *(blank)* | Blank is right for this stack — nginx fronts everything on 80/443. Add app ports only if something bypasses nginx. |
| Apply this firewall configuration? | `Yes` | Review the rule list first; the SSH port shown is read from the running sshd. |

Don't open `5432` for Postgres. In this stack it isn't published to the host
at all; if you need remote access, tunnel it: `ssh -L 5432:localhost:5432 lahat@203.0.113.10`.

### Step 7 — Scaffold the project (on your workstation)

Run this in your app's git repo — it doesn't need root and only writes files:

```bash
serv-kun springboot init
```

| Prompt | Sample answer | What to use |
|---|---|---|
| App name | `kuniel` | Lowercase, no spaces. Becomes the container prefix (`kuniel-app`, `kuniel-db`) and the certs mount path. |
| Port the app listens on inside its container | `8080` | Your `server.port`. Container-internal — it is never published to the host. |
| Max upload size | `50M` | Match your `spring.servlet.multipart.max-request-size`, or uploads fail at nginx with a confusing 413. |
| Mount `certs/` into the app container? | `Yes` | Yes if you sign JWTs with RS256. No for HS256 or no auth. |

Add a `Dockerfile` if you don't have one, then commit and push everything.

### Step 8 — Clone on the server

```bash
sudo git clone https://github.com/lahat-m/kuniel.git /opt/kuniel
cd /opt/kuniel
```

### Step 9 — Bring up the stack

```bash
sudo serv-kun springboot setup
```

| Prompt | Sample answer | What to use |
|---|---|---|
| Docker is not installed. Install it now? | `Yes` | |
| Domain this app will be served on | `kuniel.com` | Hostname only — no `https://`, no trailing slash, no port. |
| Email for Let's Encrypt renewal notices | `you@kuniel.com` | A real, monitored inbox. Expiry warnings go here. |
| App image | `ghcr.io/lahat-m/kuniel:latest` | Full registry path including tag. Must match the image your CI pushes. |
| Postgres database name | `kuniel` | Default is `app`; naming it after the project is clearer. |
| Postgres role name | `kuniel` | |
| Postgres password | *(generated)* | Generate it — don't invent one: `openssl rand -base64 24`. Stored in `.env.prod` (mode `0600`). |
| Days to keep nightly backups | `7` | Raise it if your data changes slowly and disk allows. |
| Allowed CORS origin(s) | `https://kuniel.com` | Defaults to `https://<domain>`. Add your frontend's origin if it's hosted separately. |
| Does this app use RS256 JWT keys? | `Yes` | Match what you answered in step 7. |
| → Path `certs/` is mounted at | `/etc/kuniel/certs` | Defaulted by reading the `./certs` bind mount out of your compose file, so it should already be right — accept it unless you changed the mount by hand. |
| Seed an initial admin account? | `Yes` | Only on a genuinely fresh database. |
| → Seeded admin email / password / full name | `admin@kuniel.com` / *(generated)* / `Admin` | Change the password after first login. |
| Does this app send email via Resend? | `No` | `Yes` needs an API key and a verified sender domain. |
| Does this app use S3-compatible object storage? | `No` | `Yes` needs endpoint, access key, secret key, bucket, region. |
| Proceed? | `Yes` | Review the summary — domain, image and DB name are the easy typos. |

> **If your JWT keys live in `certs/`, check the container can read them.**
> The wizard writes `certs/` as `0700` and `private.pem` as `0600`, both owned
> by root. If your `Dockerfile` sets a non-root `USER`, the app can see the key
> and not open it, and dies at startup. `sudo chown -R <uid>:<gid> certs/`
> using the container's UID fixes it; the wizard prints this reminder too.

A **proxied Cloudflare record** (orange cloud) is recognised and reported as
`→ Cloudflare edge (proxied)` rather than a failure — the HTTP-01 challenge
still reaches the origin, so issuance works normally. If instead DNS genuinely
doesn't resolve here, you'll be asked whether to proceed anyway: answer `No`,
fix the record, and re-run. Certificate issuance can't succeed until it
resolves, and repeated failures burn Let's Encrypt rate limits.

### Step 10 — Wire up continuous deployment

In your GitHub repo, add these four secrets (Settings → Secrets and variables
→ Actions) for the scaffolded `.github/workflows/deploy.yml`:

| Secret | Sample value |
|---|---|
| `HETZNER_HOST` | `203.0.113.10` |
| `HETZNER_SSH_USER` | `lahat` |
| `HETZNER_SSH_KEY` | contents of `~/.ssh/id_ed25519` (the **private** key) |
| `HETZNER_APP_DIR` | `/opt/kuniel` |

The workflow SSHes in and runs `sudo serv-kun springboot deploy`. That needs
passwordless sudo for non-interactive use, so add it on the server:

```bash
echo 'lahat ALL=(ALL) NOPASSWD: /usr/local/bin/serv-kun' | sudo tee /etc/sudoers.d/serv-kun
sudo chmod 0440 /etc/sudoers.d/serv-kun
```

Scoping `NOPASSWD` to just the `serv-kun` binary keeps the blast radius small
— interactive `sudo` for everything else still asks for the password from
step 3.

### Step 11 — Day to day

```bash
sudo serv-kun springboot status   # container status + health endpoint
sudo serv-kun springboot logs     # tail the app logs
sudo serv-kun springboot deploy   # manual redeploy, with automatic rollback
```

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

Pick **both** unless you know otherwise. The account is created with a locked
password, and the default sudoers policy requires one to run `sudo` — so a
key-only user can log in and then gets `sudo: a password is required` for
every administrative command. The key authenticates SSH; the password
authenticates `sudo`. Disabling password authentication later (`ssh harden`)
doesn't undo this: it governs SSH login, not `sudo`.

### ssh — disable root login and password authentication

```bash
sudo serv-kun ssh harden
```

Shows the current effective settings, then asks which to change: root login,
password authentication, or both. Disabling passwords also sets
`KbdInteractiveAuthentication no` — without it PAM can still offer an
interactive password prompt, and `PasswordAuthentication no` alone buys you
nothing.

#### Lockout safety

Before touching anything, it asks which user must keep access and refuses to
proceed unless that user is a non-root account **and** has at least one key in
`~/.ssh/authorized_keys` **and** can actually run `sudo`. That last check is
stronger than group membership: an account created with only an SSH key sits
in the sudo group with a locked password, and the default sudoers policy
rejects every command it runs. So the check confirms either a usable password
or a `NOPASSWD` sudoers rule — otherwise you'd be left able to log in and
unable to administer anything. If it fails you're told how to fix it
(`sudo serv-kun user`) and can only continue via an explicit opt-in that
defaults to no, for when you have console or rescue access.

#### Where it writes

On distros whose `sshd_config` has an `Include /etc/ssh/sshd_config.d/*.conf`
line (Ubuntu 22.04+, Debian 12+), it writes a drop-in named
`00-serv-kun-hardening.conf` rather than editing the main file. That naming is
load-bearing: sshd uses the **first** value it obtains for a keyword and
`Include` expands in lexical order, so a vendor file like
`50-cloud-init.conf` setting `PasswordAuthentication yes` would otherwise win
and silently undo the change. Editing the main config is the trap here, not
the fix.

Where there's no `Include` line, it edits `/etc/ssh/sshd_config` directly:
existing top-level settings for the managed keywords are commented out and the
new block is inserted at the top. Directives inside `Match` blocks are left
alone — those are deliberate per-connection overrides.

#### Validation and rollback

Either way a timestamped backup is taken first, the result is checked with
`sshd -t`, and a failure restores the backup **before** sshd is ever reloaded.
On success it reloads (not restarts, so existing sessions survive), then
re-reads the effective config to confirm the change actually took, and reminds
you to verify login from a second terminal before closing the current one.

### firewall — default-deny inbound

```bash
sudo serv-kun firewall
```

Installs and configures `ufw` (Debian/Ubuntu) or `firewalld` (RHEL/Fedora):
default deny incoming, allow outgoing, with SSH, HTTP and HTTPS allowed plus
any extra TCP ports you name. Shows the full rule list for review before
applying.

The SSH port is read from the running sshd config rather than assumed to be
`22`, and rules are added *before* default-deny takes effect — enabling a
firewall that doesn't allow the port you're connected on is the classic way to
lose a server.

Database ports are deliberately not opened. In the `springboot` stack Postgres
is never published to the host, so it needs no rule; if you do run Postgres
directly on the host, prefer an SSH tunnel over exposing `5432`.

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
