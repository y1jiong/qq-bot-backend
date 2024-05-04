package list

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"strings"
	"time"
)

func (s *sList) AddListReturnRes(ctx context.Context, listName, namespace string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 初始化 list 对象
	listE := entity.List{
		ListName:  listName,
		Namespace: namespace,
		ListJson:  "{}",
	}
	// 数据库插入
	_, err := dao.List.Ctx(ctx).
		Data(listE).
		OmitEmptyData().
		Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		retMsg = "新增 list 失败"
		return
	}
	// 同步写入
	service.Namespace().AddNamespaceList(ctx, namespace, listName)
	// 回执
	retMsg = "已新增 list(" + listName + ")"
	return
}

func (s *sList) RemoveListReturnRes(ctx context.Context, listName string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
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
	service.Namespace().RemoveNamespaceList(ctx, listE.Namespace, listName)
	// 回执
	retMsg = "已删除 list(" + listName + ")"
	return
}

func (s *sList) RecoverListReturnRes(ctx context.Context, listName string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	var listE *entity.List
	err := dao.List.Ctx(ctx).
		Fields(
			dao.List.Columns().Namespace,
			dao.List.Columns().DeletedAt,
		).
		Where(dao.List.Columns().ListName, listName).
		Unscoped().
		Scan(&listE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	if listE.DeletedAt == nil {
		retMsg = "我寻思 list(" + listName + ") 也没删除啊"
		return
	}
	// 数据库更新
	_, err = dao.List.Ctx(ctx).
		Where(dao.List.Columns().ListName, listName).
		Data(g.Map{
			dao.List.Columns().DeletedAt: nil,
			dao.List.Columns().UpdatedAt: gtime.Now(),
		}).
		Unscoped().
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 同步恢复
	service.Namespace().AddNamespaceList(ctx, listE.Namespace, listName)
	// 回执
	retMsg = "已恢复 list(" + listName + ")"
	return
}

func (s *sList) ExportListReturnRes(ctx context.Context, listName string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	var msg string
	msg = dao.List.Columns().Namespace + ": " + listE.Namespace + "\n" +
		dao.List.Columns().ListJson + ": " + listE.ListJson + "\n" +
		dao.List.Columns().UpdatedAt + ": " + listE.UpdatedAt.String()
	// 回执
	url, err := service.File().CacheFile(ctx, []byte(msg), time.Minute)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	filePath, err := service.Bot().UploadFile(ctx, url)
	if err != nil {
		retMsg = "上传文件失败"
		return
	}
	service.Bot().SendFile(ctx, filePath, "list("+listName+").txt")
	return
}

func (s *sList) QueryListLenReturnRes(ctx context.Context, listName string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(listE.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listMap := listJson.MustMap(make(map[string]any))
	// 回执
	retMsg = "list(" + listName + ") 共 " + gconv.String(len(listMap)) + " 条"
	return
}

func (s *sList) QueryListReturnRes(ctx context.Context, listName string, keys ...string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	var msg string
	if len(keys) > 0 {
		// 查询 key
		listJson, err := sj.NewJson([]byte(listE.ListJson))
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		keys[0] = service.Codec().DecodeBlank(keys[0])
		if _, ok := listJson.CheckGet(keys[0]); !ok {
			retMsg = "在 list(" + listName + ") 中未找到 key(" + keys[0] + ")"
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
		msg = dao.List.Columns().Namespace + ": " + listE.Namespace + "\n" +
			dao.List.Columns().ListJson + ": " + listE.ListJson + "\n" +
			dao.List.Columns().UpdatedAt + ": " + listE.UpdatedAt.String()
	}
	// 回执
	retMsg = msg
	return
}

func (s *sList) AddListDataReturnRes(ctx context.Context, listName, key string, value ...string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(listE.ListJson))
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
		retMsg = "已添加 key(" + key + ") value(" + value[0] + ") 到 list(" + listName + ")"
	} else {
		retMsg = "已添加 key(" + key + ") 到 list(" + listName + ")"
	}
	return
}

func (s *sList) RemoveListDataReturnRes(ctx context.Context, listName, key string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(listE.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 按照 url escape 解码空格和 %
	key = service.Codec().DecodeBlank(key)
	if _, ok := listJson.CheckGet(key); !ok {
		// 未找到 key
		retMsg = "在 list(" + listName + ") 中未找到 key(" + key + ")"
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
	retMsg = "已删除 key(" + key + ") 从 list(" + listName + ")"
	return
}

func (s *sList) ResetListDataReturnRes(ctx context.Context, listName string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
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
	retMsg = "已重置 list(" + listName + ") 的数据"
	return
}

func (s *sList) SetListDataReturnRes(ctx context.Context, listName, newListStr string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(newListStr))
	if err != nil {
		retMsg = "反序列化 JSON 失败"
		return
	}
	listM := listJson.MustMap(make(map[string]any))
	length := len(listM)
	if length < 1 {
		retMsg = "无效数据"
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
	retMsg = "已覆盖 list(" + listName + ") 的数据\n共 " + gconv.String(length) + " 条"
	return
}

func (s *sList) AppendListDataReturnRes(ctx context.Context, listName, newListStr string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(newListStr))
	if err != nil {
		retMsg = "反序列化 JSON 失败"
		return
	}
	listM := listJson.MustMap(make(map[string]any))
	appendLen := len(listM)
	if appendLen < 1 {
		retMsg = "无效数据"
		return
	}
	// 追加操作
	totalLen, err := s.AppendListData(ctx, listName, listM)
	if err != nil {
		return
	}
	// 回执
	retMsg = "已追加 list(" + listName + ") 的数据 " + gconv.String(appendLen) +
		" 条\n共 " + gconv.String(totalLen) + " 条"
	return
}

func (s *sList) GlanceListDataReturnRes(ctx context.Context, listName string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsPublicNamespace(listE.Namespace) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(listE.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listMap := listJson.MustMap(make(map[string]any))
	var msgBuilder strings.Builder
	for k := range listMap {
		msgBuilder.WriteString("`" + k + "`\n")
	}
	// 回执
	retMsg = strings.TrimRight(msgBuilder.String(), "\n")
	return
}

func (s *sList) CopyListKeyReturnRes(ctx context.Context, listName, srcKey, dstKey string) (retMsg string) {
	// 参数合法性校验
	if !legalListNameRe.MatchString(listName) {
		return
	}
	// 获取 list
	listE := getList(ctx, listName)
	if listE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, listE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	listJson, err := sj.NewJson([]byte(listE.ListJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 按照 url escape 解码空格和 %
	srcKey = service.Codec().DecodeBlank(srcKey)
	dstKey = service.Codec().DecodeBlank(dstKey)
	if _, ok := listJson.CheckGet(srcKey); !ok {
		retMsg = "在 list(" + listName + ") 中未找到 key(" + srcKey + ")"
		return
	}
	listJson.Set(dstKey, listJson.Get(srcKey).Interface())
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
	retMsg = "已复制 list(" + listName + ") 的 key(" + srcKey + ") 到 key(" + dstKey + ")"
	return
}
