package builder

import (
	"log"

	"github.com/ercole-io/ercole-agent-virtualization/config"
	"github.com/ercole-io/ercole-agent-virtualization/marshal"
	"github.com/ercole-io/ercole-agent-virtualization/model"
)

var hostDataSchemaVersion = 5

// BuildData fetch host and build HostData
func BuildData(configuration config.Configuration, version string) *model.HostData {
	log.Println("Start to build hostData...")

	out := fetcher("filesystem")
	filesystems := marshal.Filesystems(out)

	extraInfo := new(model.ExtraInfo)
	extraInfo.Filesystems = filesystems
	extraInfo.Databases = []model.Database{}
	extraInfo.Clusters = getClustersInfos(configuration)

	out = fetcher("host")
	host := marshal.Host(out)
	host.Environment = configuration.Envtype
	host.Location = configuration.Location

	hostData := new(model.HostData)
	hostData.Extra = *extraInfo
	hostData.Info = host

	hostData.Hostname = host.Hostname
	if configuration.Hostname != "default" {
		hostData.Hostname = configuration.Hostname
	}

	hostData.Environment = configuration.Envtype
	hostData.Location = configuration.Location
	hostData.HostType = configuration.HostType
	hostData.Version = version
	hostData.HostDataSchemaVersion = hostDataSchemaVersion
	hostData.Databases = ""
	hostData.Schemas = ""

	log.Println("HostData builded!")
	return hostData
}

func getClustersInfos(configuration config.Configuration) []model.ClusterInfo {
	countHypervisors := len(configuration.Hypervisors)
	clustersChan := make(chan []model.ClusterInfo, countHypervisors)
	virtualMachinesChan := make(chan []model.VMInfo, countHypervisors)

	for _, hv := range configuration.Hypervisors {
		done := make(chan bool, 1)
		go func(hv config.Hypervisor, done chan bool) {
			clustersChan <- fetchClusters(hv)
			done <- true
		}(hv, done)
		if !configuration.ParallelizeRequests {
			<-done
		}

		done = make(chan bool, 1)
		go func(hv config.Hypervisor, done chan bool) {
			virtualMachinesChan <- fetchVirtualMachines(hv)
			done <- true
		}(hv, done)
		if !configuration.ParallelizeRequests {
			<-done
		}
	}

	var clusters []model.ClusterInfo = []model.ClusterInfo{}
	var vms []model.VMInfo = []model.VMInfo{}

	for i := 0; i < countHypervisors; i++ {
		clusters = append(clusters, (<-clustersChan)...)
	}

	for i := 0; i < countHypervisors; i++ {
		vms = append(vms, (<-virtualMachinesChan)...)
	}

	setVMsInClusterInfo(clusters, vms)

	return clusters
}

func setVMsInClusterInfo(clusters []model.ClusterInfo, vms []model.VMInfo) {
	clusters = append(clusters, model.ClusterInfo{
		Name:    "not_in_cluster",
		Type:    "unknown",
		CPU:     0,
		Sockets: 0,
		VMs:     []model.VMInfo{},
	})

	clusterMap := make(map[string][]model.VMInfo)

	for _, vm := range vms {
		if vm.ClusterName == "" {
			vm.ClusterName = "not_in_cluster"
		}
		clusterMap[vm.ClusterName] = append(clusterMap[vm.ClusterName], vm)
	}

	for i := range clusters {
		if clusterMap[clusters[i].Name] != nil {
			clusters[i].VMs = clusterMap[clusters[i].Name]
		} else {
			clusters[i].VMs = []model.VMInfo{}
		}
	}
}
