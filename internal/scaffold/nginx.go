package scaffold

const nginxConfTemplate = `# Templated by the official nginx entrypoint (envsubst) at container start.
# Only ${APP_DOMAIN} is substituted — nginx's own $host/$remote_addr/etc.
# variables are left alone since they aren't environment variables.

server {
    listen 80;
    server_name ${APP_DOMAIN};

    # Let's Encrypt HTTP-01 challenge
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    http2 on;
    server_name ${APP_DOMAIN};

    ssl_certificate     /etc/letsencrypt/live/${APP_DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${APP_DOMAIN}/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    client_max_body_size {{.MaxBodySize}};

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml image/svg+xml;
    gzip_min_length 1024;

    # Docker's embedded DNS — re-resolves the app container's IP every 30s.
    # Without this, nginx caches the IP from the first resolution and keeps
    # proxying to a dead container after any redeploy that recreates "app"
    # without also restarting nginx.
    resolver 127.0.0.11 valid=30s;

    location / {
        set $upstream_app app;
        proxy_pass http://$upstream_app:{{.AppPort}};
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 60s;
        proxy_send_timeout    120s;
        proxy_read_timeout    120s;
    }
}
`

func renderNginxConf(p Params) (string, error) {
	return render("nginxConf", nginxConfTemplate, p)
}

// HardeningConf is copied verbatim — nothing in it is project-specific.
const HardeningConf = `# Global nginx hardening — loaded as-is (no envsubst), alongside the
# templated server blocks in nginx/templates/app.conf.template.

server_tokens off;

# Trust Cloudflare's proxy IPs so $remote_addr becomes the real visitor IP
# everywhere (logs, X-Forwarded-For, X-Real-IP) instead of Cloudflare's edge
# IP. Only requests actually arriving from these ranges get rewritten, so a
# request hitting the origin directly (bypassing Cloudflare) can't spoof its
# source IP via a header.
#
# List from https://www.cloudflare.com/ips/ — re-check periodically, Cloudflare
# adds ranges occasionally (rare, but this list isn't auto-updating).
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 131.0.72.0/22;
set_real_ip_from 2400:cb00::/32;
set_real_ip_from 2606:4700::/32;
set_real_ip_from 2803:f800::/32;
set_real_ip_from 2405:b500::/32;
set_real_ip_from 2405:8100::/32;
set_real_ip_from 2a06:98c0::/29;
set_real_ip_from 2c0f:f248::/32;
real_ip_header CF-Connecting-IP;
real_ip_recursive on;
`
