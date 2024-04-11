package file

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/guid"
	"net/http"
	"time"
)

const (
	fileCachePrefix = "file:cache:"
)

func (s *sFile) GetCachedFileById(ctx context.Context, id string) (content string, err error) {
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

func (s *sFile) SetCacheFile(ctx context.Context, content string, duration time.Duration) (id string, err error) {
	id = guid.S()
	err = gcache.Set(ctx, fileCachePrefix+id, content, duration)
	return
}

func (s *sFile) GetCachedFileUrl(ctx context.Context, id string) (url string, err error) {
	r := g.RequestFromCtx(ctx)
	url = r.GetSchema() + "://" + r.Host + "/file/" + id
	return
}
