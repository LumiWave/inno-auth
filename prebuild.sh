rm -rf bin/*

mkdir -p bin

cp ./etc/conf/config.$1.yml ./bin/config.yml
cp ./etc/conf/external_api.yml ./bin
cp ./etc/conf/internal_api.yml ./bin
cp ./etc/onbuffcerti.crt ./bin
cp ./etc/onbuffcerti.key ./bin
cp ./etc/IP2LOCATION-LITE-DB1.BIN ./bin
