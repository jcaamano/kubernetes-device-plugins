package macvtap

import (
	"fmt"
	"time"

	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/golang/glog"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

const (
	nicsPoolSize = 100
	tapPath      = "/dev/tap%d"
)

type MacvtapDevicePlugin struct {
	master string
}

func (mdp *MacvtapDevicePlugin) generateMacvtapDevices() []*pluginapi.Device {
	var macvtapDevs []*pluginapi.Device
	for i := 0; i < nicsPoolSize; i++ {
		// TODO use a different naming convention instead of random veth?
		name, err := ip.RandomVethName()
		if err == nil {
			continue
		}
		macvtapDevs = append(macvtapDevs, &pluginapi.Device{
			ID:     name,
			Health: pluginapi.Healthy,
		})
	}
	return macvtapDevs
}

func masterExists(master string) bool {
	_, err := netlink.LinkByName(master)
	if err != nil {
		return false
	}
	// TODO check more details about master ?
	return true
}

func createMacvtap(name string, master string) (int, error) {
	m, err := netlink.LinkByName(master)
	if err != nil {
		return 0, fmt.Errorf("failed to lookup master %q: %v", master, err)
	}

	mv := &netlink.Macvtap{
		Macvlan: netlink.Macvlan{
			LinkAttrs: netlink.LinkAttrs{
				Name:        name,
				ParentIndex: m.Attrs().Index,
			},
		},
	}

	if err := netlink.LinkAdd(mv); err != nil {
		return 0, fmt.Errorf("failed to create macvtap: %v", err)
	}

	return mv.Attrs().Index, nil
}

func (mdp *MacvtapDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	masterDevs := mdp.generateMacvtapDevices()
	noMasterDevs := make([]*pluginapi.Device, 0)
	emitResponse := func(masterExists bool) {
		if masterExists {
			glog.V(3).Info("Master exists, sending ListAndWatch response with available ports")
			s.Send(&pluginapi.ListAndWatchResponse{Devices: masterDevs})
		} else {
			glog.V(3).Info("Bridge does not exist, sending ListAndWatch response with no ports")
			s.Send(&pluginapi.ListAndWatchResponse{Devices: noMasterDevs})
		}
	}

	didMasterExist := masterExists(mdp.master)
	emitResponse(didMasterExist)

	for {
		doesMasterExist := masterExists(mdp.master)
		if didMasterExist != doesMasterExist {
			emitResponse(doesMasterExist)
			didMasterExist = doesMasterExist
		}
		time.Sleep(10 * time.Second)
	}
}

func (mdp *MacvtapDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var response pluginapi.AllocateResponse

	for _, req := range r.ContainerRequests {
		var devices []*pluginapi.DeviceSpec
		for _, vnic := range req.DevicesIDs {
			dev := new(pluginapi.DeviceSpec)
			index, err := createMacvtap(vnic, mdp.master)
			if err != nil {
				// FIXME proper error handling and cleanup
				return nil, err
			}
			devPath := fmt.Sprint(tapPath, index)
			dev.HostPath = devPath
			dev.ContainerPath = devPath
			dev.Permissions = "rw"
			devices = append(devices, dev)
		}

		response.ContainerResponses = append(response.ContainerResponses, &pluginapi.ContainerAllocateResponse{
			Devices: devices,
		})
	}

	return &response, nil
}

func (mdp *MacvtapDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (mdp *MacvtapDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return nil, nil
}
