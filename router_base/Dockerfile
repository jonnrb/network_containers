from ubuntu:xenial as cbuild
run apt-get update \
 && apt-get install -y build-essential libnfnetlink-dev bison flex
workdir /
copy iptables.tar.bz2 /
run tar xjf iptables.tar.bz2 \
 && cd iptables* \
 && ./configure --enable-static --disable-shared \
 && make \
 && make install

from quay.io/jonnrb/go as gobuild
run CGO_ENABLED=0 GOOS=linux go get go.jonnrb.io/egress/cmd/init

from gcr.io/distroless/base
copy --from=cbuild /usr/local/sbin/xtables-multi /sbin/iptables
copy --from=gobuild /go/bin/init /init
expose 8080
healthcheck --interval=10s --timeout=5s cmd ["/init", "-health_check"]
entrypoint ["/init"]
