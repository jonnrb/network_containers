#!/bin/bash
set -e
IPROUTE_VERSION=4.9.0
(git clone git://git.kernel.org/pub/scm/linux/kernel/git/shemminger/iproute2.git \
  && cd iproute2 \
  && git checkout v$IPROUTE_VERSION \
  && git archive --prefix iproute2/ HEAD . |gzip > ../iproute2.tar.gz \
  && cd .. \
  && rm -rf iproute2)
OPENVPN_VERSION=2.4.4
wget -O openvpn.tar.gz "https://swupdate.openvpn.org/community/releases/openvpn-$OPENVPN_VERSION.tar.gz"
