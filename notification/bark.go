package notification

import (
	"fmt"
	"net/http"
	"net/url"
)

// SendBarkNotification sends a notification to the Bark app.
func SendBarkNotification(barkKey, title, body string) error {
	barkURL := fmt.Sprintf("https://api.day.app/%s/%s/%s", barkKey, url.QueryEscape(title), url.QueryEscape(body))

	resp, err := http.Get(barkURL)
	if err != nil {
		return fmt.Errorf("failed to send Bark notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send Bark notification, status code: %d", resp.StatusCode)
	}
	return nil
}
