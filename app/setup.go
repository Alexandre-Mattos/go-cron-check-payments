package app

import (
	"go-cron-check-payments/config"
)

func SetupAndRunApp() error {
	err := config.LoadEnv()
	if err != nil {
		return err
	}
	return nil
}
