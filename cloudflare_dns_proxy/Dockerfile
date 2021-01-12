from alpine as build
env CLOUDFLARED_URL https://github.com/cloudflare/cloudflared/releases/download/2020.11.11/cloudflared-linux-amd64
env CLOUDFLARED_HASH 60744cb21eb5b47a38cec0646c624218b7b545507fe8d01baefcb683f29e286a
add $CLOUDFLARED_URL /cloudflared
run echo "$CLOUDFLARED_HASH  cloudflared" |sha256sum -c -

from gcr.io/distroless/base
copy --from=build /cloudflared /cloudflared
env TUNNEL_DNS_ADDRESS "0.0.0.0"
entrypoint ["/cloudflared", "proxy-dns"]
