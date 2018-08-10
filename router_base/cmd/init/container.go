package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	dockerTypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/golang/glog"
	"github.com/vishvananda/netlink"
)

var (
	lanNetwork          = flag.String("docker.lan_network", "", "Container network that this container will act as the gateway for")
	flatNetworks        = flag.String("docker.flat_networks", "", "CSV of container networks that this container will forward to (not masqueraded)")
	uplinkNetwork       = flag.String("docker.uplink_network", "", "Container network used for uplink (connections will be masqueraded)")
	uplinkInterfaceName = flag.String("docker.uplink_interface", "", "Interface used for uplink (connections will be masqueraded)")
)

func InitFromContainerEnvironment() (*RouterConfiguration, error) {
	if *lanNetwork == "" {
		return nil, errors.New("-docker.lan_network flag must be specified")
	}

	if *uplinkNetwork == "" && *uplinkInterfaceName == "" {
		return nil, errors.New("-docker.uplink_network or -docker.uplink_interface must be specified")
	}

	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("error connecting to docker: %v", err)
	}
	defer cli.Close()

	glog.V(2).Info("connected to docker")

	containerID, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("error getting hostname: %v", err)
	}

	containerJSON, err := cli.ContainerInspect(context.TODO(), containerID)
	if err != nil {
		return nil, fmt.Errorf("cannot inspect container using id %q: %v", containerID, err)
	}

	lanInterface, err := findInterfaceByDockerNetwork(*lanNetwork, containerJSON)
	if err != nil {
		return nil, err
	}

	var sr []StaticRoute
	if *flatNetworks != "" {
		for _, flatNetwork := range strings.Split(*flatNetworks, ",") {
			i, err := findInterfaceByDockerNetwork(flatNetwork, containerJSON)
			if err != nil {
				return nil, err
			}

			n, err := cli.NetworkInspect(context.TODO(), flatNetwork, dockerTypes.NetworkInspectOptions{})
			if err != nil {
				return nil, err
			}
			if len(n.IPAM.Config) != 1 {
				return nil, fmt.Errorf("expected 1 IPAM config; got: %v", n.IPAM.Config)
			}
			subnet := n.IPAM.Config[0].Subnet

			sr = append(sr, StaticRoute{
				iface:  i,
				subnet: subnet,
			})
		}
	}

	var uplinkInterface netlink.Link
	if *uplinkInterfaceName != "" {
		uplinkInterface, err = netlink.LinkByName(*uplinkInterfaceName)
		if err != nil {
			return nil, fmt.Errorf("could not get interface %q: %v", *uplinkInterfaceName, err)
		}
	} else {
		uplinkInterface, err = findInterfaceByDockerNetwork(*uplinkNetwork, containerJSON)
		if err != nil {
			return nil, fmt.Errorf("could not get interface for container network %q: %v", *uplinkNetwork, err)
		}
	}

	glog.V(2).Info("applying gateway hack")
	if err := dockerGatewayHacky(lanInterface, cli); err != nil {
		return nil, err
	}

	return &RouterConfiguration{
		lanInterface:    lanInterface,
		flatNetworks:    sr,
		uplinkInterface: uplinkInterface,
	}, nil
}

// for macvlan networks: adds the gateway ip to the lan interface
// for bridge networks: adds the "DefaultGatewayIPv4" aux-address to the lan interface
// throws an error in any other case because there is no non-hacky way to run a container as a gateway as of now
func dockerGatewayHacky(lan netlink.Link, cli *docker.Client) error {
	networkJSON, err := cli.NetworkInspect(context.TODO(), *lanNetwork, dockerTypes.NetworkInspectOptions{})

	if err != nil {
		return fmt.Errorf("error inspecting network %q: %v", *lanNetwork, err)
	}

	if networkJSON.IPAM.Driver != "default" {
		return fmt.Errorf("found unsupported ipam driver %q", networkJSON.IPAM.Driver)
	}

	switch networkJSON.Driver {
	case "bridge":
		found := false

		for _, ipam := range networkJSON.IPAM.Config {
			if gw, ok := ipam.AuxAddress["DefaultGatewayIPv4"]; ok {
				found = true
				var mask int
				if a := strings.Split(ipam.Subnet, "/"); len(a) != 2 {
					return fmt.Errorf("error parsing subnet %q: wrong format %v", ipam.Subnet, a)
				} else if n, err := fmt.Sscanf(a[1], "%d", &mask); n != 1 {
					return fmt.Errorf("error parsing subnet %q: wrong format %q", ipam.Subnet, a[1])
				} else if err != nil {
					return fmt.Errorf("error parsing subnet %q: %v", ipam.Subnet, err)
				}
				s := fmt.Sprintf("%s/%d", gw, mask)
				if addr, err := netlink.ParseAddr(s); err != nil {
					return fmt.Errorf("error parsing address %q: %v", s, err)
				} else if err = netlink.AddrAdd(lan, addr); err != nil {
					return fmt.Errorf("error adding address %q to lan: %v", s, err)
				}
				glog.V(2).Infof("added address %q to lan interface", s)
			}
		}

		if !found {
			return errors.New("did not find a suitable ipam on the bridge; DefaultGatewayIPv4 must be set as an aux-address")
		}
	case "macvlan":
		for _, ipam := range networkJSON.IPAM.Config {
			var mask int
			if a := strings.Split(ipam.Subnet, "/"); len(a) != 2 {
				return fmt.Errorf("error parsing subnet %q: wrong format %v", ipam.Subnet, a)
			} else if n, err := fmt.Sscanf(a[1], "%d", &mask); n != 1 {
				return fmt.Errorf("error parsing subnet %q: wrong format %q", ipam.Subnet, a[1])
			} else if err != nil {
				return fmt.Errorf("error parsing subnet %q: %v", ipam.Subnet, err)
			}
			s := fmt.Sprintf("%s/%d", ipam.Gateway, mask)
			if addr, err := netlink.ParseAddr(s); err != nil {
				return fmt.Errorf("error parsing address %q: %v", s, err)
			} else if err = netlink.AddrAdd(lan, addr); err != nil {
				return fmt.Errorf("error adding address %q to lan: %v", s, err)
			}
			glog.V(2).Infof("added address %q to lan interface", s)
		}
	default:
		return fmt.Errorf("found unsupported lan network driver for gateway hack: %q", networkJSON.Driver)
	}

	return nil
}

func findInterfaceByDockerNetwork(dnet string, j dockerTypes.ContainerJSON) (netlink.Link, error) {
	n, ok := j.NetworkSettings.Networks[dnet]
	if !ok {
		return nil, fmt.Errorf("network %q not found on container info", dnet)
	}

	ip := net.ParseIP(n.IPAddress)
	if ip == nil {
		return nil, fmt.Errorf("could not parse conatiner ip address %q", n.IPAddress)
	}

	return linkForIP(ip)
}

func linkForIP(ip net.IP) (netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("error listing network links: %v", err)
	}

	for _, link := range links {
		addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
		if err != nil {
			return nil, fmt.Errorf("error listing addrs on %q: %v", link.Attrs().Name, err)
		}
		for _, addr := range addrs {
			if addr.IPNet.IP.Equal(ip) {
				return link, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find link for ip %v", ip)
}
