package serverspec

import (
	"os"

	"github.com/mitchellh/packer/packer"
)

type Provisioner struct {
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
	return nil
}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
	ui.Say("Plugin test")
	return nil
}

func (p *Provisioner) Cancel() {
	os.Exit(0)
}
