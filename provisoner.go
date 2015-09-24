package serverspec

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	DestPath string `mapstructure:"dest_path"`

	ctx interpolate.Context
}

type Provisioner struct {
	config Config
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
	}, raws...)
	if err != nil {
		return err
	}

	var errs *packer.MultiError
	if p.config.DestPath == "" {
		errs = packer.MultiErrorAppend(errs,
			errors.New("`dest_path' must be specified."))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
	ui.Say("Plugin test")
	ui.Say(fmt.Sprintf("DestPath is: %s", p.config.DestPath))

	return nil
}

func (p *Provisioner) Cancel() {
	os.Exit(0)
}
