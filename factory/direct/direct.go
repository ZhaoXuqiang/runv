package direct

import (
	"github.com/Sirupsen/logrus"
	"github.com/hyperhq/runv/factory/base"
	"github.com/hyperhq/runv/hypervisor"
)

type directFactory struct {
	config hypervisor.BootConfig
}

func New(b hypervisor.BootConfig) base.Factory {
	return &directFactory{config: b}
}

func (d *directFactory) Config() *hypervisor.BootConfig {
	config := d.config
	return &config
}

func (d *directFactory) GetBaseVm() (*hypervisor.Vm, error) {
	logrus.Debugf("direct factory start create vm")
	vm, err := hypervisor.GetVm("", d.Config(), true)
	if err == nil {
		err = vm.Pause(true)
		if err != nil {
			vm.Kill()
			vm = nil
		}
	}
	if err == nil {
		logrus.Debugf("direct factory created vm: %s", vm.Id)
	} else {
		logrus.Errorf("direct factory failed to create vm")
	}
	return vm, err
}

func (d *directFactory) CloseFactory() {}
