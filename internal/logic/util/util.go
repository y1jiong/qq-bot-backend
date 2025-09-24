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

func (s *sUtil) FindBestKeywordMatch(ctx context.Context, msg string, lists map[string]any,
) (found bool, hit, value string) {
	type result struct {
		hit   string
		value string
	}

	var (
		bestPrefix    *result
		bestNonPrefix *result
	)

	for k := range lists {
		eureka, hitStr, valueStr := utility.MatchAllKeywords(msg, service.List().GetListData(ctx, k))
		if !eureka {
			continue
		}

		if strings.HasPrefix(msg, hitStr) {
			if bestPrefix == nil || len(hitStr) > len(bestPrefix.hit) {
				bestPrefix = &result{hitStr, valueStr}
			}
		} else {
			if bestNonPrefix == nil || len(hitStr) > len(bestNonPrefix.hit) {
				bestNonPrefix = &result{hitStr, valueStr}
			}
		}
	}

	if bestPrefix != nil {
		found = true
		hit = bestPrefix.hit
		value = bestPrefix.value
		return
	}
	if bestNonPrefix != nil {
		found = true
		hit = bestNonPrefix.hit
		value = bestNonPrefix.value
	}
	return
}
