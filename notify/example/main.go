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
    demoBark := os.Getenv("DEMO_BARK") == "true"

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

    if demoBark {
        fmt.Println("2. Bark Push Notification Example")
        fmt.Println("----------------------------------")
        demoBarkNotification(message)
        fmt.Println()
    } else {
        fmt.Println("2. Bark Push Notification Example (SKIPPED)")
        fmt.Println("   Set DEMO_BARK=true and configure environment variables to test:")
        fmt.Println("   - BARK_KEY (your device key)")
        fmt.Println("   - BARK_SERVER_URL (optional, defaults to https://api.day.app)")
        fmt.Println()
    }

    if demoEmail && demoBark {
        fmt.Println("3. Multi-Channel Notification Example")
        fmt.Println("--------------------------------------")
        demoMultiChannelNotification(message)
        fmt.Println()
    } else {
        fmt.Println("3. Multi-Channel Notification Example (SKIPPED)")
        fmt.Println("   Enable both DEMO_EMAIL and DEMO_BARK to test multi-channel")
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

func demoBarkNotification(message notify.Message) {
    serverURL := os.Getenv("BARK_SERVER_URL")
    if serverURL == "" {
        serverURL = "https://api.day.app"
    }

    config := notify.BarkConfig{
        ServerURL: serverURL,
        Key:       os.Getenv("BARK_KEY"),
        Sound:     "default",
    }

    barkNotifier := notify.NewBark(config)

    fmt.Println("Sending Bark push notification...")
    if err := barkNotifier.Send(message); err != nil {
        log.Printf("Bark send failed: %v\n", err)
        return
    }

    fmt.Println("✓ Bark notification sent successfully")
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

    barkServerURL := os.Getenv("BARK_SERVER_URL")
    if barkServerURL == "" {
        barkServerURL = "https://api.day.app"
    }

    barkConfig := notify.BarkConfig{
        ServerURL: barkServerURL,
        Key:       os.Getenv("BARK_KEY"),
        Sound:     "default",
    }

    emailNotifier := notify.NewEmail(emailConfig)
    barkNotifier := notify.NewBark(barkConfig)

    multiNotifier := notify.NewMultiParallel(emailNotifier, barkNotifier)

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

    fmt.Println("\nBark Example:")
    fmt.Println(`  barkNotifier := notify.NewBark(notify.BarkConfig{
    Key:   "your_device_key",  // 从 Bark App 获取
    Sound: "default",
    // ServerURL 默认为 https://api.day.app，可省略
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
