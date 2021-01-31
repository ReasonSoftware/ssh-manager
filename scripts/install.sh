#!/usr/bin/env bash
set -e

rm -rf /var/lib/ssh-manager
mkdir -p /var/lib/ssh-manager

wget $(curl -si https://api.github.com/repos/ReasonSoftware/ssh-manager/releases/latest | \
        grep browser_download_url | \
        awk -F': ' '{print $2}' | \
        tr -d '"') -O /var/lib/ssh-manager/ssh-manager.zip

unzip -j /var/lib/ssh-manager/ssh-manager.zip -d /var/lib/ssh-manager
rm -f /var/lib/ssh-manager/ssh-manager.zip

SERVICE=$(cat <<-EOF
[Unit]
Description=Central SSH Management Service
Wants=network-online.target
After=network-online.target

[Service]
Type=oneshot
ExecStart=/var/lib/ssh-manager/ssh-manager
StandardOutput=journal
User=root

[Install]
WantedBy=multi-user.target
EOF
)

echo "$SERVICE" > /etc/systemd/system/ssh-manager.service

TIMER=$(cat <<-EOF
[Unit]
Description=Central SSH Management Service
Wants=network-online.target
After=network-online.target

[Timer]
Unit=ssh-manager.service
OnBootSec=10min
OnUnitInactiveSec=60min
Persistent=true

[Install]
WantedBy=multi-user.target
EOF
)

echo "$TIMER" > /etc/systemd/system/ssh-manager.timer

systemctl daemon-reload
systemctl enable ssh-manager.service
systemctl enable --now ssh-manager.timer