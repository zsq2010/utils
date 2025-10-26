package main

import (
    "fmt"
    "log"
    "os"

    "github.com/zsq2010/utils/notify"
)

func main() {
    fmt.Println("Notify Package Example")
    fmt.Println("======================")
    fmt.Println()

    demoEmail := os.Getenv("DEMO_EMAIL") == "true"
    demoBarker := os.Getenv("DEMO_BARKER") == "true"

    message := notify.Message{
        Title:    "测试通知",
        Body:     "这是一条来自 notify 包的测试消息。",
        Priority: "high",
    }

    if demoEmail {
        fmt.Println("1. Email Notification Example")
        fmt.Println("------------------------------")
        demoEmailNotification(message)
        fmt.Println()
    } else {
        fmt.Println("1. Email Notification Example (SKIPPED)")
        fmt.Println("   Set DEMO_EMAIL=true and configure environment variables to test:")
        fmt.Println("   - EMAIL_PROVIDER (qq/outlook/gmail)")
        fmt.Println("   - EMAIL_USERNAME")
        fmt.Println("   - EMAIL_PASSWORD")
        fmt.Println("   - EMAIL_FROM")
        fmt.Println("   - EMAIL_TO")
        fmt.Println()
    }

    if demoBarker {
        fmt.Println("2. Barker Push Notification Example")
        fmt.Println("------------------------------------")
        demoBarkerNotification(message)
        fmt.Println()
    } else {
        fmt.Println("2. Barker Push Notification Example (SKIPPED)")
        fmt.Println("   Set DEMO_BARKER=true and configure environment variables to test:")
        fmt.Println("   - BARKER_SERVER_URL")
        fmt.Println("   - BARKER_KEY")
        fmt.Println()
    }

    if demoEmail && demoBarker {
        fmt.Println("3. Multi-Channel Notification Example")
        fmt.Println("--------------------------------------")
        demoMultiChannelNotification(message)
        fmt.Println()
    } else {
        fmt.Println("3. Multi-Channel Notification Example (SKIPPED)")
        fmt.Println("   Enable both DEMO_EMAIL and DEMO_BARKER to test multi-channel")
        fmt.Println()
    }

    fmt.Println("Mock Examples (always available)")
    fmt.Println("--------------------------------")
    demoMockNotifications()
}

func demoEmailNotification(message notify.Message) {
    provider := getEmailProvider(os.Getenv("EMAIL_PROVIDER"))

    config := notify.EmailConfig{
        Provider: provider,
        Username: os.Getenv("EMAIL_USERNAME"),
        Password: os.Getenv("EMAIL_PASSWORD"),
        From:     os.Getenv("EMAIL_FROM"),
        To:       []string{os.Getenv("EMAIL_TO")},
    }

    emailNotifier := notify.NewEmail(config)

    fmt.Printf("Sending email via %s...\n", getProviderName(provider))
    if err := emailNotifier.Send(message); err != nil {
        log.Printf("Email send failed: %v\n", err)
        return
    }

    fmt.Println("✓ Email sent successfully")
}

func demoBarkerNotification(message notify.Message) {
    config := notify.BarkerConfig{
        ServerURL: os.Getenv("BARKER_SERVER_URL"),
        Key:       os.Getenv("BARKER_KEY"),
        Sound:     "default",
    }

    barkerNotifier := notify.NewBarker(config)

    fmt.Println("Sending Barker push notification...")
    if err := barkerNotifier.Send(message); err != nil {
        log.Printf("Barker send failed: %v\n", err)
        return
    }

    fmt.Println("✓ Barker notification sent successfully")
}

func demoMultiChannelNotification(message notify.Message) {
    provider := getEmailProvider(os.Getenv("EMAIL_PROVIDER"))

    emailConfig := notify.EmailConfig{
        Provider: provider,
        Username: os.Getenv("EMAIL_USERNAME"),
        Password: os.Getenv("EMAIL_PASSWORD"),
        From:     os.Getenv("EMAIL_FROM"),
        To:       []string{os.Getenv("EMAIL_TO")},
    }

    barkerConfig := notify.BarkerConfig{
        ServerURL: os.Getenv("BARKER_SERVER_URL"),
        Key:       os.Getenv("BARKER_KEY"),
        Sound:     "default",
    }

    emailNotifier := notify.NewEmail(emailConfig)
    barkerNotifier := notify.NewBarker(barkerConfig)

    multiNotifier := notify.NewMultiParallel(emailNotifier, barkerNotifier)

    fmt.Println("Sending to multiple channels in parallel...")
    if err := multiNotifier.Send(message); err != nil {
        log.Printf("Multi-channel send failed: %v\n", err)
        return
    }

    fmt.Println("✓ All channels notified successfully")
}

func demoMockNotifications() {
    type mockNotifier struct {
        name string
    }

    mock := mockNotifier{name: "Mock Notifier"}

    fmt.Printf("Creating %s (for testing purposes)...\n", mock.name)
    fmt.Println("✓ Mock notifiers can be used for unit testing")

    fmt.Println("\nConfiguration Examples:")
    fmt.Println("-----------------------")

    fmt.Println("\nQQ Mail Example:")
    fmt.Println(`  emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.QQMail,
    Username: "user@qq.com",
    Password: "授权码",
    From:     "sender@qq.com",
    To:       []string{"recipient@example.com"},
  })`)

    fmt.Println("\nOutlook Example:")
    fmt.Println(`  emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Outlook,
    Username: "user@outlook.com",
    Password: "password",
    From:     "sender@outlook.com",
    To:       []string{"recipient@example.com"},
  })`)

    fmt.Println("\nGmail Example:")
    fmt.Println(`  emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Gmail,
    Username: "user@gmail.com",
    Password: "app_password",
    From:     "sender@gmail.com",
    To:       []string{"recipient@example.com"},
  })`)

    fmt.Println("\nBarker Example:")
    fmt.Println(`  barkerNotifier := notify.NewBarker(notify.BarkerConfig{
    ServerURL: "https://api.day.app",
    Key:       "your_device_key",
    Sound:     "default",
  })`)

    fmt.Println("\nAdvanced Message Example:")
    fmt.Println(`  message := notify.Message{
    Title:    "重要通知",
    Body:     "这是消息正文",
    Priority: "high",
    HTMLBody: "<p>HTML格式的<strong>消息</strong></p>",
    Extra: map[string]interface{}{
      "sound": "alarm",
      "icon":  "https://example.com/icon.png",
      "url":   "https://example.com/details",
    },
  }`)
}

func getEmailProvider(provider string) notify.EmailProvider {
    switch provider {
    case "qq":
        return notify.QQMail
    case "outlook":
        return notify.Outlook
    case "gmail":
        return notify.Gmail
    default:
        return notify.Custom
    }
}

func getProviderName(provider notify.EmailProvider) string {
    switch provider {
    case notify.QQMail:
        return "QQ Mail"
    case notify.Outlook:
        return "Outlook"
    case notify.Gmail:
        return "Gmail"
    case notify.Custom:
        return "Custom SMTP"
    default:
        return "Unknown"
    }
}
