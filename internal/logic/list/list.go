package list

import (
	"context"
	"encoding/json"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
)

type sList struct{}

func New() *sList {
	return &sList{}
}

func init() {
	service.RegisterList(New())
}

var (
	legalListNameRe = regexp.MustCompile(`^\S{1,16}$`)
)

func getList(ctx context.Context, listName string) (lEntity *entity.List) {
	// 数据库查询
	err := dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Scan(&lEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}

func (s *sList) GetListData(ctx context.Context, listName string) (listMap map[string]any) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(lEntity.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listMap = listJson.MustMap(make(map[string]any))
	return
}

func (s *sList) AppendListData(ctx context.Context, listName string, newMap map[string]any) (n int, err error) {
	listMap := s.GetListData(ctx, listName)
	if listMap == nil {
		return
	}
	// 追加数据
	for k, v := range newMap {
		listMap[k] = v
	}
	n = len(listMap)
	// 保存数据
	listBytes, err := json.Marshal(listMap)
	if err != nil {
		return
	}
	// 数据库更新
	_, err = dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Data(dao.List.Columns().ListJson, string(listBytes)).
		Update()
	return
}
