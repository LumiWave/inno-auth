# inno-auth
* inno platform auth service
 This is the content of building the api server required for development.
 
# build & run
## env info
* env : "local" | "dev" | "stage" | "live"

## windows build & run
* cmd/inno-auth> ./windows_build.sh {env}
  * ex) ./windows_build.sh stage

## linux build & run
* build
  * root@inno-auth> make env={env}
    * ex) make env=stage
* run
  * root@inno-auth> make run env={env}
    * ex) make run env=stage
