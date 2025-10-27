package notify

import (
    "context"
    "crypto/tls"
    "fmt"
    "net"
    "net/smtp"
    "path/filepath"
    "strings"
    "time"
)

// EmailProvider represents a preset email service provider.
type EmailProvider int

const (
    // QQMail is the QQ email service provider (smtp.qq.com).
    QQMail EmailProvider = iota
    // Outlook is the Outlook/Hotmail email service provider (smtp-mail.outlook.com).
    Outlook
    // Gmail is the Gmail email service provider (smtp.gmail.com).
    Gmail
    // Custom indicates a custom SMTP server configuration.
    Custom
)

// EmailConfig holds configuration for email notifications.
type EmailConfig struct {
    // Provider specifies the email service provider (QQMail, Outlook, Gmail, or Custom).
    Provider EmailProvider
    // Host is the SMTP server hostname (required for Custom provider).
    Host string
    // Port is the SMTP server port (required for Custom provider).
    Port int
    // UseTLS specifies whether to use STARTTLS for the connection.
    UseTLS bool
    // UseSSL specifies whether to use implicit TLS/SSL.
    UseSSL bool
    // Username is the SMTP authentication username.
    Username string
    // Password is the SMTP authentication password or authorization code.
    Password string
    // From is the sender email address.
    From string
    // To is the list of primary recipient email addresses.
    To []string
    // CC is the list of carbon copy recipient email addresses.
    CC []string
    // BCC is the list of blind carbon copy recipient email addresses.
    BCC []string
    // Common configuration options.
    CommonConfig
}

// EmailNotifier implements the Notifier interface for email notifications.
type EmailNotifier struct {
    config EmailConfig
}

// NewEmail creates a new EmailNotifier with the provided configuration.
func NewEmail(config EmailConfig) *EmailNotifier {
    config.CommonConfig.applyDefaults()
    config.applyProviderDefaults()
    return &EmailNotifier{config: config}
}

// applyProviderDefaults sets default host, port, and TLS settings based on the provider.
func (c *EmailConfig) applyProviderDefaults() {
    switch c.Provider {
    case QQMail:
        if c.Host == "" {
            c.Host = "smtp.qq.com"
        }
        if c.Port == 0 {
            c.Port = 587
            c.UseTLS = true
        }
    case Outlook:
        if c.Host == "" {
            c.Host = "smtp-mail.outlook.com"
        }
        if c.Port == 0 {
            c.Port = 465
            c.UseSSL = true
        }
    case Gmail:
        if c.Host == "" {
            c.Host = "smtp.gmail.com"
        }
        if c.Port == 0 {
            c.Port = 465
            c.UseSSL = true
        }
    }
}

// Send sends an email notification.
func (e *EmailNotifier) Send(message Message) error {
    ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
    defer cancel()

    var lastErr error
    attempts := e.config.RetryCount + 1

    for i := 0; i < attempts; i++ {
        if i > 0 {
            select {
            case <-ctx.Done():
                return fmt.Errorf("email send timeout: %w", ctx.Err())
            case <-time.After(e.config.RetryInterval):
            }
        }

        lastErr = e.sendOnce(message)
        if lastErr == nil {
            return nil
        }
    }

    return fmt.Errorf("email send failed after %d attempts: %w", attempts, lastErr)
}

// sendOnce performs a single email send attempt.
func (e *EmailNotifier) sendOnce(message Message) error {
    if len(e.config.To) == 0 {
        return fmt.Errorf("no recipients specified")
    }

    from := e.config.From
    if from == "" {
        from = e.config.Username
    }

    recipients := append([]string{}, e.config.To...)
    recipients = append(recipients, e.config.CC...)
    recipients = append(recipients, e.config.BCC...)

    msg := e.buildMessage(message, from)

    addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)

    if e.config.UseSSL {
        return e.sendWithSSL(addr, from, recipients, msg)
    }

    if e.config.UseTLS {
        return e.sendWithTLS(addr, from, recipients, msg)
    }

    return e.sendPlain(addr, from, recipients, msg)
}

// buildMessage constructs the email message in RFC 5322 format.
func (e *EmailNotifier) buildMessage(message Message, from string) []byte {
    var builder strings.Builder

    builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
    builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(e.config.To, ", ")))

    if len(e.config.CC) > 0 {
        builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(e.config.CC, ", ")))
    }

    builder.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Title))

    if message.HTMLBody != "" {
        boundary := "boundary-notify-email"
        builder.WriteString("MIME-Version: 1.0\r\n")
        builder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
        builder.WriteString("\r\n")

        builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
        builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
        builder.WriteString("\r\n")
        builder.WriteString(message.Body)
        builder.WriteString("\r\n\r\n")

        builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
        builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
        builder.WriteString("\r\n")
        builder.WriteString(message.HTMLBody)
        builder.WriteString("\r\n\r\n")

        builder.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
    } else {
        builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
        builder.WriteString("\r\n")
        builder.WriteString(message.Body)
    }

    if len(message.Attachments) > 0 {
        for _, att := range message.Attachments {
            builder.WriteString(fmt.Sprintf("\r\nAttachment: %s\r\n", filepath.Base(att)))
        }
    }

    return []byte(builder.String())
}

// sendWithTLS sends email using STARTTLS.
func (e *EmailNotifier) sendWithTLS(addr, from string, recipients []string, msg []byte) error {
    host, _, _ := net.SplitHostPort(addr)

    auth := smtp.PlainAuth("", e.config.Username, e.config.Password, host)

    client, err := smtp.Dial(addr)
    if err != nil {
        return fmt.Errorf("dial SMTP server: %w", err)
    }
    defer client.Close()

    if err := client.StartTLS(&tls.Config{ServerName: host}); err != nil {
        return fmt.Errorf("start TLS: %w", err)
    }

    if err := client.Auth(auth); err != nil {
        return fmt.Errorf("SMTP authentication: %w", err)
    }

    if err := client.Mail(from); err != nil {
        return fmt.Errorf("set sender: %w", err)
    }

    for _, recipient := range recipients {
        if err := client.Rcpt(recipient); err != nil {
            return fmt.Errorf("add recipient %s: %w", recipient, err)
        }
    }

    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("open data writer: %w", err)
    }

    if _, err := w.Write(msg); err != nil {
        return fmt.Errorf("write message: %w", err)
    }

    if err := w.Close(); err != nil {
        return fmt.Errorf("close data writer: %w", err)
    }

    return client.Quit()
}

// sendWithSSL sends email using implicit TLS/SSL.
func (e *EmailNotifier) sendWithSSL(addr, from string, recipients []string, msg []byte) error {
    host, _, _ := net.SplitHostPort(addr)

    tlsConfig := &tls.Config{ServerName: host}
    conn, err := tls.Dial("tcp", addr, tlsConfig)
    if err != nil {
        return fmt.Errorf("dial with SSL: %w", err)
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, host)
    if err != nil {
        return fmt.Errorf("create SMTP client: %w", err)
    }
    defer client.Close()

    auth := smtp.PlainAuth("", e.config.Username, e.config.Password, host)
    if err := client.Auth(auth); err != nil {
        return fmt.Errorf("SMTP authentication: %w", err)
    }

    if err := client.Mail(from); err != nil {
        return fmt.Errorf("set sender: %w", err)
    }

    for _, recipient := range recipients {
        if err := client.Rcpt(recipient); err != nil {
            return fmt.Errorf("add recipient %s: %w", recipient, err)
        }
    }

    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("open data writer: %w", err)
    }

    if _, err := w.Write(msg); err != nil {
        return fmt.Errorf("write message: %w", err)
    }

    if err := w.Close(); err != nil {
        return fmt.Errorf("close data writer: %w", err)
    }

    return client.Quit()
}

// sendPlain sends email without TLS (not recommended for production).
func (e *EmailNotifier) sendPlain(addr, from string, recipients []string, msg []byte) error {
    host, _, _ := net.SplitHostPort(addr)
    auth := smtp.PlainAuth("", e.config.Username, e.config.Password, host)

    return smtp.SendMail(addr, auth, from, recipients, msg)
}
