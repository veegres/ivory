#!/bin/sh
set -e

# The ssh account is created at runtime rather than baked in, so one image
# serves whatever username the Ivory vault holds.
: "${SSH_USER:=ivory}"
: "${SSH_PUBLIC_KEY:=}"

if [ "$SSH_USER" = "root" ]; then HOME_DIR=/root; else HOME_DIR="/home/$SSH_USER"; fi

if [ -n "$SSH_PUBLIC_KEY" ]; then
  if [ "$SSH_USER" != "root" ] && ! id "$SSH_USER" >/dev/null 2>&1; then
    adduser -D -s /bin/sh "$SSH_USER"
    addgroup "$SSH_USER" root
  fi
  # An alpine account created without a password is locked ("!" in /etc/shadow)
  # and sshd refuses public key auth for a locked account, with the same
  # "no supported methods remain" the client reports for a missing key.
  sed -i "s/^$SSH_USER:!:/$SSH_USER:*:/" /etc/shadow
  mkdir -p "$HOME_DIR/.ssh"
  printf '%s\n' "$SSH_PUBLIC_KEY" > "$HOME_DIR/.ssh/authorized_keys"
  chown -R "$SSH_USER" "$HOME_DIR/.ssh"
  chmod 700 "$HOME_DIR/.ssh"
  chmod 600 "$HOME_DIR/.ssh/authorized_keys"
fi

ssh-keygen -A
sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config
/usr/sbin/sshd

# The ssh account is not root and dockerd recreates the socket on every start.
( while [ ! -S /var/run/docker.sock ]; do sleep 1; done; chmod 666 /var/run/docker.sock ) &

exec dockerd-entrypoint.sh dockerd --host=unix:///var/run/docker.sock
