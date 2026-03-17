package feishu

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
)

// https://open.feishu.cn/document/server-docs/im-v1/file/create
func (c *Client) Upload() {

}

// https://open.feishu.cn/document/server-docs/im-v1/file/get
func (c *Client) Download() {

}

type UploadImageResponse struct {
	ResponseBase
	Data struct {
		ImageKey string `json:"image_key"`
	} `json:"data"`
}

// @docs https://open.feishu.cn/document/server-docs/im-v1/image/create?appId=cli_a37e09c80539d00c
// POST https://open.feishu.cn/open-apis/im/v1/images
/*
curl --location --request POST 'https://open.feishu.cn/open-apis/im/v1/images' \
--header 'Authorization: Bearer t-7f1b******8e560' \
--header 'Content-Type: multipart/form-data' \
--form 'image_type="message"' \
--form 'image=@"/xxx/二进制文件"'
*/
func (c *Client) UploadImage(imageType string, image io.Reader) (out *UploadImageResponse, err error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	err = writer.WriteField("image_type", imageType)
	if err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile("image", "image")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, image)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	data, err := c.request(
		WithPath("/im/v1/images"),
		WithAccessToken(c.AccessToken),
		WithContentType(writer.FormDataContentType()),
		WithBody(&body),
	)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &out)
	return
}

// @docs https://open.feishu.cn/document/server-docs/im-v1/image/get
// GET https://open.feishu.cn/open-apis/im/v1/images/:image_key
func (c *Client) DownloadImage() {

}
