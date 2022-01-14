set -x

rm -rf bin
mkdir bin

sh ./prebuild.sh $1

go build -o bin/inno-auth.exe rest_server/main.go

cd bin
./inno-auth.exe -c=config.yml