package thirdparty

import "qq-bot-backend/internal/service"

type sThirdParty struct{}

func New() *sThirdParty {
	return &sThirdParty{}
}

func init() {
	service.RegisterThirdParty(New())
}
