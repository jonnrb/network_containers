# network\_containers [![Build Status](https://drone.jonnrb.com/api/badges/jon/network_containers/status.svg?branch=master)](https://drone.jonnrb.com/jon/network_containers)

These are some containers I use to abuse Docker networks.

 - [hostapd](./hostapd): Runs [hostapd](https://w1.fi/hostapd/) and can plop
   clients onto a Docker bridge.

 - [router\_base](./router_base): Simple NAT that masquerades traffic from one
   network to another. Also serves as a base image for more custom router-type
   network appliances.

 - [router\_openvpn](./router_base): Builds on router\_base to provide an
   [OpenVPN](https://openvpn.net/index.php/open-source.html) client that
   masquerades a network's traffic.

And these are some one-off containers:

 - [cloudflare\_dns\_proxy](./cloudflare_dns_proxy): Proxies incoming DNS
   questions to Cloudflare's 1.1.1.1 service over TLS.

 - [reverse\_single](./reverse_single): HTTP reverse proxy for a single service.
