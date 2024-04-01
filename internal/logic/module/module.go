package module

import (
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"
)

type sModule struct{}

func New() *sModule {
	return &sModule{}
}

func init() {
	service.RegisterModule(New())
}

var (
	webhookPrefixRe = regexp.MustCompile(`^webhook:?([Gg][Ee][Tt]|[Pp][Oo][Ss][Tt])?://(.+)$`)
)

func (s *sModule) MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string) {
	for k, v := range m {
		fields := strings.Fields(k)
		fieldsLen, count := len(fields), 0
		for _, kk := range fields {
			if strings.Contains(str, kk) {
				count++
			}
		}
		if count == fieldsLen {
			contains = true
			hit = k
			if vv, ok := v.(string); ok {
				mValue = vv
			}
			return
		}
	}
	return
}
