echo "git inno-auth up-to-date"
git pull

echo "remove directory bin"
rm -rf bin/*

echo "make directory bin"
cd ../config
sh cp_inno_auth.sh

cd ../inno-auth
mkdir -p bin

echo "copy to config yml"
cp ./etc/conf/config.stage.yml ./bin/config.yml
cp ./etc/conf/external_api.yml ./bin
cp ./etc/conf/internal_api.yml ./bin
cp ./etc/onbuffcerti.crt ./bin
cp ./etc/onbuffcerti.key ./bin