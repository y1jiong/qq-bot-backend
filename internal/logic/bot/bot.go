package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/net/ghttp"
	"he3-bot/internal/service"
)

type sBot struct{}

func init() {
	service.RegisterBot(New())
}

func New() *sBot {
	return &sBot{}
}

func (s *sBot) Parse(ctx context.Context, ws *ghttp.WebSocket, msg []byte) {
	sj.New()
	req, err := sj.NewJson(msg)
	if err != nil {
		return
	}
	switch req.Get("post_type").MustString() {
	case "message":
		s.processMessage(ctx, ws, req)
	case "request":
		s.processRequest(ctx, ws, req)
	case "notice":
		s.processNotice(ctx, ws, req)
	}
}
