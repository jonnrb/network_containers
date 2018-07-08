package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type HealthChecker struct {
	c chan chan error
}

func SetupHealthCheck() io.Closer {
	hc := &HealthChecker{
		c: make(chan chan error),
	}
	go hc.loop()

	http.HandleFunc("/health", hc.Handler)

	return hc
}

func (hc *HealthChecker) Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cls := context.WithTimeout(context.Background(), 10*time.Second)
	defer cls()

	var err error
	c := make(chan error, 1)
	select {
	case hc.c <- c:
		select {
		case err = <-c:
		case <-ctx.Done():
			err = ctx.Err()
		}
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, fmt.Sprintf("%v\n", err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK\n")
	}
}

func (hc *HealthChecker) Close() error {
	close(hc.c)
	return nil
}

func (hc *HealthChecker) loop() {
	for ret := range hc.c {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := httpHeadCheck(ctx)
		ret <- err
	}
}

func httpHeadCheck(ctx context.Context) error {
	if _, err := ctxhttp.Head(ctx, nil, "https://google.com/"); err != nil {
		return err
	} else {
		return nil
	}
}
