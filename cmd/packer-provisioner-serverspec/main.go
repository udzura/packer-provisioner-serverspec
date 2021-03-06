package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/udzura/packer-provisioner-serverspec"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterProvisioner(new(serverspec.Provisioner))
	server.Serve()
}
