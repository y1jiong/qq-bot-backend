package namespace

import (
	"context"
	"github.com/gogf/gf/v2/util/grand"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/do"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sNamespace) Broadcast(ctx context.Context, namespace, message string, originGroupId int64) (err error) {
	var groups []entity.Group
	err = dao.Group.Ctx(ctx).
		Fields(dao.Group.Columns().GroupId).
		Where(do.Group{
			Namespace:       namespace,
			AcceptBroadcast: true,
		}).
		Scan(&groups)
	if err != nil {
		return
	}

	for _, group := range groups {
		if group.GroupId == originGroupId {
			continue
		}

		_, _ = service.Bot().SendMessage(ctx, "", 0, group.GroupId, message, false)

		time.Sleep(time.Duration(grand.N(1000, 10000)) * time.Millisecond)
	}

	return
}
