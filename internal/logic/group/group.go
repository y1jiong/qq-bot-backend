package group

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
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

func getGroup(ctx context.Context, groupId int64) (groupE *entity.Group) {
	// 数据库查询
	err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Scan(&groupE)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return
}

func (s *sGroup) IsBinding(ctx context.Context, groupId int64) bool {
	groupE := getGroup(ctx, groupId)
	return groupE != nil && groupE.Namespace != ""
}

func (s *sGroup) GetNamespace(ctx context.Context, groupId int64) string {
	if groupE := getGroup(ctx, groupId); groupE != nil {
		return groupE.Namespace
	}
	return ""
}
