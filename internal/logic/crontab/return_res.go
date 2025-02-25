package crontab

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/do"
	"qq-bot-backend/internal/model/entity"
	"strings"
)

func (s *sCrontab) GlanceReturnRes(ctx context.Context, creatorId int64) (retMsg string) {
	var tasks []entity.Crontab
	q := dao.Crontab.Ctx(ctx).
		Fields(
			dao.Crontab.Columns().Name,
			dao.Crontab.Columns().CreatorId,
		)
	if creatorId != 0 {
		q = q.Where(dao.Crontab.Columns().CreatorId, creatorId)
	}
	err := q.Scan(&tasks)
	if err != nil {
		return
	}

	builder := strings.Builder{}
	for _, task := range tasks {
		if creatorId != 0 {
			builder.WriteString("`" + task.Name + "`\n")
		} else {
			builder.WriteString("`" + task.Name + "` // " + gconv.String(task.CreatorId) + "\n")
		}
	}

	retMsg = builder.String()
	return
}

func (s *sCrontab) QueryReturnRes(ctx context.Context, name string, creatorId int64) (retMsg string) {
	var task *entity.Crontab
	q := dao.Crontab.Ctx(ctx).
		Fields(
			dao.Crontab.Columns().Name,
			dao.Crontab.Columns().Expression,
			dao.Crontab.Columns().Request,
			dao.Crontab.Columns().CreatorId,
			dao.Crontab.Columns().CreatedAt,
		).
		Where(dao.Crontab.Columns().Name, name)
	if creatorId != 0 {
		q = q.Where(dao.Crontab.Columns().CreatorId, creatorId)
	}
	err := q.Scan(&task)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}

	retMsg = dao.Crontab.Columns().Name + ": " + task.Name + "\n" +
		dao.Crontab.Columns().Expression + ": " + task.Expression + "\n" +
		dao.Crontab.Columns().CreatorId + ": " + gconv.String(task.CreatorId) + "\n" +
		dao.Crontab.Columns().Request + ": " + task.Request + "\n" +
		dao.Crontab.Columns().CreatedAt + ": " + task.CreatedAt.String()
	return
}

func (s *sCrontab) AddReturnRes(ctx context.Context,
	name, expr string,
	creatorId, botId int64,
	reqJSON []byte,
) (retMsg string) {
	if expr[:strings.Index(expr, " ")] == "*" {
		retMsg = "不允许设置为每分钟执行"
		return
	}

	if err := s.add(ctx, name, expr, botId, reqJSON); err != nil {
		retMsg = "加入定时任务失败"
		g.Log().Error(ctx, err)
		return
	}

	_, err := dao.Crontab.Ctx(ctx).
		Data(entity.Crontab{
			Name:       name,
			Expression: expr,
			CreatorId:  creatorId,
			BotId:      botId,
			Request:    string(reqJSON),
		}).
		OmitEmptyData().
		Insert()
	if err != nil {
		retMsg = "持久化失败"
		g.Log().Error(ctx, err)
		return
	}

	retMsg = name + "\n" + expr
	return
}

func (s *sCrontab) RemoveReturnRes(ctx context.Context, name string, creatorId int64) (retMsg string) {
	var task *entity.Crontab
	err := dao.Crontab.Ctx(ctx).
		Fields(dao.Crontab.Columns().Name).
		Where(do.Crontab{
			Name:      name,
			CreatorId: creatorId,
		}).
		Scan(&task)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}

	if task == nil {
		return
	}

	_, err = dao.Crontab.Ctx(ctx).
		Where(dao.Crontab.Columns().Name, name).
		Delete()
	if err != nil {
		retMsg = "删除失败"
		g.Log().Error(ctx, err)
		return
	}

	s.remove(name)

	retMsg = name
	return
}
