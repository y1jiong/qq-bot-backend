package util

import (
	"context"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"strings"
)

type sUtil struct{}

func New() *sUtil {
	return &sUtil{}
}

func init() {
	service.RegisterUtil(New())
}

func (s *sUtil) IsOnKeywordLists(ctx context.Context, msg string, lists map[string]any) (in bool, hit, value string) {
	for k := range lists {
		listMap := service.List().GetListData(ctx, k)
		if contains, hitStr, valueStr := s.MultiContains(msg, listMap); contains {
			in = true
			hit = hitStr
			value = valueStr
			if strings.HasPrefix(msg, hit) {
				return
			}
		}
	}
	return
}

func (s *sUtil) MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string) {
	arr := utility.ReverseSortedArrayFromMapKey(m)
	for _, k := range arr {
		fields := strings.Fields(k)
		allContains := true
		for _, kk := range fields {
			if !strings.Contains(str, kk) {
				allContains = false
				break
			}
		}
		if !allContains {
			continue
		}
		contains = true
		hit = k
		if vv, ok := m[k].(string); ok {
			mValue = vv
		}
		if strings.HasPrefix(str, k) {
			return
		}
	}
	return
}
