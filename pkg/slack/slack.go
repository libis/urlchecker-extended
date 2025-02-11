package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/libis/urlchecker-extended/pkg/config"
)

// SlackWebhookPayload represents the minimum required fields to send a webhook.
//
// https://api.slack.com/messaging/webhooks
type SlackWebhookPayload struct {
	Text   string `json:"text"`
	Mrkdwn bool   `json:"mrkdwn"`
}

// Message represents the status, url and message fields to send a webhook.
type Message struct {
	Status  int
	Url     string
	Message string
}

// SlackClient contains the Webhook URL.
type SlackClient struct {
	Webhook string
}

// SendMessage creates a SlackWebhookPayload from the concatenated messages
// and sends it to the Webhook URL.
func (c SlackClient) SendMessage(messages []Message) {
	repo := os.Getenv(config.EnvGithubRepo)

	// prepare a slice to hold all formatted messages
	var formattedMessages []string

	// loop over messages and format them
	for _, message := range messages {
		msg := fmt.Sprintf("URL: <%s>, Message: %s", message.Url, message.Message)
		formattedMessages = append(formattedMessages, msg)
	}

	// concatenate all formatted messages
	allMessages := fmt.Sprintf("Repository: %s\n\n%s", repo, strings.Join(formattedMessages, "\n\n"))

	pl := SlackWebhookPayload{
		Text: allMessages,
	}

	jsonPayload, _ := json.Marshal(pl)

	fmt.Printf("Sending message to " + c.Webhook)
	_, err := http.Post(c.Webhook, "application/json", bytes.NewBuffer(jsonPayload))

	if err != nil {
		log.Fatal(err)
	}
}
