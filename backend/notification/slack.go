package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

type Slack struct {
	api *slack.Client
}

func NewSlack(token string) *Slack {
	return &Slack{
		api: slack.New(token),
	}

}

func (s *Slack) SendMessage(channel, header, body string) error {
	attachment := slack.Attachment{
		Pretext: header,
		Text:    body,
	}

	channelID, timestamp, err := s.api.PostMessage(channel, slack.MsgOptionText(`New Notification:`, false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	return nil
}
