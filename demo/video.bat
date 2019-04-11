@echo off
if "%1" == "h" goto begin
mshta vbscript:createobject("wscript.shell").run("%~nx0 h",0)(window.close)&&exit
:begin
"PotPlayer\PotPlayerMini.exe" d:/eihoo_video/2.mp4