set GOARCH=amd64
set GOOS=linux
go build -o pokemon
cd ../build_tools/
buildtools.exe pokemon ../poke/pokemon /www/wwwroot/gopoke/cmd/poke/ reStartGo.sh
cd ../poke