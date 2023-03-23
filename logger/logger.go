package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

func Send(message string) string {
	url := os.Getenv("SLACK_WEBHOOK")

	fmt.Println(url)
	payload := slack.Payload{
		Text:      message,
		Username:  "GO-CREATE-PAYMENTS",
		Channel:   "#" + os.Getenv("SLACK_CHANNEL"),
		IconEmoji: ":rosto-saudando:",
	}
	fmt.Println(payload)

	errors := slack.Send(url, "", payload)
	if len(errors) > 0 {
		var errorStrings []string
		for _, err := range errors {
			errorStrings = append(errorStrings, err.Error())
		}
		return strings.Join(errorStrings, "\n")
	}

	return ""
}
