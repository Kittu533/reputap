#!/usr/bin/env bash
set -euo pipefail

echo "==========================================="
echo " Reputap VPS Bootstrap (Ubuntu / Debian)   "
echo " Native Go + PostgreSQL + Caddy (Auto-SSL) "
echo "==========================================="

if [ "$EUID" -ne 0 ]; then
  echo "Error: Jalankan script ini sebagai root (sudo -i)"
  exit 1
fi

DB_PASS="${1:-secret_reputap_2026}"
DOMAIN="${2:-api.kitangoding.my.id}"

# 1. Update & dependencies
apt-get update -y
apt-get install -y curl gnupg2 lsb-release debian-keyring debian-archive-keyring apt-transport-https ufw

# 2. Install PostgreSQL
apt-get install -y postgresql postgresql-contrib

systemctl enable postgresql
systemctl start postgresql

# Configure DB
sudo -u postgres psql -c "DO \$\$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'reputap') THEN CREATE ROLE reputap WITH LOGIN PASSWORD '${DB_PASS}'; END IF; END \$\$;"
sudo -u postgres psql -c "SELECT 'CREATE DATABASE reputap OWNER reputap' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'reputap')\gexec"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE reputap TO reputap;"

# 3. Install Caddy
curl -1sLF 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor --yes -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLF 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt-get update -y
apt-get install -y caddy

# 4. Setup directories
mkdir -p /var/www/reputap/migrations
mkdir -p /var/www/reputap/uploads

# 5. Default .env if not exists
if [ ! -f /var/www/reputap/.env ]; then
  JWT_SECRET=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 48)
  cat <<EOF > /var/www/reputap/.env
APP_PORT=8081
APP_ENV=production

DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=reputap
DB_USER=reputap
DB_PASSWORD=${DB_PASS}
DB_SSLMODE=disable

JWT_SECRET=${JWT_SECRET}
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

APP_BASE_URL=https://${DOMAIN}
REDIRECT_BASE_URL=https://${DOMAIN}/r
FRONTEND_URL=https://kitangoding.my.id

UPLOAD_DIR=/var/www/reputap/uploads
MAX_UPLOAD_SIZE=10485760
EOF
  echo "Created default /var/www/reputap/.env"
fi

# 6. Setup Caddyfile
cat <<EOF > /etc/caddy/Caddyfile
${DOMAIN} {
    reverse_proxy 127.0.0.1:8081
}
EOF

systemctl enable caddy
systemctl restart caddy

# 7. Setup systemd service
cat <<EOF > /etc/systemd/system/reputap.service
[Unit]
Description=Reputap Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/var/www/reputap
ExecStart=/var/www/reputap/server
Restart=always
RestartSec=3s
EnvironmentFile=/var/www/reputap/.env
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable reputap

# 8. Firewall (UFW)
ufw allow 22/tcp || true
ufw allow 80/tcp || true
ufw allow 443/tcp || true
ufw --force enable || true

echo "==========================================="
echo " VPS Setup Selesai!"
echo " Domain: https://${DOMAIN}"
echo " Path: /var/www/reputap"
echo " Systemd service 'reputap' siap dijalankan setelah binary ter-upload."
echo "==========================================="
