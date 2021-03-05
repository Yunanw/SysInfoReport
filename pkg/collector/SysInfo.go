package collector

import "log"

type SysInfo struct {
	Memory     MemoryInfo
	CPU        CPUInfo
	Partitions []Partition
	Processes  []ProcessInfo
	Services   []ServiceInfo
	Host       HostInfo
	Network    NetworkInfo
}

func ToGB(kbs uint64) float64 {
	return float64(kbs) / 1024.00 / 1024.00 / 1024.00
}

func ToMB(kbs uint64) float64 {
	return float64(kbs) / 1024.00 / 1024.00
}

func Round(value float64) float64 {
	return float64(int(value*100)) / 100
}

func CollectInfo() (*SysInfo, error) {
	sysInfo := &SysInfo{}
	collectors, err := LoadCollector()
	if err != nil {
		return nil, err
	}
	for index, _ := range *collectors {
		if err := (*collectors)[index].Collect(sysInfo); err != nil {
			log.Print(err)
		}
	}
	//if err := collectMemory(sysInfo); err != nil {
	//	return nil, err
	//}
	//
	//if err := collectCPU(sysInfo); err != nil {
	//	return nil, err
	//}
	//
	//if err := collectPartition(sysInfo); err != nil {
	//	return nil, err
	//}
	//
	//if err := collectHostName(sysInfo); err != nil {
	//	return nil, err
	//}
	//
	//if err := collectService(sysInfo); err != nil {
	//	return nil, err
	//}
	//
	//if err := collectProcess(sysInfo); err != nil {
	//	return nil, err
	//}
	//
	//if err := collectNetwork(sysInfo); err != nil {
	//	return nil, err
	//}

	return sysInfo, nil
}
