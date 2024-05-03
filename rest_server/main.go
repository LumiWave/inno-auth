package main

import (
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	if err := app.Start(); err != nil {
		log.Errorf("%v", err)
		return
	}
}
