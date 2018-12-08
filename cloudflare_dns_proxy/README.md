# cloudflare_dns_proxy [![Build Status](https://drone.jonnrb.com/api/badges/jon/network_containers/status.svg?branch=master)](https://drone.jonnrb.com/jon/network_containers)

Uses `cloudflared` to open a DNS-over-HTTPS proxy to `1.1.1.1`.

Designed to be stupid simple:
 1. Run container.
 2. Point DNS clients to that container's port 53.
