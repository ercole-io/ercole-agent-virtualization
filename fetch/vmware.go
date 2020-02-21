package fetch

import (
	"log"
	"os/exec"
	"strings"

	"github.com/ercole-io/ercole-agent-virtualization/config"
	"github.com/ercole-io/ercole-agent-virtualization/marshal"
	"github.com/ercole-io/ercole-agent-virtualization/model"
)

// VMWareClusters return VMWare clusters from the given hyperVisor
func VMWareClusters(hv config.Hypervisor) []model.ClusterInfo {
	clustersBytes := pwshFetcher("vmware.ps1", "-s", "cluster", hv.Endpoint, hv.Username, hv.Password)
	fetchedClusters := marshal.Clusters(clustersBytes)
	for _, fc := range fetchedClusters {
		fc.Type = hv.Type
	}

	return fetchedClusters
}

// VMWareVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func VMWareVirtualMachines(hv config.Hypervisor) []model.VMInfo {
	vmsBytes := pwshFetcher("vmware.ps1", "-s", "vms", hv.Endpoint, hv.Username, hv.Password)
	return marshal.VmwareVMs(vmsBytes)
}

func pwshFetcher(fetcherName string, args ...string) []byte {
	baseDir := config.GetBaseDir()

	args = append([]string{baseDir + "/fetch/" + fetcherName}, args...)
	log.Println("Pwshfetching /usr/bin/pwsh/" + " " + strings.Join(args, " "))
	out, err := exec.Command("/usr/bin/pwsh", args...).Output()
	if err != nil {
		log.Print(string(out))
		log.Fatal(err)
	}

	return out
}
