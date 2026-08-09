package scaffold

const deployDocTemplate = `# Deploying {{.AppName}} on Hetzner

Stack: Spring Boot app (Docker) + self-hosted Postgres (Docker) + nginx/certbot
for TLS, with Cloudflare in front for DNS/CDN/DDoS protection. Provisioning
and deploys are driven by [serv-kun](https://github.com/lahat-m/serv-kun)'s
` + "`springboot`" + ` wizard.

## 1. Provision the server

1. Create a Hetzner Cloud VPS (a CX22/CX32-class box is plenty to start —
   scale up later if needed). Ubuntu 24.04 recommended.
2. Add a Hetzner Cloud Firewall allowing only ` + "`22`" + ` (SSH), ` + "`80`" + `, ` + "`443`" + ` inbound.
   Postgres is never published to the host, so no firewall rule is needed for it.
3. Install serv-kun:
   ` + "```" + `
   curl -fsSL https://raw.githubusercontent.com/lahat-m/serv-kun/main/install.sh | sudo bash
   ` + "```" + `
4. In Cloudflare: add an A record for your domain pointing at the VPS IP,
   proxied (orange cloud). Set SSL/TLS mode to **Full (strict)** — the
   Let's Encrypt certificate nginx obtains satisfies that.

## 2. Set up the app directory

` + "```" + `
git clone <your-repo-url> /opt/{{.AppName}}
cd /opt/{{.AppName}}
git checkout main
` + "```" + `

## 3. First deploy

` + "```" + `
sudo serv-kun springboot setup
` + "```" + `

This installs Docker if it isn't already present, then prompts for every
` + "`.env.prod`" + ` field this project needs — domain, image, Postgres
credentials, CORS origin always; JWT keys, an admin account, email, and
object storage are each asked as a yes/no first, so declining one skips its
prompts entirely. It generates a JWT RSA keypair into ` + "`certs/`" + ` if
enabled, verifies DNS, requests the first Let's Encrypt certificate, brings
the whole stack up, and polls ` + "`/actuator/health`" + ` before declaring
success.

## 4. Continuous deployment (GitHub Actions)

` + "`.github/workflows/deploy.yml`" + ` builds the image, pushes it to GHCR,
then SSHes into the server and runs ` + "`sudo serv-kun springboot deploy`" + `
— it pulls the latest ` + "`app`" + ` image, recreates just that one
container, health-checks it, and **automatically rolls back** to the
previous image if the new one fails its health check. Set these repo
secrets:

- ` + "`HETZNER_HOST`" + ` — server IP. Prefer the IP over a hostname: a name
  proxied through Cloudflare resolves to Cloudflare's edge, and SSH never
  reaches the server.
- ` + "`HETZNER_SSH_USER`" + ` — e.g. ` + "`deploy`" + `
- ` + "`HETZNER_SSH_KEY`" + ` — private key for that user, whole file including the
  ` + "`-----BEGIN`" + `/` + "`-----END`" + ` lines (add the matching public key to the
  server's ` + "`~/.ssh/authorized_keys`" + `)
- ` + "`HETZNER_APP_DIR`" + ` — e.g. ` + "`/opt/{{.AppName}}`" + `

A secret you don't set is substituted as an empty string rather than failing,
so the workflow checks all four before connecting and names any that are
blank.

Then, on the server, grant that user passwordless sudo for the two commands
the deploy runs:

` + "```" + `
sudo serv-kun ci
` + "```" + `

Without it the deploy fails at ` + "`sudo`" + ` with *"a terminal is required to
authenticate"* — an SSH session has no terminal for sudo to prompt on. The
rule it writes covers ` + "`serv-kun`" + ` and ` + "`docker image prune -f`" + ` and
nothing else; every other ` + "`sudo`" + ` still asks for a password.

## 5. Day-to-day operations

` + "```" + `
sudo serv-kun springboot status   # docker compose ps + a health-endpoint check
sudo serv-kun springboot logs     # tail -f the app container's logs
sudo serv-kun springboot deploy   # manual redeploy, same rollback behavior as CI
` + "```" + `

## 6. Backups

The ` + "`backup`" + ` service (` + "`scripts/backup-postgres.sh`" + `) runs
` + "`pg_dump`" + ` daily into the ` + "`postgres-backups`" + ` volume and
prunes anything older than ` + "`BACKUP_KEEP_DAYS`" + ` (default 7). That
protects against app-level mistakes but not against losing the VPS itself —
for offsite copies, sync the volume to your object storage provider
periodically, e.g. via a host cron job running ` + "`docker run --rm -v {{.AppName}}_postgres-backups:/backups:ro ... s3 sync /backups s3://your-bucket`" + `.

**Restore:**
` + "```" + `
gunzip -c /backups/<db>_<timestamp>.sql.gz | \
  docker compose -f docker-compose.prod.yml --env-file .env.prod exec -T postgres psql -U <user> -d <db>
` + "```" + `
`

func renderDeployDoc(p Params) (string, error) {
	return render("deployDoc", deployDocTemplate, p)
}
