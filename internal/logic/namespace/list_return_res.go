package namespace

import (
	"context"
	"qq-bot-backend/internal/service"
)

func (s *sNamespace) LoadNamespaceListReturnRes(ctx context.Context, namespace, listName string) (retMsg string) {
	// 权限校验 owner or admin or namespace op
	if !s.IsNamespaceOwnerOrAdminOrOperator(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 获取 listName 对应的 namespace
	if listNamespace := service.List().GetListNamespace(ctx, listName); listNamespace == "" ||
		!s.IsNamespacePropertyPublic(ctx, listNamespace) {
		return
	}
	// 加载 listName 到 namespace
	s.AddNamespaceList(ctx, namespace, listName)
	return "已加载 list(" + listName + ") 到 namespace(" + namespace + ") 中"
}

func (s *sNamespace) UnloadNamespaceListReturnRes(ctx context.Context, namespace, listName string) (retMsg string) {
	// 权限校验 owner or admin or namespace op
	if !s.IsNamespaceOwnerOrAdminOrOperator(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 获取 listName 对应的 namespace
	if listNamespace := service.List().GetListNamespace(ctx, listName); listNamespace == "" ||
		listNamespace == namespace {
		return
	}
	// 从 namespace 卸载 listName
	s.RemoveNamespaceList(ctx, namespace, listName)
	return "已从 namespace(" + namespace + ") 中卸载 list(" + listName + ")"
}
