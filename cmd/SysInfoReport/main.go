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
	"strings"
	"syscall"
)

func main() {
	if !amAdmin() {
		runMeElevated()
	}

	config.InitConfig(nil)
	http.HandleFunc("/reportHtml", reportHtmlHandler)
	http.HandleFunc("/reportJson", reportJSONHandler)
	log.Fatal(http.ListenAndServe(viper.GetString("server"), nil))
}

func reportHtmlHandler(w http.ResponseWriter, r *http.Request) {

	sysInfo := SysInfo.CollectInfo()
	t := template.Must(template.ParseFiles("SysInfoReport.html"))
	t.Execute(w, sysInfo)
}

func reportJSONHandler(w http.ResponseWriter, r *http.Request) {

	sysInfo := SysInfo.CollectInfo()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sysInfo)
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
