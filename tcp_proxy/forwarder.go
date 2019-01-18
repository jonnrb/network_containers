package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Forwarder struct {
	listener       net.Listener
	forwardingAddr string
}

func NewForwarder(pair string) (*Forwarder, error) {
	s := strings.SplitN(pair, ",", 2)
	if len(s) != 2 {
		return nil, fmt.Errorf("arg should be of the form \"listen-addr,forward-addr\"; got %q", pair)
	}
	lAddr, fAddr := s[0], s[1]

	l, err := net.Listen("tcp", lAddr)
	if err != nil {
		return nil, err
	}
	if v {
		log.Printf("Listening on %q", lAddr)
	}
	return &Forwarder{
		listener:       l,
		forwardingAddr: fAddr,
	}, nil
}

func (f *Forwarder) Go(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		f.listener.Close()
	}()

	for {
		c, err := f.listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			return err
		}
		if v {
			log.Printf(
				"Connection from %q being forwarded to %q",
				c.RemoteAddr().String(),
				f.forwardingAddr)
		}
		go func() {
			f, err := net.Dial("tcp", f.forwardingAddr)
			if err != nil {
				c.Close()
				return
			}
			bidiTunnel(f, c)
		}()
	}
}

func bidiTunnel(a, b net.Conn) {
	aTCP, bTCP := a.(*net.TCPConn), b.(*net.TCPConn)
	go tunnelTCP(aTCP, bTCP)
	go tunnelTCP(bTCP, aTCP)
}

func tunnelTCP(dst, src *net.TCPConn) {
	io.Copy(dst, src)
	dst.CloseWrite()
	src.CloseRead()
}
