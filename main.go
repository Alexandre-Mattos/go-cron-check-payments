package main

import (
	"go-cron-check-payments/app"
	"go-cron-check-payments/kernel"
	//"go-cron-check-payments/kernel"
)

func main() {
	//setup and run the app
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}

	/* db, err := database.Connect()
	if err != nil {
		logger.Send(err.Error(), "debug")
	}

	var conta models.Multa
	db.Model(&models.Multa{}).
		Where("id = ?", 1).
		Find(&conta)

	fmt.Println(conta) */
	/* err = logger.Send("Se você conseguiu ler isso, quer dizer que o Alexandre conseguiu integrar o golang com o slack", "success")

	if err != nil {
		panic(err)
	}
	*/
	kernel.Run()
}
