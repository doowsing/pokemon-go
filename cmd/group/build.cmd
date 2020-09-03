set GOARCH=amd64
set GOOS=linux
go build -o group
cd ../build_tools/
buildtools.exe group ../group/group /www/wwwroot/gopoke/cmd/group/ startGroup.sh
cd ../group