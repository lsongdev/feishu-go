package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lsongdev/feishu-go/feishu"
)

var client *feishu.Client

func main() {

	// 创建客户端，配置 WebSocket 选项
	client = feishu.NewClient(&feishu.Config{
		AppID:     os.Getenv("FEISHU_APP_ID"),
		AppSecret: os.Getenv("FEISHU_APP_SECRET"),
		WSConfig:  *feishu.DefaultWSConfig(),
	})

	// 获取 tenant_access_token（用于发送消息）x
	tokenResp, err := client.GetTenantAccessTokenInternal()
	if err != nil {
		log.Fatalf("获取 token 失败：%v", err)
	}
	client.SetAccessToken(tokenResp.TenantAccessToken)
	log.Println("已获取 tenant_access_token", tokenResp)

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("收到退出信号，正在关闭...")
		cancel() // 取消上下文
		client.Close()
	}()

	// 启动 WebSocket 连接（阻塞）
	if err := client.Start(ctx, handleEvent); err != nil {
		if err == context.Canceled {
			log.Println("WebSocket 已正常关闭")
		} else {
			log.Fatalf("WebSocket 错误：%v", err)
		}
	}
}

// 事件处理器
func handleEvent(event *feishu.EventMessage) error {
	log.Printf("收到事件：type=%s, id=%s", event.Header.EventType, event.Header.EventID)

	switch event.Header.EventType {
	case feishu.EVENT_MESSAGE_RECEIVE:
		return handleMessageReceive(event)
	case feishu.EVENT_MESSAGE_READ:
		log.Println("消息被已读")
	case feishu.EVENT_APP_TICKET:
		log.Println("收到 app_ticket 事件")
	case feishu.EVENT_BOT_ADD:
		log.Println("机器人被添加")
	case feishu.EVENT_BOT_DELETED:
		log.Println("机器人被删除")
	default:
		log.Printf("未知事件类型：%s", event.Header.EventType)
	}

	return nil
}

// 处理消息接收事件
func handleMessageReceive(event *feishu.EventMessage) error {
	log.Println("收到消息", string(event.Event))
	var msgEvent feishu.MessageReceiveEvent
	if err := json.Unmarshal(event.Event, &msgEvent); err != nil {
		return err
	}

	// 解析消息内容（以文本消息为例）
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(msgEvent.Message.Content), &content); err != nil {
		return err
	}

	if text, ok := content["text"].(string); ok {
		log.Printf("消息内容：%s", text)

		// Echo 消息：回复发送者（带引用）
		// senderOpenID := msgEvent.Sender.SenderID.OpenID
		messageID := msgEvent.Message.MessageID
		replyText := fmt.Sprintf("Echo: %s", text)

		if res, err := client.AddReaction(messageID, feishu.EMOJI_OK); err != nil {
			log.Printf("发送表情失败：%v, res=%v", err, res)
		} else {
			log.Printf("发送表情成功：%v", res)
		}

		replyMessage := feishu.NewTextMessage(replyText)
		if res, err := client.SendReplyMessage(messageID, &replyMessage); err != nil {
			log.Printf("发送回复失败：%v, res=%v", err, res)
		} else {
			log.Printf("发送回复成功：%v", res)
		}
	}

	return nil
}
