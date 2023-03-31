package netUtils

import "net"

func ipToAddr(ip net.IP) net.Addr {
	if ip == nil {
		return nil
	}
	return &net.IPAddr{IP: ip}
}
