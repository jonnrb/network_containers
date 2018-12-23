from alpine:3.8 as build
workdir /
run apk add --no-cache gcc libc-dev openssl-dev lz4-dev lzo-dev libmnl-dev \
  linux-headers bash bison flex make pkgconfig 

from build as iproute2_build
env IPROUTE2_SRC https://git.kernel.org/pub/scm/network/iproute2/iproute2.git/snapshot/iproute2-4.19.0.tar.gz
env IPROUTE2_HASH 1f31fc30d030c6e532ff6a59f7bbe3a6fa929f790c600f05c966f358f693a0c7
add $IPROUTE2_SRC /iproute2.tar.gz
run echo "$IPROUTE2_HASH  iproute2.tar.gz" |sha256sum -c - \
 && tar xf iproute2.tar.gz \
 && cd iproute2-* \
 && ./configure \
 && make \
 && make install

from build as openvpn_build
env OPENVPN_SRC https://swupdate.openvpn.net/community/releases/openvpn-2.4.6.tar.gz
env OPENVPN_HASH 738dbd37fcf8eb9382c53628db22258c41ba9550165519d9200e8bebaef4cbe2
add $OPENVPN_SRC /openvpn.tar.gz
run echo "$OPENVPN_HASH  openvpn.tar.gz" |sha256sum -c - \
 && tar xf openvpn.tar.gz \
 && cd openvpn* \
 && ./configure --disable-plugin-auth-pam --enable-iproute2 --disable-server \
 && make \
 && make install

from jonnrb/router_base
copy --from=iproute2_build /etc/iproute2 /etc/iproute2
copy --from=iproute2_build /sbin/ip /sbin/ip
copy --from=iproute2_build /usr/lib/libmnl.so* /lib/
copy --from=openvpn_build /usr/local/sbin/openvpn /sbin/openvpn
copy --from=openvpn_build /lib/libcrypto.so* /lib/libssl.so* \
     /usr/lib/liblzo2.so* /usr/lib/liblz4.so* /lib/libz.so* /lib/
entrypoint ["/init", "-c", "openvpn /data/openvpn.conf"]
