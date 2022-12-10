::SET output="tank.exe"
SET output="C:\Users\xiangqian\Desktop\tmp\tank\tank.exe"
go build -ldflags="-s -w" -o %output%
::upx -9 --brute %output%
pause