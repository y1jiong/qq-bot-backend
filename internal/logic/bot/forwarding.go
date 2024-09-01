package bot

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"qq-bot-backend/internal/consts"
	"sync"
)

var (
	forwardClient *http.Client
)

var initForwardClient = sync.OnceFunc(func() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	// No validation for https certification of the server in default.
	t.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	forwardClient = &http.Client{
		Transport: t,
	}
})

func (s *sBot) Forward(ctx context.Context, url, authorization string) error {
	payload, err := s.reqJsonFromCtx(ctx).MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", consts.ProjName+"/"+consts.Version)
	if authorization != "" {
		req.Header.Set("Authorization", "Bearer "+authorization)
	}

	initForwardClient()
	resp, err := forwardClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		s.SendMsg(ctx, string(body))
	}

	return resp.Body.Close()
}
