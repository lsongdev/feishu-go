package feishu

import (
	"encoding/json"
	"strings"
)

type User struct {
	UserID string `json:"user_id"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
}

type FeishuContactResponseData struct {
	UserList []User `json:"user_list"`
}

type FeishuContactResponse struct {
	ResponseBase
	Data FeishuContactResponseData `json:"data"`
}

// 通过手机号或邮箱获取用户 ID
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/user/batch_get_id
func (c *Client) GetUsersBy(userIdType string, query []string) ([]User, error) {
	var emails []string
	var mobiles []string
	for _, item := range query {
		if strings.Contains(item, "@") {
			emails = append(emails, item)
		} else {
			mobiles = append(mobiles, item)
		}
	}
	input := map[string][]string{
		"emails":  emails,
		"mobiles": mobiles,
	}
	data, err := c.RequestWithAccessToken("/contact/v3/users/batch_get_id?user_id_type="+userIdType, input)
	if err != nil {
		return nil, err
	}
	var out FeishuContactResponse
	err = json.Unmarshal(data, &out)
	return out.Data.UserList, err
}
