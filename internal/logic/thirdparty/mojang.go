package thirdparty

import (
	"context"
	"io"
	"net/http"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/utility"
	"qq-bot-backend/utility/codec"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *sThirdParty) QueryMinecraftGenuineUser(ctx context.Context, name string,
) (genuine bool, realName, uuid string, err error) {
	url := "https://api.mojang.com/users/profiles/minecraft/" + name

	ctx, span := gtrace.NewSpan(ctx, codec.GetAbsoluteURL(url))
	defer span.End()
	span.SetAttributes(attribute.String("http.url", url))
	span.SetAttributes(attribute.String("minecraft.name", name))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// GET 请求 mojang api
	var resp *gclient.Response
	_ = utility.RetryWithBackoff(ctx, func() bool {
		resp, err = gclient.New().SetTimeout(30*time.Second).SetAgent(consts.ProjName+"/"+consts.Version).Get(ctx, url)
		if err != nil {
			span.RecordError(err)
			return false
		}
		return true
	}, 3, utility.ExponentialBackoffWithJitter(ctx))
	if err != nil || resp == nil {
		return
	}
	defer resp.Close()

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
