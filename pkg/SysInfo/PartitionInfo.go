package SysInfo

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/thoas/go-funk"
)

type Partition struct {
	Device      string
	Fstype      string
	Mountpoint  string
	Total       float64
	Free        float64
	Used        float64
	UsedPercent float64
}

func collectPartition(sysInfo *SysInfo) {
	partitionStat, _ := disk.Partitions(true)
	sysInfo.Partitions = funk.Map(partitionStat, func(p disk.PartitionStat) Partition {
		ret := Partition{}
		ret.Device = p.Device
		ret.Fstype = p.Fstype
		ret.Mountpoint = p.Mountpoint
		diskUsed, _ := disk.Usage(p.Mountpoint)
		ret.Total = Round(ToGB(diskUsed.Total))
		ret.Free = Round(ToGB(diskUsed.Free))
		ret.Used = Round(ToGB(diskUsed.Used))
		ret.UsedPercent = Round(diskUsed.UsedPercent)
		return ret
	}).([]Partition)
}
