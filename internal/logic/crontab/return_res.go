package crontab

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"strings"
)

func (s *sCrontab) GlanceReturnRes(ctx context.Context) (retMsg string) {
	tasks, err := s.getTasks(ctx)
	if err != nil {
		return
	}

	builder := strings.Builder{}
	for _, task := range tasks {
		builder.WriteString("`" + task.Name + "`\n")
	}

	retMsg = builder.String()
	return
}

func (s *sCrontab) QueryReturnRes(ctx context.Context, name string) (retMsg string) {
	var task *entity.Crontab
	err := dao.Crontab.Ctx(ctx).
		Where(dao.Crontab.Columns().Name, name).
		Scan(&task)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}

	retMsg = task.Name + "\n" + task.Expression + "\n" + task.Request + "\n" + task.CreatedAt.String()
	return
}

func (s *sCrontab) AddReturnRes(ctx context.Context, name, expr string, botId int64, reqJSON []byte,
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

func (s *sCrontab) RemoveReturnRes(ctx context.Context, name string) (retMsg string) {
	_, err := dao.Crontab.Ctx(ctx).
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
