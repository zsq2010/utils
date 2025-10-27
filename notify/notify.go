// Package notify provides a unified interface for sending notifications across multiple channels.
package notify

import (
	"time"
)

// Notifier is the core interface that all notification channels must implement.
type Notifier interface {
	Send(message Message) error
}

// Message represents a notification message with common fields supported by most channels.
type Message struct {
	// Title is the subject or headline of the notification.
	Title string
	// Body is the main content of the notification.
	Body string
	// Priority indicates the importance of the notification (e.g., "high", "normal", "low").
	// Different channels may interpret this differently.
	Priority string
	// HTMLBody is an optional HTML-formatted version of the body (for channels that support it).
	HTMLBody string
	// Attachments is a list of file paths to attach (for channels that support attachments).
	Attachments []string
	// Extra holds channel-specific additional data.
	Extra map[string]interface{}
}

// CommonConfig holds common configuration options for all notifiers.
type CommonConfig struct {
	// Timeout specifies the maximum duration for send operations.
	// If zero, a default timeout of 30 seconds is used.
	Timeout time.Duration
	// RetryCount specifies how many times to retry on failure.
	RetryCount int
	// RetryInterval specifies the delay between retry attempts.
	RetryInterval time.Duration
}

// applyDefaults sets default values for common configuration.
func (c *CommonConfig) applyDefaults() {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	if c.RetryInterval == 0 {
		c.RetryInterval = 2 * time.Second
	}
}
