package collector

type ICollector interface {
	Collect(sysInfo *SysInfo) error
}

type Collector struct {
	Enable bool
}

type FileMessageCollector struct {
	FileType string
}

type MQMessageCollector struct {
	MQType    string
	QueueName string
}

func LoadCollector() (*[]ICollector, error) {
	collectors := make([]ICollector, 0)

	cpu := &CPUCollector{PreCPU: true, Interval: 0}
	collectors = append(collectors, cpu)

	memory := &MemoryCollector{}
	collectors = append(collectors, memory)

	host := &HostCollector{}
	collectors = append(collectors, host)

	partition := &PartitionCollector{}
	collectors = append(collectors, partition)

	network := &NetworkCollector{
		CollectionInfo: true,
		IP:             true,
	}
	collectors = append(collectors, network)

	process := &ProcessCollector{
		Target: []ProcessTarget{
			{
				Exec:           "C:\\Windows\\System32\\notepad.exe",
				ShowCPUPercent: true,
				ShowMemory:     true,
			},
		},
	}
	collectors = append(collectors, process)

	service := &ServiceCollector{
		Target: []ServiceTarget{
			{
				Name:           "edgeupdate",
				ShowCPUPercent: true,
				ShowMemory:     true,
			},
		},
	}
	collectors = append(collectors, service)

	return &collectors, nil
}
