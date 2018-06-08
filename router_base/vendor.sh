#!/bin/bash
# Pull a version of iptables into the repo.
set -e
VERSION=1.4.21
curl -q -L -o iptables.tar.bz2 \
  ftp://ftp.netfilter.org/pub/iptables/iptables-$VERSION.tar.bz2
