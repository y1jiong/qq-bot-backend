package list

import (
	"context"
	"github.com/bytedance/sonic"
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

func getList(ctx context.Context, listName string) (listE *entity.List) {
	// 数据库查询
	err := dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Scan(&listE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}

// GetListData 获取 list 数据，返回的 map 一定不为 nil
func (s *sList) GetListData(ctx context.Context, listName string) (listMap map[string]any) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 数据处理
	listJson, err := sonic.GetFromString(listE.ListJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listMap, _ = listJson.Map()
	if listMap == nil {
		listMap = make(map[string]any)
	}
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
	listBytes, err := sonic.ConfigStd.Marshal(listMap)
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
