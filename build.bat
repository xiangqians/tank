::set output="tank.exe"
set output="C:\Users\xiangqian\Desktop\tmp\tank\tank.exe"
cd ./src && go build -ldflags="-s -w" -o %output%
upx -9 --brute %output%
pause