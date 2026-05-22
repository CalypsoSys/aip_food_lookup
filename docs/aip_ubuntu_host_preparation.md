# AIP Food Lookup Ubuntu host preparation

This is the one-time Ubuntu host bootstrap checklist for the lab machine that will run the AIP Food Lookup API stack in
Docker.

Routine deployment and refresh steps belong in:

- [aip_production_runbook.md](aip_production_runbook.md)

## Goal

Prepare the Ubuntu host to run:

- the Dockerized `aip-food-lookup-api` Go API
- file-backed AIP food data and runtime suggestion/feedback files
- `cloudflared` for the tunnel-backed API origin
- host-installed Caddy as the local reverse proxy behind Cloudflare Tunnel
- the rendered-config deployment wrapper used by the AIP Docker stack

Use `aip-food-lookup` as the server-side identity for directories, stack names, logs, and service names.

## 1. Confirm baseline host details

Run on the Ubuntu host:

```bash
uname -a
lsblk
free -h
ip addr
```

Recommended baseline:

- Ubuntu Server
- `x86_64` / `amd64`
- enough disk for Docker images, food data, feedback JSONL files, logs, and backups
- enough RAM for the Go API and the host services

## 2. Install official Docker Engine

Run on the Ubuntu host:

```bash
sudo apt update
sudo apt install -y ca-certificates curl smartmontools

sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

sudo tee /etc/apt/sources.list.d/docker.sources >/dev/null <<EOF_DOCKER
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF_DOCKER

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl enable --now docker
sudo docker run hello-world
```

## 3. Allow the deployment user to run Docker

Run on the Ubuntu host:

```bash
sudo usermod -aG docker $USER
newgrp docker
docker version
docker compose version
systemctl status docker --no-pager
```

## 4. Set up SSH key access

If needed, create a key on the Windows workstation in PowerShell:

```powershell
$serverIp = "replace_with_server_ip"
$serverUser = "replace_with_server_user"

ssh-keygen -t ed25519 -C "$serverUser@aip-food-lookup-lab"
```

Copy the public key to the Ubuntu host:

```powershell
type $env:USERPROFILE\.ssh\id_ed25519.pub | ssh ${serverUser}@${serverIp} "mkdir -p ~/.ssh && chmod 700 ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
```

Verify login:

```powershell
ssh ${serverUser}@${serverIp}
```

## 5. Harden SSH server settings

Edit:

```bash
sudo vi /etc/ssh/sshd_config
```

Recommended settings:

```text
PermitRootLogin no
PubkeyAuthentication yes
PasswordAuthentication no
KbdInteractiveAuthentication no
```

Reload SSH:

```bash
sudo systemctl restart ssh
sudo systemctl status ssh --no-pager
```

## 6. Enable firewall baseline

Run on the Ubuntu host:

```bash
sudo apt update
sudo apt install -y ufw
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status verbose
```

If the API will only be reached through Cloudflare Tunnel and no direct ingress will be used, you can tighten this
later.

## 7. Enable automatic updates

Run on the Ubuntu host:

```bash
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

## 8. Disable lid-close sleep behavior if this is laptop hardware

Edit:

```bash
sudo vi /etc/systemd/logind.conf
```

Recommended:

```text
HandleLidSwitch=ignore
HandleLidSwitchExternalPower=ignore
HandleLidSwitchDocked=ignore
```

Apply:

```bash
sudo systemctl restart systemd-logind
sudo systemd-analyze cat-config systemd/logind.conf | grep -E 'HandleLidSwitch|HandleLidSwitchExternalPower|HandleLidSwitchDocked'
```

## 9. Install cloudflared

If this host will publish the API through Cloudflare Tunnel, install `cloudflared`:

```bash
curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null
echo 'deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared any main' | sudo tee /etc/apt/sources.list.d/cloudflared.list
sudo apt update
sudo apt install -y cloudflared
cloudflared --version
```

Then install the tunnel as a systemd service using the token from Cloudflare Zero Trust. Treat the tunnel token as a
secret: do not commit it, paste it into shared notes, or leave it in reusable scripts.

```bash
sudo cloudflared service install <paste-tunnel-token-from-cloudflare>
sudo systemctl enable --now cloudflared
systemctl is-active --quiet cloudflared && echo "cloudflared running"
```

If a direct service status check reports `Unit cloudflared.service could not be found`, the package may be installed but
the tunnel service has not been installed yet. Run the service install step above before continuing.

In the recommended steady state, Cloudflare Tunnel fronts Caddy rather than pointing directly at the API container.

## 10. Prepare host-installed Caddy

Prepare the long-term local ingress layer behind Cloudflare Tunnel.

Recommended request path:

- browser or Flutter client -> `https://hashimojoe.com/api`
- `/api/*` -> Cloudflare Pages Functions
- Pages Functions -> Cloudflare Tunnel hostname for the API origin
- Cloudflare Tunnel -> host-installed Caddy on the Ubuntu host
- Caddy -> `aip-food-lookup-api`

Recommended server layout:

```text
/srv/logs/caddy
  caddy.log

/etc/caddy
  Caddyfile
```

Install Caddy:

```bash
sudo apt update
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install -y caddy
```

Prepare the host log directory:

```bash
sudo mkdir -p /srv/logs/caddy
sudo chown -R caddy:caddy /srv/logs/caddy
sudo chmod 755 /srv/logs/caddy
```

Recommended Caddyfile:

```caddy
{
    auto_https off

    log {
        output file /srv/logs/caddy/caddy.log
        format console
    }
}

http://api.hashimojoe.com {
    reverse_proxy 127.0.0.1:8084
}
```

Validate and start Caddy:

```bash
sudo vi /etc/caddy/Caddyfile
sudo caddy validate --config /etc/caddy/Caddyfile
sudo systemctl enable --now caddy
sudo systemctl restart caddy
sudo systemctl status caddy --no-pager
```

Install the checked-in Caddy logrotate policy after this repo is available on the host:

```bash
export AIP_REPO_ROOT=/absolute/path/to/aip_food_lookup
sudo cp "$AIP_REPO_ROOT/scripts/caddy/caddy.logrotate" /etc/logrotate.d/caddy
sudo chmod 644 /etc/logrotate.d/caddy
sudo logrotate -d /etc/logrotate.d/caddy
```

See [aip_caddy_host_setup.md](aip_caddy_host_setup.md) for the dedicated Caddy setup and validation guide.

## 11. Verify prepared host state

Run on the Ubuntu host:

```bash
docker version
docker compose version
systemctl status docker --no-pager
systemctl status ssh --no-pager
systemctl status caddy --no-pager
systemctl is-active --quiet cloudflared && echo "cloudflared running"
sudo ufw status verbose
cloudflared --version
sudo caddy validate --config /etc/caddy/Caddyfile
```

## 12. Continue with the production runbook

The remaining AIP-specific deployment steps belong in:

- [aip_production_runbook.md](aip_production_runbook.md)

That runbook covers:

- server directory layout
- rendered-config binary placement
- internal runtime ports
- Docker stack bring-up
- Cloudflare wiring and validation
- file-data backup

## Deployment direction

Recommended operating model:

- build the API image locally in WSL/Linux
- stage artifacts through `C:\transfer`
- copy artifacts and stack files to the Ubuntu host
- keep runtime secrets and server-local config on the host
- run the API in Docker on the host
- run host-installed Caddy as the local reverse proxy
- run `cloudflared` on the host as the public tunnel into Caddy

That keeps the lab box as the runtime target without turning it into the primary build machine.
