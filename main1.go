package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/shirou/gopsutil/v3/winservices"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

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

type ServiceStatus struct {
}

type Service struct {
	Name          string
	Config        mgr.Config
	Status        ServiceStatus
	State         svc.State
	Accepts       svc.Accepted
	Pid           uint32
	Win32ExitCode uint32
	// contains filtered or unexported fields
}

type SysInfo struct {
	Memory     MemoryInfo
	CPU        CPUInfo
	Partitions []Partition
	Host       HostInfo
}

type ProcessInfo struct {
	PID        uint32
	exec       string
	running    bool
	CpuPercent float64
	MemUsage   float64
}

type ServiceMonitor struct {
	Name            string
	WarningOnStop   bool
	ShowProcessInfo bool
}

func SetDefault() {
	viper.Set("Server", ":9090")
	viper.Set("ShowCPU", "true")
	viper.Set("ShowHost", "true")
	viper.Set("ShowMemory", "true")
	viper.Set("ShowPartitions", "true")
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
	if !amAdmin() {
		runMeElevated()
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			SetDefault()
			viper.WriteConfig()
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
	sysInfo.CPU.LoadPercent = funk.Map(sysInfo.CPU.LoadPercent, sysInfo.round).([]float64)
	sysInfo.CPU.CpuCount = len(sysInfo.CPU.LoadPercent)
}

func collectDisk(sysInfo *SysInfo) {
	partitionStat, _ := disk.Partitions(true)
	sysInfo.Partitions = funk.Map(partitionStat, func(p disk.PartitionStat) Partition {
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
	}).([]Partition)
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

func collectService(sysInfo *SysInfo) {
	serviceList, err := winservices.ListServices()
	serviceMonitor := make([]ServiceMonitor, 1)
	viper.UnmarshalKey("ServiceMonitor", &serviceMonitor)
	serviceMonitorMap := funk.ToMap(serviceMonitor, "Name").(map[string]ServiceMonitor)

	serviceList = funk.Filter(serviceList, func(s winservices.Service) bool {
		_, ok := serviceMonitorMap[s.Name]
		return ok
	}).([]winservices.Service)

	for index, _ := range serviceList {
		s, _ := winservices.NewService(serviceList[index].Name)
		s.GetServiceDetail()
		serviceList[index] = *s

		if s.Status.Pid != 0 {
			proc, _ := process.NewProcess(int32(s.Status.Pid))
			fmt.Println(proc)
		}

	}

	fmt.Println(serviceList)
	fmt.Println(err)
}

func collectProcess(sysInfo *SysInfo) {
	processInfo, _ := process.Processes()
	for index, _ := range processInfo {
		p, _ := process.NewProcess(processInfo[index].Pid)
		if exe, ok := p.Exe(); ok == nil {
			fmt.Println(exe)
			m, _ := p.MemoryInfo()
			fmt.Println(m)
		}
		fmt.Println(p)

	}

}

func collectInfo() *SysInfo {
	sysInfo := &SysInfo{}
	collectMemory(sysInfo)
	collectCPU(sysInfo)
	collectDisk(sysInfo)
	collectHostName(sysInfo)
	collectService(sysInfo)
	collectProcess(sysInfo)
	return sysInfo
}

func reportHandler(w http.ResponseWriter, r *http.Request) {

	sysInfo := collectInfo()
	t := template.Must(template.ParseFiles("SysInfoReport.html"))
	t.Execute(w, sysInfo)
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
}
