package module

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"qq-bot-backend/internal/consts"
)

func (s *sModule) WebhookGetHeadConnectOptionsTrace(ctx context.Context, header, method, url string) (
	statusCode int, contentType string, body []byte, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	if header != "" {
		c.SetHeaderRaw(header)
	}
	resp, err := c.DoRequest(ctx, method, url)
	if err != nil || resp == nil {
		return
	}
	defer resp.Close()
	statusCode = resp.StatusCode
	contentType = resp.Header.Get("Content-Type")
	body = resp.ReadAll()
	return
}

func (s *sModule) WebhookPostPutPatchDelete(ctx context.Context, header, method, url string, payload any) (
	statusCode int, contentType string, body []byte, err error) {
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
	statusCode = resp.StatusCode
	contentType = resp.Header.Get("Content-Type")
	body = resp.ReadAll()
	return
}
