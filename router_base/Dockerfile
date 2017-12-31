from alpine:latest

# Get required packages
run apk update && apk add iptables

add iptables.tmpl /data/iptables.tmpl
add init.sh /

healthcheck --timeout=3s --interval=15s --retries=4 \
  cmd ["ping", "-c", "1", "8.8.8.8"]

entrypoint ["/init.sh"]
