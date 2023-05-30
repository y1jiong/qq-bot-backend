package file

import (
	"qq-bot-backend/internal/service"
)

type sFile struct{}

func New() *sFile {
	return &sFile{}
}

func init() {
	service.RegisterFile(New())
}
