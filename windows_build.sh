set -x

sh ./prebuild.sh $1

rm -rf bin/inno-auth

go build -o bin/inno-auth.exe main.go

cd bin
./inno-auth.exe -c=config.yml