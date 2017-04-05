from golang:alpine

run apk update && \
  apk add bash iptables git alpine-sdk iproute2

workdir /app

# Get docker-gen
env DOCKER_GEN_VERSION 0.7.3
run wget https://github.com/jwilder/docker-gen/releases/download/$DOCKER_GEN_VERSION/docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz \
 && tar -C /usr/local/bin -xvzf docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz \
 && rm docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz

# Get golang deps
run go get github.com/jonnrb/docker_utils/get_bridge_name && \
    go get github.com/jonnrb/docker_utils/get_container_ip

add netter.sh /app/netter.sh
add docker-gen.tmpl /app/docker-gen.tmpl
ENTRYPOINT [ "/usr/local/bin/docker-gen" \
           , "-watch", "-notify-output"  \
           , "-notify", "/app/netter.sh" \
           , "/app/docker-gen.tmpl"      \
           , "/data/netter-gen.conf"     \
           ]
