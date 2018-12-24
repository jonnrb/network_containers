from quay.io/jonnrb/go as build
add main.go .
run CGO_ENABLED=0 go build -o reverse_single .

from gcr.io/distroless/static
copy --from=build /src/reverse_single /reverse_single
expose 8080
entrypoint ["/reverse_single"]
