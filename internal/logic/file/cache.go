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
	fileCachePrefix = "file_cache_"
)

func getFileCacheKey(id string) string {
	return fileCachePrefix + id
}

func (s *sFile) GetCachedFileById(ctx context.Context, id string) (content []byte, err error) {
	v, err := gcache.Get(ctx, getFileCacheKey(id))
	if err != nil {
		return
	}
	if v == nil {
		err = gerror.NewCode(gcode.New(http.StatusNotFound, "", nil), "file not found")
		return
	}
	content = v.Bytes()
	return
}

func (s *sFile) getCachedFileId(ctx context.Context, content []byte, duration time.Duration) (id string, err error) {
	id = guid.S()
	err = gcache.Set(ctx, getFileCacheKey(id), content, duration)
	return
}

func (s *sFile) getCachedFileUrl(ctx context.Context, id string) (url string, err error) {
	r := g.RequestFromCtx(ctx)
	url = r.GetSchema() + "://" + r.Host + "/file/" + id
	return
}

func (s *sFile) CacheFile(ctx context.Context, content []byte, duration time.Duration) (url string, err error) {
	id, err := s.getCachedFileId(ctx, content, duration)
	if err != nil {
		return
	}
	url, err = s.getCachedFileUrl(ctx, id)
	return
}
