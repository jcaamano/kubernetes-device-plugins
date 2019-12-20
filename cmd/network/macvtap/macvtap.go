package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"github.com/kubevirt/kubernetes-device-plugins/pkg/network/macvtap"
)

func main() {
	flag.Parse()

	_, mastersListDefined := os.LookupEnv(macvtap.MastersListEnvironmentVariable)
	if !mastersListDefined {
		glog.Exit("MASTERS environment variable must be set in format MASTER[,MASTER[...]]")
	}

	manager := dpm.NewManager(macvtap.MacvtapLister{})
	manager.Run()
}
