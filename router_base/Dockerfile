from alpine:latest

# Get required packages
run apk update && apk add iptables

add iptables.tmpl /data/iptables.tmpl
add init.sh /

entrypoint ["/init.sh"]
