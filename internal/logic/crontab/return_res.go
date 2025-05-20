package crontab

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"strings"
)

func (s *sCrontab) ShowReturnRes(ctx context.Context, creatorId int64) (retMsg string) {
	var (
		tasks []entity.Crontab
	)

	q := dao.Crontab.Ctx(ctx).
		Fields(
			dao.Crontab.Columns().Name,
			dao.Crontab.Columns().Expression,
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
		builder.WriteString(task.Expression + " " + task.Name + "\n")
	}

	retMsg = strings.TrimSuffix(builder.String(), "\n")
	return
}

func (s *sCrontab) QueryReturnRes(ctx context.Context, name string, creatorId int64) (retMsg string) {
	var (
		task *entity.Crontab
	)

	q := dao.Crontab.Ctx(ctx).
		Fields(
			dao.Crontab.Columns().Name,
			dao.Crontab.Columns().Expression,
			dao.Crontab.Columns().CreatorId,
			dao.Crontab.Columns().BotId,
			dao.Crontab.Columns().Request,
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
		dao.Crontab.Columns().BotId + ": " + gconv.String(task.BotId) + "\n" +
		dao.Crontab.Columns().Request + ": " + task.Request + "\n" +
		dao.Crontab.Columns().CreatedAt + ": " + task.CreatedAt.String()
	return
}

func (s *sCrontab) AddReturnRes(ctx context.Context,
	name, expr string,
	creatorId, botId int64,
	reqJSON []byte,
) (retMsg string) {
	if strings.HasPrefix(expr, "* ") {
		retMsg = "不允许设置为每分钟执行"
		return
	}

	if err := s.add(ctx, name, expr, botId, reqJSON); err != nil {
		g.Log().Error(ctx, err)
		retMsg = "加入定时任务失败"
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
		g.Log().Error(ctx, err)
		s.remove(name)
		retMsg = "持久化失败"
		return
	}

	retMsg = expr + " " + name
	return
}

func (s *sCrontab) RemoveReturnRes(ctx context.Context, name string, creatorId int64) (retMsg string) {
	q := dao.Crontab.Ctx(ctx).
		Where(dao.Crontab.Columns().Name, name)
	if creatorId != 0 {
		q = q.Where(dao.Crontab.Columns().CreatorId, creatorId)
	}
	result, err := q.Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "删除失败"
		return
	}

	if aff, _ := result.RowsAffected(); aff == 0 {
		return
	}

	s.remove(name)

	retMsg = name
	return
}

func (s *sCrontab) ChangeExpressionReturnRes(ctx context.Context,
	expr, name string, creatorId int64,
) (retMsg string) {
	if strings.HasPrefix(expr, "* ") {
		retMsg = "不允许设置为每分钟执行"
		return
	}

	q := dao.Crontab.Ctx(ctx).
		Data(dao.Crontab.Columns().Expression, expr).
		Where(dao.Crontab.Columns().Name, name)
	if creatorId != 0 {
		q = q.Where(dao.Crontab.Columns().CreatorId, creatorId)
	}
	result, err := q.Update()
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "修改失败"
		return
	}

	if aff, _ := result.RowsAffected(); aff == 0 {
		return
	}

	task, err := s.getTask(ctx, name, creatorId)
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "获取失败"
		return
	}
	if task == nil {
		retMsg = "未找到任务"
		return
	}

	s.remove(name)

	if err = s.add(ctx, task.Name, task.Expression, task.BotId, []byte(task.Request)); err != nil {
		g.Log().Error(ctx, err)
		retMsg = "重新加入定时任务失败"
		return
	}

	retMsg = expr + " " + name
	return
}

func (s *sCrontab) ChangeBotIdReturnRes(ctx context.Context,
	botId int64, name string, creatorId int64,
) (retMsg string) {
	if botId == 0 {
		return
	}

	q := dao.Crontab.Ctx(ctx).
		Data(dao.Crontab.Columns().BotId, botId).
		Where(dao.Crontab.Columns().Name, name)
	if creatorId != 0 {
		q = q.Where(dao.Crontab.Columns().CreatorId, creatorId)
	}
	result, err := q.Update()
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "修改失败"
		return
	}

	if aff, _ := result.RowsAffected(); aff == 0 {
		return
	}

	task, err := s.getTask(ctx, name, creatorId)
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "获取失败"
		return
	}
	if task == nil {
		retMsg = "未找到任务"
		return
	}

	s.remove(name)

	if err = s.add(ctx, task.Name, task.Expression, task.BotId, []byte(task.Request)); err != nil {
		g.Log().Error(ctx, err)
		retMsg = "重新加入定时任务失败"
		return
	}

	retMsg = name + " -> " + gconv.String(botId)
	return
}

func (s *sCrontab) OneshotReturnRes(ctx context.Context, name string, creatorId int64) (retMsg string) {
	task, err := s.getTask(ctx, name, creatorId)
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "获取失败"
		return
	}
	if task == nil {
		retMsg = "未找到任务"
		return
	}

	if err = s.oneshot(task.BotId, []byte(task.Request)); err != nil {
		g.Log().Error(ctx, err)
		retMsg = "执行失败"
		return
	}

	retMsg = name
	return
}
