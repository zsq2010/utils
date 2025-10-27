package notify

import (
    "strings"
    "testing"
)

func TestEmailConfig_applyProviderDefaults(t *testing.T) {
    tests := []struct {
        name         string
        provider     EmailProvider
        expectedHost string
        expectedPort int
        expectedTLS  bool
    }{
        {
            name:         "QQMail defaults",
            provider:     QQMail,
            expectedHost: "smtp.qq.com",
            expectedPort: 587,
            expectedTLS:  true,
        },
        {
            name:         "Outlook defaults",
            provider:     Outlook,
            expectedHost: "smtp-mail.outlook.com",
            expectedPort: 587,
            expectedTLS:  true,
        },
        {
            name:         "Gmail defaults",
            provider:     Gmail,
            expectedHost: "smtp.gmail.com",
            expectedPort: 587,
            expectedTLS:  true,
        },
        {
            name:         "Custom no defaults",
            provider:     Custom,
            expectedHost: "",
            expectedPort: 0,
            expectedTLS:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            config := EmailConfig{Provider: tt.provider}
            config.applyProviderDefaults()

            if config.Host != tt.expectedHost {
                t.Errorf("Host = %v, want %v", config.Host, tt.expectedHost)
            }
            if config.Port != tt.expectedPort {
                t.Errorf("Port = %v, want %v", config.Port, tt.expectedPort)
            }
            if config.UseTLS != tt.expectedTLS {
                t.Errorf("UseTLS = %v, want %v", config.UseTLS, tt.expectedTLS)
            }
        })
    }
}

func TestEmailConfig_customHostPort(t *testing.T) {
    config := EmailConfig{
        Provider: QQMail,
        Host:     "custom.smtp.com",
        Port:     465,
    }
    config.applyProviderDefaults()

    if config.Host != "custom.smtp.com" {
        t.Errorf("Host = %v, want %v", config.Host, "custom.smtp.com")
    }
    if config.Port != 465 {
        t.Errorf("Port = %v, want %v", config.Port, 465)
    }
}

func TestNewEmail(t *testing.T) {
    config := EmailConfig{
        Provider: QQMail,
        Username: "user@qq.com",
        Password: "password",
        From:     "sender@qq.com",
        To:       []string{"recipient@example.com"},
    }

    notifier := NewEmail(config)
    if notifier == nil {
        t.Fatal("NewEmail returned nil")
    }

    if notifier.config.Host != "smtp.qq.com" {
        t.Errorf("Host = %v, want %v", notifier.config.Host, "smtp.qq.com")
    }
    if notifier.config.Port != 587 {
        t.Errorf("Port = %v, want %v", notifier.config.Port, 587)
    }
}

func TestEmailNotifier_buildMessage(t *testing.T) {
    notifier := NewEmail(EmailConfig{
        Provider: Gmail,
        From:     "sender@gmail.com",
        To:       []string{"recipient1@example.com", "recipient2@example.com"},
        CC:       []string{"cc@example.com"},
    })

    tests := []struct {
        name     string
        message  Message
        contains []string
    }{
        {
            name: "plain text message",
            message: Message{
                Title: "Test Subject",
                Body:  "Test Body Content",
            },
            contains: []string{
                "From: sender@gmail.com",
                "To: recipient1@example.com, recipient2@example.com",
                "Cc: cc@example.com",
                "Subject: Test Subject",
                "Test Body Content",
            },
        },
        {
            name: "HTML message",
            message: Message{
                Title:    "HTML Test",
                Body:     "Plain text version",
                HTMLBody: "<p>HTML version</p>",
            },
            contains: []string{
                "Subject: HTML Test",
                "multipart/alternative",
                "Plain text version",
                "<p>HTML version</p>",
            },
        },
        {
            name: "message with attachments",
            message: Message{
                Title:       "With Attachments",
                Body:        "Body",
                Attachments: []string{"/path/to/file.txt", "/path/to/doc.pdf"},
            },
            contains: []string{
                "Subject: With Attachments",
                "Attachment: file.txt",
                "Attachment: doc.pdf",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            msg := notifier.buildMessage(tt.message, notifier.config.From)
            msgStr := string(msg)

            for _, substr := range tt.contains {
                if !strings.Contains(msgStr, substr) {
                    t.Errorf("message does not contain %q", substr)
                }
            }
        })
    }
}

func TestEmailNotifier_sendOnce_noRecipients(t *testing.T) {
    notifier := NewEmail(EmailConfig{
        Provider: Gmail,
        Username: "user@gmail.com",
        Password: "password",
        From:     "sender@gmail.com",
        To:       []string{},
    })

    message := Message{
        Title: "Test",
        Body:  "Test Body",
    }

    err := notifier.sendOnce(message)
    if err == nil {
        t.Error("expected error for no recipients, got nil")
    }
    if !strings.Contains(err.Error(), "no recipients") {
        t.Errorf("error message = %v, want to contain 'no recipients'", err)
    }
}

func TestEmailNotifier_Send_integration(t *testing.T) {
    t.Skip("Skipping integration test - requires real SMTP server")

    notifier := NewEmail(EmailConfig{
        Provider: Custom,
        Host:     "localhost",
        Port:     1025,
        Username: "test",
        Password: "test",
        From:     "test@example.com",
        To:       []string{"recipient@example.com"},
    })

    message := Message{
        Title: "Test Email",
        Body:  "This is a test email",
    }

    err := notifier.Send(message)
    if err != nil {
        t.Errorf("Send failed: %v", err)
    }
}

func TestEmailNotifier_buildMessage_noCc(t *testing.T) {
    notifier := NewEmail(EmailConfig{
        Provider: Gmail,
        From:     "sender@gmail.com",
        To:       []string{"recipient@example.com"},
    })

    message := Message{
        Title: "Test",
        Body:  "Body",
    }

    msg := notifier.buildMessage(message, notifier.config.From)
    msgStr := string(msg)

    if strings.Contains(msgStr, "Cc:") {
        t.Error("message should not contain Cc header when no CC recipients")
    }
}

func TestEmailNotifier_recipients(t *testing.T) {
    notifier := NewEmail(EmailConfig{
        Provider: Gmail,
        From:     "sender@gmail.com",
        To:       []string{"to1@example.com", "to2@example.com"},
        CC:       []string{"cc@example.com"},
        BCC:      []string{"bcc@example.com"},
    })

    message := Message{
        Title: "Test",
        Body:  "Body",
    }

    msg := notifier.buildMessage(message, notifier.config.From)
    msgStr := string(msg)

    if !strings.Contains(msgStr, "To: to1@example.com, to2@example.com") {
        t.Error("message should contain all To recipients")
    }
    if !strings.Contains(msgStr, "Cc: cc@example.com") {
        t.Error("message should contain CC recipients")
    }
}

func TestEmailNotifier_fromDefault(t *testing.T) {
    notifier := NewEmail(EmailConfig{
        Provider: Gmail,
        Username: "user@gmail.com",
        To:       []string{"recipient@example.com"},
    })

    message := Message{
        Title: "Test",
        Body:  "Body",
    }

    msg := notifier.buildMessage(message, notifier.config.Username)
    msgStr := string(msg)

    if !strings.Contains(msgStr, "From: user@gmail.com") {
        t.Error("should use Username as From when From is empty")
    }
}

func TestEmailProvider_constants(t *testing.T) {
    providers := []EmailProvider{QQMail, Outlook, Gmail, Custom}
    if len(providers) != 4 {
        t.Error("should have 4 email providers")
    }
}

func TestEmailConfig_SSL(t *testing.T) {
    config := EmailConfig{
        Provider: QQMail,
        Port:     465,
        UseSSL:   true,
    }
    config.applyProviderDefaults()

    if !config.UseSSL {
        t.Error("UseSSL should be preserved")
    }
    if config.Port != 465 {
        t.Error("Port 465 should be preserved for SSL")
    }
}
