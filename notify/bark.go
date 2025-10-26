package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BarkConfig holds configuration for Bark push notifications.
// Bark is an iOS notification service (https://api.day.app/).
type BarkConfig struct {
	// ServerURL is the Bark server endpoint URL (default: https://api.day.app).
	ServerURL string
	// Key is the device key for routing the notification to your iOS device.
	Key string
	// Sound is the notification sound name (optional).
	Sound string
	// Icon is the notification icon URL (optional).
	Icon string
	// Group is the notification group name for organizing messages (optional).
	Group string
	// URL is a URL to open when the notification is tapped (optional).
	URL string
	// Common configuration options.
	CommonConfig
}

// BarkNotifier implements the Notifier interface for Bark push notifications.
type BarkNotifier struct {
	config BarkConfig
	client *http.Client
}

// NewBark creates a new BarkNotifier with the provided configuration.
// If ServerURL is empty, it defaults to https://api.day.app.
func NewBark(config BarkConfig) *BarkNotifier {
	config.CommonConfig.applyDefaults()

	if config.ServerURL == "" {
		config.ServerURL = "https://api.day.app"
	}

	return &BarkNotifier{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// barkRequest represents the JSON payload for Bark API.
type barkRequest struct {
	Title    string `json:"title,omitempty"`
	Body     string `json:"body"`
	Sound    string `json:"sound,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Group    string `json:"group,omitempty"`
	URL      string `json:"url,omitempty"`
	Level    string `json:"level,omitempty"`
	Badge    int    `json:"badge,omitempty"`
	AutoCopy string `json:"autoCopy,omitempty"`
	Copy     string `json:"copy,omitempty"`
	IsArchive int   `json:"isArchive,omitempty"`
}

// barkResponse represents the JSON response from Bark API.
type barkResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Send sends a push notification via Bark.
func (b *BarkNotifier) Send(message Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.config.Timeout)
	defer cancel()

	var lastErr error
	attempts := b.config.RetryCount + 1

	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("bark send timeout: %w", ctx.Err())
			case <-time.After(b.config.RetryInterval):
			}
		}

		lastErr = b.sendOnce(ctx, message)
		if lastErr == nil {
			return nil
		}
	}

	return fmt.Errorf("bark send failed after %d attempts: %w", attempts, lastErr)
}

// sendOnce performs a single Bark send attempt.
func (b *BarkNotifier) sendOnce(ctx context.Context, message Message) error {
	if b.config.Key == "" {
		return fmt.Errorf("bark key is required")
	}

	req := barkRequest{
		Title: message.Title,
		Body:  message.Body,
		Sound: b.config.Sound,
		Icon:  b.config.Icon,
		Group: b.config.Group,
		URL:   b.config.URL,
	}

	if message.Priority != "" {
		switch message.Priority {
		case "high", "urgent":
			req.Level = "timeSensitive"
		case "low":
			req.Level = "passive"
		default:
			req.Level = "active"
		}
	}

	if message.Extra != nil {
		if sound, ok := message.Extra["sound"].(string); ok && sound != "" {
			req.Sound = sound
		}
		if icon, ok := message.Extra["icon"].(string); ok && icon != "" {
			req.Icon = icon
		}
		if group, ok := message.Extra["group"].(string); ok && group != "" {
			req.Group = group
		}
		if url, ok := message.Extra["url"].(string); ok && url != "" {
			req.URL = url
		}
		if badge, ok := message.Extra["badge"].(int); ok {
			req.Badge = badge
		}
		if autoCopy, ok := message.Extra["autoCopy"].(string); ok && autoCopy != "" {
			req.AutoCopy = autoCopy
		}
		if copy, ok := message.Extra["copy"].(string); ok && copy != "" {
			req.Copy = copy
		}
		if isArchive, ok := message.Extra["isArchive"].(int); ok {
			req.IsArchive = isArchive
		}
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/%s", b.config.ServerURL, b.config.Key)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := b.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark server returned status %d: %s", resp.StatusCode, string(body))
	}

	var barkResp barkResponse
	if err := json.Unmarshal(body, &barkResp); err != nil {
		return nil
	}

	if barkResp.Code != 200 {
		return fmt.Errorf("bark API error (code %d): %s", barkResp.Code, barkResp.Message)
	}

	return nil
}
