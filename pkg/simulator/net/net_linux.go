//go:build linux
// +build linux

package net

import (
	"net"

	"github.com/vishvananda/netlink"
)

//AddIPAddress add ip addeess
//
//example:name:eth0 prefixip:10.0.50.161/24
func AddIPAddress(netname string, prefixip string) error {
	adapter, err := GetAdaptersAddresses(netname)
	if err != nil {
		return err
	}
	ip, ipnet, err := net.ParseCIDR(prefixip)
	if err != nil {
		return err
	}
	ipnet.IP = ip
	addr := &netlink.Addr{IPNet: ipnet}

	err = netlink.AddrAdd(adapter, addr)
	if err != nil {
		return err
	}

	return nil
}

//DeleteIPAddress delete ip addeess
//
//example:name:eth0 prefixip:10.0.50.161/24
func DeleteIPAddress(netname string, prefixip string) error {
	adapter, err := GetAdaptersAddresses(netname)
	if err != nil {
		return err
	}
	ip, ipnet, err := net.ParseCIDR(prefixip)
	if err != nil {
		return err
	}
	ipnet.IP = ip
	addr := &netlink.Addr{IPNet: ipnet}

	err = netlink.AddrDel(adapter, addr)
	if err != nil {
		return err
	}

	return nil
}

func GetAdaptersAddresses(netname string) (netlink.Link, error) {
	localInterface, err := netlink.LinkByName(netname)
	if err != nil {
		return nil, err
	}
	return localInterface, err
}
