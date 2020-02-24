package builder;

// BuildData fetch host and build HostData
func BuildData(configuration config.Configuration) model.HostData {
	out := fetcher("host")
	host := marshal.Host(out)

	host.Environment = configuration.Envtype
	host.Location = configuration.Location
	out = fetcher("filesystem")
	filesystems := marshal.Filesystems(out)

	var clusters []model.ClusterInfo = []model.ClusterInfo{}
	var vms []model.VMInfo = []model.VMInfo{}

	countHypervisors := len(configuration.Hypervisors)
	clustersChannel := make(chan []model.ClusterInfo, countHypervisors)
	virtualMachinesChannel := make(chan []model.VMInfo, countHypervisors)

	for _, hv := range configuration.Hypervisors {
		go func(hv config.Hypervisor) {
			clustersChannel <- fetch.GetClusters(hv)
		}(hv)

		go func(hv config.Hypervisor) {
			virtualMachinesChannel <- fetch.GetVirtualMachines(hv)
		}(hv)
	}

	for i := 0; i < countHypervisors; i++ {
		clusters = append(clusters, (<-clustersChannel)...)
	}

	for i := 0; i < countHypervisors; i++ {
		vms = append(vms, (<-virtualMachinesChannel)...)
	}

	clusterMap := make(map[string][]model.VMInfo)
	clusters = append(clusters, model.ClusterInfo{
		Name:    "not_in_cluster",
		Type:    "unknown",
		CPU:     0,
		Sockets: 0,
		VMs:     []model.VMInfo{},
	})
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
	hostData := new(model.HostData)
	extraInfo := new(model.ExtraInfo)
	extraInfo.Filesystems = filesystems
	extraInfo.Databases = []model.Database{}
	extraInfo.Clusters = clusters
	hostData.Extra = *extraInfo
	hostData.Info = host
	hostData.Hostname = host.Hostname
	// override host name with the one in config if != default
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

	return hostData
}

func fetcher(fetcherName string, args ...string) []byte {
	var (
		cmd    *exec.Cmd
		err    error
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	baseDir := config.GetBaseDir()
	log.Println("Fetching " + baseDir + "/fetch/" + fetcherName + " " + strings.Join(args, " "))

	cmd = exec.Command(baseDir+"/fetch/"+fetcherName, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	// log.Println(stderr)
	if len(stderr.Bytes()) > 0 {
		log.Print(string(stderr.Bytes()))
	}

	if err != nil {
		log.Fatal(err)
	}

	return stdout.Bytes()
}
