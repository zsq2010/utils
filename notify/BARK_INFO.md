# Bark 推送通知服务集成说明

## Bark 是什么？

Bark 是一个开源的 iOS 推送通知服务，允许开发者通过简单的 API 调用向 iOS 设备发送推送通知。

- **官方仓库**: https://github.com/Finb/Bark
- **官方 API**: https://api.day.app/
- **App Store**: 搜索 "Bark" 下载

## 特点

- 🚀 简单易用：只需一个 Key 即可发送通知
- 🔒 隐私保护：通知内容端到端加密
- 🎨 高度自定义：支持自定义图标、声音、分组等
- 🏠 支持自建：可以部署自己的 Bark 服务器
- 📱 iOS 原生：完美支持 iOS 通知特性

## 快速开始

### 1. 安装 Bark App

在 iOS 设备上从 App Store 下载并安装 Bark。

### 2. 获取设备 Key

打开 Bark App，会显示一个类似这样的 URL：
```
https://api.day.app/YourDeviceKey
```

其中 `YourDeviceKey` 就是你的设备 Key。

### 3. 使用 notify 包发送通知

```go
package main

import (
    "log"
    "github.com/zsq2010/utils/notify"
)

func main() {
    // 创建 Bark 通知器
    barkNotifier := notify.NewBark(notify.BarkConfig{
        Key: "YourDeviceKey",  // 替换为你的实际 Key
    })

    // 发送通知
    err := barkNotifier.Send(notify.Message{
        Title: "Hello from notify!",
        Body:  "这是一条测试通知",
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
```

## API 参数详解

### 基本参数

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| Key | string | ✅ | 设备 Key，从 Bark App 获取 |
| Title | string | ❌ | 通知标题 |
| Body | string | ✅ | 通知内容 |

### 扩展参数

通过 `Message.Extra` 字段可以设置以下参数：

| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| sound | string | 通知声音 | "alarm", "bell", "default" |
| icon | string | 通知图标 URL | "https://example.com/icon.png" |
| group | string | 通知分组 | "MyApp" |
| url | string | 点击通知时打开的 URL | "https://example.com" |
| badge | int | App 角标数字 | 5 |
| autoCopy | string | 自动复制到剪贴板的内容 | "copy this" |
| copy | string | 点击通知时复制的内容 | "copy this" |
| isArchive | int | 是否自动保存 (1=是, 0=否) | 1 |

### 优先级 (Priority)

通过 `Message.Priority` 字段设置：

| 值 | Bark Level | 说明 |
|----|------------|------|
| "high" / "urgent" | timeSensitive | 时效性通知，即使在勿扰模式也会展示 |
| "low" | passive | 被动通知，不会立即展示 |
| "normal" 或其他 | active | 普通通知 |

## 使用示例

### 基本通知

```go
barkNotifier.Send(notify.Message{
    Title: "提醒",
    Body:  "该吃饭了！",
})
```

### 带声音的通知

```go
barkNotifier.Send(notify.Message{
    Title: "重要提醒",
    Body:  "会议即将开始",
    Extra: map[string]interface{}{
        "sound": "alarm",
    },
})
```

### 带跳转链接的通知

```go
barkNotifier.Send(notify.Message{
    Title: "新消息",
    Body:  "点击查看详情",
    Extra: map[string]interface{}{
        "url": "https://example.com/message/123",
    },
})
```

### 高优先级通知

```go
barkNotifier.Send(notify.Message{
    Title:    "紧急警告",
    Body:     "服务器宕机！",
    Priority: "urgent",
    Extra: map[string]interface{}{
        "sound": "alarm",
    },
})
```

### 完整配置示例

```go
barkNotifier.Send(notify.Message{
    Title:    "系统通知",
    Body:     "这是一个完整配置的通知",
    Priority: "high",
    Extra: map[string]interface{}{
        "sound":     "bell",
        "icon":      "https://example.com/icon.png",
        "group":     "System",
        "url":       "https://example.com/details",
        "badge":     3,
        "isArchive": 1,
    },
})
```

## 自建 Bark 服务器

如果你想使用自己的 Bark 服务器：

```go
barkNotifier := notify.NewBark(notify.BarkConfig{
    ServerURL: "https://your-bark-server.com",
    Key:       "YourDeviceKey",
})
```

部署 Bark 服务器请参考官方文档：
https://github.com/Finb/Bark/blob/master/README.md

## 声音列表

Bark 支持的声音（部分）：
- `alarm` - 闹钟
- `anticipate` - 预期
- `bell` - 铃声
- `birdsong` - 鸟鸣
- `bloom` - 绽放
- `calypso` - 卡利普索
- `chime` - 钟声
- `choo` - 火车
- `descent` - 下降
- `electronic` - 电子
- `fanfare` - 号角
- `glass` - 玻璃
- `gotosleep` - 去睡觉
- `healthnotification` - 健康通知
- `horn` - 喇叭
- `ladder` - 梯子
- `mailsent` - 邮件发送
- `minuet` - 小步舞曲
- `multiwayinvitation` - 多方邀请
- `newmail` - 新邮件
- `newsflash` - 新闻快讯
- `noir` - 黑色
- `paymentsuccess` - 支付成功
- `shake` - 摇动
- `sherwoodforest` - 舍伍德森林
- `silence` - 静音
- `spell` - 咒语
- `suspense` - 悬念
- `telegraph` - 电报
- `tiptoes` - 蹑手蹑脚
- `typewriters` - 打字机
- `update` - 更新

## 故障排查

### 通知没有收到？

1. 检查设备 Key 是否正确
2. 确认 iOS 设备上的 Bark App 正在运行
3. 检查网络连接
4. 查看 Bark App 中的历史记录

### API 返回错误？

常见错误码：
- `400` - 参数错误（检查 Key 和内容）
- `404` - Key 不存在
- `500` - 服务器错误

### 自定义服务器无法连接？

1. 确认服务器 URL 正确
2. 检查服务器是否正常运行
3. 确认防火墙规则允许访问

## 最佳实践

1. **保护你的 Key**: Key 相当于密码，不要泄露给他人
2. **合理使用声音**: 不要在夜间使用吵闹的声音
3. **分组管理**: 使用 `group` 参数对不同类型的通知分组
4. **设置重试**: 对于重要通知，配置重试机制
5. **监控配额**: 注意官方 API 的使用限制

## 更多信息

- Bark GitHub: https://github.com/Finb/Bark
- Bark API 文档: https://github.com/Finb/Bark/blob/master/API.md
- notify 包文档: [README.md](README.md)

## 社区与支持

如有问题或建议，欢迎提交 Issue 或 Pull Request。
