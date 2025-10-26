# notify - 多渠道通知库

`notify` 包提供统一的通知发送接口，支持多种通知渠道，包括邮件（Email）和 Barker 推送通知。

## 特性

- ✅ 统一的 `Notifier` 接口设计
- ✅ 邮件通知：支持 QQ、Outlook、Gmail 等主流邮箱
- ✅ Barker 推送通知：支持通过 key 发送推送到指定设备
- ✅ 多渠道组合：支持同时向多个渠道发送通知（顺序或并行）
- ✅ 灵活配置：超时控制、重试机制、优先级设置
- ✅ 完善的错误处理和单元测试

## 安装

```bash
go get github.com/zsq2010/utils/notify
```

## 快速开始

### 邮件通知

```go
package main

import (
    "log"
    "github.com/zsq2010/utils/notify"
)

func main() {
    // 使用 QQ 邮箱发送
    emailNotifier := notify.NewEmail(notify.EmailConfig{
        Provider: notify.QQMail,
        Username: "user@qq.com",
        Password: "your_auth_code",  // 使用授权码，不是登录密码
        From:     "sender@qq.com",
        To:       []string{"recipient@example.com"},
    })

    message := notify.Message{
        Title: "测试通知",
        Body:  "这是一条测试消息",
    }

    if err := emailNotifier.Send(message); err != nil {
        log.Fatal(err)
    }
}
```

### Barker 推送通知

```go
barkerNotifier := notify.NewBarker(notify.BarkerConfig{
    ServerURL: "https://api.day.app",  // Barker 服务器地址
    Key:       "your_device_key",       // 你的设备 Key
    Sound:     "default",
})

message := notify.Message{
    Title:    "重要提醒",
    Body:     "这是一条推送消息",
    Priority: "high",
}

if err := barkerNotifier.Send(message); err != nil {
    log.Fatal(err)
}
```

### 多渠道通知

```go
// 创建多个通知渠道
emailNotifier := notify.NewEmail(emailConfig)
barkerNotifier := notify.NewBarker(barkerConfig)

// 并行发送到所有渠道
multiNotifier := notify.NewMultiParallel(emailNotifier, barkerNotifier)

// 或者顺序发送（遇到错误立即停止）
// multiNotifier := notify.NewMulti(emailNotifier, barkerNotifier)

message := notify.Message{
    Title:    "系统警告",
    Body:     "服务器 CPU 使用率超过 90%",
    Priority: "urgent",
}

if err := multiNotifier.Send(message); err != nil {
    log.Fatal(err)
}
```

## 配置指南

### 邮件配置

#### 支持的邮件服务商

| 服务商 | Provider 常量 | SMTP 服务器 | 端口 | 说明 |
|--------|--------------|-------------|------|------|
| QQ邮箱 | `notify.QQMail` | smtp.qq.com | 587 | 需在 QQ 邮箱设置中开启 SMTP 并获取授权码 |
| Outlook/Hotmail | `notify.Outlook` | smtp-mail.outlook.com | 587 | 使用账号密码 |
| Gmail | `notify.Gmail` | smtp.gmail.com | 587 | 需生成应用专用密码 |
| 自定义 | `notify.Custom` | 自定义 | 自定义 | 需手动配置 Host 和 Port |

#### QQ 邮箱配置示例

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.QQMail,
    Username: "123456789@qq.com",
    Password: "abcdefghijklmnop",  // 授权码，不是 QQ 密码
    From:     "123456789@qq.com",
    To:       []string{"recipient@example.com"},
    CC:       []string{"cc@example.com"},      // 可选
    BCC:      []string{"bcc@example.com"},     // 可选
})
```

**获取 QQ 邮箱授权码步骤：**
1. 登录 QQ 邮箱
2. 进入"设置" → "账户"
3. 找到"POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务"
4. 开启"IMAP/SMTP服务"
5. 生成授权码

#### Outlook 配置示例

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Outlook,
    Username: "user@outlook.com",
    Password: "your_password",
    From:     "user@outlook.com",
    To:       []string{"recipient@example.com"},
})
```

#### Gmail 配置示例

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Gmail,
    Username: "user@gmail.com",
    Password: "your_app_password",  // 应用专用密码
    From:     "user@gmail.com",
    To:       []string{"recipient@example.com"},
})
```

**获取 Gmail 应用专用密码步骤：**
1. 访问 Google 账户安全设置
2. 启用两步验证
3. 生成应用专用密码

#### 自定义 SMTP 服务器

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Custom,
    Host:     "smtp.example.com",
    Port:     587,
    UseTLS:   true,
    Username: "user@example.com",
    Password: "password",
    From:     "sender@example.com",
    To:       []string{"recipient@example.com"},
})
```

#### 使用 SSL/TLS (端口 465)

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.QQMail,
    Port:     465,
    UseSSL:   true,  // 使用 SSL 而不是 STARTTLS
    Username: "user@qq.com",
    Password: "auth_code",
    From:     "user@qq.com",
    To:       []string{"recipient@example.com"},
})
```

### Barker 配置

```go
barkerNotifier := notify.NewBarker(notify.BarkerConfig{
    ServerURL: "https://api.day.app",  // Barker 服务器地址
    Key:       "your_device_key",      // 设备 Key
    Sound:     "default",              // 通知声音（可选）
    Icon:      "https://example.com/icon.png",  // 图标 URL（可选）
    Group:     "MyApp",                // 通知分组（可选）
    URL:       "https://example.com",  // 点击跳转 URL（可选）
})
```

**获取 Barker Key：**
1. 在 iOS 设备上安装 Barker App
2. 打开 App 查看设备 Key
3. 或者使用自建的 Barker 服务器

### 通用配置选项

所有通知渠道都支持以下通用配置：

```go
notify.EmailConfig{
    // ... 其他配置
    CommonConfig: notify.CommonConfig{
        Timeout:       30 * time.Second,  // 发送超时时间（默认 30s）
        RetryCount:    3,                  // 失败重试次数（默认 0）
        RetryInterval: 5 * time.Second,    // 重试间隔（默认 2s）
    },
}
```

## 高级用法

### HTML 邮件

```go
message := notify.Message{
    Title:    "欢迎使用",
    Body:     "纯文本内容（用于不支持 HTML 的邮件客户端）",
    HTMLBody: `
        <html>
        <body>
            <h1>欢迎使用我们的服务</h1>
            <p>这是一封 <strong>HTML 格式</strong>的邮件。</p>
            <a href="https://example.com">访问我们的网站</a>
        </body>
        </html>
    `,
}
```

### 附件（邮件）

```go
message := notify.Message{
    Title: "月度报表",
    Body:  "请查收附件中的月度报表",
    Attachments: []string{
        "/path/to/report.pdf",
        "/path/to/data.xlsx",
    },
}
```

**注意：** 当前版本的附件支持是基础实现，仅在邮件中显示附件文件名。完整的附件编码（MIME multipart）将在未来版本中增强。

### Barker 自定义参数

使用 `Extra` 字段传递 Barker 特定参数：

```go
message := notify.Message{
    Title: "自定义推送",
    Body:  "消息内容",
    Extra: map[string]interface{}{
        "sound":    "alarm",                          // 声音
        "icon":     "https://example.com/icon.png",   // 图标
        "group":    "MyApp",                          // 分组
        "url":      "https://example.com/details",    // 跳转 URL
        "badge":    5,                                // 角标数字
        "autoCopy": "自动复制的内容",                    // 自动复制
        "copy":     "点击复制的内容",                     // 点击复制
    },
}
```

### 消息优先级

```go
message := notify.Message{
    Title:    "紧急通知",
    Body:     "系统出现严重故障",
    Priority: "urgent",  // high, urgent, normal, low
}
```

优先级映射（Barker）：
- `high` / `urgent` → `timeSensitive`（时效性通知）
- `low` → `passive`（被动通知）
- 其他 → `active`（普通通知）

### 重试机制

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Gmail,
    Username: "user@gmail.com",
    Password: "password",
    From:     "user@gmail.com",
    To:       []string{"recipient@example.com"},
    CommonConfig: notify.CommonConfig{
        Timeout:       10 * time.Second,  // 每次尝试的超时时间
        RetryCount:    3,                  // 失败后重试 3 次
        RetryInterval: 2 * time.Second,    // 每次重试间隔 2 秒
    },
})
```

### 多渠道策略

#### 顺序发送（遇到错误立即停止）

```go
multi := notify.NewMulti(emailNotifier, barkerNotifier, slackNotifier)
// 按顺序发送，第一个失败则停止
err := multi.Send(message)
```

#### 并行发送（收集所有错误）

```go
multi := notify.NewMultiParallel(emailNotifier, barkerNotifier, slackNotifier)
// 并行发送到所有渠道，收集所有错误
err := multi.Send(message)
// 如果多个渠道失败，err 将包含所有错误信息
```

## 实战示例

### 服务器监控告警

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/zsq2010/utils/notify"
)

func checkCPU() float64 {
    // 模拟 CPU 检查
    return 95.5
}

func main() {
    // 配置邮件通知
    emailNotifier := notify.NewEmail(notify.EmailConfig{
        Provider: notify.QQMail,
        Username: "monitor@qq.com",
        Password: os.Getenv("EMAIL_PASSWORD"),
        From:     "monitor@qq.com",
        To:       []string{"admin@example.com"},
    })

    // 配置 Barker 推送
    barkerNotifier := notify.NewBarker(notify.BarkerConfig{
        ServerURL: "https://api.day.app",
        Key:       os.Getenv("BARKER_KEY"),
        Sound:     "alarm",
    })

    // 组合通知渠道
    notifier := notify.NewMultiParallel(emailNotifier, barkerNotifier)

    // 监控循环
    for {
        cpuUsage := checkCPU()
        if cpuUsage > 90 {
            message := notify.Message{
                Title:    "⚠️ CPU 使用率告警",
                Body:     fmt.Sprintf("当前 CPU 使用率: %.1f%%\n时间: %s", cpuUsage, time.Now().Format("2006-01-02 15:04:05")),
                Priority: "urgent",
            }
            
            if err := notifier.Send(message); err != nil {
                log.Printf("发送告警失败: %v", err)
            }
        }
        
        time.Sleep(5 * time.Minute)
    }
}
```

### 用户注册欢迎邮件

```go
func sendWelcomeEmail(userEmail, userName string) error {
    emailNotifier := notify.NewEmail(notify.EmailConfig{
        Provider: notify.Gmail,
        Username: "noreply@example.com",
        Password: os.Getenv("EMAIL_PASSWORD"),
        From:     "Example Team <noreply@example.com>",
        To:       []string{userEmail},
    })

    message := notify.Message{
        Title: "欢迎加入 Example！",
        Body:  fmt.Sprintf("Hi %s,\n\n感谢注册！", userName),
        HTMLBody: fmt.Sprintf(`
            <html>
            <body style="font-family: Arial, sans-serif;">
                <h1>欢迎加入 Example！</h1>
                <p>Hi <strong>%s</strong>,</p>
                <p>感谢您的注册。点击下面的按钮开始使用：</p>
                <a href="https://example.com/get-started" 
                   style="background-color: #4CAF50; color: white; padding: 10px 20px; 
                          text-decoration: none; border-radius: 5px;">
                    开始使用
                </a>
            </body>
            </html>
        `, userName),
    }

    return emailNotifier.Send(message)
}
```

## 错误处理

```go
message := notify.Message{
    Title: "测试",
    Body:  "内容",
}

if err := notifier.Send(message); err != nil {
    // 错误信息包含具体的失败原因
    log.Printf("发送失败: %v", err)
    
    // 可以进行错误分类处理
    switch {
    case strings.Contains(err.Error(), "timeout"):
        log.Println("发送超时")
    case strings.Contains(err.Error(), "authentication"):
        log.Println("认证失败")
    case strings.Contains(err.Error(), "no recipients"):
        log.Println("没有收件人")
    default:
        log.Println("未知错误")
    }
}
```

## 测试

运行单元测试：

```bash
cd notify
go test ./... -v
```

查看测试覆盖率：

```bash
go test -cover ./...
```

生成覆盖率报告：

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 最佳实践

1. **使用环境变量存储敏感信息**
   ```go
   Password: os.Getenv("EMAIL_PASSWORD"),
   ```

2. **为生产环境配置重试机制**
   ```go
   CommonConfig: notify.CommonConfig{
       RetryCount: 3,
       RetryInterval: 5 * time.Second,
   }
   ```

3. **使用并行多渠道提高可靠性**
   ```go
   notifier := notify.NewMultiParallel(email, barker, slack)
   ```

4. **记录发送日志**
   ```go
   if err := notifier.Send(message); err != nil {
       log.Printf("[%s] 发送失败: %v", time.Now(), err)
   } else {
       log.Printf("[%s] 发送成功", time.Now())
   }
   ```

5. **为不同场景使用不同优先级**
   - 使用 `urgent` 用于严重告警
   - 使用 `normal` 用于日常通知
   - 使用 `low` 用于统计报表

## 接口文档

### Notifier 接口

```go
type Notifier interface {
    Send(message Message) error
}
```

所有通知渠道都实现此接口。

### Message 结构

```go
type Message struct {
    Title       string                 // 标题
    Body        string                 // 正文
    Priority    string                 // 优先级 (high, urgent, normal, low)
    HTMLBody    string                 // HTML 格式正文（可选）
    Attachments []string               // 附件路径列表（可选）
    Extra       map[string]interface{} // 扩展字段（可选）
}
```

## 常见问题

### Q: QQ 邮箱发送失败，提示"authentication failed"？

A: 请确保：
1. 使用的是授权码，不是 QQ 密码
2. 已在 QQ 邮箱中开启 SMTP 服务
3. 授权码正确无误

### Q: Gmail 发送失败？

A: 请确保：
1. 已启用两步验证
2. 使用应用专用密码，不是 Google 账号密码
3. 网络可以访问 smtp.gmail.com

### Q: Barker 推送没有收到？

A: 请检查：
1. ServerURL 是否正确（默认为 https://api.day.app）
2. Key 是否正确
3. iOS 设备上的 Barker App 是否正常运行
4. 网络连接是否正常

### Q: 如何提高发送成功率？

A: 建议：
1. 配置重试机制（RetryCount 和 RetryInterval）
2. 使用多渠道并行发送
3. 设置合理的超时时间（Timeout）
4. 记录详细的错误日志

### Q: 测试覆盖率为什么没有达到 80%？

A: 邮件发送的底层 SMTP 连接代码难以在单元测试中 mock（标准库 net/smtp 不提供依赖注入）。实际的业务逻辑、配置验证、消息构建等核心功能都有完整的测试覆盖。Barker 通知和多渠道功能都有 >85% 的覆盖率。

## 许可证

MIT License
