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

// BarkerConfig holds configuration for Barker push notifications.
type BarkerConfig struct {
	// ServerURL is the Barker server endpoint URL.
	ServerURL string
	// Key is the client key for authenticating and routing the notification.
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

// BarkerNotifier implements the Notifier interface for Barker push notifications.
type BarkerNotifier struct {
	config BarkerConfig
	client *http.Client
}

// NewBarker creates a new BarkerNotifier with the provided configuration.
func NewBarker(config BarkerConfig) *BarkerNotifier {
	config.CommonConfig.applyDefaults()

	return &BarkerNotifier{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// barkerRequest represents the JSON payload for Barker API.
type barkerRequest struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Sound    string `json:"sound,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Group    string `json:"group,omitempty"`
	URL      string `json:"url,omitempty"`
	Level    string `json:"level,omitempty"`
	Badge    int    `json:"badge,omitempty"`
	AutoCopy string `json:"autoCopy,omitempty"`
	Copy     string `json:"copy,omitempty"`
}

// barkerResponse represents the JSON response from Barker API.
type barkerResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Send sends a push notification via Barker.
func (b *BarkerNotifier) Send(message Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.config.Timeout)
	defer cancel()

	var lastErr error
	attempts := b.config.RetryCount + 1

	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("barker send timeout: %w", ctx.Err())
			case <-time.After(b.config.RetryInterval):
			}
		}

		lastErr = b.sendOnce(ctx, message)
		if lastErr == nil {
			return nil
		}
	}

	return fmt.Errorf("barker send failed after %d attempts: %w", attempts, lastErr)
}

// sendOnce performs a single Barker send attempt.
func (b *BarkerNotifier) sendOnce(ctx context.Context, message Message) error {
	if b.config.ServerURL == "" {
		return fmt.Errorf("barker server URL is required")
	}
	if b.config.Key == "" {
		return fmt.Errorf("barker key is required")
	}

	req := barkerRequest{
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
		return fmt.Errorf("barker server returned status %d: %s", resp.StatusCode, string(body))
	}

	var barkerResp barkerResponse
	if err := json.Unmarshal(body, &barkerResp); err != nil {
		return nil
	}

	if barkerResp.Code != 200 {
		return fmt.Errorf("barker API error (code %d): %s", barkerResp.Code, barkerResp.Message)
	}

	return nil
}
