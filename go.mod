module github.com/kubevirt/kubernetes-device-plugins

go 1.12

require (
	github.com/Microsoft/go-winio v0.4.14
	github.com/containernetworking/plugins v0.8.3
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/kubevirt/device-plugin-manager v1.14.0
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/pkg/errors v0.8.1
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20190625233234-7109fa855b0f
	golang.org/x/net v0.0.0-20191014212845-da9a3fd4c582
	golang.org/x/sys v0.0.0-20191018095205-727590c5006e
	golang.org/x/text v0.3.2
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03
	google.golang.org/grpc v1.24.0
	k8s.io/kubernetes v1.14.8
)
