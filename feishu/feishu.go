package feishu

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	AppID       string `json:"app_id,omitempty"`
	AppSecret   string `json:"app_secret,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	// WebSocket 配置
	WSConfig *WSConfig `json:"ws_config,omitempty"`
}

type Client struct {
	*Config
	*http.Client
	ws        *websocket.Conn
	wsURL     string
	wsHandler EventHandler
	mu        sync.Mutex
}

// WSConfig WebSocket 客户端配置
type WSConfig struct {
	AutoReconnect     bool          `json:"auto_reconnect,omitempty"`
	ReconnectCount    int           `json:"reconnect_count,omitempty"`     // 重连次数，-1 为无限
	ReconnectInterval time.Duration `json:"reconnect_interval,omitempty"`  // 重连间隔
	ReconnectNonce    int           `json:"reconnect_nonce,omitempty"`     // 首次重连随机抖动（秒）
	PingInterval      time.Duration `json:"ping_interval,omitempty"`       // Ping 间隔
}

// DefaultWSConfig 返回默认 WebSocket 配置
func DefaultWSConfig() *WSConfig {
	return &WSConfig{
		AutoReconnect:     true,
		ReconnectCount:    -1, // 无限重连
		ReconnectInterval: 2 * time.Minute,
		ReconnectNonce:    30,
		PingInterval:      2 * time.Minute,
	}
}

// ClientConfig 服务端返回的客户端配置（用于合并到 WSConfig）
type ClientConfig struct {
	ReconnectCount    int `json:"reconnect_count,omitempty"`
	ReconnectInterval int `json:"reconnect_interval,omitempty"`
	ReconnectNonce    int `json:"reconnect_nonce,omitempty"`
	PingInterval      int `json:"ping_interval,omitempty"`
}

type ResponseBase struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func NewClient(config *Config) *Client {
	return &Client{
		Config: config,
		Client: http.DefaultClient,
	}
}

// RequestOption 请求选项
type RequestOption func(*requestOptions)

type requestOptions struct {
	method      string
	baseURL     string
	contentType string
	headers     map[string]string
	body        interface{}
}

func WithMethod(method string) RequestOption {
	return func(opts *requestOptions) {
		opts.method = method
	}
}

func WithBaseURL(baseURL string) RequestOption {
	return func(opts *requestOptions) {
		opts.baseURL = baseURL
	}
}

func WithContentType(contentType string) RequestOption {
	return func(opts *requestOptions) {
		opts.contentType = contentType
	}
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(opts *requestOptions) {
		opts.headers = headers
	}
}

func WithBody(body interface{}) RequestOption {
	return func(opts *requestOptions) {
		opts.body = body
	}
}

// Request 发送 HTTP 请求（通用版本）
func (c *Client) Request(path string, opts ...RequestOption) (out []byte, err error) {
	options := &requestOptions{
		method:      http.MethodPost,
		baseURL:     "https://open.feishu.cn/open-apis",
		contentType: "application/json; charset=utf-8",
		headers:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(options)
	}

	var bodyReader io.Reader
	if options.body != nil {
		payload, err := json.Marshal(options.body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(options.method, options.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", options.contentType)
	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (c *Client) RequestWithAppSecret(path string) (out []byte, err error) {
	return c.Request(path, WithBody(map[string]string{
		"app_id":     c.Config.AppID,
		"app_secret": c.Config.AppSecret,
	}))
}

func (c *Client) SetAccessToken(accessToken string) {
	// log.Println("accessToken:", accessToken)
	c.Config.AccessToken = accessToken
}

func (c *Client) RequestWithAccessToken(path string, data interface{}) (out []byte, err error) {
	return c.Request(path,
		WithHeaders(map[string]string{
			"Authorization": "Bearer " + c.Config.AccessToken,
		}),
		WithBody(data),
	)
}
