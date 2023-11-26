package list

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
)

func isListOpLegal(ctx context.Context, A, B, C string) (legal bool) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(A) || !legalListNameRe.MatchString(B) || !legalListNameRe.MatchString(C) {
		return
	}
	// 局部变量
	userId := service.Bot().GetUserId(ctx)
	// 权限校验
	AE := getList(ctx, A)
	BE := getList(ctx, B)
	CE := getList(ctx, C)
	if AE == nil || BE == nil || CE == nil {
		return
	}
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, AE.Namespace, userId) ||
		!service.Namespace().IsNamespaceOwnerOrAdmin(ctx, BE.Namespace, userId) ||
		!service.Namespace().IsNamespaceOwnerOrAdmin(ctx, CE.Namespace, userId) {
		return
	}
	legal = true
	return
}

func (s *sList) UnionListReturnRes(ctx context.Context, A, B, C string) (retMsg string) {
	// 校验
	if !isListOpLegal(ctx, A, B, C) {
		return
	}
	// 数据处理 并集运算
	n, err := s.UnionOp(ctx, A, B, C)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = fmt.Sprintf("list %s union %s equals %s, %d items", A, B, C, n)
	return
}

func (s *sList) IntersectListReturnRes(ctx context.Context, A, B, C string) (retMsg string) {
	// 校验
	if !isListOpLegal(ctx, A, B, C) {
		return
	}
	// 数据处理 交集运算
	n, err := s.IntersectOp(ctx, A, B, C)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = fmt.Sprintf("list %s intersect %s equals %s, %d items", A, B, C, n)
	return
}

func (s *sList) DifferenceListReturnRes(ctx context.Context, A, B, C string) (retMsg string) {
	// 校验
	if !isListOpLegal(ctx, A, B, C) {
		return
	}
	// 数据处理 差集运算
	n, err := s.DifferenceOp(ctx, A, B, C)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = fmt.Sprintf("list %s difference %s equals %s, %d items", A, B, C, n)
	return
}
