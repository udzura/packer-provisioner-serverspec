package serverspec

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	SourcePath string `mapstructure:"source_path"`

	ctx interpolate.Context
}

type Provisioner struct {
	config   Config
	destPath string
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
	if p.config.SourcePath == "" {
		errs = packer.MultiErrorAppend(errs,
			errors.New("`source_path' must be specified."))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
	// Upload serverspec source
	dest := fmt.Sprintf("/tmp/packer-provisioner-serverspec.%d.%d", time.Now().Unix(), syscall.Getpid)
	ui.Say(fmt.Sprintf("Uploading %s => %s", p.config.SourcePath, dest))

	info, err := os.Stat(p.config.SourcePath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source_path %s is not a directory.", p.config.SourcePath)
	}

	err = comm.UploadDir(dest, p.config.SourcePath, nil)
	if err != nil {
		return err
	}
	p.destPath = dest

	return nil
}

func (p *Provisioner) Cancel() {
	os.Exit(0)
}
