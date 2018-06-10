from ubuntu:trusty as build

run apt-get update \
 # For OpenVPN
 && apt-get install -y build-essential libssl-dev liblz4-dev liblzo2-dev \
 # For iproute2
 && apt-get install -y pkg-config bison flex

workdir /
copy openvpn.tar.gz iproute2.tar.gz /
run set -x \
 # Build OpenVPN
 && tar xvf openvpn.tar.gz \
 && cd openvpn* \
 && ./configure --disable-shared --enable-static --disable-plugin-auth-pam \
        --disable-plugin-down-root --enable-iproute2 --disable-server \
 && make LIBS="-all-static" \
 && make install \
 && cd .. \\
 # Build iproute2
 && tar xvf iproute2.tar.gz \
 && cd iproute2 \
 && ./configure \
 && make \
 && make install

from jonnrb/router_base
copy --from=build /usr/local/sbin/openvpn /sbin/openvpn
copy --from=build /etc/iproute2 /etc/iproute2
copy --from=build /sbin/ip /sbin/ip
entrypoint ["/init", "-c", "openvpn /data/openvpn.conf"]
