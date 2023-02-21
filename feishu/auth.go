package feishu

type UserAccessTokenResponse struct {
	Data map[string]interface{} `json:"data"`
}

type AppAccessTokenResponse struct {
	Expire         int    `json:"expire"`
	AppAccessToken string `json:"app_access_token"`
}
type TenantAccessTokenResponse struct {
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// 获取 user_access_token
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/access_token/create
func (c *Client) GetUserAccessToken() (out *UserAccessTokenResponse, err error) {
	err = c.RequestWithAppSecret("/authen/v1/access_token", &out)
	return
}

// 商店应用获取 app_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/app_access_token
func (c *Client) GetAppAccessToken() (out *AppAccessTokenResponse, err error) {
	err = c.RequestWithAppSecret("/auth/v3/app_access_token", &out)
	return
}

// 自建应用获取 app_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/app_access_token_internal
func (c *Client) GetAppAccessTokenInternal() (out *AppAccessTokenResponse, err error) {
	err = c.RequestWithAppSecret("/auth/v3/app_access_token/internal", &out)
	return out, err
}

// 商店应用获取 tenant_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token
func (c *Client) GetTenantAccessToken() (out *TenantAccessTokenResponse, err error) {
	err = c.RequestWithAppSecret("/auth/v3/tenant_access_token", &out)
	return
}

// 自建应用获取 tenant_access_token
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token_internal
func (c *Client) GetTenantAccessTokenInternal() (out *TenantAccessTokenResponse, err error) {
	err = c.RequestWithAppSecret("/auth/v3/tenant_access_token/internal", &out)
	return
}

// 刷新 user_access_token
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/refresh_access_token/create
func (c *Client) RefreshUserAccessToken(refreshToken string, grantType string) (out *UserAccessTokenResponse, err error) {
	params := map[string]string{
		"refresh_token": refreshToken,
		"grant_type":    grantType,
	}
	err = c.RequestWithAccessToken("/authen/v1/refresh_access_token", params, &out)
	return
}

// 重新获取 app_ticket
// https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/app_ticket_resend
func (c *Client) ResendAppTicket() (err error) {
	err = c.RequestWithAppSecret("/auth/v3/app_ticket/resend", nil)
	return
}
