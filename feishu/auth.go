package feishu

import (
	"encoding/json"
)

type UserAccessTokenResponse struct {
	ResponseBase
	Data map[string]interface{} `json:"data"`
}

type AppAccessTokenResponse struct {
	ResponseBase
	Expire         int    `json:"expire"`
	AppAccessToken string `json:"app_access_token"`
}
type TenantAccessTokenResponse struct {
	ResponseBase
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// 获取 user_access_token
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/access_token/create
func (c *Client) GetUserAccessToken() (out *UserAccessTokenResponse, err error) {
	data, err := c.RequestWithAppSecret("/authen/v1/access_token")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return out, err
}

// 刷新 user_access_token
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/refresh_access_token/create
func (c *Client) RefreshUserAccessToken(refreshToken string, grantType string) (out *UserAccessTokenResponse, err error) {
	data, err := c.RequestWithAccessToken("/authen/v1/refresh_access_token", map[string]string{
		"refresh_token": refreshToken,
		"grant_type":    grantType,
	})
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return
}

// 商店应用获取 app_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/app_access_token
func (c *Client) GetAppAccessToken() (out *AppAccessTokenResponse, err error) {
	data, err := c.RequestWithAppSecret("/auth/v3/app_access_token")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return
}

// 自建应用获取 app_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/app_access_token_internal
func (c *Client) GetAppAccessTokenInternal() (out *AppAccessTokenResponse, err error) {
	data, err := c.RequestWithAppSecret("/auth/v3/app_access_token/internal")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return out, err
}

// 商店应用获取 tenant_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token
func (c *Client) GetTenantAccessToken() (out *TenantAccessTokenResponse, err error) {
	data, err := c.RequestWithAppSecret("/auth/v3/tenant_access_token")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return
}

// 自建应用获取 tenant_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token_internal
func (c *Client) GetTenantAccessTokenInternal() (out *TenantAccessTokenResponse, err error) {
	data, err := c.RequestWithAppSecret("/auth/v3/tenant_access_token/internal")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return out, err
}

// 重新获取 app_ticket
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/app_ticket_resend
func (c *Client) ResendAppTicket() (out *ResponseBase, err error) {
	data, err := c.RequestWithAppSecret("/auth/v3/app_ticket/resend")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &out)
	return
}
