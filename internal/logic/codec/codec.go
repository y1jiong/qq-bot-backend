package codec

import "qq-bot-backend/internal/service"

type sCodec struct{}

func New() *sCodec {
	return &sCodec{}
}

func init() {
	service.RegisterCodec(New())
}
