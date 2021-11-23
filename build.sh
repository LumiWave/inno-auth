set -x

sh ./prebuild.sh

go build -o bin/inno-auth rest_server/main.go

cd bin
./inno-auth -c=config.yml