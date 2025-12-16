package crontab

import (
	"context"
	"errors"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcron"
	"go.opentelemetry.io/otel/trace"
)

type sCrontab struct{}

func New() *sCrontab {
	return &sCrontab{}
}

func init() {
	service.RegisterCrontab(New())
}

var (
	crontabMu sync.Mutex
	crontab   *gcron.Cron
)

func (s *sCrontab) Run(ctx context.Context) {
	var tasks []entity.Crontab

	if err := utility.RetryWithBackoff(ctx, func() bool {
		var err error
		tasks, err = s.getTasks(ctx)
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		return true
	}, 4, utility.ExponentialBackoffWithJitter(ctx)); err != nil {
		g.Log().Error(ctx, err)
		return
	}

	s.renew()

	count := len(tasks)
	for _, task := range tasks {
		if err := s.add(ctx, task.Name, task.Expression, task.BotId, []byte(task.Request)); err != nil {
			g.Log().Error(ctx, err)
			count--
			continue
		}
	}

	g.Log().Info(ctx, "crontab loaded", count, "tasks")
}

func (s *sCrontab) getTasks(ctx context.Context) (tasks []entity.Crontab, err error) {
	err = dao.Crontab.Ctx(ctx).
		Fields(
			dao.Crontab.Columns().Name,
			dao.Crontab.Columns().Expression,
			dao.Crontab.Columns().BotId,
			dao.Crontab.Columns().Request,
		).
		Scan(&tasks)
	return
}

func (s *sCrontab) getTask(ctx context.Context, name string, creatorId int64) (task *entity.Crontab, err error) {
	q := dao.Crontab.Ctx(ctx).
		Fields(
			dao.Crontab.Columns().Name,
			dao.Crontab.Columns().Expression,
			dao.Crontab.Columns().BotId,
			dao.Crontab.Columns().Request,
		).
		Where(dao.Crontab.Columns().Name, name)
	if creatorId != 0 {
		q = q.Where(dao.Crontab.Columns().CreatorId, creatorId)
	}
	err = q.Scan(&task)
	return
}

func (s *sCrontab) oneshot(botId int64, reqJSON []byte) (err error) {
	ctx := service.Bot().LoadConnection(botId)
	if ctx == nil {
		return errors.New("bot not found")
	}
	// new trace
	ctx, span := gtrace.NewSpan(trace.ContextWithSpanContext(ctx, trace.SpanContext{}), "crontab.oneshot")
	defer span.End()
	service.Bot().Process(ctx, reqJSON, service.Process().Process)
	return
}

func (s *sCrontab) lazyInit() {
	if crontab != nil {
		return
	}

	crontabMu.Lock()
	defer crontabMu.Unlock()

	if crontab == nil {
		crontab = gcron.New()
	}
}

func (s *sCrontab) add(ctx context.Context, name, expr string, botId int64, reqJSON []byte) (err error) {
	s.lazyInit()

	if len(strings.Fields(expr)) == 5 {
		expr = "# " + expr
	}

	_, err = crontab.AddSingleton(ctx, expr, func(ctx context.Context) {
		if err := s.oneshot(botId, reqJSON); err != nil {
			g.Log().Notice(ctx, err)
		}
	}, name)

	return
}

func (s *sCrontab) remove(name string) {
	s.lazyInit()
	crontab.Remove(name)
}

func (s *sCrontab) renew() {
	crontabMu.Lock()
	defer crontabMu.Unlock()

	if crontab != nil {
		crontab.Close()
	}
	crontab = gcron.New()
}
