@echo off
if "%1" == "h" goto begin
mshta vbscript:createobject("wscript.shell").run("%~nx0 h",0)(window.close)&&exit
:begin
"C:\Program Files\DAUM\PotPlayer\PotPlayerMini64.exe" d:/eihoo_video/2.mp4