from jonnrb/router_base

run apk update \
 && apk add ca-certificates tor wget \
 && update-ca-certificates

# Get docker-gen
env DOCKER_GEN_VERSION 0.7.3
run wget https://github.com/jwilder/docker-gen/releases/download/$DOCKER_GEN_VERSION/docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz \
 && tar -C /usr/local/bin -xvzf docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz \
 && rm docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz

add netresolve.tmpl /data/netresolve.tmpl
add tor.sh /tor.sh
add torrc /data/torrc

cmd ["/tor.sh"]
