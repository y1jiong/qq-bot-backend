package bot

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"qq-bot-backend/utility/codec"
	"qq-bot-backend/utility/segment"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

func (s *sBot) Forward(ctx context.Context, url, key string) (err error) {
	ctx, span := gtrace.NewSpan(ctx, codec.GetAbsoluteURL(url))
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if len(body) == 0 {
		return
	}

	const kind = "forwarding"
	if limited, _ := utility.AutoLimit(ctx, kind, url, consts.MaxSendMessageCount, time.Minute); limited {
		return
	}

	if msg, ok := medium(ctx, resp.Header.Get("Content-Type"), body); ok {
		// 如果是图片、音频、视频，直接发送
		s.SendMsg(ctx, msg)
		return
	}
	if len(body) > consts.MaxMessageLength*3 &&
		utf8.RuneCount(segment.FilterCQCode(body)) > consts.MaxMessageLength {
		s.SendForwardMsg(ctx, string(body))
		return
	}
	s.SendMsg(ctx, string(body))

	return
}

func medium(ctx context.Context, contentType string, body []byte) (string, bool) {
	// 如果是图片
	if strings.HasPrefix(contentType, "image/") {
		mediumURL, err := service.File().CacheFile(ctx, body, 5*time.Minute)
		if err != nil {
			return "Image cache failed", true
		}
		return "[CQ:image,file=" + codec.EncodeCQCode(mediumURL) + "]", true
	}
	// 如果是音频
	if strings.HasPrefix(contentType, "audio/") {
		mediumURL, err := service.File().CacheFile(ctx, body, 5*time.Minute)
		if err != nil {
			return "Audio cache failed", true
		}
		return "[CQ:record,file=" + codec.EncodeCQCode(mediumURL) + "]", true
	}
	// 如果是视频
	if strings.HasPrefix(contentType, "video/") {
		mediumURL, err := service.File().CacheFile(ctx, body, 5*time.Minute)
		if err != nil {
			return "Video cache failed", true
		}
		return "[CQ:video,file=" + codec.EncodeCQCode(mediumURL) + "]", true
	}
	return "", false
}
