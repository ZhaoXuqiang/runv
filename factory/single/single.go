package single

import (
	"github.com/Sirupsen/logrus"
	"github.com/hyperhq/runv/factory/base"
	"github.com/hyperhq/runv/hypervisor"
)

type Factory struct{ base.Factory }

func New(b base.Factory) Factory {
	return Factory{Factory: b}
}

func (f Factory) GetVm(cpu, mem int) (*hypervisor.Vm, error) {
	// check if match the base
	config := f.Config()
	if config.CPU > cpu || config.Memory > mem {
		boot := *config
		return Dummy(boot).GetVm(cpu, mem)
	}

	vm, err := f.GetBaseVm()
	if err != nil {
		return nil, err
	}

	// unpause
	vm.Pause(false)

	// hotplug add cpu and memory
	var needOnline bool = false
	if vm.Cpu < cpu {
		needOnline = true
		logrus.Debug("HotAddCpu for cached Vm")
		err = vm.SetCpus(cpu)
		logrus.Debugf("HotAddCpu result %v", err)
	}
	if vm.Mem < mem {
		needOnline = true
		logrus.Debug("HotAddMem for cached Vm")
		err = vm.AddMem(mem)
		logrus.Debugf("HotAddMem result %v", err)
	}
	if needOnline {
		logrus.Debug("OnlineCpuMem for cached Vm")
		vm.OnlineCpuMem()
	}
	if err != nil {
		vm.Kill()
		vm = nil
	}
	return vm, err
}

type Dummy hypervisor.BootConfig

func (f Dummy) GetVm(cpu, mem int) (*hypervisor.Vm, error) {
	config := hypervisor.BootConfig(f)
	config.CPU = cpu
	config.Memory = mem
	return hypervisor.GetVm("", &config, false)
}
func (f Dummy) CloseFactory() {}
