package main

import (
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Error(err)
		return
	}

	if err := app.Start(); err != nil {
		log.Error(err)
		return
	}
}
