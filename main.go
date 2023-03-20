package main

import (
	"go-cron-check-payments/app"
	"go-cron-check-payments/kernel"
)

func main() {
	//setup and run the app
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}

	kernel.Run()
}
