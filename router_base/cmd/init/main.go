package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/vishvananda/netlink"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

var (
	healthCheck   = flag.Bool("health_check", false, "If set, connects to the internal healthcheck endpoint and exits.")
	tunCreateName = flag.String("create_tun", "", "If set, creates a tun interface with the specified name (to be used with -docker.uplink_interface and probably a VPN client")
	cmd           = flag.String("c", "", "Command to run after initialization")
	httpAddr      = flag.String("http.addr", "0.0.0.0:8080", "Port to serve metrics and health status on")
)

type RouterConfiguration struct {
	lanInterface    netlink.Link
	uplinkInterface netlink.Link
}

func main() {
	flag.Parse()

	args := flag.Args()

	if *healthCheck {
		client := &http.Client{}
		_, port, err := net.SplitHostPort(*httpAddr)
		if err != nil {
			fmt.Printf("bad address %q: %v\n", *httpAddr, err)
			os.Exit(1)
		}
		resp, err := client.Get(fmt.Sprintf("http://localhost:%v/health", port))
		if err != nil {
			fmt.Printf("error connecting to healthcheck: %v\n", err)
			os.Exit(1)
		}
		io.Copy(os.Stdout, resp.Body)
		if resp.StatusCode != http.StatusOK {
			os.Exit(resp.StatusCode)
		}
		return
	}

	var err error
	if *cmd != "" {
		if len(args) > 0 {
			glog.Exit("-c or an exec line; pick one")
		}
		args, err = shellquote.Split(*cmd)
		if err != nil {
			glog.Exitf("error parsing shell command %q: %v", *cmd, err)
		}
	}

	glog.V(2).Infof("MaybeCreateNetworks()")
	if err := MaybeCreateNetworks(); err != nil {
		glog.Exitf("error creating networks: %v", err)
	}

	glog.V(2).Infof("InitFromContainerEnvironment()")
	conf, err := InitFromContainerEnvironment()
	if err != nil {
		glog.Exitf("error initializing network configuration: %v", err)
	}

	glog.V(2).Infof("PatchIPTables()")
	if err := conf.PatchIPTables(); err != nil {
		glog.Exitf("error patching iptables: %v", err)
	}

	scraper, err := SetupMetrics(conf)
	if err != nil {
		glog.Exitf("error setting up metrics: %v", err)
	}
	defer scraper.Close()

	hc := SetupHealthCheck()
	defer hc.Close()

	var e errgroup.Group

	e.Go(func() error {
		l, err := net.Listen("tcp", *httpAddr)
		defer l.Close()

		_, port, err := net.SplitHostPort(l.Addr().String())
		if err != nil {
			return fmt.Errorf("wtf: %v", err)
		}
		OpenPort("tcp", port)

		glog.Infof("listening on %q", *httpAddr)

		return http.Serve(l, nil)
	})

	e.Go(func() error {
		if len(args) > 0 {
			glog.Infof("running %q", strings.Join(args, " "))
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("error starting subprocess: %v", err)
			}
			if err := ReapChildren(cmd.Process); err != nil {
				return fmt.Errorf("error waiting for subprocess: %v", err)
			}
		} else {
			glog.Info("sleeping forever")
			for {
				time.Sleep(time.Duration(9223372036854775807))
			}
		}
		return nil
	})

	if err := e.Wait(); err != nil {
		glog.Exit(err)
	}
}

func MaybeCreateNetworks() error {
	if *tunCreateName == "" {
		return nil
	}

	if err := maybeCreateDevNetTun(); err != nil {
		return fmt.Errorf("error creating /dev/net/tun: %v", err)
	}

	la := netlink.NewLinkAttrs()
	la.Name = *tunCreateName

	link := &netlink.Tuntap{
		LinkAttrs: la,
		Mode:      netlink.TUNTAP_MODE_TUN,
		Flags:     netlink.TUNTAP_DEFAULTS,
	}

	err := netlink.LinkAdd(link)
	if err != nil {
		return fmt.Errorf("error creating tun %q: %v", *tunCreateName, err)
	}

	return nil
}

func maybeCreateDevNetTun() error {
	if err := os.Mkdir("/dev/net", os.FileMode(0755)); !os.IsExist(err) && err != nil {
		return err
	}
	tunMode := uint32(020666)
	if err := unix.Mknod("/dev/net/tun", tunMode, int(unix.Mkdev(10, 200))); !os.IsExist(err) && err != nil {
		return err
	}
	return nil
}
