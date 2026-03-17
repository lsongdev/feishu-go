package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Config struct {
	AppID       string `json:"app_id,omitempty"`
	AppSecret   string `json:"app_secret,omitempty"`
	AccessToken string `json:"access_token,omitempty"`

	WSConfig
}

type Client struct {
	*Config
	*http.Client
	ws              *websocket.Conn
	mu              sync.Mutex
	IncomingMessage chan *EventMessage
}

type ResponseBase struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func NewClient(config *Config) *Client {
	return &Client{
		Config:          config,
		Client:          http.DefaultClient,
		IncomingMessage: make(chan *EventMessage),
	}
}

// RequestOption 请求选项
type RequestOption func(*requestOptions)

type requestOptions struct {
	method      string
	base        string
	path        string
	contentType string
	headers     map[string]string
	body        any
}

func WithMethod(method string) RequestOption {
	return func(opts *requestOptions) {
		opts.method = method
	}
}

func WithURL(baseURL string) RequestOption {
	return func(opts *requestOptions) {
		opts.base = baseURL
	}
}

func WithPath(path string) RequestOption {
	return func(opts *requestOptions) {
		opts.path = path
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

func WithAccessToken(token string) RequestOption {
	return func(opts *requestOptions) {
		opts.headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
	}
}

func WithBody(body interface{}) RequestOption {
	return func(opts *requestOptions) {
		opts.body = body
	}
}

func (c *Client) request(opts ...RequestOption) (out []byte, err error) {
	options := &requestOptions{
		method:      http.MethodPost,
		path:        "",
		base:        "https://open.feishu.cn/open-apis",
		contentType: "application/json; charset=utf-8",
		headers:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(options)
	}

	var bodyReader io.Reader
	if options.body != nil {
		// 如果 body 是 io.Reader 类型，直接使用（用于 multipart/form-data）
		if reader, ok := options.body.(io.Reader); ok {
			bodyReader = reader
		} else {
			// 否则序列化为 JSON
			payload, err := json.Marshal(options.body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewBuffer(payload)
		}
	}

	url := options.base + options.path
	req, err := http.NewRequest(options.method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if options.contentType != "" {
		req.Header.Set("content-type", options.contentType)
	}
	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp ResponseBase
	err = json.Unmarshal(data, &resp)
	if resp.Code != 0 {
		return nil, fmt.Errorf("response error: %s", resp.Msg)
	}
	// println(string(data))
	return data, err
}

func (c *Client) RequestWithAppSecret(path string, out any) (err error) {
	data, err := c.request(
		WithPath(path),
		WithBody(map[string]string{
			"app_id":     c.Config.AppID,
			"app_secret": c.Config.AppSecret,
		}))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, out)
	return
}

func (c *Client) RequestWithAccessToken(path string, params any, out any) (err error) {
	data, err := c.request(
		WithPath(path),
		WithAccessToken(c.AccessToken),
		WithBody(params),
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, out)
	return
}

func (c *Client) SetAccessToken(accessToken string) {
	// log.Println("accessToken:", accessToken)
	c.Config.AccessToken = accessToken
}
