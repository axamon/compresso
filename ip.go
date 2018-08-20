package main

import (
	"net"
)

//IP è l'indirizzo ip
type IP struct {
	ip   net.IP
	ipv6 bool
}

func makeIPv4(a, b, c, d byte) IP {
	ip := IP{net.IPv4(a, b, c, d), false}
	return ip
}

func iPv4ToInt(ip net.IP) uint64 {
	switch len(ip) {
	case 4: // IPv4
		return uint64(ip[0])<<24 | uint64(ip[1])<<16 | uint64(ip[2])<<8 | uint64(ip[3])
	case 16: // IPv6
		return uint64(ip[12])<<24 | uint64(ip[13])<<16 | uint64(ip[14])<<8 | uint64(ip[15])
	default:
		return 0
	}
}

func intToIPv4(ipInt uint64) net.IP {
	return net.IPv4(byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))
}

/* func main() {
	clientip := net.ParseIP(os.Args[1])
	l := IPv4ToInt(clientip)
	fmt.Println(l)
	m := IntToIPv4(l)
	fmt.Println(m)
} */
