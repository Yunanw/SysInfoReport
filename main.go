package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"html/template"
	"log"
	"net/http"
)

import . "github.com/ahmetb/go-linq/v3"

type MemoryInfo struct {
	Total float64
	Used float64
	UsedPercent float64
}

type CPUInfo struct {
	CpuCount int
	LoadPercent []float64

}

type Partition  struct {
	Device     string
	Fstype     string
	Mountpoint string
	Total      float64
	Free       float64
	Used       float64
	UsedPercent float64

}

type SysInfo struct {
	Memory MemoryInfo
	CPU CPUInfo
	Partitions []Partition
}



func (sysInfo *SysInfo) ToGB(kbs uint64) float64 {
	return float64(kbs) / 1024.00 / 1024.00 / 1024.00

}

func (sysInfo *SysInfo) ToMB(kbs uint64) float64 {

	return float64(kbs) / 1024.00 / 1024.00
}

func  (sysInfo *SysInfo) round(value float64) float64 {
	return float64(int(value * 100)) / 100
}




func main() {


	http.HandleFunc("/report", reportHandler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	v, _ := mem.VirtualMemory()
	sysInfo:= SysInfo{}
	sysInfo.Memory.Total = sysInfo.round(sysInfo.ToGB(v.Total))
	sysInfo.Memory.Used =  sysInfo.round(sysInfo.ToMB(v.Used))
	sysInfo.Memory.UsedPercent = v.UsedPercent

	sysInfo.CPU.LoadPercent, _ = cpu.Percent(0, true)
	From(sysInfo.CPU.LoadPercent).SelectT(sysInfo.round).ToSlice(&sysInfo.CPU.LoadPercent)
	sysInfo.CPU.CpuCount = len(sysInfo.CPU.LoadPercent)
	partitionStat,_ := disk.Partitions(true)
	From(partitionStat).SelectT(func(p disk.PartitionStat) Partition {
		ret := Partition{}
		ret.Device = p.Device
		ret.Fstype = p.Fstype
		ret.Mountpoint = p.Mountpoint
		diskUsed,_ := disk.Usage(p.Mountpoint)
		ret.Total =  sysInfo.round(sysInfo.ToGB(diskUsed.Total))
		ret.Free = sysInfo.round(sysInfo.ToGB(diskUsed.Free))
		ret.Used = sysInfo.round(sysInfo.ToGB(diskUsed.Used))
		ret.UsedPercent = sysInfo.round(diskUsed.UsedPercent)
		return ret
	}).ToSlice(&sysInfo.Partitions)


	t := template.Must(template.ParseFiles("SysInfoReport.html"))
	t.Execute(w, sysInfo)
}

