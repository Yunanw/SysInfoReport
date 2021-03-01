package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"time"
)

import . "github.com/ahmetb/go-linq/v3"

type MemoryInfo struct {
	Total       float64
	Used        float64
	UsedPercent float64
}

type CPUInfo struct {
	CpuCount    int
	LoadPercent []float64
}

type Partition struct {
	Device      string
	Fstype      string
	Mountpoint  string
	Total       float64
	Free        float64
	Used        float64
	UsedPercent float64
}

type HostInfo struct {
	Hostname        string
	Uptime          uint64
	BootTime        time.Time
	Procs           uint64
	OS              string
	Platform        string
	PlatformVersion string
}

type Network struct {
	Ipv4 []string
}

type SysInfo struct {
	Memory     MemoryInfo
	CPU        CPUInfo
	Partitions []Partition
	Host       HostInfo
}

func (sysInfo *SysInfo) ToGB(kbs uint64) float64 {
	return float64(kbs) / 1024.00 / 1024.00 / 1024.00

}

func (sysInfo *SysInfo) ToMB(kbs uint64) float64 {

	return float64(kbs) / 1024.00 / 1024.00
}

func (sysInfo *SysInfo) round(value float64) float64 {
	return float64(int(value*100)) / 100
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			viper.Set("server", ":9090")
			viper.WriteConfig()
			viper.SafeWriteConfig()
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	http.HandleFunc("/report", reportHandler)
	log.Fatal(http.ListenAndServe(viper.GetString("server"), nil))
}

func collectMemory(sysInfo *SysInfo) {
	v, _ := mem.VirtualMemory()

	sysInfo.Memory.Total = sysInfo.round(sysInfo.ToGB(v.Total))
	sysInfo.Memory.Used = sysInfo.round(sysInfo.ToMB(v.Used))
	sysInfo.Memory.UsedPercent = v.UsedPercent
}

func collectCPU(sysInfo *SysInfo) {
	sysInfo.CPU.LoadPercent, _ = cpu.Percent(0, true)
	From(sysInfo.CPU.LoadPercent).SelectT(sysInfo.round).ToSlice(&sysInfo.CPU.LoadPercent)
	sysInfo.CPU.CpuCount = len(sysInfo.CPU.LoadPercent)
}

func collectDisk(sysInfo *SysInfo) {
	partitionStat, _ := disk.Partitions(true)
	From(partitionStat).SelectT(func(p disk.PartitionStat) Partition {
		ret := Partition{}
		ret.Device = p.Device
		ret.Fstype = p.Fstype
		ret.Mountpoint = p.Mountpoint
		diskUsed, _ := disk.Usage(p.Mountpoint)
		ret.Total = sysInfo.round(sysInfo.ToGB(diskUsed.Total))
		ret.Free = sysInfo.round(sysInfo.ToGB(diskUsed.Free))
		ret.Used = sysInfo.round(sysInfo.ToGB(diskUsed.Used))
		ret.UsedPercent = sysInfo.round(diskUsed.UsedPercent)
		return ret
	}).ToSlice(&sysInfo.Partitions)
}

func collectHostName(sysInfo *SysInfo) {
	info, _ := host.Info()
	sysInfo.Host.BootTime = time.Unix(int64(info.BootTime), 8)
	sysInfo.Host.Hostname = info.Hostname
	sysInfo.Host.OS = info.OS
	sysInfo.Host.Platform = info.Platform
	sysInfo.Host.PlatformVersion = info.PlatformVersion
	sysInfo.Host.Procs = info.Procs
	sysInfo.Host.Uptime = info.Uptime
}

func collectInfo() *SysInfo {
	sysInfo := &SysInfo{}
	collectMemory(sysInfo)
	collectCPU(sysInfo)
	collectDisk(sysInfo)
	collectHostName(sysInfo)
	return sysInfo
}

func reportHandler(w http.ResponseWriter, r *http.Request) {

	sysInfo := collectInfo()
	t := template.Must(template.ParseFiles("SysInfoReport.html"))
	t.Execute(w, sysInfo)
}
