#!/bin/bash
set -e
IPROUTE_VERSION=4.9.0
(git clone git://git.kernel.org/pub/scm/linux/kernel/git/shemminger/iproute2.git \
  && cd iproute2 \
  && git config advice.detachedHead false \
  && echo "iproute '$(git tag -l |tail -1)' available" \
  && git checkout v$IPROUTE_VERSION \
  && git archive --prefix iproute2/ HEAD . |gzip -n > ../iproute2.tar.gz \
  && cd .. \
  && rm -rf iproute2)
OPENVPN_VERSION=2.4.5
wget -O openvpn.tar.gz "https://swupdate.openvpn.org/community/releases/openvpn-$OPENVPN_VERSION.tar.gz"
