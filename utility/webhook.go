package utility

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"qq-bot-backend/internal/consts"
)

func SendWebhookRequest(ctx context.Context, header, method, url string, payload ...any,
) (statusCode int, contentType string, body []byte, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	if header != "" {
		c.SetHeaderRaw(header)
	}

	resp, err := c.DoRequest(ctx, method, url, payload)
	if err != nil || resp == nil {
		return
	}
	defer resp.Close()

	return resp.StatusCode, resp.Header.Get("Content-Type"), resp.ReadAll(), nil
}
