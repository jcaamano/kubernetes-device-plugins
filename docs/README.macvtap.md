# Network Macvtap Device Plugin

Creates a macvtap interface in the host and makes the associated tap device
available in the pod. Use in combination with multus + macvtap cni to make the
macvtap interface available in the pod. Main use case is to provide tha
interface as a non managed ethernet device to a libvirt VM hosted in the pod.


## Usage

Specify list of host interfaces that should be exposed to the device plugin as
available parents/masters of the macvtap interfaces:

```bash
$ kubectl create configmap device-plugin-network-macvtap --from-literal=masters="eth0"
configmap "device-plugin-network-macvtap" created
```

**Note:** Multiple masters can be exposed by using a comma separated list of
names: `eth0,eth1,eth2`

Once this is done, the device plugin can be deployed using daemon set:

```
$ kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubernetes-device-plugins/master/manifests/macvtap-ds.yml
daemonset "device-plugin-network-macvtap" created

$ kubectl get pods
NAME                                 READY     STATUS    RESTARTS   AGE
device-plugin-network-macvtap-745x4    1/1    Running           0    5m
```

If the daemonset is running the user can define a pod requesting a macvtap tap
device over host interface `eth0`.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: bridge-consumer
spec:
  containers:
  - name: busybox
    image: busybox
    command: ["/bin/sleep", "123"]
    resources:
      limits:
        macvtap.network.kubevirt.io/eth0: 1
```

```bash
$ kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubernetes-device-plugins/master/examples/macvtap-consumer.yml
pod "macvtap-consumer" created
```

Once the pod is created it should have a tap device corresponding to the macvtap
interface.

```bash
$ ls /dev
/dev/tap4
```

## Pending

* FIXME: proper error handling and cleanup
* FIXME: add tests
* TODO: use predictable interface/device names instead of randomized
* TODO: better health check over parent interfaces