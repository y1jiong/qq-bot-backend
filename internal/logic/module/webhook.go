package module

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"qq-bot-backend/internal/consts"
)

func (s *sModule) WebhookGetHeadConnectOptionsTrace(ctx context.Context, method, url string) (statusCode int, body string, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	resp, err := c.DoRequest(ctx, method, url)
	if err != nil || resp == nil {
		return
	}
	defer resp.Close()
	statusCode = resp.StatusCode
	body = resp.ReadAllString()
	return
}

func (s *sModule) WebhookPostPutPatchDelete(ctx context.Context, method, url string, payload any) (statusCode int, body string, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	resp, err := c.DoRequest(ctx, method, url, payload)
	if err != nil || resp == nil {
		return
	}
	defer resp.Close()
	statusCode = resp.StatusCode
	body = resp.ReadAllString()
	return
}
