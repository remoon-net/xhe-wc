package main

import (
	"fmt"
	"net"
	"net/netip"
)

func parseIP(s string) (addr netip.Addr, err error) {
	var ip net.IP
	if ip, _, err = net.ParseCIDR(s); err == nil {
		if ip.To4() != nil {
			ip = ip.To4()
		}
	} else if ip = net.ParseIP(s); ip == nil {
		return netip.Addr{}, fmt.Errorf("ip parse failed")
	}
	addr, _ = netip.AddrFromSlice(ip)
	return
}
