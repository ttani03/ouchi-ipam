package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func parseCIDR(cidr string) (int64, int32, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, 0, err
	}

	ip := ipnet.IP.To4()
	mask, _ := ipnet.Mask.Size()

	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3]), int32(mask), nil
}

func parseIPv4ToInt(ip net.IP) int64 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}

	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

func ipv4IntToString(i int64) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, uint32(i))

	return ip.String()
}

func stringToIPv4Int(s string) (int64, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return 0, fmt.Errorf("%s is not a valid IPv4 address", s)
	}
	ip = ip.To4()

	ipInt := int64(ip[0])<<24 + int64(ip[1])<<16 + int64(ip[2])<<8 + int64(ip[3])

	return ipInt, nil
}

func isCIDR(s string) bool {
	ip, ipnet, err := net.ParseCIDR(s)

	return err == nil && ip.String() == ipnet.IP.String()
}
