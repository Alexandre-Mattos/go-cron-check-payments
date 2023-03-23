package main

import (
	"go-cron-check-payments/app"
	"go-cron-check-payments/logger"
	//"go-cron-check-payments/kernel"
)

func main() {
	//setup and run the app
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}

	errors := logger.Send("GOLANG ESTEVE AQUI")
	if errors != "" {
		panic(errors)
	}
	//kernel.Run()
}
