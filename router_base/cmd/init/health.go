package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sparrc/go-ping"
)

type HealthChecker struct {
	c chan chan error
}

func SetupHealthCheck() io.Closer {
	hc := &HealthChecker{
		c: make(chan chan error),
	}
	go hc.loop()

	http.HandleFunc("/health", hc.Handler())

	return hc
}

func (hc *HealthChecker) Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cls := context.WithTimeout(context.Background(), 10*time.Second)
		defer cls()

		var err error
		c := make(chan error)
		hc.c <- c

		select {
		case err = <-c:
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
}

func (hc *HealthChecker) Close() error {
	hc.c <- nil
	return nil
}

func (hc *HealthChecker) loop() {
	for {
		c := <-hc.c
		if c == nil {
			return
		}

		ctx, cls := context.WithTimeout(context.Background(), 10*time.Second)
		err := PingCheck(ctx)
		cls()

		c <- err

		for more := true; more; {
			select {
			case c := <-hc.c:
				if c == nil {
					return
				}
				c <- err
			default:
				more = false
			}
		}
	}
}

func PingCheck(ctx context.Context) error {
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		return fmt.Errorf("error creating pinger: %v", err)
	}

	const numPackets = 4

	pinger.Count = numPackets
	pinger.SetPrivileged(true)
	if deadline, ok := ctx.Deadline(); ok {
		pinger.Timeout = deadline.Sub(time.Now())
		if pinger.Timeout <= 0 {
			return fmt.Errorf("deadline exceeded")
		}
	}

	pinger.Run()

	stats := pinger.Statistics()
	if stats.PacketsSent != numPackets {
		return fmt.Errorf("failed to send packets %v < %v", stats.PacketsSent, numPackets)
	}
	if stats.PacketsSent != stats.PacketsRecv {
		return fmt.Errorf("lost %v packets", stats.PacketsSent-stats.PacketsSent)
	}

	return nil
}
