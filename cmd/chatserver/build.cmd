set GOARCH=amd64
set GOOS=linux
go build -o chatserver
cd ../build_tools/
buildtools.exe chatserver ../chatserver/chatserver /www/wwwroot/gopoke/cmd/chatserver/ reStartChat.sh
cd ../chatserver