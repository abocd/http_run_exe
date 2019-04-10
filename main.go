package main

import "os/exec"
import (
	"encoding/json"
	// "errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func runExe(exeAdress string) (err2 error, msg string) {
	dir := path.Dir(exeAdress)
	os.Chdir(dir)
	cmd := exec.Command("cmd.exe", "/c", "start "+exeAdress)

	err := cmd.Run()
	if err != nil {
		log.Println("启动失败:", err)
		return err, "启动失败"
	} else {
		log.Println("启动成功!")
		return nil, "启动成功"
	}
}

func closeExe(exeAdress string) (err2 error, msg string) {
	cmd := exec.Command("taskkill", "/f", "/t", "/im", exeAdress)
	err := cmd.Run()
	if err != nil {
		log.Println("结束失败:", err)
		return err, "结束失败"
	} else {
		log.Println("结束成功!")
		return nil, "结束成功"
	}
}

var port int
var exe string
var exeMaps map[string]string

func main() {
	exeMaps = make(map[string]string, 3)
	//go run "d:\go\src\github.com\abocd\test\exec.go" -exe="D:/unity/xiaohu/build/xiaohu/xiaohu.exe|e:/xiaohu/xiaohu2.exe" -port=8081
	flag.IntVar(&port, "port", 8080, "监听的端口号")
	flag.StringVar(&exe, "exe", "D:/unity/xiaohu/build/xiaohu/xiaohu.exe", "监听的程序，多个用|隔开")
	flag.Parse()
	_exe := strings.Split(exe, "|")
	for i := 0; i < len(_exe); i++ {
		_exe[i] = path.Clean(_exe[i])
		_, err := os.Stat(_exe[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		exeMaps[path.Base(_exe[i])] = _exe[i]
	}
	fmt.Println(port, exeMaps)
	startServer()
}

func startServer() {
	http.HandleFunc("/", web)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}

type Msg struct {
	Error   string
	Success string
}

func web(w http.ResponseWriter, r *http.Request) {
	var msg Msg
	fmt.Println(r.RequestURI)
	status := r.URL.Query().Get("status")
	exe := r.URL.Query().Get("exe")
	fmt.Println(status, exe)
	if _, ok := exeMaps[exe]; !ok {
		fmt.Println(exe, "程序不存在")
		msg.Error = fmt.Sprintf("%s程序不存在", exe)
		jsonData, _ := json.Marshal(msg)
		w.Write(jsonData)
		return
	}
	fmt.Println(exeMaps[exe])
	var resultStatus error
	var resultMsg string
	if status == "open" {
		resultStatus, resultMsg = runExe(exeMaps[exe])
	} else {
		resultStatus, resultMsg = closeExe(exe)
	}
	if resultStatus == nil {
		msg.Success = resultMsg
	} else {
		msg.Error = resultMsg
	}

	jsonData, _ := json.Marshal(msg)
	w.Write(jsonData)
}
