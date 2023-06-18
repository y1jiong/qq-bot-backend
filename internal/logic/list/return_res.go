package list

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sList) AddListReturnRes(ctx context.Context, listName, namespace string) {
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

func (s *sList) RemoveListReturnRes(ctx context.Context, listName string) {
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

func (s *sList) ExportListReturnRes(ctx context.Context, listName string) {
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
	var msg string
	msg = dao.List.Columns().Namespace + ": " + lEntity.Namespace + "\n" +
		dao.List.Columns().ListJson + ": " + lEntity.ListJson + "\n" +
		dao.List.Columns().UpdatedAt + ": " + lEntity.UpdatedAt.String()
	// 回执
	id, err := service.File().SetCachedFile(ctx, msg, time.Minute)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	url, err := service.File().GetCachedFileUrl(ctx, id)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	service.Bot().SendFile(ctx, "list("+listName+").txt", url)
}

func (s *sList) QueryListLenReturnRes(ctx context.Context, listName string) {
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
	listMap := listJson.MustMap(make(map[string]any))
	// 回执
	service.Bot().SendPlainMsg(ctx, "list("+listName+") 共 "+gconv.String(len(listMap))+" 条")
}

func (s *sList) QueryListReturnRes(ctx context.Context, listName string, keys ...string) {
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
	var msg string
	if len(keys) > 0 {
		// 查询 key
		listJson, err := sj.NewJson([]byte(lEntity.ListJson))
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		keys[0] = service.Codec().DecodeBlank(keys[0])
		if _, ok := listJson.CheckGet(keys[0]); !ok {
			service.Bot().SendPlainMsg(ctx, "在 list("+listName+") 中未找到 key("+keys[0]+")")
			return
		}
		viewJson := sj.New()
		viewJson.Set(keys[0], listJson.Get(keys[0]))
		msgBytes, err := viewJson.Encode()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		msg = string(msgBytes)
	} else {
		msg = dao.List.Columns().Namespace + ": " + lEntity.Namespace + "\n" +
			dao.List.Columns().ListJson + ": " + lEntity.ListJson + "\n" +
			dao.List.Columns().UpdatedAt + ": " + lEntity.UpdatedAt.String()
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, msg)
}

func (s *sList) AddListDataReturnRes(ctx context.Context, listName, key string, value ...string) {
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

func (s *sList) RemoveListDataReturnRes(ctx context.Context, listName, key string) {
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
		// 未找到 key
		service.Bot().SendPlainMsg(ctx, "在 list("+listName+") 中未找到 key("+key+")")
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

func (s *sList) ResetListDataReturnRes(ctx context.Context, listName string) {
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

func (s *sList) SetListDataReturnRes(ctx context.Context, listName, newListStr string) {
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
	listJson, err := sj.NewJson([]byte(newListStr))
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

func (s *sList) AppendListDataReturnRes(ctx context.Context, listName, newListStr string) {
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
	listJson, err := sj.NewJson([]byte(newListStr))
	if err != nil {
		service.Bot().SendPlainMsg(ctx, "反序列化 json 失败")
		return
	}
	listM := listJson.MustMap(make(map[string]any))
	appendLen := len(listM)
	if appendLen < 1 {
		service.Bot().SendPlainMsg(ctx, "无效数据")
		return
	}
	// 追加操作
	totalLen, err := s.AppendListData(ctx, listName, listM)
	if err != nil {
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已追加 list("+listName+") 的数据 "+gconv.String(appendLen)+
		" 条\n共 "+gconv.String(totalLen)+" 条")
}
