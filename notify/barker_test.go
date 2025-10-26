package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewBarker(t *testing.T) {
	config := BarkerConfig{
		ServerURL: "https://barker.example.com",
		Key:       "test-key",
		Sound:     "default",
	}

	notifier := NewBarker(config)
	if notifier == nil {
		t.Fatal("NewBarker returned nil")
	}

	if notifier.config.ServerURL != "https://barker.example.com" {
		t.Errorf("ServerURL = %v, want %v", notifier.config.ServerURL, "https://barker.example.com")
	}
	if notifier.config.Key != "test-key" {
		t.Errorf("Key = %v, want %v", notifier.config.Key, "test-key")
	}
	if notifier.config.Timeout == 0 {
		t.Error("Timeout should be set to default")
	}
}

func TestBarkerNotifier_Send_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Method = %v, want POST", r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("Content-Type = %v, want application/json", contentType)
		}

		var req barkerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if req.Title != "Test Title" {
			t.Errorf("Title = %v, want %v", req.Title, "Test Title")
		}
		if req.Body != "Test Body" {
			t.Errorf("Body = %v, want %v", req.Body, "Test Body")
		}

		resp := barkerResponse{
			Code:      200,
			Message:   "Success",
			Timestamp: time.Now().Unix(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	notifier := NewBarker(BarkerConfig{
		ServerURL: server.URL,
		Key:       "test-key",
	})

	message := Message{
		Title: "Test Title",
		Body:  "Test Body",
	}

	err := notifier.Send(message)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestBarkerNotifier_Send_withPriority(t *testing.T) {
	tests := []struct {
		name          string
		priority      string
		expectedLevel string
	}{
		{
			name:          "high priority",
			priority:      "high",
			expectedLevel: "timeSensitive",
		},
		{
			name:          "urgent priority",
			priority:      "urgent",
			expectedLevel: "timeSensitive",
		},
		{
			name:          "low priority",
			priority:      "low",
			expectedLevel: "passive",
		},
		{
			name:          "normal priority",
			priority:      "normal",
			expectedLevel: "active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req barkerRequest
				json.NewDecoder(r.Body).Decode(&req)

				if req.Level != tt.expectedLevel {
					t.Errorf("Level = %v, want %v", req.Level, tt.expectedLevel)
				}

				resp := barkerResponse{Code: 200, Message: "Success"}
				json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			notifier := NewBarker(BarkerConfig{
				ServerURL: server.URL,
				Key:       "test-key",
			})

			message := Message{
				Title:    "Test",
				Body:     "Body",
				Priority: tt.priority,
			}

			err := notifier.Send(message)
			if err != nil {
				t.Errorf("Send failed: %v", err)
			}
		})
	}
}

func TestBarkerNotifier_Send_withExtra(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req barkerRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Sound != "custom-sound" {
			t.Errorf("Sound = %v, want %v", req.Sound, "custom-sound")
		}
		if req.Icon != "custom-icon" {
			t.Errorf("Icon = %v, want %v", req.Icon, "custom-icon")
		}
		if req.Group != "custom-group" {
			t.Errorf("Group = %v, want %v", req.Group, "custom-group")
		}
		if req.URL != "https://example.com" {
			t.Errorf("URL = %v, want %v", req.URL, "https://example.com")
		}
		if req.Badge != 5 {
			t.Errorf("Badge = %v, want %v", req.Badge, 5)
		}

		resp := barkerResponse{Code: 200, Message: "Success"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	notifier := NewBarker(BarkerConfig{
		ServerURL: server.URL,
		Key:       "test-key",
	})

	message := Message{
		Title: "Test",
		Body:  "Body",
		Extra: map[string]interface{}{
			"sound": "custom-sound",
			"icon":  "custom-icon",
			"group": "custom-group",
			"url":   "https://example.com",
			"badge": 5,
		},
	}

	err := notifier.Send(message)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestBarkerNotifier_Send_serverError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	notifier := NewBarker(BarkerConfig{
		ServerURL: server.URL,
		Key:       "test-key",
	})

	message := Message{
		Title: "Test",
		Body:  "Body",
	}

	err := notifier.Send(message)
	if err == nil {
		t.Error("expected error for server error, got nil")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("error message = %v, want to contain 'status 500'", err)
	}
}

func TestBarkerNotifier_Send_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := barkerResponse{
			Code:    400,
			Message: "Invalid request",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	notifier := NewBarker(BarkerConfig{
		ServerURL: server.URL,
		Key:       "test-key",
	})

	message := Message{
		Title: "Test",
		Body:  "Body",
	}

	err := notifier.Send(message)
	if err == nil {
		t.Error("expected error for API error, got nil")
	}
	if !strings.Contains(err.Error(), "code 400") {
		t.Errorf("error message = %v, want to contain 'code 400'", err)
	}
}

func TestBarkerNotifier_Send_missingConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      BarkerConfig
		wantErrText string
	}{
		{
			name: "missing server URL",
			config: BarkerConfig{
				Key: "test-key",
			},
			wantErrText: "server URL is required",
		},
		{
			name: "missing key",
			config: BarkerConfig{
				ServerURL: "https://example.com",
			},
			wantErrText: "key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := NewBarker(tt.config)
			message := Message{Title: "Test", Body: "Body"}

			err := notifier.Send(message)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErrText) {
				t.Errorf("error message = %v, want to contain %q", err, tt.wantErrText)
			}
		})
	}
}

func TestBarkerNotifier_Send_withRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		resp := barkerResponse{Code: 200, Message: "Success"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	notifier := NewBarker(BarkerConfig{
		ServerURL: server.URL,
		Key:       "test-key",
		CommonConfig: CommonConfig{
			RetryCount:    3,
			RetryInterval: 10 * time.Millisecond,
		},
	})

	message := Message{Title: "Test", Body: "Body"}

	err := notifier.Send(message)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}
