package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/sync/errgroup"
)

var v bool

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage %v: [-v] listen-addr,forward-addr ...\n", os.Args[0])
	flag.PrintDefaults()
}

func parseArgs() (forwarders []*Forwarder) {
	flag.Usage = usage
	flag.BoolVar(&v, "v", false, "log connections")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for i, arg := range flag.Args() {
		f, err := NewForwarder(arg)
		if err != nil {
			log.Fatalf("Could not create forwarder for arg %d: %v", i, err)
		}
		forwarders = append(forwarders, f)
	}
	return
}

func main() {
	forwarders := parseArgs()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grp, ctx := errgroup.WithContext(ctx)
	for _, f := range forwarders {
		f := f
		grp.Go(func() error {
			return f.Go(ctx)
		})
	}
	grp.Wait()
}
