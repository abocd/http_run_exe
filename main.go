package main

import (
	"io/ioutil"
	"os/exec"
)
import (
	"encoding/json"
	// "errors"
	"flag"
	"fmt"

	// "log"
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
	fmt.Println("正在启动程序", exeAdress)
	if err != nil {
		fmt.Println("启动失败:", err)
		return err, "启动失败"
	} else {
		fmt.Println("启动成功!")
		return nil, "启动成功"
	}
}

func closeExe(exeAdress string) (err2 error, msg string) {
	if strings.Index(exeAdress, "video") == 0 {
		exeAdress = videoPlayer
	}
	cmd := exec.Command("taskkill", "/f", "/t", "/im", exeAdress)
	err := cmd.Run()
	fmt.Println("正在结束程序", exeAdress)
	if err != nil {
		fmt.Println("结束失败:", err)
		return err, "结束失败"
	} else {
		fmt.Println("结束成功!")
		return nil, "结束成功"
	}
}

/**
 *  列表程序
 */
func listExe() (err2 error, msg map[string]Exe) {
	exeList := make(map[string]Exe)
	for name, exe := range ExeList {
		if exe.Show {
			exeList[name] = exe
		}
	}
	return nil, exeList
}

var port int
var exe string
var name string
var videoPlayer string

type Exe struct {
	Name string
	Path string
	Ico  string
	Show bool
}

var ExeList map[string]Exe

func main() {

	ExeList = make(map[string]Exe, 3)
	//go run "d:\go\src\github.com\abocd\test\exec.go" -exe="D:/unity/xiaohu/build/xiaohu/xiaohu.exe|e:/xiaohu/xiaohu2.exe" -port=8081
	flag.IntVar(&port, "port", 8080, "监听的端口号")
	flag.StringVar(&exe, "exe", "D:/unity/xiaohu/build/xiaohu/xiaohu.exe", "监听的程序，多个用|隔开")
	flag.StringVar(&name, "name", "语音精灵", "和程序配套的对应的程序名称，多个用|隔开")
	flag.StringVar(&videoPlayer, "player", "PotPlayerMini64.exe", "视频播放器名称")
	flag.Parse()
	_exe := strings.Split(exe, "|")
	_name := strings.Split(name, "|")
	if len(_exe) != len(_name) {
		fmt.Println("exe和name参数数量不一致")
		return
	}
	for i := 0; i < len(_exe); i++ {
		_exe[i] = path.Clean(_exe[i])
		_, err := os.Stat(_exe[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		var name = path.Base(_exe[i])
		if i < len(_name) {
			name = _name[i]
		}
		ExeList[path.Base(_exe[i])] = Exe{name, _exe[i], fmt.Sprintf("%s.jpg", _exe[i]), true}
	}
	ExeList[videoPlayer] = Exe{videoPlayer, videoPlayer, "", false}
	fmt.Println(port, ExeList)
	startServer()
}

func startServer() {
	http.HandleFunc("/", web)
	fmt.Println("服务器启动成功")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}

type resultMsgType interface{}

type Msg struct {
	Error   resultMsgType
	Success resultMsgType
}

func ico(exeAdress string, w http.ResponseWriter, r *http.Request) {
	icoPath := ExeList[exeAdress].Ico
	icoPath = path.Clean(icoPath)
	_, err := os.Stat(icoPath)
	//fmt.Println("找图标 ", exeAdress, icoPath)
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile(icoPath)
	if err != nil {
		return
	}
	w.Write(data)
}

func web(w http.ResponseWriter, r *http.Request) {
	var msg Msg
	fmt.Println(r.RequestURI, r.URL.Path)
	//status := r.URL.Query().Get("status")
	// exe := r.URL.Query().Get("exe")
	exe := r.PostFormValue("exe")
	if exe == "" {
		exe = r.URL.Query().Get("exe")
	}
	//fmt.Println(status, exe)
	if r.URL.Path != "/listapp" {
		if _, ok := ExeList[exe]; !ok {
			fmt.Println(exe, "程序不存在")
			msg.Error = fmt.Sprintf("%s程序不存在", exe)
			jsonData, _ := json.Marshal(msg)
			w.Write(jsonData)
			return
		}
	}
	//fmt.Println(ExeList[exe])
	var resultStatus error
	var resultMsg resultMsgType
	if r.URL.Path == "/controller/openapp" {
		resultStatus, resultMsg = runExe(ExeList[exe].Path)
	} else if r.URL.Path == "/controller/closeapp" {
		resultStatus, resultMsg = closeExe(exe)
	} else if r.URL.Path == "/icoapp" {
		ico(exe, w, r)
		return
	} else {
		resultStatus, resultMsg = listExe()
	}
	if resultStatus == nil {
		msg.Success = resultMsg
	} else {
		msg.Error = resultMsg
	}

	jsonData, _ := json.Marshal(msg)
	w.Write(jsonData)
}
