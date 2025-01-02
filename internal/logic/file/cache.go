package file

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/guid"
	"qq-bot-backend/internal/consts/errcode"
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
		err = gerror.NewCode(errcode.FileNotFound)
		return
	}
	return v.Bytes(), nil
}

func (s *sFile) getCachedFileId(ctx context.Context, content []byte, duration time.Duration) (id string, err error) {
	id = guid.S()
	return id, gcache.Set(ctx, getFileCacheKey(id), content, duration)
}

func (s *sFile) getCachedFileURL(ctx context.Context, id string) (url string, err error) {
	r := g.RequestFromCtx(ctx)
	return r.GetSchema() + "://" + r.Host + "/file/" + id, nil
}

func (s *sFile) CacheFile(ctx context.Context, content []byte, duration time.Duration) (url string, err error) {
	id, err := s.getCachedFileId(ctx, content, duration)
	if err != nil {
		return
	}
	return s.getCachedFileURL(ctx, id)
}
