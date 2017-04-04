#!/bin/sh

echo "[init.sh] starting router_base"

# iptables rules cannot be just iptables-restored thanks to docker's dns server
if [ -f /data/iptables.tmpl ]; then
	echo "[init.sh] adding iptables rules"
	# strip comments and empty lines
	cat /data/iptables.tmpl |grep -Ev '^\w*(\#|$)' |while read rule; do
		# add rules in order from rule template file
		iptables ${rule}
	done
	# TODO: cry about docker mangling my networking. i am a but a router
else
	echo "[init.sh] no iptables template found"
fi

# hi. run me with privileged=true or i will say nasty error. ty
sysctl -w net.ipv4.ip_forward=1

if [ -n "$*" ]; then
	echo "[init.sh] running $*"
	exec "$@"
else
	echo "[init.sh] no docker cmd; sleeping forever"
	trap exit SIGINT SIGHUP
	while true; do sleep 1; done
fi
