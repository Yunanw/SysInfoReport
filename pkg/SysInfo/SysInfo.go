package SysInfo

type SysInfo struct {
	Memory     MemoryInfo
	CPU        CPUInfo
	Partitions []Partition
	Processes  []ProcessInfo
	Host       HostInfo
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

func CollectInfo() *SysInfo {
	sysInfo := &SysInfo{}
	collectMemory(sysInfo)
	collectCPU(sysInfo)
	collectPartition(sysInfo)
	collectHostName(sysInfo)
	collectService(sysInfo)
	collectProcess(sysInfo)
	return sysInfo
}
