package SysInfo

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"strconv"
)

type NetworkInfo struct {
	ConnectionInfo []ConnectionInfo
	IP             []string
}

type NetAddress struct {
	IP   string
	Port uint32
}

type ConnectionInfo struct {
	Fd     uint32
	Family uint32
	Type   uint32
	Laddr  NetAddress
	Raddr  NetAddress
	Status string
	Pid    int32
	Exec   string
}

func collectNetwork(sysInfo *SysInfo) error {
	network := &NetworkInfo{}
	if err := collectLocalIP(network); err != nil {
		return errors.Wrap(err, "读取本机IP失败")
	}
	connsStat, err := net.Connections("inet")
	if err != nil {
		return errors.Wrap(err, "读取连接失败")
	}
	for _, conn := range connsStat {

		info := &ConnectionInfo{
			Fd:     conn.Fd,
			Family: conn.Family,
			Type:   conn.Type,
			Laddr: NetAddress{
				IP:   conn.Laddr.IP,
				Port: conn.Laddr.Port,
			},
			Raddr: NetAddress{
				IP:   conn.Raddr.IP,
				Port: conn.Raddr.Port,
			},
			Status: conn.Status,
			Pid:    conn.Pid,
		}
		exec, err := fetchProcessExec(info.Pid)
		if err == nil {
			info.Exec = exec.Exec
		} else {
			log.Printf("%+v\n", errors.Errorf("取得进程信息失败:%s", strconv.Itoa(int(info.Pid))))
		}
		network.ConnectionInfo = append(network.ConnectionInfo, *info)
	}
	sysInfo.Network = *network
	return err
}

func collectLocalIP(network *NetworkInfo) error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range ifaces {
		network.IP = append(network.IP, i.String())
	}
	return nil
}
