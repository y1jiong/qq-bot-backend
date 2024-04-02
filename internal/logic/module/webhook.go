package module

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"qq-bot-backend/internal/consts"
)

func (s *sModule) WebhookGetHeadConnectOptionsTrace(ctx context.Context, method, url string) (body string, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	resp, err := c.DoRequest(ctx, method, url)
	if err != nil || resp == nil {
		return
	}
	body = resp.ReadAllString()
	return
}

func (s *sModule) WebhookPostPutPatchDelete(ctx context.Context, method, url string, payload any) (body string, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	resp, err := c.DoRequest(ctx, method, url, payload)
	if err != nil || resp == nil {
		return
	}
	body = resp.ReadAllString()
	return
}
