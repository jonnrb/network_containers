from alpine:3.12.0 as build
workdir /
run apk add --no-cache gcc libc-dev openssl-dev lz4-dev lzo-dev libmnl-dev \
  linux-headers bash bison flex make pkgconfig 

from build as iproute2_build
env IPROUTE2_SRC https://git.kernel.org/pub/scm/network/iproute2/iproute2.git/snapshot/iproute2-5.7.0.tar.gz
env IPROUTE2_HASH 12a3861f463c6bbd1bb3b213ac734f75c89172b74104140dd0bbfcb1e13ee798
add $IPROUTE2_SRC /iproute2.tar.gz
run echo "$IPROUTE2_HASH  iproute2.tar.gz" |sha256sum -c - \
 && tar xf iproute2.tar.gz \
 && cd iproute2-* \
 && (echo "#include <signal.h>"; cat devlink/devlink.c) >/tmp/devlink.c \
 && mv /tmp/devlink.c devlink/devlink.c \
 && ./configure \
 && make \
 && make install

from build as openvpn_build
env OPENVPN_SRC https://swupdate.openvpn.net/community/releases/openvpn-2.4.9.tar.gz
env OPENVPN_HASH 46b268ef88e67ca6de2e9f19943eb9e5ac8544e55f5c1f3af677298d03e64b6e
add $OPENVPN_SRC /openvpn.tar.gz
run echo "$OPENVPN_HASH  openvpn.tar.gz" |sha256sum -c - \
 && tar xf openvpn.tar.gz \
 && cd openvpn* \
 && ./configure --disable-plugin-auth-pam --enable-iproute2 --disable-server \
 && make \
 && make install

from gcr.io/distroless/base
copy --from=iproute2_build /etc/iproute2 /etc/iproute2
copy --from=iproute2_build /sbin/ip /sbin/ip
copy --from=iproute2_build /usr/lib/libmnl.so* /lib/
copy --from=openvpn_build /usr/local/sbin/openvpn /sbin/openvpn
copy --from=openvpn_build /lib/libc* /lib/ld* /lib/libcrypto.so* \
     /lib/libssl.so* /usr/lib/liblzo2.so* /usr/lib/liblz4.so* /lib/libz.so* \
     /lib/
entrypoint ["/sbin/openvpn"]
cmd ["/data/openvpn.conf"]
