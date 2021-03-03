package config

import (
	"github.com/spf13/viper"
)

type ServiceMonitor struct {
	Name           string
	ProcessMonitor ProcessMonitor
}

type ProcessMonitor struct {
	Exec           string
	ShowCPUPercent bool
	ShowMemory     bool
}

type SysInfoReportConfig struct {
	Server         string
	ShowCPU        bool
	ShowHost       bool
	ShowMemory     bool
	ShowPartitions bool
	ServiceMonitor []ServiceMonitor
	ProcessMonitor []ProcessMonitor
}

func setToViper(config *SysInfoReportConfig) {
	viper.Set("Server", config.Server)
	viper.Set("ShowCPU", config.ShowCPU)
	viper.Set("ShowHost", config.ShowHost)
	viper.Set("ShowMemory", config.ShowMemory)
	viper.Set("ShowPartitions", config.ShowPartitions)
	viper.Set("ServiceMonitor", config.ServiceMonitor)
	viper.Set("ProcessMonitor", config.ProcessMonitor)
}

func getFromViper() *SysInfoReportConfig {

	config := &SysInfoReportConfig{
		Server:         viper.GetString("Server"),
		ShowCPU:        viper.GetBool("ShowCPU"),
		ShowHost:       viper.GetBool("ShowHost"),
		ShowMemory:     viper.GetBool("ShowMemory"),
		ShowPartitions: viper.GetBool("ShowPartitions"),
	}
	serviceMonitor := make([]ServiceMonitor, 1)
	viper.UnmarshalKey("ServiceMonitor", &serviceMonitor)
	config.ServiceMonitor = serviceMonitor

	processMonitor := make([]ProcessMonitor, 1)
	viper.UnmarshalKey("ProcessMonitor", &processMonitor)
	config.ProcessMonitor = processMonitor
	return config
}

func DefaultConfig() *SysInfoReportConfig {

	serviceMonitor := make([]ServiceMonitor, 0)
	serviceMonitor = append(serviceMonitor, ServiceMonitor{
		Name: "edgeupdate",
		ProcessMonitor: ProcessMonitor{
			Exec:           "",
			ShowCPUPercent: true,
			ShowMemory:     true,
		},
	})

	processMonitor := make([]ProcessMonitor, 0)
	processMonitor = append(processMonitor, ProcessMonitor{
		Exec:           "C:\\Windows\\explorer.exe",
		ShowCPUPercent: true,
		ShowMemory:     true,
	})

	ret := SysInfoReportConfig{
		Server:         ":9090",
		ShowCPU:        true,
		ShowHost:       true,
		ShowMemory:     true,
		ShowPartitions: true,
		ServiceMonitor: serviceMonitor,
		ProcessMonitor: processMonitor,
	}

	return &ret
}

func InitConfig(config *SysInfoReportConfig) error {

	if config == nil {
		config = DefaultConfig()
	}
	viper.AddConfigPath(".")
	viper.SetConfigName("SysInfoReport")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			SaveToLocalConfig(config)
		} else {
			return err
		}
	}
	return nil
}

func SaveToLocalConfig(config *SysInfoReportConfig) {
	setToViper(config)
	viper.SafeWriteConfig()
	viper.WriteConfig()
}

func GetSysInfoReportConfig() *SysInfoReportConfig {
	return getFromViper()
}
