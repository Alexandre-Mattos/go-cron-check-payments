package kernel

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func Run() error {
	cron := gocron.NewScheduler(time.UTC)
	var count int = 0
	cron.Every(10).Seconds().Do(func() {
		fmt.Println("Vezes rodado: ", count)
		count++
	})

	cron.StartBlocking()

	return nil
}
