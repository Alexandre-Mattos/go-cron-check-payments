package kernel

import (
	"go-cron-check-payments/commands"
	"time"

	"github.com/go-co-op/gocron"
)

func Run() error {
	cron := gocron.NewScheduler(time.UTC)
	cron.Every(10).Minutes().Do(commands.CreatePayments)

	cron.StartBlocking()

	return nil
}
