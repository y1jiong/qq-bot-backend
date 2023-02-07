package list

import (
	"context"
	"encoding/json"
	"errors"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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

func (s *sList) AddList(ctx context.Context, listName, namespace string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 初始化 list 对象
	lEntity := entity.List{
		ListName:  listName,
		Namespace: namespace,
		ListJson:  "{}",
	}
	// 数据库插入
	_, err := dao.List.Ctx(ctx).
		Data(lEntity).
		OmitEmpty().
		Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		service.Bot().SendPlainMsg(ctx, "新增 list 失败")
		return
	}
	// 同步写入
	service.Namespace().AddNamespaceList(ctx, namespace, listName)
	// 回执
	service.Bot().SendPlainMsg(ctx, "已新增 list("+listName+")")
}

func (s *sList) RemoveList(ctx context.Context, listName string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, lEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据库软删除
	_, err := dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 同步删除
	service.Namespace().RemoveNamespaceList(ctx, lEntity.Namespace, listName)
	// 回执
	service.Bot().SendPlainMsg(ctx, "已删除 list("+listName+")")
}

func (s *sList) QueryList(ctx context.Context, listName string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, lEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	msg := dao.List.Columns().Namespace + ": " + lEntity.Namespace + "\n" +
		dao.List.Columns().ListJson + ": " + lEntity.ListJson + "\n" +
		dao.List.Columns().UpdatedAt + ": " + lEntity.UpdatedAt.String()
	// 回执
	service.Bot().SendPlainMsg(ctx, msg)
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

func (s *sList) AddListData(ctx context.Context, listName, key string, value ...string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, lEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(lEntity.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 按照 url escape 解码空格和 %
	key = service.Codec().DecodeBlank(key)
	if len(value) > 0 {
		// 按照 url escape 解码空格和 %
		value[0] = service.Codec().DecodeBlank(value[0])
		listJson.Set(key, value[0])
	} else {
		listJson.Set(key, nil)
	}
	// 保存数据
	listBytes, err := listJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Data(dao.List.Columns().ListJson, string(listBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if len(value) > 0 {
		service.Bot().SendPlainMsg(ctx, "已添加 key("+key+") value("+value[0]+") 到 list("+listName+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已添加 key("+key+") 到 list("+listName+")")
	}
}

func (s *sList) RemoveListData(ctx context.Context, listName, key string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, lEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(lEntity.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 按照 url escape 解码空格和 %
	key = service.Codec().DecodeBlank(key)
	if _, ok := listJson.CheckGet(key); !ok {
		// 不存在 key
		service.Bot().SendPlainMsg(ctx, key+" 不存在")
		return
	}
	listJson.Del(key)
	// 保存数据
	listBytes, err := listJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Data(dao.List.Columns().ListJson, string(listBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已删除 key("+key+") 从 list("+listName+")")
}

func (s *sList) ResetListData(ctx context.Context, listName string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, lEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据库更新
	_, err := dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Data(dao.List.Columns().ListJson, "{}").
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已重置 list("+listName+") 的数据")
}

func (s *sList) SetListData(ctx context.Context, listName, newListBytes string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	lEntity := getList(ctx, listName)
	if lEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, lEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(newListBytes))
	if err != nil {
		service.Bot().SendPlainMsg(ctx, "反序列化 json 失败")
		return
	}
	listM := listJson.MustMap(make(map[string]any))
	length := len(listM)
	if length < 1 {
		service.Bot().SendPlainMsg(ctx, "无效数据")
		return
	}
	// 保存数据
	listBytes, err := listJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Data(dao.List.Columns().ListJson, string(listBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已覆盖 list("+listName+") 的数据\n共 "+gconv.String(length)+" 条")
}

func (s *sList) AppendListData(ctx context.Context, listName string, newMap map[string]any) (err error) {
	listMap := s.GetListData(ctx, listName)
	if listMap == nil {
		err = errors.New("未找到 list(" + listName + ")")
		return
	}
	// 追加数据
	for k, v := range newMap {
		listMap[k] = v
	}
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
