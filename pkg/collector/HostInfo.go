package collector

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/host"
	"time"
)

type HostInfo struct {
	Hostname        string
	Uptime          uint64
	BootTime        time.Time
	Procs           uint64
	OS              string
	Platform        string
	PlatformVersion string
}

type HostCollector struct {
}

func (collector *HostCollector) Collect(sysInfo *SysInfo) error {
	info, err := host.Info()
	if err != nil {
		return errors.Wrap(err, "读取Host信息失败")
	}
	sysInfo.Host.BootTime = time.Unix(int64(info.BootTime), 8)
	sysInfo.Host.Hostname = info.Hostname
	sysInfo.Host.OS = info.OS
	sysInfo.Host.Platform = info.Platform
	sysInfo.Host.PlatformVersion = info.PlatformVersion
	sysInfo.Host.Procs = info.Procs
	sysInfo.Host.Uptime = info.Uptime
	return nil
}
