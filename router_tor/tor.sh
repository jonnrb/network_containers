#!/bin/sh

docker-gen /data/netresolve.tmpl /data/nets.all

local_net=$(grep "${TOR_NET}" /data/nets.all |cut -d ' ' -f 2)
trans_port="9040"
dns_port="5353"
virt_addr="10.192.0.0/10"

# redirect nicey tor traffic
iptables -t nat -A PREROUTING -s "${local_net}" -p udp -m udp --dport 53 -j REDIRECT --to-ports "${dns_port}"
iptables -t nat -A PREROUTING -s "${local_net}" -p udp -m udp --dport 5353 -j REDIRECT --to-ports "${dns_port}"
iptables -t nat -A PREROUTING -s "${local_net}" -d "${virt_addr}" -p tcp --syn -j REDIRECT --to-ports "${trans_port}"


exec tor -f /data/torrc
