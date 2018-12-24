# quay.io/jonnrb/cloudflare\_dns\_proxy [![Build Status](https://drone.jonnrb.com/api/badges/jon/network_containers/status.svg?branch=master)](https://drone.jonnrb.com/jon/network_containers) [![Docker Repository on Quay](https://quay.io/repository/jonnrb/cloudflare_dns_proxy/status "Docker Repository on Quay")](https://quay.io/repository/jonnrb/cloudflare_dns_proxy)

Uses `cloudflared` to open a DNS-over-HTTPS proxy to `1.1.1.1`.

Designed to be stupid simple:
 1. Run container.
 2. Point DNS clients to that container's port 53.
