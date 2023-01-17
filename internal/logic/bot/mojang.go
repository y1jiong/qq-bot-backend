package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"he3-bot/internal/service"
	"io"
	"net/http"
	"time"
)

func (s *sBot) queryMinecraftGenuineUser(ctx context.Context, id string) (genuine bool, realName, uuid string, err error) {
	get := func() (*gclient.Response, error) {
		return g.Client().Get(ctx, "https://api.mojang.com/users/profiles/minecraft/"+id)
	}
	res, err := get()
	// 失败重试
	if err != nil {
		for i := 0; i < 2; i++ {
			time.Sleep(service.Cfg().GetRetryIntervalMilliseconds(ctx) * time.Millisecond)
			res, err = get()
			if err == nil {
				break
			}
		}
		if err != nil {
			return
		}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return
	}
	rawJson, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	resJson, err := sj.NewJson(rawJson)
	if err != nil {
		return
	}
	// 导出正确的 name 和 uuid
	realName = resJson.Get("name").MustString()
	uuid = resJson.Get("id").MustString()
	if uuid != "" {
		genuine = true
	}
	return
}
