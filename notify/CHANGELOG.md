# Changelog - notify 包

## 更新说明

### 2024-10-26 - Bark 服务集成

#### 重要变更
- **重命名**: 将 `Barker` 相关的所有内容重命名为 `Bark`
  - `BarkerConfig` → `BarkConfig`
  - `BarkerNotifier` → `BarkNotifier`
  - `NewBarker()` → `NewBark()`
  - 文件名: `barker.go` → `bark.go`
  - 测试文件: `barker_test.go` → `bark_test.go`

#### 原因
Bark 是一个已有的 iOS 推送通知服务，官方 API 地址为 https://api.day.app/。
之前的命名 "Barker" 是误用，现已更正为正确的服务名称 "Bark"。

#### Bark 服务说明
- **官方网站**: https://github.com/Finb/Bark
- **API 地址**: https://api.day.app/
- **使用方式**:
  1. 在 iOS 设备上安装 Bark App
  2. 从 App 获取设备 Key
  3. 使用 API 发送推送通知到该设备

#### API 使用示例

```go
// 创建 Bark 通知器
barkNotifier := notify.NewBark(notify.BarkConfig{
    Key:   "your_device_key",  // 从 Bark App 获取
    Sound: "default",
    // ServerURL 默认为 https://api.day.app，可省略
})

// 发送通知
message := notify.Message{
    Title: "测试通知",
    Body:  "这是一条测试消息",
    Priority: "high",
}

err := barkNotifier.Send(message)
```

#### 支持的参数
- `Title`: 通知标题
- `Body`: 通知正文
- `Sound`: 通知声音 (可选)
- `Icon`: 通知图标 URL (可选)
- `Group`: 通知分组 (可选)
- `URL`: 点击跳转 URL (可选)
- `Level`: 优先级 (通过 Priority 字段自动映射)
  - `high`/`urgent` → `timeSensitive`
  - `low` → `passive`
  - 其他 → `active`
- `Badge`: 角标数字 (可选)
- `AutoCopy`: 自动复制内容 (可选)
- `Copy`: 点击复制内容 (可选)
- `IsArchive`: 是否自动保存 (可选)

#### 配置选项
```go
type BarkConfig struct {
    ServerURL    string        // Bark 服务器地址 (默认: https://api.day.app)
    Key          string        // 设备 Key (必需)
    Sound        string        // 通知声音 (可选)
    Icon         string        // 通知图标 URL (可选)
    Group        string        // 通知分组 (可选)
    URL          string        // 点击跳转 URL (可选)
    CommonConfig              // 通用配置 (超时、重试等)
}
```

#### 默认行为
- 如果不指定 `ServerURL`，自动使用 `https://api.day.app`
- 支持自建 Bark 服务器，只需设置自定义的 `ServerURL`

#### 测试覆盖
- 完整的单元测试覆盖
- 使用 `httptest` 模拟 Bark API 服务器
- 测试包含成功场景、错误处理、重试机制等

#### 文档更新
- 更新主 README.md
- 更新 notify/README.md
- 更新示例程序 notify/example/main.go
- 添加手动测试示例 notify/example/bark_test_example.go

#### 迁移指南
如果您之前使用了 `Barker` 相关的 API，请按以下方式更新代码：

```go
// 旧代码
barkerNotifier := notify.NewBarker(notify.BarkerConfig{
    ServerURL: "https://api.day.app",
    Key:       "your_key",
})

// 新代码
barkNotifier := notify.NewBark(notify.BarkConfig{
    Key: "your_key",  // ServerURL 默认为 https://api.day.app，可省略
})
```

环境变量也相应更新：
- `BARKER_SERVER_URL` → `BARK_SERVER_URL`
- `BARKER_KEY` → `BARK_KEY`
- `DEMO_BARKER` → `DEMO_BARK`
