.PHONY: build init config docs clean run test stop stat
default: clean init config docs build

# set variables for build information.
buildTime=$(shell date -u "+%Y%m%d%I%M")
revVersion=$(shell git rev-parse --short HEAD)

server_int_api=internal_api.yml
server_ext_api=external_api.yml

# set variables for build
output=bin
app=inno-auth-$(env)
src=cmd/$(app)

ifeq ("$(env)", "")
env=local
endif

# set variables for dependencies

build: init
	# build $(app)
	go build -ldflags "-X github.com/ONBUFF-IP-TOKEN/baseapp/base.AppVersion=$(revVersion).$(buildTime)" -o $(output)/$(app) rest_server/main.go

init:
	# initialize output directory.
	@if [ -d $(output) ]; then rm -rf $(output); fi;
	@if [ ! -d $(output) ]; then mkdir -p $(output); mkdir -p $(output)/logs; fi;
	@if [ ! -d $(output) ]; then mkdir -p $(output); mkdir -p $(output)/docs; fi;
	@if [ ! -d $(output)/docs ]; then mkdir -p $(output); mkdir -p $(output)/docs/ext; fi;

run: stop
	# run in background mode
	- cd bin && ./$(app) -c=config.yml &

stat:
	- @pgrep $(app)
stop:
	- @pkill -9 -f $(app)

test: 
	# test
	# If the test succeeds, kill the test process and then do the next build,
	# if it fails, kill the test process and exit the build.
	go test github.com/ONBUFF-IP-TOKEN/inno-auth -v -timeout 30s \
		&& { echo test success; make stop; } \
		|| { echo test failure; make stop; exit 1;}

deploy: test
	# If the test is successful, deploy it to the development environment.

config: init
	# config file copy
	cp etc/conf/config.$(env).yml $(output)/config.yml
	cp etc/conf/$(server_int_api) $(output)/$(server_int_api)
	cp etc/conf/$(server_ext_api) $(output)/$(server_ext_api)
	cp etc/onbuffcerti.crt bin/onbuffcerti.crt
	cp etc/onbuffcerti.key bin/onbuffcerti.key
	cp etc/swagger/ext/*.* bin/docs/ext
	
docs: init
	# copy the API docs to the output directory.
	# @if [ ! -d $(output)/docs/int ]; then mkdir -p $(output)/docs/int; fi;
	# @if [ ! -d $(output)/docs/ext ]; then mkdir -p $(output)/docs/ext; fi;
	# cp -R etc/swagger/int/* bin/docs/int/
	# cp -R etc/swagger/ext/* bin/docs/ext/

clean:
	# clean output
	rm -rf $(output)

