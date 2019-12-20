package macvtap

import (
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

const (
	resourceNamespace              = "bmacvtapridge.network.kubevirt.io"
	MastersListEnvironmentVariable = "MASTERS"
)

type MacvtapLister struct{}

func (ml MacvtapLister) GetResourceNamespace() string {
	return resourceNamespace
}

func (ml MacvtapLister) Discover(pluginListCh chan dpm.PluginNameList) {
	var plugins = make(dpm.PluginNameList, 0)

	mastersListRaw := os.Getenv(MastersListEnvironmentVariable)
	masters := strings.Split(mastersListRaw, ",")
	plugins = append(plugins, masters...)

	pluginListCh <- plugins
}

func (ml MacvtapLister) NewPlugin(master string) dpm.PluginInterface {
	glog.V(3).Infof("Creating device plugin %s", master)
	return &MacvtapDevicePlugin{
		master: master,
	}
}
