package notify

import (
	"testing"
	"time"
)

func TestCommonConfig_applyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    CommonConfig
		expected CommonConfig
	}{
		{
			name:  "empty config gets defaults",
			input: CommonConfig{},
			expected: CommonConfig{
				Timeout:       30 * time.Second,
				RetryInterval: 2 * time.Second,
			},
		},
		{
			name: "custom values preserved",
			input: CommonConfig{
				Timeout:       10 * time.Second,
				RetryCount:    3,
				RetryInterval: 5 * time.Second,
			},
			expected: CommonConfig{
				Timeout:       10 * time.Second,
				RetryCount:    3,
				RetryInterval: 5 * time.Second,
			},
		},
		{
			name: "partial custom values with defaults",
			input: CommonConfig{
				Timeout: 15 * time.Second,
			},
			expected: CommonConfig{
				Timeout:       15 * time.Second,
				RetryInterval: 2 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.input
			config.applyDefaults()

			if config.Timeout != tt.expected.Timeout {
				t.Errorf("Timeout = %v, want %v", config.Timeout, tt.expected.Timeout)
			}
			if config.RetryCount != tt.expected.RetryCount {
				t.Errorf("RetryCount = %v, want %v", config.RetryCount, tt.expected.RetryCount)
			}
			if config.RetryInterval != tt.expected.RetryInterval {
				t.Errorf("RetryInterval = %v, want %v", config.RetryInterval, tt.expected.RetryInterval)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Title:    "Test Title",
		Body:     "Test Body",
		Priority: "high",
		HTMLBody: "<p>Test HTML</p>",
		Attachments: []string{
			"/path/to/file1.txt",
			"/path/to/file2.pdf",
		},
		Extra: map[string]interface{}{
			"custom_field": "custom_value",
		},
	}

	if msg.Title != "Test Title" {
		t.Errorf("Title = %v, want %v", msg.Title, "Test Title")
	}
	if msg.Body != "Test Body" {
		t.Errorf("Body = %v, want %v", msg.Body, "Test Body")
	}
	if msg.Priority != "high" {
		t.Errorf("Priority = %v, want %v", msg.Priority, "high")
	}
	if msg.HTMLBody != "<p>Test HTML</p>" {
		t.Errorf("HTMLBody = %v, want %v", msg.HTMLBody, "<p>Test HTML</p>")
	}
	if len(msg.Attachments) != 2 {
		t.Errorf("len(Attachments) = %v, want %v", len(msg.Attachments), 2)
	}
	if msg.Extra["custom_field"] != "custom_value" {
		t.Errorf("Extra[custom_field] = %v, want %v", msg.Extra["custom_field"], "custom_value")
	}
}
