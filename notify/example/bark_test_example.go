// +build manual

// 这是一个手动测试示例，用于验证 Bark 通知功能
// 使用方式：
// 1. 在 iOS 设备上安装 Bark App
// 2. 从 App 获取你的设备 Key
// 3. 设置环境变量 BARK_KEY
// 4. 运行: go run -tags manual bark_test_example.go
//
// 或者直接运行:
// BARK_KEY=your_key go run -tags manual bark_test_example.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zsq2010/utils/notify"
)

func main() {
	key := os.Getenv("BARK_KEY")
	if key == "" {
		log.Fatal("请设置 BARK_KEY 环境变量")
	}

	// 创建 Bark 通知器
	// ServerURL 默认为 https://api.day.app，无需显式设置
	barkNotifier := notify.NewBark(notify.BarkConfig{
		Key:   key,
		Sound: "bell",
	})

	// 发送简单通知
	fmt.Println("发送简单通知...")
	err := barkNotifier.Send(notify.Message{
		Title: "测试通知",
		Body:  "这是来自 notify 包的 Bark 测试消息",
	})
	if err != nil {
		log.Fatalf("发送失败: %v", err)
	}
	fmt.Println("✓ 通知发送成功！请查看您的 iOS 设备")

	// 发送带优先级的通知
	fmt.Println("\n发送高优先级通知...")
	err = barkNotifier.Send(notify.Message{
		Title:    "紧急通知",
		Body:     "这是一条高优先级消息",
		Priority: "urgent",
	})
	if err != nil {
		log.Fatalf("发送失败: %v", err)
	}
	fmt.Println("✓ 高优先级通知发送成功！")

	// 发送带自定义参数的通知
	fmt.Println("\n发送带自定义参数的通知...")
	err = barkNotifier.Send(notify.Message{
		Title: "自定义通知",
		Body:  "点击可以跳转到指定网址",
		Extra: map[string]interface{}{
			"url":   "https://github.com/zsq2010/utils",
			"sound": "alarm",
			"icon":  "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png",
			"group": "notify-test",
		},
	})
	if err != nil {
		log.Fatalf("发送失败: %v", err)
	}
	fmt.Println("✓ 自定义通知发送成功！")

	fmt.Println("\n所有测试完成！")
}
