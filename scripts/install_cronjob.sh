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

(echo "0 * * * * bash -lc '/var/lib/ssh-manager/ssh-manager > /var/lib/ssh-manager/execution.log'") | crontab -
