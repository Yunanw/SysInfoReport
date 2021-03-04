package main

import (
	"SysInfoReport/pkg/SysInfo"
	"SysInfoReport/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/sys/windows"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	if !amAdmin() {
		runMeElevated()
		return
	}

	if err := config.InitConfig(nil); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/reportHtml", reportHtmlHandler)
	http.HandleFunc("/reportJson", reportJSONHandler)
	log.Fatal(http.ListenAndServe(viper.GetString("server"), nil))
}

func reportHtmlHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	sysInfo, err := SysInfo.CollectInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	t := template.Must(template.ParseFiles("SysInfoReport.html"))
	err = t.Execute(w, sysInfo)
	fmt.Println(err)
}

func reportJSONHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	sysInfo, err := SysInfo.CollectInfo()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if err = json.NewEncoder(w).Encode(sysInfo); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

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
	if runtime.GOOS != "windows" {
		return true
	}
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("Windows下需要管理员权限运行本程序")
		return false
	}
	fmt.Println("admin yes")
	return true
}
