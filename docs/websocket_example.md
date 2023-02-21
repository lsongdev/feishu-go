# WebSocket 消息接收示例

## 功能说明

本实现提供了飞书 WebSocket 消息接收的完整功能，包括：

- ✅ 自动获取 WebSocket 连接端点
- ✅ WebSocket 连接建立
- ✅ Protobuf 消息解析
- ✅ 事件分发处理
- ✅ Ping/Pong 心跳保活
- ✅ 自动重连机制
- ✅ 优雅退出

## 使用方法

### 1. 配置环境变量

```bash
export FEISHU_APP_ID="your_app_id"
export FEISHU_APP_SECRET="your_app_secret"
```

### 2. 运行示例

```bash
go run examples/main.go
```

### 3. 代码示例

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "os"

    "github.com/lsongdev/feishu-go/feishu"
)

func main() {
    appID := os.Getenv("FEISHU_APP_ID")
    appSecret := os.Getenv("FEISHU_APP_SECRET")

    // 创建 WebSocket 客户端
    wsClient := feishu.NewWSClient(appID, appSecret,
        feishu.WithEventHandler(handleEvent),
        feishu.WithAutoReconnect(true),
    )

    // 启动 WebSocket 连接（阻塞）
    ctx := context.Background()
    if err := wsClient.Start(ctx); err != nil {
        log.Fatalf("WebSocket 启动失败：%v", err)
    }
}

// 事件处理器
func handleEvent(ctx context.Context, event *feishu.EventMessage) error {
    log.Printf("收到事件：%s", event.Header.EventType)

    switch event.Header.EventType {
    case feishu.EVENT_MESSAGE_RECEIVE:
        return handleMessageReceive(event)
    default:
        log.Printf("未知事件类型：%s", event.Header.EventType)
    }

    return nil
}

// 处理消息接收事件
func handleMessageReceive(event *feishu.EventMessage) error {
    var msgEvent feishu.MessageReceiveEvent
    if err := json.Unmarshal(event.Event, &msgEvent); err != nil {
        return err
    }

    log.Printf("收到消息：%s", msgEvent.Message.MessageID)

    // 解析消息内容
    var content map[string]interface{}
    if err := json.Unmarshal([]byte(msgEvent.Message.Content), &content); err != nil {
        return err
    }

    if text, ok := content["text"].(string); ok {
        log.Printf("消息内容：%s", text)
    }

    return nil
}
```

## 配置选项

```go
// 创建客户端时可配置选项
wsClient := feishu.NewWSClient(appID, appSecret,
    feishu.WithEventHandler(handler),      // 设置事件处理器
    feishu.WithAutoReconnect(true),        // 设置自动重连
)
```

## 支持的事件类型

```go
feishu.EVENT_MESSAGE_RECEIVE  // 消息接收
feishu.EVENT_APP_TICKET       // 应用票据
feishu.EVENT_BOT_ADD          // 机器人添加
feishu.EVENT_BOT_DELETED      // 机器人删除
feishu.EVENT_GROUP_ADDED      // 群组添加
feishu.EVENT_USER_ADDED       // 用户添加
// ... 更多事件类型
```

## 注意事项

1. **应用权限**：确保你的飞书应用已开通 WebSocket 接收消息的权限
2. **事件订阅**：在飞书开放平台后台配置需要订阅的事件
3. **网络环境**：确保服务器可以访问 `open.feishu.cn` 的 WebSocket 服务
4. **API 路径**：WebSocket 端点 API 路径为 `/callback/ws/endpoint`（不包含 `/open-apis` 前缀）

## 故障排查

### 获取端点失败

如果看到 "获取端点失败" 错误，检查：
- AppID 和 AppSecret 是否正确
- 应用是否已发布
- 是否开通了 WebSocket 权限

### 连接断开后重连

客户端会自动重连，重连配置：
- 默认无限次重连
- 重连间隔：2 分钟
- 首次重连随机抖动：0-30 秒
