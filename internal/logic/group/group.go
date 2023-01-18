package group

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
)

type sGroup struct{}

func New() *sGroup {
	return &sGroup{}
}

func init() {
	service.RegisterGroup(New())
}

func (s *sGroup) BindNamespace(ctx context.Context, groupId int64, namespace string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	if !service.Namespace().IsNamespaceExist(ctx, namespace) {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) ||
		!service.Namespace().IsNamespaceOwnerOrAdmin(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据库计数
	n, err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if n > 0 {
		// 数据库更新
		_, err = dao.Group.Ctx(ctx).
			Where(dao.Group.Columns().GroupId, groupId).
			Data(dao.Group.Columns().Namespace, namespace).
			Update()
	} else {
		// 数据库插入
		gEntity := entity.Group{
			GroupId:   groupId,
			Namespace: namespace,
		}
		_, err = dao.Group.Ctx(ctx).
			Data(gEntity).
			OmitEmpty().
			Insert()
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已绑定当前 group("+gconv.String(groupId)+") 到 namespace("+namespace+")")
}

func (s *sGroup) Unbind(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !s.IsGroupBindNamespaceOwnerOrAdmin(ctx, groupId, service.Bot().GetUserId(ctx)) {
		return
	}
	// 过程
	n, err := dao.Group.Ctx(ctx).Where(dao.Group.Columns().GroupId, groupId).Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if n < 1 {
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().Namespace, "").
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已解除 group("+gconv.String(groupId)+") 的 namespace 绑定")
}

func (s *sGroup) BindQuery(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 过程
	var gEntity *entity.Group
	err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Scan(&gEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if gEntity == nil || gEntity.Namespace == "" {
		service.Bot().SendMsg(ctx, "group("+gconv.String(groupId)+") 没有绑定任何 namespace")
		return
	}
	service.Bot().SendMsg(ctx, "group("+gconv.String(groupId)+") 绑定了 namespace("+gEntity.Namespace+")")
}

func (s *sGroup) IsGroupBindNamespaceOwnerOrAdmin(ctx context.Context, groupId, userId int64) (yes bool) {
	// 参数合法性校验
	if groupId < 1 || userId < 1 {
		return
	}
	// 获取 group 绑定的 namespace
	var gEntity *entity.Group
	err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Scan(&gEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if gEntity == nil {
		return
	}
	// 过程
	return service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, userId)
}
