# http_run_exe
通过http运行主机上的exe




### 启动方法
http_run_exe.exe -exe="c:/a.exe|d:/b.exe" -name="程序名称|程序名称" -port=8081  -single=1

`多个程序使用|隔开，需要将\换成/`

### 命令调用方法

`打开`
http://localhost:8081/open?exe=a.exe


`关闭`
http://localhost:8081/close?exe=a.exe


`列表`
http://localhost:8081/list


`缩略图`
http://localhost:8081/ico?exe=a.exe
