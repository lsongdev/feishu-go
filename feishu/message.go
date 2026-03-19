package feishu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReceiveIdType string

const (
	OPEN_ID  ReceiveIdType = "open_id"
	USER_ID  ReceiveIdType = "user_id"
	UNION_ID ReceiveIdType = "union_id"
	EMAIL    ReceiveIdType = "email"
	CHAT_ID  ReceiveIdType = "chat_id"
)

// https://open.feishu.cn/document/server-docs/im-v1/message/create
type Message struct {
	ReceiveIdType string `json:"-"`
	ReceiveID     string `json:"receive_id"`
	Type          string `json:"msg_type"`
	Content       string `json:"content"`
	UUID          string `json:"uuid,omitempty"`
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

// https://open.feishu.cn/document/server-docs/im-v1/message-content-description/create_json#7111df05
func NewImageMessage(imageKey string) (message Message) {
	var data, _ = json.Marshal(map[string]string{
		"image_key": imageKey,
	})
	message.Type = "image"
	message.Content = string(data)
	return
}

// https://open.feishu.cn/document/server-docs/im-v1/message-content-description/create_json#45e0953e
func NewPostMessage(lang, title string, paragraphs ...[]PostTag) (message Message) {
	var data, _ = json.Marshal(map[string]any{
		lang: map[string]any{
			"title":   title,
			"content": paragraphs,
		},
	})
	message.Type = "post"
	message.Content = string(data)
	return
}

type PostTag = map[string]any

func NewPostParagraph(lines ...PostTag) []PostTag {
	return lines
}

func NewMarkdownTag(content string) PostTag {
	return PostTag{
		"tag":  "md",
		"text": content,
	}
}

func NewMarkdownMessage(title, content string) Message {
	return NewPostMessage("zh_cn", title,
		NewPostParagraph(
			NewMarkdownTag(content),
		),
	)
}

// https://open.feishu.cn/document/server-docs/im-v1/message/create
func (c *Client) SendMessage(message *Message) (out *MessageResponse, err error) {
	path := fmt.Sprintf("/im/v1/messages?receive_id_type=%s", message.ReceiveIdType)
	err = c.RequestWithAccessToken(path, message, &out)
	return
}

// SendReplyMessage 回复指定消息（引用回复）
// https://open.feishu.cn/document/server-docs/im-v1/message/reply
func (c *Client) SendReplyMessage(messageID string, message *Message) (out *MessageResponse, err error) {
	path := fmt.Sprintf("/im/v1/messages/%s/reply", messageID)
	err = c.RequestWithAccessToken(path, message, &out)
	return
}

// ReactionType 表情回应类型
type EmojiType string

// https://open.feishu.cn/document/server-docs/im-v1/message-reaction/emojis-introduce
// copy([...document.querySelectorAll('td > p')].map(x => x.textContent).filter(x => x).map(x => `EMOJI_${x}    EmojiType = "${x}"`).join('\n'))
const (
	EMOJI_OK                       EmojiType = "OK"
	EMOJI_THUMBSUP                 EmojiType = "THUMBSUP"
	EMOJI_THANKS                   EmojiType = "THANKS"
	EMOJI_MUSCLE                   EmojiType = "MUSCLE"
	EMOJI_FINGERHEART              EmojiType = "FINGERHEART"
	EMOJI_APPLAUSE                 EmojiType = "APPLAUSE"
	EMOJI_FISTBUMP                 EmojiType = "FISTBUMP"
	EMOJI_JIAYI                    EmojiType = "JIAYI"
	EMOJI_DONE                     EmojiType = "DONE"
	EMOJI_SMILE                    EmojiType = "SMILE"
	EMOJI_BLUSH                    EmojiType = "BLUSH"
	EMOJI_LAUGH                    EmojiType = "LAUGH"
	EMOJI_SMIRK                    EmojiType = "SMIRK"
	EMOJI_LOL                      EmojiType = "LOL"
	EMOJI_FACEPALM                 EmojiType = "FACEPALM"
	EMOJI_LOVE                     EmojiType = "LOVE"
	EMOJI_WINK                     EmojiType = "WINK"
	EMOJI_PROUD                    EmojiType = "PROUD"
	EMOJI_WITTY                    EmojiType = "WITTY"
	EMOJI_SMART                    EmojiType = "SMART"
	EMOJI_SCOWL                    EmojiType = "SCOWL"
	EMOJI_THINKING                 EmojiType = "THINKING"
	EMOJI_SOB                      EmojiType = "SOB"
	EMOJI_CRY                      EmojiType = "CRY"
	EMOJI_ERROR                    EmojiType = "ERROR"
	EMOJI_NOSEPICK                 EmojiType = "NOSEPICK"
	EMOJI_HAUGHTY                  EmojiType = "HAUGHTY"
	EMOJI_SLAP                     EmojiType = "SLAP"
	EMOJI_SPITBLOOD                EmojiType = "SPITBLOOD"
	EMOJI_TOASTED                  EmojiType = "TOASTED"
	EMOJI_GLANCE                   EmojiType = "GLANCE"
	EMOJI_DULL                     EmojiType = "DULL"
	EMOJI_INNOCENTSMILE            EmojiType = "INNOCENTSMILE"
	EMOJI_JOYFUL                   EmojiType = "JOYFUL"
	EMOJI_WOW                      EmojiType = "WOW"
	EMOJI_TRICK                    EmojiType = "TRICK"
	EMOJI_YEAH                     EmojiType = "YEAH"
	EMOJI_ENOUGH                   EmojiType = "ENOUGH"
	EMOJI_TEARS                    EmojiType = "TEARS"
	EMOJI_EMBARRASSED              EmojiType = "EMBARRASSED"
	EMOJI_KISS                     EmojiType = "KISS"
	EMOJI_SMOOCH                   EmojiType = "SMOOCH"
	EMOJI_DROOL                    EmojiType = "DROOL"
	EMOJI_OBSESSED                 EmojiType = "OBSESSED"
	EMOJI_MONEY                    EmojiType = "MONEY"
	EMOJI_TEASE                    EmojiType = "TEASE"
	EMOJI_SHOWOFF                  EmojiType = "SHOWOFF"
	EMOJI_COMFORT                  EmojiType = "COMFORT"
	EMOJI_CLAP                     EmojiType = "CLAP"
	EMOJI_PRAISE                   EmojiType = "PRAISE"
	EMOJI_STRIVE                   EmojiType = "STRIVE"
	EMOJI_XBLUSH                   EmojiType = "XBLUSH"
	EMOJI_SILENT                   EmojiType = "SILENT"
	EMOJI_WAVE                     EmojiType = "WAVE"
	EMOJI_WHAT                     EmojiType = "WHAT"
	EMOJI_FROWN                    EmojiType = "FROWN"
	EMOJI_SHY                      EmojiType = "SHY"
	EMOJI_DIZZY                    EmojiType = "DIZZY"
	EMOJI_LOOKDOWN                 EmojiType = "LOOKDOWN"
	EMOJI_CHUCKLE                  EmojiType = "CHUCKLE"
	EMOJI_WAIL                     EmojiType = "WAIL"
	EMOJI_CRAZY                    EmojiType = "CRAZY"
	EMOJI_WHIMPER                  EmojiType = "WHIMPER"
	EMOJI_HUG                      EmojiType = "HUG"
	EMOJI_BLUBBER                  EmojiType = "BLUBBER"
	EMOJI_WRONGED                  EmojiType = "WRONGED"
	EMOJI_HUSKY                    EmojiType = "HUSKY"
	EMOJI_SHHH                     EmojiType = "SHHH"
	EMOJI_SMUG                     EmojiType = "SMUG"
	EMOJI_ANGRY                    EmojiType = "ANGRY"
	EMOJI_HAMMER                   EmojiType = "HAMMER"
	EMOJI_SHOCKED                  EmojiType = "SHOCKED"
	EMOJI_TERROR                   EmojiType = "TERROR"
	EMOJI_PETRIFIED                EmojiType = "PETRIFIED"
	EMOJI_SKULL                    EmojiType = "SKULL"
	EMOJI_SWEAT                    EmojiType = "SWEAT"
	EMOJI_SPEECHLESS               EmojiType = "SPEECHLESS"
	EMOJI_SLEEP                    EmojiType = "SLEEP"
	EMOJI_DROWSY                   EmojiType = "DROWSY"
	EMOJI_YAWN                     EmojiType = "YAWN"
	EMOJI_SICK                     EmojiType = "SICK"
	EMOJI_PUKE                     EmojiType = "PUKE"
	EMOJI_BETRAYED                 EmojiType = "BETRAYED"
	EMOJI_HEADSET                  EmojiType = "HEADSET"
	EMOJI_EatingFood               EmojiType = "EatingFood"
	EMOJI_MeMeMe                   EmojiType = "MeMeMe"
	EMOJI_Sigh                     EmojiType = "Sigh"
	EMOJI_Typing                   EmojiType = "Typing"
	EMOJI_Lemon                    EmojiType = "Lemon"
	EMOJI_Get                      EmojiType = "Get"
	EMOJI_LGTM                     EmojiType = "LGTM"
	EMOJI_OnIt                     EmojiType = "OnIt"
	EMOJI_OneSecond                EmojiType = "OneSecond"
	EMOJI_VRHeadset                EmojiType = "VRHeadset"
	EMOJI_YouAreTheBest            EmojiType = "YouAreTheBest"
	EMOJI_SALUTE                   EmojiType = "SALUTE"
	EMOJI_SHAKE                    EmojiType = "SHAKE"
	EMOJI_HIGHFIVE                 EmojiType = "HIGHFIVE"
	EMOJI_UPPERLEFT                EmojiType = "UPPERLEFT"
	EMOJI_ThumbsDown               EmojiType = "ThumbsDown"
	EMOJI_SLIGHT                   EmojiType = "SLIGHT"
	EMOJI_TONGUE                   EmojiType = "TONGUE"
	EMOJI_EYESCLOSED               EmojiType = "EYESCLOSED"
	EMOJI_RoarForYou               EmojiType = "RoarForYou"
	EMOJI_CALF                     EmojiType = "CALF"
	EMOJI_BEAR                     EmojiType = "BEAR"
	EMOJI_BULL                     EmojiType = "BULL"
	EMOJI_RAINBOWPUKE              EmojiType = "RAINBOWPUKE"
	EMOJI_ROSE                     EmojiType = "ROSE"
	EMOJI_HEART                    EmojiType = "HEART"
	EMOJI_PARTY                    EmojiType = "PARTY"
	EMOJI_LIPS                     EmojiType = "LIPS"
	EMOJI_BEER                     EmojiType = "BEER"
	EMOJI_CAKE                     EmojiType = "CAKE"
	EMOJI_GIFT                     EmojiType = "GIFT"
	EMOJI_CUCUMBER                 EmojiType = "CUCUMBER"
	EMOJI_Drumstick                EmojiType = "Drumstick"
	EMOJI_Pepper                   EmojiType = "Pepper"
	EMOJI_CANDIEDHAWS              EmojiType = "CANDIEDHAWS"
	EMOJI_BubbleTea                EmojiType = "BubbleTea"
	EMOJI_Coffee                   EmojiType = "Coffee"
	EMOJI_Yes                      EmojiType = "Yes"
	EMOJI_No                       EmojiType = "No"
	EMOJI_OKR                      EmojiType = "OKR"
	EMOJI_CheckMark                EmojiType = "CheckMark"
	EMOJI_CrossMark                EmojiType = "CrossMark"
	EMOJI_MinusOne                 EmojiType = "MinusOne"
	EMOJI_Hundred                  EmojiType = "Hundred"
	EMOJI_AWESOMEN                 EmojiType = "AWESOMEN"
	EMOJI_Pin                      EmojiType = "Pin"
	EMOJI_Alarm                    EmojiType = "Alarm"
	EMOJI_Loudspeaker              EmojiType = "Loudspeaker"
	EMOJI_Trophy                   EmojiType = "Trophy"
	EMOJI_Fire                     EmojiType = "Fire"
	EMOJI_BOMB                     EmojiType = "BOMB"
	EMOJI_Music                    EmojiType = "Music"
	EMOJI_XmasTree                 EmojiType = "XmasTree"
	EMOJI_Snowman                  EmojiType = "Snowman"
	EMOJI_XmasHat                  EmojiType = "XmasHat"
	EMOJI_FIREWORKS                EmojiType = "FIREWORKS"
	EMOJI_2022                     EmojiType = "2022"
	EMOJI_REDPACKET                EmojiType = "REDPACKET"
	EMOJI_FORTUNE                  EmojiType = "FORTUNE"
	EMOJI_LUCK                     EmojiType = "LUCK"
	EMOJI_FIRECRACKER              EmojiType = "FIRECRACKER"
	EMOJI_StickyRiceBalls          EmojiType = "StickyRiceBalls"
	EMOJI_HEARTBROKEN              EmojiType = "HEARTBROKEN"
	EMOJI_POOP                     EmojiType = "POOP"
	EMOJI_StatusFlashOfInspiration EmojiType = "StatusFlashOfInspiration"
	EMOJI_18X                      EmojiType = "18X"
	EMOJI_CLEAVER                  EmojiType = "CLEAVER"
	EMOJI_Soccer                   EmojiType = "Soccer"
	EMOJI_Basketball               EmojiType = "Basketball"
	EMOJI_GeneralDoNotDisturb      EmojiType = "GeneralDoNotDisturb"
	EMOJI_Status_PrivateMessage    EmojiType = "Status_PrivateMessage"
	EMOJI_GeneralInMeetingBusy     EmojiType = "GeneralInMeetingBusy"
	EMOJI_StatusReading            EmojiType = "StatusReading"
	EMOJI_StatusInFlight           EmojiType = "StatusInFlight"
	EMOJI_GeneralBusinessTrip      EmojiType = "GeneralBusinessTrip"
	EMOJI_GeneralWorkFromHome      EmojiType = "GeneralWorkFromHome"
	EMOJI_StatusEnjoyLife          EmojiType = "StatusEnjoyLife"
	EMOJI_GeneralTravellingCar     EmojiType = "GeneralTravellingCar"
	EMOJI_StatusBus                EmojiType = "StatusBus"
	EMOJI_GeneralSun               EmojiType = "GeneralSun"
	EMOJI_GeneralMoonRest          EmojiType = "GeneralMoonRest"
	EMOJI_MoonRabbit               EmojiType = "MoonRabbit"
	EMOJI_Mooncake                 EmojiType = "Mooncake"
	EMOJI_JubilantRabbit           EmojiType = "JubilantRabbit"
	EMOJI_TV                       EmojiType = "TV"
	EMOJI_Movie                    EmojiType = "Movie"
	EMOJI_Pumpkin                  EmojiType = "Pumpkin"
	EMOJI_BeamingFace              EmojiType = "BeamingFace"
	EMOJI_Delighted                EmojiType = "Delighted"
	EMOJI_ColdSweat                EmojiType = "ColdSweat"
	EMOJI_FullMoonFace             EmojiType = "FullMoonFace"
	EMOJI_Partying                 EmojiType = "Partying"
	EMOJI_GoGoGo                   EmojiType = "GoGoGo"
	EMOJI_ThanksFace               EmojiType = "ThanksFace"
	EMOJI_SaluteFace               EmojiType = "SaluteFace"
	EMOJI_Shrug                    EmojiType = "Shrug"
	EMOJI_ClownFace                EmojiType = "ClownFace"
	EMOJI_HappyDragon              EmojiType = "HappyDragon"
)

// Operator 操作人信息
type Operator struct {
	OperatorID   string `json:"operator_id"`
	OperatorType string `json:"operator_type"`
}

// ReactionTypeObject 表情类型对象
type ReactionType struct {
	EmojiType EmojiType `json:"emoji_type"`
}

// AddReactionRequest 添加表情回应的请求体
type ReactionRequest struct {
	ReactionType ReactionType `json:"reaction_type"`
}

// ReactionResponse 表情回应响应
type ReactionResponse struct {
	Data struct {
		ReactionID string `json:"reaction_id"`
	} `json:"data"`
}

func (c *Client) ListReaction(messageID string) (err error) {
	path := fmt.Sprintf("/im/v1/messages/%s/reactions", messageID)
	data, err := c.request(
		WithMethod(http.MethodGet),
		WithPath(path),
		WithAccessToken(c.AccessToken),
	)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	return
}

// AddMessageReaction 为消息添加表情回应
// https://open.feishu.cn/document/server-docs/im-v1/message-reaction/create
func (c *Client) AddReaction(messageID string, emoji EmojiType) (out *ReactionResponse, err error) {
	path := fmt.Sprintf("/im/v1/messages/%s/reactions", messageID)
	body := ReactionRequest{
		ReactionType: ReactionType{
			EmojiType: emoji,
		},
	}
	err = c.RequestWithAccessToken(path, body, &out)
	return
}

func (c *Client) RemoveReaction(messageID string) (err error) {
	path := fmt.Sprintf("/im/v1/messages/%s/reactions", messageID)
	data, err := c.request(
		WithMethod(http.MethodDelete),
		WithPath(path),
		WithAccessToken(c.AccessToken),
	)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	return
}
