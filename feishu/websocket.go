package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// EventType 定义事件类型
type EventType string

const (
	// 消息相关事件
	EVENT_MESSAGE_RECEIVE EventType = "im.message.receive_v1"
	EVENT_MESSAGE_READ    EventType = "im.message.message_read_v1"

	// 应用相关事件
	EVENT_APP_TICKET      EventType = "app_ticket"
	EVENT_APP_STATUS      EventType = "application_status_v2"
	EVENT_APP_OPEN        EventType = "application_open_v2"
	EVENT_APP_CANCEL      EventType = "application_cancel_v2"
	EVENT_APP_SUSPEND     EventType = "application_suspend_v2"
	EVENT_APP_RESTORE     EventType = "application_restore_v2"
	EVENT_APP_REDIRECT    EventType = "application_redirect_v2"
	EVENT_APP_ACCOUNT     EventType = "application_account_v2"
	EVENT_APP_DEMO        EventType = "application_demo_v2"
	EVENT_APP_CHARGE      EventType = "application_charge_v2"
	EVENT_APP_REFUND      EventType = "application_refund_v2"
	EVENT_APP_TRIAL       EventType = "application_trial_v2"
	EVENT_APP_SALES       EventType = "application_sales_v2"
	EVENT_APP_SUBSCRIBE   EventType = "application_subscribe_v2"
	EVENT_APP_UNSUBSCRIBE EventType = "application_unsubscribe_v2"

	// 机器人相关事件
	EVENT_BOT_ADD     EventType = "bot.add"
	EVENT_BOT_DELETED EventType = "bot.deleted"

	// 群组相关事件
	EVENT_GROUP_ADDED   EventType = "group.added"
	EVENT_GROUP_DELETED EventType = "group.deleted"
	EVENT_GROUP_UPDATED EventType = "group.updated"

	// 成员相关事件
	EVENT_USER_ADDED   EventType = "user.added"
	EVENT_USER_DELETED EventType = "user.deleted"
	EVENT_USER_UPDATED EventType = "user.updated"
)

// FrameType 定义帧类型
type FrameType int32

const (
	FrameTypeControl FrameType = 0
	FrameTypeData    FrameType = 1
)

// MessageType 定义消息类型
type MessageType string

const (
	MessageTypeEvent MessageType = "event"
	MessageTypePing  MessageType = "ping"
	MessageTypePong  MessageType = "pong"
)

// FramePB 是 protobuf 生成的 Frame 类型的别名
type FramePB = Frame

// EndpointResponse WebSocket 端点响应
type EndpointResponse struct {
	ResponseBase
	Data struct {
		URL          string    `json:"url"`
		ClientConfig *WSConfig `json:"client_config,omitempty"`
	} `json:"data"`
}

// EventMessage 事件消息结构
type EventMessage struct {
	Schema string          `json:"schema"`
	Header EventHeader     `json:"header"`
	Event  json.RawMessage `json:"event"`
	Token  string          `json:"token,omitempty"`
}

// EventHeader 事件头信息
type EventHeader struct {
	EventID    string    `json:"event_id"`
	EventType  EventType `json:"event_type"`
	CreateTime string    `json:"create_time"`
	Token      string    `json:"token"`
	AppID      string    `json:"app_id"`
	TenantKey  string    `json:"tenant_key"`
}

// {
//     "message": {
//         "chat_id": "oc_a7b4b44e69866856c9510d3cd6ed3b5d",
//         "chat_type": "p2p",
//         "content": "{\"text\":\"www\"}",
//         "create_time": "1773377115036",
//         "message_id": "om_x100b540d2b63cca0b32c1273314ffc7",
//         "message_type": "text",
//         "update_time": "1773377115036"
//     },
//     "sender": {
//         "sender_id": {
//             "open_id": "ou_0fec6d92414810535e1e3a5c96a021b8",
//             "union_id": "on_ac42fc6683e1ae22eab11707db6c5fcd",
//             "user_id": "553d5845"
//         },
//         "sender_type": "user",
//         "tenant_key": "10913ce36a825740"
//     }
// }

// MessageReceiveEvent 消息接收事件数据
type MessageReceiveEvent struct {
	Message struct {
		ChatID        string `json:"chat_id"`
		ChatType      string `json:"chat_type"`
		MessageID     string `json:"message_id"`
		MessageType   string `json:"msg_type"`
		Content       string `json:"content"`
		CreateTime    string `json:"create_time"`
		UpdateTime    string `json:"update_time"`
		MentionedAll  bool   `json:"mentioned_all"`
		MentionedInfo struct {
			MentionedList []string `json:"mention_list"`
		} `json:"mentioned_info"`
	} `json:"message"`
	Sender struct {
		SenderID struct {
			OpenID  string `json:"open_id"`
			UnionID string `json:"union_id"`
			UserID  string `json:"user_id"`
		} `json:"sender_id"`
		SenderType string `json:"sender_type"`
		TenantKey  string `json:"tenant_key"`
	} `json:"sender"`
}

// EventHandler 事件处理函数类型
type EventHandler func(event *EventMessage) error

// getWebSocketEndpoint 获取 WebSocket 端点
func (c *Client) getWebSocketEndpoint() (endpointResp EndpointResponse, err error) {
	body := map[string]string{
		"AppID":     c.Config.AppID,
		"AppSecret": c.Config.AppSecret,
	}
	respBody, err := c.request(
		WithURL("https://open.feishu.cn"),
		WithPath("/callback/ws/endpoint"),
		WithBody(body),
	)
	err = json.Unmarshal(respBody, &endpointResp)
	if endpointResp.Code != 0 {
		err = fmt.Errorf("获取端点失败：code=%d, msg=%s", endpointResp.Code, endpointResp.Msg)
	}
	return
}

// connectWebSocket 建立 WebSocket 连接
func (c *Client) connectWebSocket(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ws != nil {
		return nil
	}

	endpoint, err := c.getWebSocketEndpoint()
	if err != nil {
		return fmt.Errorf("获取端点失败：%w", err)
	}

	c.ws, _, err = websocket.DefaultDialer.Dial(endpoint.Data.URL, nil)
	if err != nil {
		return fmt.Errorf("WebSocket 拨号失败：%w", err)
	}
	go c.receiveMessageLoop(ctx)
	return nil
}

// receiveMessageLoop 接收消息循环
func (c *Client) receiveMessageLoop(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WebSocket] receiveMessageLoop panic: %v", r)
		}
		c.disconnectWebSocket()
		if err := c.reconnectWebSocket(ctx); err != nil {
			log.Printf("[WebSocket] 重连失败：%v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		c.mu.Lock()
		conn := c.ws
		c.mu.Unlock()

		if conn == nil {
			log.Printf("[WebSocket] 连接已关闭，退出接收循环")
			return
		}

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WebSocket] 读取消息失败：%v", err)
			return
		}

		if messageType != websocket.BinaryMessage {
			log.Printf("[WebSocket] 收到未知消息类型：%d", messageType)
			continue
		}

		go c.handleMessage(message)
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(msg []byte) {
	frame := &FramePB{}
	if err := frame.Unmarshal(msg); err != nil {
		log.Printf("[WebSocket] 解析 Frame 失败：%v", err)
		return
	}

	if FrameType(frame.Method) == FrameTypeControl {
		c.handleControlFrame(frame)
	} else if FrameType(frame.Method) == FrameTypeData {
		c.handleDataFrame(frame)
	}
}

// ClientConfig 服务端返回的客户端配置（用于合并到 WSConfig）
type WSConfig struct {
	ReconnectCount    int `json:"reconnect_count,omitempty"`
	ReconnectInterval int `json:"reconnect_interval,omitempty"`
	ReconnectNonce    int `json:"reconnect_nonce,omitempty"`
	PingInterval      int `json:"ping_interval,omitempty"`
}

// DefaultWSConfig 返回默认 WebSocket 配置
func DefaultWSConfig() *WSConfig {
	return &WSConfig{
		ReconnectCount:    3, // 无限重连
		ReconnectInterval: 2,
		ReconnectNonce:    30,
		PingInterval:      2,
	}
}

// handleControlFrame 处理控制帧
func (c *Client) handleControlFrame(frame *FramePB) {
	hs := headersToMap(frame.Headers)
	t := hs["type"]
	if t == string(MessageTypePong) {
		log.Printf("[WebSocket] 收到 Pong")
		if len(frame.Payload) == 0 {
			return
		}
		var conf WSConfig
		if err := json.Unmarshal(frame.Payload, &conf); err != nil {
			return
		}
		c.WSConfig = conf
	}
}

// handleDataFrame 处理数据帧
func (c *Client) handleDataFrame(frame *FramePB) {
	hs := headersToMap(frame.Headers)
	msgType := hs["type"]
	if msgType != string(MessageTypeEvent) {
		log.Printf("[WebSocket] 未知数据帧类型：%s", msgType)
		return
	}

	var eventMsg EventMessage
	if err := json.Unmarshal(frame.Payload, &eventMsg); err != nil {
		log.Printf("[WebSocket] 解析事件失败：%v", err)
		return
	}

	// log.Printf("[WebSocket] 收到事件：%s", eventMsg.Header.EventType)

	// 调用事件处理器
	c.mu.Lock()
	handler := c.wsHandler
	c.mu.Unlock()

	if handler != nil {
		if err := handler(&eventMsg); err != nil {
			log.Printf("[WebSocket] 事件处理失败：%v", err)
		}
	}

	c.sendResponse(frame, http.StatusOK)
}

// sendResponse 发送响应
func (c *Client) sendResponse(frame *FramePB, statusCode int) {
	response := map[string]interface{}{
		"code": statusCode,
	}

	payload, _ := json.Marshal(response)
	frame.Payload = payload

	data, err := frame.Marshal()
	if err != nil {
		log.Printf("[WebSocket] 序列化响应失败：%v", err)
		return
	}

	c.mu.Lock()
	conn := c.ws
	c.mu.Unlock()

	if conn == nil {
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("[WebSocket] 发送响应失败：%v", err)
	}
}

// pingLoop Ping 循环
func (c *Client) pingLoop(ctx context.Context) {
	d := time.Duration(c.WSConfig.PingInterval) * time.Minute
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.sendPing()
		}
	}
}

// sendPing 发送 Ping 帧
func (c *Client) sendPing() {
	pingFrame := &FramePB{
		Method:  int32(FrameTypeControl),
		Service: 1,
		Headers: []Header{{Key: "type", Value: string(MessageTypePing)}},
	}
	data, _ := pingFrame.Marshal()
	c.mu.Lock()
	conn := c.ws
	c.mu.Unlock()
	if conn == nil {
		return
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("[WebSocket] 发送 Ping 失败：%v", err)
	}
}

// reconnectWebSocket 重连逻辑
func (c *Client) reconnectWebSocket(ctx context.Context) error {
	// 首次重连随机抖动
	if c.WSConfig.ReconnectNonce > 0 {
		rand.Seed(time.Now().UnixNano())
		delay := time.Duration(rand.Intn(c.WSConfig.ReconnectNonce*1000)) * time.Millisecond
		log.Printf("[WebSocket] %v 后开始重连", delay)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	var i = 0
	for {
		log.Printf("[WebSocket] 尝试重连：%d/%d", i+1, c.WSConfig.ReconnectCount)
		err := c.connectWebSocket(ctx)
		if err == nil {
			// success
			return nil
		}
		select {
		case <-time.After(time.Duration(c.WSConfig.ReconnectInterval) * time.Minute):
		case <-ctx.Done():
			return ctx.Err()
		}
		// 如果 -1 表示无限重试
		if c.WSConfig.ReconnectCount == -1 {
			continue
		}
		i++
		if i >= c.WSConfig.ReconnectCount {
			break
		}
	}
	return fmt.Errorf("重连 %d 次后仍失败", c.WSConfig.ReconnectCount)
}

// disconnectWebSocket 断开连接
func (c *Client) disconnectWebSocket() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ws == nil {
		return
	}
	c.ws.Close()
	c.ws = nil
}

// headersToMap 将 Header 列表转换为 map
func headersToMap(headers []Header) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		result[h.Key] = h.Value
	}
	return result
}

// StartWithContext 启动 WebSocket 连接并接收事件（阻塞）
// 可通过 ctx 控制退出
func (c *Client) Start(ctx context.Context, fn EventHandler) error {
	c.mu.Lock()
	c.wsHandler = fn
	c.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := c.connectWebSocket(ctx)
	if err != nil {
		log.Printf("[WebSocket] 连接失败：%v", err)
		if err = c.reconnectWebSocket(ctx); err != nil {
			// 如果是 context 取消导致的错误，直接返回
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
	}

	go c.pingLoop(ctx)

	// 等待上下文取消
	<-ctx.Done()
	return ctx.Err()
}

// Close 关闭 WebSocket 连接
func (c *Client) Close() error {
	c.disconnectWebSocket()
	return nil
}
