package namespace

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"
)

type sNamespace struct{}

func New() *sNamespace {
	return &sNamespace{}
}

func init() {
	service.RegisterNamespace(New())
}

var (
	legalNamespaceNameRe = regexp.MustCompile(`^\S{1,16}$`)
)

// setting json key
const (
	adminArrKey     = "admins"
	whitelistArrKey = "whitelists"
)

func getNamespace(ctx context.Context, namespace string) (nEntity *entity.Namespace) {
	// 数据库查询
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Scan(&nEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 没找到
	if nEntity == nil {
		service.Bot().SendMsg(ctx, "没找到 namespace "+namespace)
		return
	}
	return
}

func isNamespaceOwner(userId int64, nEntity *entity.Namespace) (yes bool) {
	// 判断 owner
	if userId != nEntity.OwnerId {
		return
	}
	yes = true
	return
}

func isNamespaceOwnerOrAdmin(ctx context.Context, userId int64, nEntity *entity.Namespace) (yes bool) {
	// 判断 owner
	if userId == nEntity.OwnerId {
		yes = true
		return
	}
	// 解析 setting json
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin array
	var adminArr []any
	if adminJson, ok := settingJson.CheckGet(adminArrKey); ok {
		adminArr = adminJson.MustArray()
	} else {
		return
	}
	// 判断 admin
	for _, v := range adminArr {
		if userId == gconv.Int64(v) {
			yes = true
		}
	}
	return
}

func (s *sNamespace) IsNamespaceExist(ctx context.Context, namespace string) (yes bool) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 过程
	n, err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return n > 0
}

func (s *sNamespace) IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 过程
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	return isNamespaceOwnerOrAdmin(ctx, userId, nEntity)
}

func (s *sNamespace) AddNewNamespace(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 初始化 setting json
	settingJson := sj.New()
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 初始化 namespace 对象
	namespaceEntity := entity.Namespace{
		Namespace:   namespace,
		OwnerId:     service.Bot().GetUserId(ctx),
		SettingJson: string(settingBytes),
	}
	// 数据库插入
	_, err = dao.Namespace.Ctx(ctx).
		Data(namespaceEntity).
		OmitEmpty().
		Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		// 返回错误
		service.Bot().SendMsg(ctx, "新增 namespace 失败")
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已新增 namespace("+namespace+")")
}

func (s *sNamespace) RemoveNamespace(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据库软删除
	_, err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已删除 namespace("+namespace+")")
}

func (s *sNamespace) QueryNamespace(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断 owner or admin
	if !isNamespaceOwnerOrAdmin(ctx, service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 回执
	msg := dao.Namespace.Columns().Namespace + ": " + nEntity.Namespace + "\n" +
		dao.Namespace.Columns().OwnerId + ": " + gconv.String(nEntity.OwnerId) + "\n" +
		dao.Namespace.Columns().SettingJson + ": " + nEntity.SettingJson
	service.Bot().SendMsg(ctx, msg)
}

func (s *sNamespace) QueryOwnNamespace(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 创建数组指针
	var nEntities []*entity.Namespace
	// 数据库查询
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().OwnerId, userId).
		Scan(&nEntities)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 判断空
	if len(nEntities) == 0 {
		return
	}
	// 回执
	var msg strings.Builder
	for _, v := range nEntities {
		msg.WriteString(dao.Namespace.Columns().Namespace)
		msg.WriteString(": ")
		msg.WriteString(v.Namespace)
		msg.WriteString("\n")
		msg.WriteString(dao.Namespace.Columns().CreatedAt)
		msg.WriteString(": ")
		msg.WriteString(v.CreatedAt.String())
		msg.WriteString("\n---\n")
	}
	service.Bot().SendMsg(ctx, msg.String())
}

func (s *sNamespace) AddNamespaceAdmin(ctx context.Context, namespace string, userId int64) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	if userId < 1 {
		return
	}
	// 获取 namespace 对象
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin array
	var adminArr []any
	if adminJson, ok := settingJson.CheckGet(adminArrKey); ok {
		adminArr = adminJson.MustArray()
	}
	// 添加 userId 的 admin 权限
	adminArr = append(adminArr, userId)
	settingJson.Set(adminArrKey, adminArr)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新 添加 admin 权限
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已添加 user("+gconv.String(userId)+") 的 namespace("+namespace+") admin 权限")
}

func (s *sNamespace) RemoveNamespaceAdmin(ctx context.Context, namespace string, userId int64) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 数据库查询
	var nEntity *entity.Namespace
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Scan(&nEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 没找到
	if nEntity == nil {
		service.Bot().SendMsg(ctx, "没找到 namespace "+namespace)
		return
	}
	// 判断是否是 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin array
	var adminArr []any
	if adminJson, ok := settingJson.CheckGet(adminArrKey); ok {
		adminArr = adminJson.MustArray()
	}
	// 判断空
	lengthOfAdminArr := len(adminArr)
	if lengthOfAdminArr < 1 {
		service.Bot().SendMsg(ctx, "admin 数组为空")
		return
	}
	// 查找 userId 是否拥有 admin 权限
	exist, position := false, -1
	for i, v := range adminArr {
		adminId := gconv.Int64(v)
		if userId == adminId {
			exist = true
			position = i
			break
		}
	}
	if !exist {
		service.Bot().SendMsg(ctx, gconv.String(userId)+"不存在")
		return
	}
	// 删除 userId 的 admin 权限
	if position == lengthOfAdminArr-1 {
		adminArr = adminArr[:position]
	} else {
		adminArr = append(adminArr[:position], adminArr[position+1:]...)
	}
	settingJson.Set(adminArrKey, adminArr)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新 移除 admin 权限
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已取消 user("+gconv.String(userId)+") 的 namespace("+namespace+") admin 权限")
}

func (s *sNamespace) ResetNamespace(ctx context.Context, namespace, option string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 数据库查询
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断是否是 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 选择重置
	switch option {
	case "admin":
		settingJson.Set(adminArrKey, []any{})
	case "whitelist":
		settingJson.Set(whitelistArrKey, []any{})
	case "all":
		settingJson = sj.New()
	}
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已重置 namespace("+namespace+") 的 "+option)
}
