package utility

import (
	"context"
	"fmt"
	"net/http"
	"qq-bot-backend/internal/consts"
	"runtime"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/guid"
	"golang.org/x/sync/semaphore"
	"golang.org/x/sync/singleflight"
)

var (
	sem   *semaphore.Weighted
	group singleflight.Group

	initSem = sync.OnceFunc(func() {
		// A maximum of runtime.NumCPU() concurrent requests can be made
		sem = semaphore.NewWeighted(int64(runtime.NumCPU()))
	})
)

func SendWebhookRequest(ctx context.Context, header, method, url string, payload ...any,
) (statusCode int, contentType string, body []byte, err error) {
	initSem()

	type data struct {
		statusCode  int
		contentType string
		body        []byte
	}

	key := url
	if method != http.MethodGet && method != http.MethodHead {
		key = guid.S()
	}

	v, err, _ := group.Do(key, func() (any, error) {
		const cost = 1
		if err := sem.Acquire(ctx, cost); err != nil {
			return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
		}
		defer sem.Release(cost)

		c := gclient.New().SetTimeout(30 * time.Second).SetAgent(consts.ProjName + "/" + consts.Version)
		if header != "" {
			c.SetHeaderRaw(header)
		}

		resp, err := c.DoRequest(ctx, method, url, payload...)
		if err != nil || resp == nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Close()

		return &data{
			statusCode:  resp.StatusCode,
			contentType: resp.Header.Get("Content-Type"),
			body:        resp.ReadAll(),
		}, nil
	})
	if err != nil {
		return
	}

	dat, ok := v.(*data)
	if !ok {
		err = fmt.Errorf("invalid response type")
		return
	}

	return dat.statusCode, dat.contentType, dat.body, nil
}
