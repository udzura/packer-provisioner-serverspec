package serverspec

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

	var cmd *packer.RemoteCmd
	// Preparing serverspec
	installer := []string{
		"if which apt-get; then",
		"  FILE=`mktemp`;",
		"  curl -qL https://packagecloud.io/omnibus-serverspec/serverspec/packages/ubuntu/trusty/serverspec_2.19.0+20150626234406-198_amd64.deb/download > $FILE;",
		"  sudo dpkg -i $FILE;",
		"  rm -f $FILE;",
		"else",
		"  sudo rpm -Uvh https://packagecloud.io/omnibus-serverspec/serverspec/packages/el/6/serverspec-2.19.0+20150626234135-198.el6.x86_64.rpm/download;",
		"fi",
	}

	ui.Say("Preparing serverspec via omnubus package...")
	cmd = &packer.RemoteCmd{Command: strings.Join(installer, "")}
	err = cmd.StartWithUi(comm, ui)
	if err != nil {
		return err
	}

	// Running serverspec
	runner := fmt.Sprintf("cd %s && sudo /usr/local/bin/rake spec", dest)

	ui.Say("Running serverspec via `rake spec'...")
	cmd = &packer.RemoteCmd{Command: runner}
	err = cmd.StartWithUi(comm, ui)
	if err != nil {
		return err
	}

	// Cleanup processes
	cleaner := []string{
		"if which apt-get; then",
		"  sudo apt-get -y remove serverspec;",
		"else",
		"  sudo yum -y remove serverspec;",
		"fi;",
		fmt.Sprintf("sudo rm -rf %s", dest),
	}
	ui.Say("Cleaning up serverspec packages and files...")
	cmd = &packer.RemoteCmd{Command: strings.Join(cleaner, "")}
	err = cmd.StartWithUi(comm, ui)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provisioner) Cancel() {
	os.Exit(0)
}
