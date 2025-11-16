package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility/codec"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	ocrReq struct {
		Base64  string    `json:"base64"`
		Options ocrOption `json:"options,omitempty"`
	}
	ocrResp struct {
		Code      int     `json:"code"`
		Data      string  `json:"data"`
		Time      float64 `json:"time"`
		Timestamp float64 `json:"timestamp"`
	}
)

// https://github.com/hiroi-sora/Umi-OCR/blob/main/docs/http/api_ocr.md#/api/ocr/get_options
type ocrOption struct {
	OcrLimitSideLen int    `json:"ocr.limit_side_len,omitempty"`
	DataFormat      string `json:"data.format,omitempty"`
}

func (s *sUtil) OCR(ctx context.Context, image []byte) (string, error) {
	url := service.Cfg().GetOcrURL(ctx)
	if url == "" {
		return "", errors.New("empty ocr url")
	}

	ctx, span := gtrace.NewSpan(ctx, "util.OCR")
	defer span.End()

	req := &ocrReq{
		Base64: base64.StdEncoding.EncodeToString(image),
		Options: ocrOption{
			OcrLimitSideLen: 4320,
			DataFormat:      "text",
		},
	}

	resp, err := s.ocr(ctx, url, req)
	if err != nil {
		return "", err
	}
	if resp.Code != 100 {
		return "", nil
	}

	span.AddEvent("ocr.result", trace.WithAttributes(
		attribute.Int("ocr.result.code", resp.Code),
		attribute.String("ocr.result.data", resp.Data),
	))

	return resp.Data, nil
}

func (s *sUtil) ocr(ctx context.Context, url string, req *ocrReq) (respObj *ocrResp, err error) {
	ctx, span := gtrace.NewSpan(ctx, codec.GetAbsoluteURL(url))
	defer span.End()
	span.SetAttributes(attribute.String("http.url", url))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	payload, err := sonic.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", consts.ProjName+"/"+consts.Version)
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, errors.New(resp.Status)
	}

	respJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = sonic.Unmarshal(respJson, &respObj); err != nil {
		return nil, err
	}

	return respObj, nil
}

func (s *sUtil) httpGetQQImage(ctx context.Context, url string) (respBody []byte, err error) {
	ctx, span := gtrace.NewSpan(ctx, codec.GetAbsoluteURL(url))
	defer span.End()
	span.SetAttributes(attribute.String("http.url", url))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, errors.New("image: " + resp.Status)
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
