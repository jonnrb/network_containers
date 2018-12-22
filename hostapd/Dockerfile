from alpine:3.8 as build
env HOSTAPD_SRC https://w1.fi/releases/hostapd-2.7.tar.gz
env HOSTAPD_HASH 21b0dda3cc3abe75849437f6b9746da461f88f0ea49dd621216936f87440a141
add $HOSTAPD_SRC hostapd.tar.gz
run apk add --no-cache alpine-sdk git iw linux-headers libnl3-dev openssl-dev \
 && echo "$HOSTAPD_HASH  hostapd.tar.gz" |sha256sum -c - \
 && tar xf hostapd.tar.gz \
 && cd hostapd-2.7 \
 && cd hostapd \
 && cp defconfig .config \
 && sed -i '/CONFIG_ACS=y/s/^#//' .config \
 && sed -i '/CONFIG_TLSV11=y/s/^#//' .config \
 && sed -i '/CONFIG_TLSV12=y/s/^#//' .config \
 && sed -i '/CONFIG_IEEE80211AC=y/s/^#//' .config \
 && sed -i '/CONFIG_IEEE80211N=y/s/^#//' .config \
 && make \
 && make install

from gcr.io/distroless/static
copy --from=build /usr/local/bin/hostapd* /usr/local/bin/
copy --from=build /lib/libc* /lib/ld* /usr/lib/libnl* /lib/libcrypto* \
     /lib/libssl* /lib/libz* /lib/
copy --from=build /usr/sbin/iw /usr/sbin/
add init.sh /
entrypoint ["/usr/local/bin/hostapd"]
cmd ["/etc/hostapd.conf"]
