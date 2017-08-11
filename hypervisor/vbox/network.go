package vbox

import (
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/hyperhq/runv/api"
	"github.com/hyperhq/runv/hypervisor/network"
)

func (vd *VBoxDriver) BuildinNetwork() bool {
	return true
}

func (vd *VBoxDriver) InitNetwork(bIface, bIP string, disableIptables bool) error {
	var i = 0

	if bIP == "" {
		network.BridgeIP = network.DefaultBridgeIP
	} else {
		network.BridgeIP = bIP
	}

	bip, ipnet, err := net.ParseCIDR(network.BridgeIP)
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}

	gateway := bip.Mask(ipnet.Mask)
	inc(gateway, 2)

	if !ipnet.Contains(gateway) {
		logrus.Errorf(err.Error())
		return fmt.Errorf("get Gateway from BridgeIP %s failed", network.BridgeIP)
	}
	prefixSize, _ := ipnet.Mask.Size()
	_, network.BridgeIPv4Net, err = net.ParseCIDR(gateway.String() + fmt.Sprintf("/%d", prefixSize))
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	network.BridgeIPv4Net.IP = gateway
	logrus.Warnf(network.BridgeIPv4Net.String())
	/*
	 * Filter the IPs which can not be used for VMs
	 */
	bip = bip.Mask(ipnet.Mask)
	for inc(bip, 1); ipnet.Contains(bip) && i < 2; inc(bip, 1) {
		i++
		logrus.Debugf("Try %s", bip.String())
		_, err = network.IpAllocator.RequestIP(network.BridgeIPv4Net, bip)
		if err != nil {
			logrus.Errorf(err.Error())
			return err
		}
	}

	return nil
}

func (vc *VBoxContext) ConfigureNetwork(config *api.InterfaceDescription) (*network.Settings, error) {
	ip, ipnet, err := net.ParseCIDR(config.Ip)
	if err != nil {
		logrus.Errorf("Parse interface IP failed %s", err)
		return nil, err
	}

	maskSize, _ := ipnet.Mask.Size()

	return &network.Settings{
		Mac:         config.Mac,
		IPAddress:   ip.String(),
		Gateway:     config.Gw,
		Bridge:      "",
		IPPrefixLen: maskSize,
		Device:      "",
		File:        nil,
	}, nil
}

// Release an interface for a select ip
func (vc *VBoxContext) ReleaseNetwork(releasedIP string) error {
	if err := network.IpAllocator.ReleaseIP(network.BridgeIPv4Net, net.ParseIP(releasedIP)); err != nil {
		return err
	}

	return nil
}

func inc(ip net.IP, count int) {
	for j := len(ip) - 1; j >= 0; j-- {
		for i := 0; i < count; i++ {
			ip[j]++
		}
		if ip[j] > 0 {
			break
		}
	}
}
