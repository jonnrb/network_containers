/*
Listens on port 8080 and proxies requests to the URL passed as the first
argument.
*/
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s https://backend-url:1234/some/path\n", os.Args[0])
		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("invalid url %q: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	(&http.Server{
		Addr:    ":8080",
		Handler: httputil.NewSingleHostReverseProxy(u),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}).ListenAndServe()
}
