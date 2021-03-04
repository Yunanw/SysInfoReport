package SysInfo

import (
	. "github.com/ahmetb/go-linq/v3"
	"github.com/shirou/gopsutil/v3/disk"
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

func collectPartition(sysInfo *SysInfo) error {
	partitionStat, err := disk.Partitions(true)
	if err != nil {
		return err
	}
	From(partitionStat).SelectT(func(p disk.PartitionStat) Partition {
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
	}).ToSlice(&sysInfo.Partitions)
	return nil
}
