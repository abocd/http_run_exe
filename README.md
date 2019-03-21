# http_run_exe
通过http运行主机上的exe




### 启动方法
http_run_exe.exe -exe="c:/a.exe|d:/b.exe" -port=8081

`多个程序使用|隔开，需要将\换成/`

### 命令调用方法

`打开`
http://localhost:8081/?status=open&exe=a.exe


`关闭`
http://localhost:8081/?status=close&exe=a.exe
