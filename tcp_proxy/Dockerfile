from quay.io/jonnrb/go as build
add go.* ./
run go mod download
add *.go ./
run CGO_ENABLED=0 go get .

from gcr.io/distroless/static
copy --from=build /go/bin/tcp_proxy /tcp_proxy
expose 8080
entrypoint ["/tcp_proxy"]
