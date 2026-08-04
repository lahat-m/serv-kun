package scaffold

// BackupScript is copied verbatim — it's driven entirely by env vars the
// compose file already injects (POSTGRES_HOST/DB/USER/PASSWORD,
// BACKUP_KEEP_DAYS), so nothing here is project-specific.
const BackupScript = `#!/bin/sh
# Nightly pg_dump loop for the "backup" sidecar in docker-compose.prod.yml.
# Runs inside the same postgres:18-alpine image as the DB, so pg_dump always
# matches the server version. Dumps go to the postgres-backups volume; sync
# them offsite periodically if you want protection against losing the VPS
# itself, not just app-level mistakes.
set -eu

# Without pipefail, "pg_dump | gzip" reports gzip's exit status — so a failed
# dump still looks successful and leaves behind a valid, empty .gz. BusyBox
# ash (this image's /bin/sh) and bash both support it; the guard keeps the
# script working on a shell that doesn't.
if (set -o pipefail) 2>/dev/null; then
  set -o pipefail
fi

BACKUP_DIR=/backups
KEEP_DAYS="${BACKUP_KEEP_DAYS:-7}"

mkdir -p "$BACKUP_DIR"

while true; do
  timestamp=$(date +%Y%m%d_%H%M%S)
  file="$BACKUP_DIR/${POSTGRES_DB}_${timestamp}.sql.gz"
  backup_ok=0

  echo "[$(date -Iseconds)] backing up $POSTGRES_DB to $file"
  if PGPASSWORD="$POSTGRES_PASSWORD" pg_dump --clean --if-exists --no-owner \
       -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" | gzip > "$file"; then
    # Verify the archive in a single pass: under pipefail a corrupt archive
    # fails the pipeline, and a zero byte count catches a dump that produced
    # nothing (even an empty database writes SQL headers). A separate gzip -t
    # would decompress the whole file a second time for no extra signal.
    if bytes=$(gzip -dc "$file" | wc -c) && [ "$bytes" -gt 0 ]; then
      backup_ok=1
      echo "[$(date -Iseconds)] backup ok ($(wc -c < "$file") bytes compressed, $bytes uncompressed)"
    else
      echo "[$(date -Iseconds)] backup FAILED verification — discarding" >&2
      rm -f "$file"
    fi
  else
    echo "[$(date -Iseconds)] backup FAILED" >&2
    rm -f "$file"
  fi

  # Only prune after a good dump. Otherwise a run of failures would quietly
  # age out the last known-good backup and leave nothing behind.
  if [ "$backup_ok" -eq 1 ]; then
    find "$BACKUP_DIR" -name '*.sql.gz' -mtime "+$KEEP_DAYS" -delete
  else
    echo "[$(date -Iseconds)] skipping prune — keeping older backups" >&2
  fi

  sleep 86400
done
`

// DeployWorkflow is copied verbatim — it uses ${{ github.repository }} and
// ${{ github.actor }} rather than any hardcoded project name, so it's
// already reusable as-is.
const DeployWorkflow = `name: Build and deploy

on:
  push:
    branches: [main]

permissions:
  contents: read
  packages: write

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract image metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/${{ github.repository }}
          tags: |
            type=sha,format=long
            type=raw,value=latest

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    steps:
      - name: Deploy over SSH
        uses: appleboy/ssh-action@v1.2.0
        with:
          host: ${{ secrets.HETZNER_HOST }}
          username: ${{ secrets.HETZNER_SSH_USER }}
          key: ${{ secrets.HETZNER_SSH_KEY }}
          script: |
            cd ${{ secrets.HETZNER_APP_DIR }}
            git pull --ff-only
            sudo serv-kun springboot deploy
            docker image prune -f
`
