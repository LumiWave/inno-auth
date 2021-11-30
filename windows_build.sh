# set -x

sh prebuild.sh

go build -o bin/inno-auth.exe rest_server/main.go

cd bin
./inno-auth.exe -c=config.yml