#!/bin/sh

rm -f /etc/dhcp/dhcpd.conf
docker-gen /etc/dhcp/dhcpd.conf.tmpl /etc/dhcp/dhcpd.conf

exec /usr/sbin/dhcpd -4 -f -d --no-pid -cf /etc/dhcp/dhcpd.conf
