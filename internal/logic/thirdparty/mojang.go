package thirdparty

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"io"
	"net/http"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility/codec"
	"regexp"
	"time"
)

var (
	// 匹配 Minecraft 玩家名称
	legalMinecraftNameRe = regexp.MustCompile(`(?:^|[^\w#])(\w{3,16})$`)
)

func (s *sThirdParty) QueryMinecraftGenuineUser(ctx context.Context, name string) (genuine bool, realName, uuid string, err error) {
	// 全字匹配
	if !legalMinecraftNameRe.MatchString(name) {
		return
	}
	name = legalMinecraftNameRe.FindStringSubmatch(name)[1]

	url := "https://api.mojang.com/users/profiles/minecraft/" + name

	ctx, span := gtrace.NewSpan(ctx, codec.GetRouteURL(url))
	defer span.End()
	span.SetAttributes(attribute.String("http.url", url))
	span.SetAttributes(attribute.String("minecraft.name", name))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// GET 请求 mojang api
	var resp *http.Response
	for range 3 {
		resp, err = http.DefaultClient.Get(url)
		if err == nil {
			break
		}
		time.Sleep(service.Cfg().GetRetryIntervalSeconds(ctx) * time.Second)
	}
	if err != nil || resp == nil {
		return
	}
	defer resp.Body.Close()

	// 判断是否正版
	if resp.StatusCode != http.StatusOK {
		return
	}

	// 解析响应
	rawJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	resJson, err := sonic.Get(rawJson)
	if err != nil {
		return
	}

	// 导出正确的 name 和 uuid
	realName, _ = resJson.Get("name").StrictString()
	uuid, _ = resJson.Get("id").StrictString()

	// 判断正版
	if uuid != "" {
		genuine = true
	}

	return
}
