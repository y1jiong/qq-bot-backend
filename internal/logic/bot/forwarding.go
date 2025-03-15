package bot

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/utility/codec"
	"qq-bot-backend/utility/segment"
	"sync"
	"unicode/utf8"
)

var (
	forwarding *http.Client
)

var initForwarding = sync.OnceFunc(func() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	// No validation for https certification of the server in default.
	t.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	forwarding = &http.Client{
		Transport: t,
	}
})

func (s *sBot) Forward(ctx context.Context, url, key string) (err error) {
	ctx, span := gtrace.NewSpan(ctx, codec.GetRouteURL(url))
	defer span.End()
	span.SetAttributes(attribute.String("http.url", url))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	payload, err := s.reqNodeFromCtx(ctx).MarshalJSON()
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return
	}

	// Inject trace
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	req.Header.Set("User-Agent", consts.ProjName+"/"+consts.Version)
	if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	req.Header.Set("Content-Type", "application/json")

	initForwarding()
	resp, err := forwarding.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if len(body) != 0 {
		if len(body) > consts.MaxMessageLength*3 &&
			utf8.RuneCount(segment.FilterCQCode(body)) > consts.MaxMessageLength {
			s.SendForwardMsg(ctx, string(body))
		} else {
			s.SendMsg(ctx, string(body))
		}
	}

	return resp.Body.Close()
}
