package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"github.com/vishvananda/netlink"
)

var (
	iptablesBin = flag.String("iptables.bin", "/sbin/iptables", "Path to iptables binary")
)

func (r RouterConfiguration) PatchIPTables() error {
	glog.V(2).Info("setting base rules")
	if err := setBaseRules(); err != nil {
		return fmt.Errorf("error setting base rules set: %v", err)
	}

	for _, s := range r.flatNetworks {
		if err := forward(r.lanInterface, s.iface, s.subnet); err != nil {
			return err
		}
	}

	if err := forward(r.lanInterface, r.uplinkInterface, ""); err != nil {
		return err
	}

	uplinkAttrs := r.uplinkInterface.Attrs()
	glog.V(2).Infof("setting masquerading out of %q", uplinkAttrs.Name)
	masqueradeCmd := fmt.Sprintf("-t nat -A POSTROUTING -j MASQUERADE -o %v", uplinkAttrs.Name)
	glog.V(4).Infof("applying rule %q", masqueradeCmd)
	if err := iptablesRaw(masqueradeCmd); err != nil {
		return fmt.Errorf("error applying masquerade rule: %v", err)
	}

	return nil
}

func forward(a, b netlink.Link, dst string) error {
	aAttr, bAttr := a.Attrs(), b.Attrs()

	var forwardingCmd string
	if dst == "" {
		glog.V(2).Infof("allowing forwarding from %q to %q", aAttr.Name, bAttr.Name)
		forwardingCmd = fmt.Sprintf("-t filter -A fw-interfaces -j ACCEPT -i %v -o %v", aAttr.Name, bAttr.Name)
	} else {
		glog.V(2).Infof("allowing forwarding from %q to %q with destination %q", aAttr.Name, bAttr.Name, dst)
		forwardingCmd = fmt.Sprintf("-t filter -A fw-interfaces -j ACCEPT -d %v -i %v -o %v", dst, aAttr.Name, bAttr.Name)
	}
	glog.V(4).Infof("applying rule %q", forwardingCmd)

	if err := iptablesRaw(forwardingCmd); err != nil {
		return fmt.Errorf("error applying forwarding rule: %v", err)
	} else {
		return nil
	}
}

func OpenPort(proto, port string) error {
	switch proto {
	case "tcp":
	case "udp":
	default:
		return fmt.Errorf("invalid proto: %q", proto)
	}

	glog.V(2).Infof("opening %s port %s", proto, port)

	rule := fmt.Sprintf("-I in-%s -j ACCEPT -p %s --dport %s", proto, proto, port)
	err := iptablesRaw(rule)
	if err != nil {
		return fmt.Errorf("error applying rule %q: %v", rule, err)
	}

	return nil
}

var POLICY_RULES = []string{
	"-t filter -P INPUT DROP",
	"-t filter -P FORWARD DROP",
	"-t filter -N in-tcp",
	"-t filter -N in-udp",
	"-t filter -N fw-interfaces",
	"-t filter -N fw-open",
}

var BASE_CHAINS = []string{
	"-t filter -A INPUT -j DROP -m state --state INVALID",
	"-t filter -A INPUT -j ACCEPT -m conntrack --ctstate RELATED,ESTABLISHED",
	"-t filter -A INPUT -j ACCEPT -i lo",
	"-t filter -A INPUT -j ACCEPT -p icmp --icmp-type 8 -m conntrack --ctstate NEW",
	"-t filter -A INPUT -j in-tcp -p tcp --tcp-flags FIN,SYN,RST,ACK SYN -m conntrack --ctstate NEW",
	"-t filter -A INPUT -j in-udp -p udp -m conntrack --ctstate NEW",
	"-t filter -A INPUT -j REJECT --reject-with icmp-proto-unreachable",
	"-t filter -A FORWARD -j ACCEPT -m conntrack --ctstate ESTABLISHED,RELATED",
	"-t filter -A FORWARD -j fw-interfaces",
	"-t filter -A FORWARD -j fw-open",
	"-t filter -A FORWARD -j REJECT --reject-with icmp-host-unreach",
	"-t nat -A PREROUTING -j ACCEPT -m conntrack --ctstate RELATED,ESTABLISHED",
}

var REJECTIONS = []string{
	"-t filter -A in-tcp -j REJECT -p tcp --reject-with tcp-reset",
	"-t filter -A in-udp -j REJECT -p udp --reject-with icmp-port-unreachable",
}

func setBaseRules() error {
	for _, rule := range POLICY_RULES {
		glog.V(3).Infof("applying iptables rule %q", rule)
		if err := iptablesRaw(rule); err != nil {
			return fmt.Errorf("error applying rule: %q", rule)
		}
	}

	for _, rule := range BASE_CHAINS {
		glog.V(3).Infof("applying iptables rule %q", rule)
		if err := iptablesRaw(rule); err != nil {
			return fmt.Errorf("error applying rule: %q", rule)
		}
	}

	// TODO: insert user rules here maybe

	for _, rule := range REJECTIONS {
		glog.V(3).Infof("applying iptables rule %q", rule)
		if err := iptablesRaw(rule); err != nil {
			return fmt.Errorf("error applying rule: %q", rule)
		}
	}

	return nil
}

func iptablesRaw(cmd string) error {
	return iptablesRawArg(strings.Split(cmd, " ")...)
}

func iptablesRawArg(arg ...string) error {
	cmd := exec.Command(*iptablesBin, arg...)
	return cmd.Run()
}
