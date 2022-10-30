package id

import (
	"bytes"
	"encoding/binary"
	"gitee.com/sy_183/common/uns"
	"net"
	"sync/atomic"
	"unsafe"
)

const ptrSize = 8

func getTypeId(i interface{}) uintptr {
	return *(*uintptr)(unsafe.Pointer(&i))
}

func getTypeIdString(i interface{}) string {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(getTypeId(i)))
	return uns.BytesToString(bs)
}

func getIpId(ip net.IP) []byte {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return ip
	}
	return ipv4
}

func GenerateIpAddrId(addr *net.IPAddr) string {
	return uns.BytesToString(bytes.Join([][]byte{getIpId(addr.IP), uns.StringToBytes(addr.Zone)}, nil))
}

func GenerateIpNetId(ipNet *net.IPNet) string {
	return uns.BytesToString(bytes.Join([][]byte{getIpId(ipNet.IP), ipNet.Mask}, nil))
}

func GenerateTcpAddrId(addr *net.TCPAddr) string {
	return uns.BytesToString(bytes.Join([][]byte{
		getIpId(addr.IP),
		{byte(addr.Port >> 8), byte(addr.Port)},
		uns.StringToBytes(addr.Zone),
	}, nil))
}

func GenerateUdpAddrId(addr *net.UDPAddr) string {
	return GenerateTcpAddrId((*net.TCPAddr)(addr))
}

func GenerateUnixAddrId(addr *net.UnixAddr) string {
	return addr.Net + ":" + addr.Name
}

func GeneratePairTcpAddrId(localAddr, remoteAddr *net.TCPAddr) string {
	return GenerateTcpAddrId(localAddr) + GenerateTcpAddrId(remoteAddr)
}

func GeneratePairUdpAddrId(localAddr, remoteAddr *net.UDPAddr) string {
	return GenerateUdpAddrId(localAddr) + GenerateUdpAddrId(remoteAddr)
}

func GenerateTcpConnId(conn *net.TCPConn) string {
	return GeneratePairTcpAddrId(conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr))
}

func GenerateUdpConnId(conn *net.UDPConn) string {
	return GeneratePairUdpAddrId(conn.LocalAddr().(*net.UDPAddr), conn.RemoteAddr().(*net.UDPAddr))
}

var (
	ipAddrTypeId  = getTypeIdString(&net.IPAddr{})
	ipNetTypeId   = getTypeIdString(&net.IPAddr{})
	tcpAddrTypeId = getTypeIdString(&net.TCPAddr{})
	udpAddrTypeId = getTypeIdString(&net.UDPAddr{})
	unixAddrId    = getTypeIdString(&net.UnixAddr{})
)

func GenerateAddrId(addr net.Addr) string {
	switch rawAddr := addr.(type) {
	case *net.IPAddr:
		return ipAddrTypeId + GenerateIpAddrId(rawAddr)
	case *net.IPNet:
		return ipNetTypeId + GenerateIpNetId(rawAddr)
	case *net.TCPAddr:
		return tcpAddrTypeId + GenerateTcpAddrId(rawAddr)
	case *net.UDPAddr:
		return udpAddrTypeId + GenerateUdpAddrId(rawAddr)
	case *net.UnixAddr:
		return unixAddrId + GenerateUnixAddrId(rawAddr)
	default:
		return getTypeIdString(addr) + addr.String()
	}
}

func Uint64Id(context *uint64) uint64 {
	return atomic.AddUint64(context, 1)
}
