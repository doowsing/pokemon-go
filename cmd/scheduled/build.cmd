set GOARCH=amd64
set GOOS=linux
go build -o scheduled
cd ../build_tools/
buildtools.exe scheduled ../scheduled/scheduled /www/wwwroot/gopoke/cmd/scheduled/ startTimer.sh
cd ../scheduled