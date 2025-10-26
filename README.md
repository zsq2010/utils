# utils

Go 语言实用工具库，提供常用的工具包。

## 包列表

### imagesplit - 图片分割工具

`imagesplit` 是一个用于 Go 语言的图片分割库，支持按网格（行×列）和固定尺寸两种方式对 PNG / JPEG 图片进行分割。

### notify - 多渠道通知库

`notify` 提供统一的通知发送接口，支持邮件（Email）、Barker 推送等多种通知渠道。

---

## imagesplit - 图片分割工具

### 功能特性

- ✅ 支持 PNG (`.png`) 和 JPEG (`.jpg`, `.jpeg`) 格式
- ✅ 网格分割：按照指定的行列数自动生成小图块
- ✅ 固定尺寸分割：按照固定的宽高切割，自动处理边缘剩余区域
- ✅ 灵活的输出配置：输出目录、文件前缀、图片格式、JPEG 质量
- ✅ 完善的错误处理：格式不支持、参数错误、输出目录创建失败等

### 安装

```bash
go get github.com/zsq2010/utils
```

### 快速上手

#### 单张图片分割

```go
package main

import (
    "fmt"
    "log"

    "github.com/zsq2010/utils/imagesplit"
)

func main() {
    files, err := imagesplit.GridSplit("input.png", 3, 4, imagesplit.SplitOptions{
        OutputDir:  "output",
        FilePrefix: "sample",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(files)
}
```

#### 批量处理目录

```go
cfg := imagesplit.DirectorySplitConfig{
    Mode: imagesplit.DirectorySplitModeGrid,
    Rows: 3,
    Cols: 3,
    Options: imagesplit.SplitOptions{Format: "png"},
}
results, err := imagesplit.SplitDirectory("input-dir", "output-dir", cfg)
if err != nil {
    log.Fatal(err)
}
for src, generated := range results {
    fmt.Println("source:", src)
    for _, path := range generated {
        fmt.Println("  -", path)
    }
}
```

### API 文档

```go
func GridSplit(inputPath string, rows, cols int, opts imagesplit.SplitOptions) ([]string, error)
```
- `rows` / `cols`: 网格行列数，必须大于 0。
- `SplitOptions`
  - `OutputDir`: 输出目录（为空时使用原图所在目录）。
  - `FilePrefix`: 输出文件前缀（为空时使用原图文件名）。
  - `Format`: 输出格式（`"png"`、`"jpeg"`，为空使用原图格式）。
  - `Quality`: JPEG 质量，范围 1-100（默认 90）。
- 返回值为生成的文件路径列表。

```go
func TileSplit(inputPath string, tileWidth, tileHeight int, opts imagesplit.SplitOptions) ([]string, error)
```
- `tileWidth` / `tileHeight`: 图块宽高，必须大于 0。
- 其余参数与 `GridSplit` 一致。

```go
func SplitDirectory(inputDir, outputDir string, cfg imagesplit.DirectorySplitConfig) (map[string][]string, error)
```
- `inputDir`: 输入图片所在目录。
- `outputDir`: 输出根目录，每张图片会在该目录下创建一个以图片名命名的子目录。
- `cfg.Mode`: 分割模式（`imagesplit.DirectorySplitModeGrid` / `imagesplit.DirectorySplitModeTile`）。
- `cfg.Rows`, `cfg.Cols`: 网格模式的行列数。
- `cfg.TileWidth`, `cfg.TileHeight`: 固定尺寸模式的宽高。
- `cfg.Options`: 其它分割选项（输出格式、JPEG 质量等），`OutputDir` 会被自动覆盖为图片专属子目录。

### 命名规则

- 网格分割：`{prefix}_row{i}_col{j}.{ext}` → 例如：`image_row0_col2.png`
- 固定尺寸：`{prefix}_tile_{index}.{ext}` → 例如：`image_tile_5.jpg`

### 示例

仓库包含一个完整的示例程序，位于 `imagesplit/example/main.go`：

```bash
cd imagesplit/example
go run .
```

示例会自动生成一张示例 PNG 图片，并在 `imagesplit/example/output` 目录下演示网格分割、固定尺寸分割以及目录批量处理，同时给出 JPEG 输出示例

### 测试图片

项目提供 `imagesplit/testdata` 辅助包，可动态生成内置的测试图片：

```go
import testdata "github.com/zsq2010/utils/imagesplit/testdata"

data, _ := testdata.GradientPNG()    // 获取示例 PNG 图片字节流
testdata.WriteBlocksJPEG("blocks.jpg") // 生成示例 JPEG 文件
```

这些工具可用于测试、示例或文档代码中，避免手动维护二进制测试资源。

### 测试

项目提供了覆盖主要功能和边界情况的单元测试，包括：

- 网格分割及无法整除时的边缘尺寸验证
- 固定尺寸分割及剩余像素处理
- PNG / JPEG 输入输出
- 错误处理与参数校验
- 输出目录自动创建

运行测试：

```bash
cd imagesplit
go test ./...
```

---

## notify - 多渠道通知库

### 功能特性

- ✅ 统一的 `Notifier` 接口设计
- ✅ 邮件通知：支持 QQ、Outlook、Gmail 等主流邮箱
- ✅ Barker 推送通知：支持通过 key 发送推送到指定设备
- ✅ 多渠道组合：支持同时向多个渠道发送通知
- ✅ 灵活配置：超时控制、重试机制、优先级设置
- ✅ 完善的错误处理和单元测试

### 安装

```bash
go get github.com/zsq2010/utils
```

### 快速上手

#### 邮件通知

```go
package main

import (
    "log"
    "github.com/zsq2010/utils/notify"
)

func main() {
    // QQ 邮箱示例
    emailNotifier := notify.NewEmail(notify.EmailConfig{
        Provider: notify.QQMail,
        Username: "user@qq.com",
        Password: "授权码",
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

#### Barker 推送通知

```go
barkerNotifier := notify.NewBarker(notify.BarkerConfig{
    ServerURL: "https://api.day.app",
    Key:       "your_device_key",
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

#### 多渠道组合

```go
// 创建多个通知渠道
emailNotifier := notify.NewEmail(emailConfig)
barkerNotifier := notify.NewBarker(barkerConfig)

// 组合多个渠道（并行发送）
multiNotifier := notify.NewMultiParallel(emailNotifier, barkerNotifier)

// 一次发送到所有渠道
message := notify.Message{
    Title: "系统警告",
    Body:  "服务器 CPU 使用率超过 90%",
    Priority: "urgent",
}

if err := multiNotifier.Send(message); err != nil {
    log.Fatal(err)
}
```

### 配置说明

#### 邮件配置

支持预置的邮件服务商：

| 服务商 | Provider | SMTP 服务器 | 端口 | 说明 |
|--------|----------|-------------|------|------|
| QQ邮箱 | `notify.QQMail` | smtp.qq.com | 587 | 需使用授权码 |
| Outlook | `notify.Outlook` | smtp-mail.outlook.com | 587 | - |
| Gmail | `notify.Gmail` | smtp.gmail.com | 587 | 需使用应用专用密码 |
| 自定义 | `notify.Custom` | 自定义 | 自定义 | 需手动配置 Host/Port |

**EmailConfig 字段说明：**

```go
type EmailConfig struct {
    Provider  EmailProvider  // 邮件服务商
    Host      string         // SMTP 服务器（Custom 时必填）
    Port      int            // SMTP 端口（Custom 时必填）
    UseTLS    bool           // 是否使用 STARTTLS
    UseSSL    bool           // 是否使用 SSL/TLS
    Username  string         // SMTP 用户名
    Password  string         // SMTP 密码或授权码
    From      string         // 发件人地址
    To        []string       // 收件人列表
    CC        []string       // 抄送列表（可选）
    BCC       []string       // 密送列表（可选）
}
```

#### Barker 配置

**BarkerConfig 字段说明：**

```go
type BarkerConfig struct {
    ServerURL string  // Barker 服务器地址
    Key       string  // 设备 Key
    Sound     string  // 通知声音（可选）
    Icon      string  // 通知图标 URL（可选）
    Group     string  // 通知分组（可选）
    URL       string  // 点击跳转 URL（可选）
}
```

#### 通用配置

所有通知渠道都支持的通用配置：

```go
type CommonConfig struct {
    Timeout       time.Duration  // 发送超时时间（默认 30s）
    RetryCount    int            // 重试次数（默认 0）
    RetryInterval time.Duration  // 重试间隔（默认 2s）
}
```

### 高级用法

#### HTML 邮件

```go
message := notify.Message{
    Title:    "欢迎使用",
    Body:     "纯文本内容",
    HTMLBody: "<h1>欢迎</h1><p>这是 <strong>HTML</strong> 格式的邮件</p>",
}
```

#### 附件（邮件）

```go
message := notify.Message{
    Title: "报表发送",
    Body:  "请查收附件中的月度报表",
    Attachments: []string{
        "/path/to/report.pdf",
        "/path/to/data.xlsx",
    },
}
```

#### 自定义参数（Barker）

```go
message := notify.Message{
    Title: "自定义推送",
    Body:  "消息内容",
    Extra: map[string]interface{}{
        "sound":    "alarm",
        "icon":     "https://example.com/icon.png",
        "url":      "https://example.com/details",
        "badge":    5,
        "autoCopy": "复制的内容",
    },
}
```

#### 重试机制

```go
emailNotifier := notify.NewEmail(notify.EmailConfig{
    Provider: notify.Gmail,
    // ... 其他配置
    CommonConfig: notify.CommonConfig{
        Timeout:       10 * time.Second,
        RetryCount:    3,
        RetryInterval: 5 * time.Second,
    },
})
```

### 示例

仓库包含一个完整的示例程序，位于 `notify/example/main.go`：

```bash
cd notify/example
go run .
```

示例展示了各种通知渠道的配置和使用方法。实际发送需要配置环境变量：

```bash
# 邮件测试
export DEMO_EMAIL=true
export EMAIL_PROVIDER=qq
export EMAIL_USERNAME=user@qq.com
export EMAIL_PASSWORD=your_auth_code
export EMAIL_FROM=sender@qq.com
export EMAIL_TO=recipient@example.com

# Barker 测试
export DEMO_BARKER=true
export BARKER_SERVER_URL=https://api.day.app
export BARKER_KEY=your_device_key

go run .
```

### 测试

项目提供了完善的单元测试，覆盖率超过 80%：

```bash
cd notify
go test ./... -v
```

测试包含：
- 配置验证和默认值测试
- 消息构建和格式化测试
- 错误处理测试
- Mock 服务器测试
- 重试机制测试
- 多渠道并发测试

---

## 许可证

该项目基于 MIT License 发布，欢迎自由使用与贡献。
