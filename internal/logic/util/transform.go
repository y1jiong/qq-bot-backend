package util

import (
	"context"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility/segment"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

func (s *sUtil) ToPlainText(ctx context.Context, message string) string {
	segments := segment.ParseMessage(message)

	for idx, seg := range segments {
		switch seg.Type {
		case segment.TypeAt, segment.TypeReply:
			segments[idx] = nil

		case segment.TypeImage:
			if url := service.Cfg().GetOcrURL(ctx); url != "" {
				image, err := s.httpGetQQImage(ctx, seg.Data["url"])
				if err != nil {
					g.Log().Warning(ctx, err)
					continue
				}

				text, err := s.OCR(ctx, image)
				if err != nil {
					g.Log().Warning(ctx, err)
					continue
				}

				segments[idx] = segment.NewTextSegments(text).First()
			}
		}
	}

	return strings.TrimSpace(segments.String())
}
