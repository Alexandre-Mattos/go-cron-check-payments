package logger

import (
	"os"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
)

func WarningLevelEmoji(warning string) string {
	switch warning {
	case "warning":
		return ":weary:"
	case "error":
		return ":rage:"
	case "debug":
		return ":sleeping:"
	default:
		return ":sunglasses:"

	}
}

func WarningLevelColor(warning string) *string {
	switch warning {
	case "warning":
		color := "#FBFF00"
		return &color
	case "error":
		color := "#B60000"
		return &color
	case "debug":
		color := "#FFFFFF"
		return &color
	default:
		color := "#2AFF00"
		return &color
	}
}

func Send(text string, warning string) error {
	//slack incoming webhook url
	webhookUrl := os.Getenv("SLACK_WEBHOOK")

	footer := os.Getenv("GO_ENV") + " Golang Log | " + time.Now().Format(time.ANSIC)

	attachment := slack.Attachment{
		Color:  WarningLevelColor(warning),
		Footer: &footer,
	}

	attachment.AddField(slack.Field{Title: "Message", Value: "[" + os.Getenv("GO_ENV") + "] " + text}).AddField(slack.Field{Title: "Status", Value: warning})
	payload := slack.Payload{
		Username:    os.Getenv("SLACK_USERNAME"),
		IconEmoji:   WarningLevelEmoji(warning),
		Attachments: []slack.Attachment{attachment},
	}

	errors := slack.Send(webhookUrl, "", payload)
	if (errors) != nil {
		for err := range errors {
			panic(err)
		}
	}

	return nil
}
