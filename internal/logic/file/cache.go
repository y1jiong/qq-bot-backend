package file

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/guid"
	"net/http"
	"qq-bot-backend/internal/service"
	"time"
)

const (
	fileCachePrefix = "file:cache:"
)

func (s *sFile) GetCachedFileFromId(ctx context.Context, id string) (content string, err error) {
	v, err := gcache.Get(ctx, fileCachePrefix+id)
	if err != nil {
		return
	}
	if v == nil {
		err = gerror.NewCode(gcode.New(http.StatusNotFound, "", nil), "file not found")
		return
	}
	content = v.String()
	return
}

func (s *sFile) SetCachedFile(ctx context.Context, content string, duration time.Duration) (id string, err error) {
	id = guid.S()
	err = gcache.Set(ctx, fileCachePrefix+id, content, duration)
	return
}

func (s *sFile) GetCachedFileUrl(ctx context.Context, id string) (url string, err error) {
	u, err := service.Cfg().GetUrlPrefix(ctx)
	if err != nil {
		return
	}
	url = u + "/file/" + id
	return
}
