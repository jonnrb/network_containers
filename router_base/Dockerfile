from alpine:3.8 as cbuild
workdir /
env SRC https://www.netfilter.org/projects/iptables/files/iptables-1.8.2.tar.bz2
env SRC_HASH a3778b50ed1a3256f9ca975de82c2204e508001fc2471238c8c97f3d1c4c12af
add $SRC /iptables.tar.bz2
run echo "$SRC_HASH  iptables.tar.bz2" |sha256sum -c - \
 && apk add --no-cache gcc musl-dev make libmnl-dev bison flex bzip2 \
 && tar xjf iptables.tar.bz2 \
 && cd iptables* \
 && ./configure --enable-static --disable-shared --disable-nftables \
 && make \
 && make install

from quay.io/jonnrb/go as gobuild
run git clone https://git.jonnrb.com/jon/egress . \
 && CGO_ENABLED=0 GOOS=linux go get ./cmd/init

from gcr.io/distroless/static
copy --from=cbuild /usr/local/sbin/xtables-legacy-multi /sbin/iptables
copy --from=cbuild /lib/libc* /lib/ld* /lib/
copy --from=gobuild /go/bin/init /init
expose 8080
healthcheck --interval=10s --timeout=5s cmd ["/init", "-health_check"]
entrypoint ["/init"]
