package icmp

import (
	"log"
	"net"
)

// Return true if address is IPv4.
func IsIPv4(address string) bool {
	ip := net.ParseIP(GetIP(address))
	return ip.To4() != nil
}

// Return true if address is IPv6.
func IsIPv6(address string) bool {
	ip := net.ParseIP(GetIP(address))
	return ip.To4() == nil
}

// Split uint16 in to two uint8, return (hign uint8, low uint8).
func SplitInt16(i uint16) (uint8, uint8) {
	high := uint8(i >> 8) // Get the high byte
	low := uint8(i)       // Get the low byte
	return high, low
}

// Combine (hign uint8, low uint8) and return uint16.
func CombineInt8(high, low uint8) uint16 {
	return uint16(high)<<8 | uint16(low)
}

// Get IP address of hostname
func GetIP(hostname string) string {
	ips, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatalf("Error looking up IP addresses for %s: %v\n", hostname, err)
	}
	return ips[len(ips)-1]
}

// Set the value of TargetIP.
func SetTargetIP(hostname string) {
	dstAddr, err := net.ResolveIPAddr("ip", GetIP(hostname))
	if err != nil {
		log.Fatal("net.ResolveIPAddr:", err)
	}
	TargetIP = dstAddr
}

// Set the value of ListenIP.
func SetListenIP(hostname string) {
	if hostname == "" {
		ListenIP = ""
	} else {
		ListenIP = GetIP(hostname)
	}
}
