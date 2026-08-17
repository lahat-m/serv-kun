
# serv-kun

A CLI toolbox of guided setup wizards for provisioning a Linux server — no
config file editing required.

- **nginx** — reverse proxy config + SSL via Let's Encrypt
- **user** — create a Linux user with sudo privileges
- **ssh** — disable root login and password auth with a lockout safety check, and audit login attempts from the logs
- **firewall** — default-deny inbound via ufw/firewalld, allowing SSH, HTTP and HTTPS, with optional SSH rate limiting and source allow-listing
- **postgres** — install PostgreSQL, create a deployment role/database, and back up/restore
- **node** — install Node.js via nvm for a user and pick a default version
- **springboot** — set up and deploy a Dockerized Spring Boot app shaped like [Kuniel](https://github.com/lahat-m/kuniel)'s `docker-compose.prod.yml`

> Built on [**server-wiz**](https://github.com/hardope/server-wiz) by
> [hardope](https://github.com/hardope) — see [Acknowledgements](#acknowledgements).

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
| Rate-limit inbound SSH on port 22? | `Yes` | Blocks a source IP that opens 6+ connections in 30 seconds. Answer `No` if CI deploys over SSH from one fixed address, or if your office shares a public IP behind NAT. |
| Restrict SSH to specific source addresses? | `No` | Say `Yes` only if you can list every address that will ever need SSH — office, VPN, CI runner. It's the one control that stops a distributed attack, and the one that can lock you out. Your live sessions are checked against the list either way. |
| Additional TCP ports to open | *(blank)* | Blank is right for this stack — nginx fronts everything on 80/443. Add app ports only if something bypasses nginx. |
| Apply this firewall configuration? | `Yes` | Review the rule list first; the SSH port shown is read from the running sshd. |

Don't open `5432` for Postgres. In this stack it isn't published to the host
at all; if you need remote access, tunnel it: `ssh -L 5432:localhost:5432 lahat@203.0.113.10`.

Once the server has been up for a day, check what's been knocking:

```bash
sudo serv-kun ssh audit
```

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
| *(project-specific variables)* | *(varies)* | Asked only when this project needs them — see below. |
| Proceed? | `Yes` | Review the summary — domain, image and DB name are the easy typos. |

> **The wizard asks what your project's own configuration requires.**
> The questions above produce a fixed set of variable names, which is only ever
> right for the project they were modelled on. Before writing anything, setup
> reads `src/main/resources/application.properties` and
> `application-prod.properties` for every property that resolves from an
> environment variable *with no fallback*, and asks for the ones the generated
> file wouldn't cover — using the comment above each property as the prompt's
> help text.
>
> It counts a property as covered if the file defines either the variable the
> author named **or** the property's own relaxed-binding name, because Spring
> accepts both and the environment outranks the properties file. That's why a
> project reading `${JWT_PUBLIC_KEY_PATH}` is silently satisfied by the
> `JWT_PUBLIC_KEY` this wizard writes, and isn't asked about — while
> `app.storage.endpoint=${S3_ENDPOINT}` is not covered by `APP_S3_ENDPOINT`
> (that maps to `app.s3.endpoint`) and *is* asked for. Without the distinction
> the check would flag eight satisfied variables next to every real one and
> quickly become noise.
>
> A blank answer is allowed and warned about: the app will see the literal
> `${VAR}` and fail on it. A project deployed from a published image, with no
> sources on the box, is skipped silently; one configured through YAML is
> skipped with a note, since that isn't parsed yet.
>
> **Re-running setup no longer discards variables it didn't write.** Anything
> the previous `.env.prod` defined that the new one doesn't is carried across
> and reported, rather than surviving only in the timestamped backup.

> **JWT keys in `certs/` are handed to the container's user for you.**
> The wizard writes `certs/` as `0700` and `private.pem` as `0600`, both owned
> by root — which an image with a non-root `USER` can see but not open, killing
> the app at startup on its own signing key. Once the image exists, setup
> resolves the numeric UID that `USER` stands for and chowns `certs/` to it,
> reporting `✓ certs/ handed to 1001:1001`. Permissions are untouched: the keys
> become readable to that one container user, not to the host at large. On an
> image whose `USER` can't be resolved — a name, with no `id(1)` to resolve it
> against, as in a distroless base — setup says so and prints the
> `sudo chown -R <uid>:<gid> certs/` to run by hand.

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

A secret you never set is substituted as an empty string rather than failing,
so the workflow checks all four before it connects and names any that are
blank. Without that check the first symptom is `missing server host` from deep
inside the SSH action, which tells you neither which secret nor how many.

The workflow SSHes in and runs `sudo serv-kun springboot deploy`. An SSH
session has no terminal, so sudo has nowhere to prompt for a password and
fails with *"a terminal is required to authenticate"*. Grant the deploy user
passwordless sudo on the server:

```bash
sudo serv-kun ci
```

It asks which user CI logs in as (defaulting to you), shows the rule before
writing anything, validates it with `visudo` while it is still staged under a
name sudo ignores, then installs it as `/etc/sudoers.d/serv-kun` and confirms
with sudo that it took.

The rule covers the `serv-kun` binary and `docker image prune -f` — the two
commands the workflow actually runs — and nothing else. Interactive `sudo` for
everything else still asks for the password from step 3. The docker entry pins
its arguments deliberately: unrestricted `docker` is a root shell in one hop
(`docker run -v /:/host`), which is also why the deploy user should not be in
the `docker` group.

> The rule is per-user. If `HETZNER_SSH_USER` holds a different account than
> the one you answered here, it won't match and the deploy still fails.

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

### ci — let CI deploy without a password prompt

```bash
sudo serv-kun ci
```

The workflow `springboot init` scaffolds SSHes in and runs `sudo serv-kun
springboot deploy`. That session has no terminal, so sudo cannot prompt and
the deploy dies with `sudo: a terminal is required to authenticate` — a
message that names neither sudoers nor the user it means.

This writes `/etc/sudoers.d/serv-kun` granting that user the two commands the
deploy runs, without a password:

- the `serv-kun` binary, with any arguments — its subcommands are the point
- `docker image prune -f`, with those arguments exactly

The docker entry is pinned because unrestricted `docker` is a larger grant
than passwordless `serv-kun`: `docker run -v /:/host` is a root shell. For the
same reason, don't put the deploy user in the `docker` group.

Three details it handles that a hand-written rule usually doesn't:

- **The file is validated before it can apply.** It's staged inside
  `/etc/sudoers.d` under a name containing a dot, which sudo skips, and only
  renamed into place once `visudo` accepts it. A malformed sudoers file that
  sudo *does* read breaks `sudo` for everyone at once.
- **The filename has no dot in it.** sudo ignores `serv-kun.conf` entirely —
  the rule would be valid, root-owned, and inert.
- **It asks sudo whether the rule took**, rather than trusting the write. A
  rule can be installed and still not apply if something in `/etc/sudoers`
  overrides it.

`springboot setup` warns when it finds a scaffolded `deploy.yml` and the
invoking user has no such rule, so this surfaces at setup rather than on the
first push.

### ssh harden — disable root login and password authentication

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

### ssh audit — who has been trying to log in

```bash
sudo serv-kun ssh audit
sudo serv-kun ssh audit --since 7d --top 20
sudo serv-kun ssh audit --top 0            # every address, nothing truncated
sudo serv-kun ssh audit --file /var/log/auth.log.1
```

`--top 0` lifts the row limit on every table at once — sources, blocks and
accounts — for when you want the full list rather than a summary, to build
deny rules from or to check whether a particular address appears at all. A
busy public server sees a four-figure number of addresses a week, so the
default truncates and tells you how many it held back.

Read-only. Reads the systemd journal (`-u ssh -u sshd`), falling back to
`/var/log/auth.log` or `/var/log/secure`, and summarises the authentication
activity in the window: the busiest source addresses, the accounts they named,
and every successful login with its method.

Root is required for a reason — both the journal and `auth.log` are root-only
on a stock install, so a non-root run would report zero attempts rather than
fail, which is the most misleading output a security check can produce.

Counts are split along two axes, because a single login attempt writes several
log lines:

| Counter | What it means |
|---|---|
| Rejected credentials | a password or key was offered for a real account and refused |
| Probes and scans | a username that doesn't exist here, or a connection dropped before authenticating |
| Successful logins | someone got in |

The findings are then weighed against the running sshd config, since the same
numbers mean very different things on different servers — a thousand rejected
passwords is background noise where `PasswordAuthentication no` is set, and an
active problem where it isn't. Two findings are called out by name:

- an address that **failed and then succeeded**, which is the shape of a
  brute-force that worked and the one result here that can mean a live
  intrusion
- addresses making 50+ attempts, with a pointer to `serv-kun firewall` for
  rate limiting

#### Distributed attacks

A per-address table can look completely clean during exactly the attack worth
worrying about. Split a run across two thousand hosts and every row reads "3
attempts", no rate limit ever fires, and nothing is out of the ordinary at the
level anyone is looking. So the report also totals across addresses:

- **Volume that no rate limit would have caught** — how many addresses stayed
  under 10 attempts each, and how many attempts they made between them. Twenty
  or more such addresses is reported as a distributed attack, along with the
  two things that do stop one: keys-only authentication, and restricting SSH
  to known source ranges.
- **Credential spraying** — the "most-targeted accounts" table counts distinct
  source addresses per account alongside the attempt total. One account tried
  from ten or more addresses is a coordinated attempt on a login someone
  believes exists here, which reads very differently from one bot retrying.
- **Clustered network blocks** — sources grouped into /24 (IPv4) or /48
  (IPv6). Botnet nodes and rented ranges cluster, so several addresses in one
  block is usually one operator rather than several.

  Blocks are ranked by **attempts per address**, not by how many addresses
  they hold, and the `PER ADDR` column shows it. Address count alone cannot
  tell the two kinds of cluster apart: seven addresses averaging 800 attempts
  each is an operator working on you, while seven averaging 1 is an
  internet-wide sweep that touched your block once on its way through every
  address that exists. They printed as neighbours under the old ordering and
  they call for opposite responses — the sweep's hosts have already moved on,
  so blocking them buys nothing but a permanent firewall rule.

  A deny is therefore suggested only for a block that is both clustered (5+
  addresses) and still working (10+ attempts per address):
  `ufw prepend deny from <block>`. The `prepend` matters: ufw takes the first
  matching rule, so a plain `ufw deny from` is appended *after* the SSH allow
  and never reached. (`ufw insert 1` puts it in the same place but errors on
  an empty rule list, so `prepend` is the form that works either way. It
  needs ufw ≥ 0.36, which is what every currently supported Debian and Ubuntu
  ships.)

Parsing is deliberately tolerant: sshd's wording varies across OpenSSH
releases and distro patches, so an unrecognised line is skipped rather than
aborting the run. Every pattern is anchored at the start of the message, so a
username crafted to look like a log line can't attribute attempts to an
address of the attacker's choosing.

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
lose a server. If the port can't be read, you're told so and asked to confirm
it rather than being shown a guess presented as fact.

#### Rate-limiting SSH

On ufw you're offered `ufw limit` for the SSH rule instead of `ufw allow`: the
port opens as normal, but a source address that starts 6 or more connections
within 30 seconds is denied. The limit is **per source**, so an attacker who
trips it can't lock anyone else out.

It's opt-in because it counts connections, not failures — anything that opens
SSH sessions in bursts from one address can trip it too. Say no if CI deploys
over SSH from a fixed IP, or if a NAT'd office shares one public address.

It's offered on ufw only. firewalld's rich-rule `limit` counts matches against
the rule rather than per source address, so one attacker saturating the budget
would send every other address — yours included — to the default deny. A
lockout is a worse outcome than the noise it would suppress.

Rate limiting slows one busy address. It does nothing about the same attack
split across many addresses, since none of them ever reaches the limit.

#### Restricting SSH to known sources

```
Restrict SSH to specific source addresses? Yes

Connected over SSH right now:
  203.0.113.7
Every one of these must be covered, or its session ends when the policy applies.

Addresses allowed to reach SSH (space-separated CIDR or IP): 203.0.113.0/24
```

This is the only firewall control that stops a distributed attack: it refuses
every address not on the list, however few attempts each one makes. It's also
the fastest way to lose a server, so it defaults to no and is checked before
it's applied.

The check asks the kernel — `ss -Htn state established sport = :<sshport>` —
for every address currently connected on the SSH port, and refuses to accept a
list that leaves any of them out without an explicit opt-in that defaults to
no. `SSH_CLIENT` would be the obvious source for this and is the wrong one:
sudo's default `env_reset` drops it, and this command always runs under sudo.
Asking the kernel also covers a colleague's session rather than only yours.

Include everything that legitimately connects — your office, your VPN, and any
CI runner that deploys — remembering that a NAT'd office appears as one shared
address and that a home connection's address usually changes. If you can't
confirm the full list, answering no leaves SSH open, which is the right
outcome: a firewall nobody can get past is worse than one an attacker still
has to guess a key for.

On ufw each source becomes its own rule (`ufw limit from 203.0.113.0/24 to any
port 22 proto tcp`), since ufw takes one source per rule. On firewalld it
becomes a rich rule *instead of* an open port — the two aren't additive in the
helpful direction, and adding both would leave the port reachable from
everywhere.

Neither of these removes the attack. Keys-only authentication (`ssh harden`)
is what does that — a distributed guess at a key is not feasible — and `ssh
audit` tells you whether anything is getting through.

Database ports are deliberately not opened. In the `springboot` stack Postgres
is never published to the host, so it needs no rule; if you do run Postgres
directly on the host, prefer an SSH tunnel over exposing `5432`.

> **ufw does not filter Docker-published ports.** Docker installs its own
> iptables rules, and traffic to a published container port is forwarded
> before it reaches the `INPUT` chain where ufw's default-deny applies. A
> container run with `-p 5432:5432` is therefore reachable from the internet
> even while `ufw status` reports `Default: deny (incoming)`.
>
> The `springboot` stack is unaffected: it publishes only 80 and 443, which
> are allowed anyway, and Postgres is never published. But don't read
> "default-deny" as covering containers. Bind anything not meant to be public
> to localhost — `ports: ["127.0.0.1:5432:5432"]` — rather than relying on the
> firewall to hide it.

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

The generated compose file names its own project — `name: <directory>-prod`
at the top. Left to default, compose names a project after the directory it
sits in, so a production stack and a dev `docker-compose.yml` in the same
repo end up sharing one container namespace. Compose then matches containers
by project and service name rather than by `container_name`, and an ordinary
`docker compose up -d` against the dev file recreates the *production*
database container as if it were its own; `--remove-orphans` deletes the rest
of the stack. Declaring the name in the file rather than passing `-p` means
hand-typed `docker compose -f docker-compose.prod.yml` commands get it too,
which matters — those are the ones that do the damage.

Pick that name before the first `up` and leave it. The project name prefixes
every volume, so changing it later points the stack at `<newname>_postgres-data`,
starts it on an empty database, and strands the real one under the old prefix.
`init` warns before overwriting a compose file that is already running under a
different name.

```bash
sudo serv-kun springboot setup
```

Both `setup` and `deploy` check first for another compose file in the same
directory that resolves to the same project name — the collision described
above — and name it, along with the services the two files have in common,
before doing anything else. They only report it. The obvious repair is to
rename this project, and that is the one thing not to do blindly: the fix
belongs on the *other* stack, whose volumes are usually disposable. Existing
deployments predate the generated `name:` key, so this is how they find out.

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

If the project mounts a keypair from `certs/`, ownership of it is reasserted
against the image that is about to run — on the way out *and* on the way back.
An image that renumbers its user (a base-image bump is enough) would otherwise
crash-loop on keys chowned for its predecessor; and once the new uid owns them,
a rollback to the earlier image would fail on the same keys, taking down the
one mechanism there to catch a bad deploy.

The container is always recreated, even when the image and config are
unchanged. Left to itself `docker compose up` skips such a service and
reports it as `Running`, which reads as a completed deploy while nothing was
deployed — and the health check that follows would be judging a container
this command never touched.

If the app turns out to be running the same image the pull just produced —
usual when CI has already deployed it, or when it's crash-looping — there is
no earlier image to return to, and you're told exactly that instead of
watching a "rollback" redeploy the identical build. A failure in that state
is in the image or in `.env.prod`, not in the upgrade.

```bash
sudo serv-kun springboot status   # docker compose ps + a health-endpoint check
sudo serv-kun springboot logs     # tail -f the app container's logs
```

## Acknowledgements

serv-kun is built on [**server-wiz**](https://github.com/hardope/server-wiz) by
[hardope](https://github.com/hardope). That's where this project starts from,
not something it happens to resemble.

The premise it's named for is hardope's: that provisioning a Linux box should
be a sequence of answered questions rather than a pile of hand-edited config
files, with each concern getting its own wizard. So is the shape that carries
it — the Go CLI layout, the per-distro package-manager handling, and the
`nginx`, `user`, `postgres` and `node` wizards, which between them still do
most of the work here.

What serv-kun adds on top:

- `ssh harden` — disabling root login and password authentication, with a
  lockout check that refuses to run unless a named non-root user has both a
  working key and usable `sudo`, and drop-in placement that survives
  `50-cloud-init.conf`
- `ssh audit` — a read-only summary of who has been trying to log in, judged
  against the running sshd config, calling out any address that failed and
  then succeeded, and totalling activity across addresses so a distributed
  attack doesn't read as a clean per-address table
- `firewall` — default-deny inbound via ufw/firewalld, reading the real SSH
  port from the running sshd and adding allow rules *before* the default
  policy takes effect, with opt-in per-source rate limiting on SSH and an
  opt-in source allow-list checked against live sessions before it applies
- `springboot` — `init`/`setup`/`deploy`/`status`/`logs` for a Dockerized
  Spring Boot + Postgres stack behind nginx, including Let's Encrypt
  bootstrap and health-checked deploys with automatic rollback

Thanks to hardope for building the thing that was worth extending.

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
