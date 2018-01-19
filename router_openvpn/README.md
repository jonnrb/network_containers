# jonnrb/router\_openvpn [![Docker Automated Build](https://img.shields.io/docker/automated/jonnrb/router_openvpn.svg)](https://hub.docker.com/r/jonnrb/router_openvpn/) [![Docker Build Status](https://img.shields.io/docker/build/jonnrb/router_openvpn.svg)](https://hub.docker.com/r/jonnrb/router_openvpn/)

### Usage

You'll need an uplink network, for OpenVPN to tunnel through, created somewhat
like

```bash
docker network create net_uplink
```

but any named Docker network that can connect to the internet should do.

You'll also need an internal network set up a bit more carefully:

```bash
docker network create net_vpn \
  -o com.docker.network.bridge.enable_ip_masquerade=false \
  -o com.docker.network.bridge.enable_icc=false \
  --subnet=10.55.55.0/24 --gateway=10.55.55.2 \
  --aux-address DefaultGatewayIPv4=10.55.55.1
```

This creates a bridge network but disables communication to other Docker
containers via the gateway IP (10.55.55.2) on the Docker host. It also uses the
secret `DefaultGatewayIPv4` option that works on bridge networks to set the
default route to something _other_ than the Docker host's gateway IP.

Finally, the router:

```bash
docker create --name router --cap-add NET_ADMIN \
  -v /path/to/openvpn.conf:/data/openvpn.conf:ro \
  -v /path/to/keys/maybe:/data/keys:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  jonnrb/router_openvpn -logtostderr -v 2 \
  -create_tun vpntun -docker.uplink_interface vpntun \
  -docker.lan_network net_vpn

docker network connect net_uplink router
docker network connect net_vpn router

docker start -a router
```

Your `openvpn.conf` will need to reference the tunnel `vpntun` or you can live
on the edge and change the name given to `-create_tun` and
`-docker.uplink_interface`.

(The container needs `CAP_NET_ADMIN` for iptables and the raw sockets ping used
for the healthcheck.)
