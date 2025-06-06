package utility

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"qq-bot-backend/internal/consts"
	"sync"
)

var (
	sem chan struct{}

	initSem = sync.OnceFunc(func() {
		// A maximum of 32 concurrent requests can be made
		sem = make(chan struct{}, 32)
	})
)

func SendWebhookRequest(ctx context.Context, header, method, url string, payload ...any,
) (statusCode int, contentType string, body []byte, err error) {
	initSem()
	sem <- struct{}{}
	defer func() { <-sem }()

	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	if header != "" {
		c.SetHeaderRaw(header)
	}

	resp, err := c.DoRequest(ctx, method, url, payload...)
	if err != nil || resp == nil {
		return
	}
	defer resp.Close()

	return resp.StatusCode, resp.Header.Get("Content-Type"), resp.ReadAll(), nil
}
