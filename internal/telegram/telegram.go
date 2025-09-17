package telegram

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	token  string
	chatID string
}

func NewClient(token, chatID string) *Client {
	return &Client{token: token, chatID: chatID}
}

func (c *Client) SendMessage(message string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	resp, err := http.PostForm(endpoint, url.Values{
		"chat_id":    {c.chatID},
		"text":       {message},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", resp.Status)
	}
	return nil
}
