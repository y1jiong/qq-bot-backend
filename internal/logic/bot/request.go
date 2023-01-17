package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gorilla/websocket"
	"regexp"
)

var (
	// 匹配 Minecraft 玩家名称
	genuineNameRegexp = regexp.MustCompile(`[A-Za-z\d_]{3,}`)
)

func (s *sBot) processRequest(ctx context.Context, ws *ghttp.WebSocket, req *sj.Json) {
	// 加群申请
	if req.Get("request_type").MustString() == "group" &&
		req.Get("sub_type").MustString() == "add" {
		name := genuineNameRegexp.FindString(req.Get("comment").MustString())
		var genuine bool
		var uuid string
		if name != "" {
			// TODO: 第二个返回值 uuid 记得处理
			var err error
			genuine, name, uuid, err = s.queryMinecraftGenuineUser(ctx, name)
			if err != nil {
				glog.Noticef(ctx, "an error occurred while queryMinecraftGenuineUser")
			}
		}
		// 执行最后的审批
		s.approve(ctx, ws, req, genuine)
		// 打印通过的日志
		if genuine {
			glog.Infof(ctx, "approve user(%v) join group(%v) with %v(%v) in %v",
				req.Get("user_id").MustInt64(),
				req.Get("group_id").MustInt64(),
				name, uuid,
				req.Get("comment").MustString())
		}
	}
}

func (s *sBot) approve(ctx context.Context, ws *ghttp.WebSocket, req *sj.Json, agree bool) {
	reqJson := sj.New()
	reqJson.Set("action", "set_group_add_request")
	params := make(map[string]any)
	params["flag"] = req.Get("flag").MustString()
	params["sub_type"] = req.Get("sub_type").MustString()
	params["approve"] = agree
	// 当不予通过时，给出理由
	if !agree {
		params["reason"] = "invalid name"
	}
	reqJson.Set("params", params)
	res, err := reqJson.Encode()
	if err != nil {
		glog.Warningf(ctx, "an error occurred while reqJson.Encode, %v", err)
		return
	}
	err = ws.WriteMessage(websocket.TextMessage, res)
	if err != nil {
		glog.Warningf(ctx, "an error occurred while ws.WriteMessage, %v", err)
	}
}
