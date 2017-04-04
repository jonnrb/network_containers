from jonnrb/router_base

# Get required packages
run apk update && apk add openvpn

add iptables.tmpl /data/

cmd ["openvpn", "/data/openvpn.conf"]
