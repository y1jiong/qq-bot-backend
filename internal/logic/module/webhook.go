package module

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"qq-bot-backend/internal/consts"
)

func (s *sModule) WebhookGet(ctx context.Context, url string) (body string, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	resp, err := c.Get(ctx, url)
	if err != nil || resp == nil {
		return
	}
	body = resp.ReadAllString()
	return
}

func (s *sModule) WebhookPost(ctx context.Context, url string, payload any) (body string, err error) {
	c := gclient.New()
	c.SetAgent(consts.ProjName + "/" + consts.Version)
	resp, err := c.Post(ctx, url, payload)
	if err != nil || resp == nil {
		return
	}
	body = resp.ReadAllString()
	return
}
