package driverloader

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/hyperhq/runv/hypervisor"
	"github.com/hyperhq/runv/hypervisor/kvmtool"
	"github.com/hyperhq/runv/hypervisor/libvirt"
	"github.com/hyperhq/runv/hypervisor/qemu"
	"github.com/hyperhq/runv/hypervisor/xen"
	"github.com/hyperhq/runv/lib/vsock"
)

func Probe(driver string) (hd hypervisor.HypervisorDriver, err error) {
	defer func() {
		if hd != nil && hypervisor.VsockCidManager == nil {
			hypervisor.VsockCidManager = vsock.NewDefaultVsockCidAllocator()
		}
	}()

	driver = strings.ToLower(driver)
	switch driver {
	case "libvirt":
		ld := libvirt.InitDriver()
		if ld != nil {
			logrus.Infof("Driver %q loaded", driver)
			return ld, nil
		}
	case "kvm", "qemu-kvm":
		if _, err := os.Stat("/dev/kvm"); os.IsNotExist(err) {
			return nil, fmt.Errorf("Driver %q is unavailable", driver)
		}
		qd := qemu.InitDriver()
		if qd != nil {
			logrus.Infof("Driver %q loaded", driver)
			return qd, nil
		}
	case "xen", "":
		xd := xen.InitDriver()
		if xd != nil {
			logrus.Infof("Driver \"xen\" loaded")
			return xd, nil
		}
		if driver == "xen" {
			return nil, fmt.Errorf("Driver %q is unavailable", driver)
		}
		fallthrough // only for ""
	case "qemu": // "qemu" or "", kvm will be enabled if the system enables kvm
		qd := qemu.InitDriver()
		if qd != nil {
			logrus.Infof("Driver \"qemu\" loaded")
			return qd, nil
		}
	case "kvmtool":
		kd := kvmtool.InitDriver()
		if kd != nil {
			logrus.Infof("Driver \"kvmtool\" loaded")
			return kd, nil
		}
	}

	return nil, fmt.Errorf("Unsupported driver %q", driver)
}
