# WebSocket 使用指南

## 快速开始

### 1. 创建客户端

```go
client := feishu.NewClient(&feishu.Config{
    AppID:     "your_app_id",
    AppSecret: "your_app_secret",
})
```

### 2. 启动 WebSocket 连接

```go
// 简单方式 - 阻塞接收事件
if err := client.Start(handleEvent); err != nil {
    log.Fatalf("WebSocket 启动失败：%v", err)
}
```

### 3. 事件处理函数

```go
func handleEvent(event *feishu.EventMessage) error {
    log.Printf("收到事件：%s", event.Header.EventType)
    
    switch event.Header.EventType {
    case feishu.EVENT_MESSAGE_RECEIVE:
        return handleMessageReceive(event)
    case feishu.EVENT_APP_TICKET:
        log.Println("收到 app_ticket 事件")
    default:
        log.Printf("未知事件类型：%s", event.Header.EventType)
    }
    
    return nil
}
```

### 4. 处理消息事件

```go
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

### 5. 优雅退出

```go
// 监听退出信号
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigChan
    log.Println("正在关闭...")
    client.Close()
}()

// 启动 WebSocket
client.Start(handleEvent)
```

### 6. 使用 Context 控制

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 启动 WebSocket（可通过 ctx 控制退出）
if err := client.StartWithContext(ctx, handleEvent); err != nil {
    log.Fatalf("WebSocket 启动失败：%v", err)
}
```

## 完整示例

```go
package main

import (
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/lsongdev/feishu-go/feishu"
)

func main() {
    client := feishu.NewClient(&feishu.Config{
        AppID:     os.Getenv("FEISHU_APP_ID"),
        AppSecret: os.Getenv("FEISHU_APP_SECRET"),
    })

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Println("正在关闭...")
        client.Close()
    }()

    log.Println("正在连接飞书 WebSocket...")
    if err := client.Start(handleEvent); err != nil {
        log.Fatalf("WebSocket 启动失败：%v", err)
    }
}

func handleEvent(event *feishu.EventMessage) error {
    log.Printf("收到事件：%s", event.Header.EventType)
    
    if event.Header.EventType == feishu.EVENT_MESSAGE_RECEIVE {
        var msgEvent feishu.MessageReceiveEvent
        if err := json.Unmarshal(event.Event, &msgEvent); err != nil {
            return err
        }
        
        var content map[string]interface{}
        json.Unmarshal([]byte(msgEvent.Message.Content), &content)
        
        if text, ok := content["text"].(string); ok {
            log.Printf("消息内容：%s", text)
        }
    }
    
    return nil
}
```

## API 说明

### `client.Start(fn EventHandler) error`

启动 WebSocket 连接并接收事件（阻塞方法）

- `fn`: 事件处理函数，每次收到事件时会被调用
- 自动重连（默认开启）
- 自动心跳保活

### `client.StartWithContext(ctx context.Context, fn EventHandler) error`

启动 WebSocket 连接并接收事件（可通过 ctx 控制退出）

- `ctx`: 上下文，用于控制生命周期
- `fn`: 事件处理函数

### `client.Close() error`

关闭 WebSocket 连接

## 事件类型

```go
feishu.EVENT_MESSAGE_RECEIVE  // 消息接收
feishu.EVENT_APP_TICKET       // 应用票据
feishu.EVENT_BOT_ADD          // 机器人添加
feishu.EVENT_BOT_DELETED      // 机器人删除
feishu.EVENT_GROUP_ADDED      // 群组添加
feishu.EVENT_GROUP_DELETED    // 群组删除
feishu.EVENT_GROUP_UPDATED    // 群组更新
feishu.EVENT_USER_ADDED       // 用户添加
feishu.EVENT_USER_DELETED     // 用户删除
feishu.EVENT_USER_UPDATED     // 用户更新
```

## 配置说明

默认配置（通过 `DefaultWSConfig()` 查看）：

- `AutoReconnect`: true（自动重连）
- `ReconnectCount`: -1（无限重连）
- `ReconnectInterval`: 2 分钟
- `ReconnectNonce`: 30 秒（首次重连随机抖动）
- `PingInterval`: 2 分钟

## 注意事项

1. **应用权限**：确保应用已开通 WebSocket 接收消息权限
2. **事件订阅**：在飞书开放平台配置需要订阅的事件
3. **错误处理**：事件处理函数中的错误会被记录但不会中断连接
