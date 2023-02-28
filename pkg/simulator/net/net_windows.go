//go:build windows
// +build windows

package net

import (
	"fmt"
	"net/netip"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

func GetAdaptersAddresses(name string) (*winipcfg.IPAdapterAddresses, error) {
	ifcs, err := winipcfg.GetAdaptersAddresses(windows.AF_INET, winipcfg.GAAFlagDefault)
	if err != nil {
		return nil, err
	}

	for _, v := range ifcs {
		if v.FriendlyName() == name {
			return v, nil
		}

	}
	return nil, fmt.Errorf("not found %v", name)
}

//AddIPAddress add ip addeess
//
//example:name:eth0 prefixip:10.0.50.161/24
func AddIPAddress(netname, prefixip string) error {
	adapter, err := GetAdaptersAddresses(netname)
	if err != nil {
		return err
	}
	addr, err := netip.ParsePrefix(prefixip)
	if err != nil {
		return err
	}
	err = adapter.LUID.AddIPAddress(addr)
	if err != nil {
		return err
	}
	return nil
}

//DeleteIPAddress delete ip addeess
//
//example:name:eth0 prefixip:10.0.50.161/24
func DeleteIPAddress(netname, prefixip string) error {
	adapter, err := GetAdaptersAddresses(netname)
	if err != nil {
		return err
	}
	addr, err := netip.ParsePrefix(prefixip)
	if err != nil {
		return err
	}
	err = adapter.LUID.DeleteIPAddress(addr)
	if err != nil {
		return err
	}
	return nil
}
