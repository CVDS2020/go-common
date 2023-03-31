package netUtils

import (
	"fmt"
	"io"
	"net"
	"syscall"
)

type fdSocket int

func (s fdSocket) Close() error {
	return syscall.Close(int(s))
}

func AllocPort(network string, ip net.IP) (socket io.Closer, port int, err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
	case "udp", "udp4", "udp6":
	default:
		return nil, 0, &net.OpError{Op: "dial", Net: network, Source: ipToAddr(ip), Addr: ipToAddr(ip), Err: net.UnknownNetworkError(network)}
	}
	var sockaddr syscall.Sockaddr
retry:
	switch len(ip) {
	case 0:
		switch network {
		case "tcp", "tcp4", "udp", "udp4":
			sockaddr = &syscall.SockaddrInet4{}
		case "tcp6", "udp6":
			sockaddr = &syscall.SockaddrInet6{}
		}
	case net.IPv4len:
		if network == "tcp6" || network == "udp6" {
			ip = ip.To16()
			goto retry
		}
		sockaddrInet4 := &syscall.SockaddrInet4{}
		copy(sockaddrInet4.Addr[:], ip)
		sockaddr = sockaddrInet4
	case net.IPv6len:
		if network == "tcp4" || network == "udp4" {
			ipv4 := ip.To4()
			if ipv4 == nil {
				return nil, 0, &net.AddrError{Err: "错误的IPV4地址", Addr: ip.String()}
			}
			ip = ipv4
			goto retry
		}
		sockaddrInet6 := &syscall.SockaddrInet6{}
		copy(sockaddrInet6.Addr[:], ip)
		sockaddr = sockaddrInet6
	default:
		return nil, 0, &net.AddrError{Err: "错误的IP地址", Addr: ip.String()}
	}

	var fd int
	switch network[:3] {
	case "tcp":
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	case "udp":
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	}
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err != nil {
			syscall.Close(fd)
		}
	}()

	if err = syscall.Bind(fd, sockaddr); err != nil {
		return nil, 0, err
	}

	sockaddr, err = syscall.Getsockname(fd)
	if err != nil {
		return nil, 0, err
	}

	switch sa := sockaddr.(type) {
	case *syscall.SockaddrInet4:
		return fdSocket(fd), sa.Port, nil
	case *syscall.SockaddrInet6:
		return fdSocket(fd), sa.Port, nil
	default:
		return nil, 0, &net.AddrError{Err: "无效的sockaddr类型", Addr: fmt.Sprint(sa)}
	}
}
