from golang:alpine

run apk update && apk add alpine-sdk hostapd git iw

run go get github.com/jonnrb/docker_utils/get_bridge_name

add init.sh /
cmd ["/init.sh"]
