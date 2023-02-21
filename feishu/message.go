package feishu

import (
	"encoding/json"
	"fmt"
)

type ReceiveIdType string

const (
	OPEN_ID  ReceiveIdType = "open_id"
	USER_ID  ReceiveIdType = "user_id"
	UNION_ID ReceiveIdType = "union_id"
	EMAIL    ReceiveIdType = "email"
	CHAT_ID  ReceiveIdType = "chat_id"
)

type Message struct {
	ReceiveIdType ReceiveIdType
	ReceiveId     string `json:"receive_id"`
	Type          string `json:"msg_type"`
	Content       string `json:"content"`
	UUID          string `json:"uuid,omitempty"`
	// 回复消息相关
	RootID   string `json:"root_id,omitempty"`
	ReplyID  string `json:"reply_id,omitempty"`
}

type MessageResponseData struct {
	MessageID string `json:"message_id"`
	RootID    string `json:"root_id,omitempty"`
}

type MessageResponse struct {
	ResponseBase
	Data MessageResponseData `json:"data"`
}

func NewTextMessage(content string) (message Message) {
	var data, _ = json.Marshal(map[string]string{
		"text": content,
	})
	message.Type = "text"
	message.Content = string(data)
	return
}

// https://open.feishu.cn/document/server-docs/im-v1/message/create
func (c *Client) SendMessage(message *Message) (out *MessageResponse, err error) {
	if message.ReceiveIdType == "" {
		message.ReceiveIdType = OPEN_ID
	}
	path := fmt.Sprintf("/im/v1/messages?receive_id_type=%s", message.ReceiveIdType)
	data, err := c.RequestWithAccessToken(path, message)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return
}

// SendReplyMessage 回复指定消息（引用回复）
func (c *Client) SendReplyMessage(message *Message, messageID string) (out *MessageResponse, err error) {
	// 使用专门的回复消息 API：POST /im/v1/messages/:message_id/reply
	path := fmt.Sprintf("/im/v1/messages/%s/reply", messageID)
	
	data, err := c.RequestWithAccessToken(path, message)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &out)
	return
}

func (c *Client) SendTextMessage(receiveId string, content string) (out *MessageResponse, err error) {
	return c.SendTextMessageWithType(receiveId, content, OPEN_ID)
}

func (c *Client) SendTextMessageWithType(receiveId string, content string, idType ReceiveIdType) (out *MessageResponse, err error) {
	message := NewTextMessage(content)
	message.ReceiveId = receiveId
	message.ReceiveIdType = idType
	return c.SendMessage(&message)
}

// ReactionType 表情回应类型
type ReactionType string

const (
	REACTION_THUMBS_UP   ReactionType = "LIKE"      // 点赞 👍
	REACTION_HEART       ReactionType = "LOVE"      // 爱心 ❤️
	REACTION_LAUGH       ReactionType = "LAUGH"     // 大笑 😄
	REACTION_SURPRISED   ReactionType = "SURPRISED" // 惊讶 😮
	REACTION_SAD         ReactionType = "SAD"       // 难过 😢
	REACTION_DISLIKE     ReactionType = "DISLIKE"   // 不喜欢 👎
	REACTION_CELEBRATION ReactionType = "CELEBRATION" // 庆祝 🎉
	REACTION_SMILE       ReactionType = "SMILE"     // 微笑 😊
)

// Operator 操作人信息
type Operator struct {
	OperatorID   string `json:"operator_id"`
	OperatorType string `json:"operator_type"`
}

// ReactionTypeObject 表情类型对象
type ReactionTypeObject struct {
	EmojiType ReactionType `json:"emoji_type"`
}

// AddReactionRequest 添加表情回应的请求体
type AddReactionRequest struct {
	ReactionID  string             `json:"reaction_id,omitempty"`
	Operator    Operator           `json:"operator,omitempty"`
	ActionTime  string             `json:"action_time,omitempty"`
	ReactionType ReactionTypeObject `json:"reaction_type"`
}

// ReactionResponse 表情回应响应
type ReactionResponse struct {
	ResponseBase
	Data struct {
		ReactionID string `json:"reaction_id"`
	} `json:"data"`
}

// AddMessageReaction 为消息添加表情回应
func (c *Client) AddMessageReaction(messageID string, reactionType ReactionType) (out *ReactionResponse, err error) {
	path := fmt.Sprintf("/im/v1/messages/%s/reactions", messageID)
	body := AddReactionRequest{
		ReactionType: ReactionTypeObject{
			EmojiType: reactionType,
		},
	}
	data, err := c.RequestWithAccessToken(path, body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &out)
	return
}
