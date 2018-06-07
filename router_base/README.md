# [jonnrb/router\_base](https://hub.docker.com/r/jonnrb/router_base) [![Build Status](https://drone.jonnrb.com/api/badges/jon/network_containers/status.svg?branch=master)](https://drone.jonnrb.com/jon/network_containers)

Base image for creating a more interesting Docker router. (Can be used by
itself too.) This image provides an entrypoint (init program) that sets up
routing rules between Docker networks within the container's network namespace.

### Details

The entrypoint uses the Docker socket to read the networking configuration so
this must be mounted to `/var/run/docker.sock` or the usual environment
variables should be set.

Docker doesn't really support containers acting as gateways for its networks,
so some trickery (_hacks_) are involved and some special configuration is
needed depending on what you're trying to do. The way I see it, there are two
scenarios where you'd want a Docker container running as a basic NAT (and
Docker managing everything won't cut it):

 1. You have a mixed network where Docker containers and other endpoints are
    supposed to be happy and living together. You might want to run a
    [DHCP server](https://github.com/JonNRb/etcdhcp) on this network to handle
    the other endpoints. In this case, it's a pretty good assumption you are
    (or want to be) using the `macvlan` network driver. You also might want
    to use some other uplink (VPN) for this network because you value privacy.

 2. You want a network of all containers not connecting to the internet via
    your usual uplink. You might want some privacy by using a VPN or Tor or
    some other private network.

### Usage

You'll need an uplink network created something like

```bash
docker network create net_uplink
```

but any named Docker network that can connect to the internet should do.

You'll also need an internal network set up a bit more carefully:

```bash
docker network create net_internal \
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
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  jonnrb/router_base -logtostderr -v 2 \
  -docker.uplink_network net_uplink -docker.lan_network net_internal

docker network connect net_uplink router
docker network connect net_internal router

docker start -a router
```

(The container needs `CAP_NET_ADMIN` for iptables and the raw sockets ping used
for the healthcheck.)
