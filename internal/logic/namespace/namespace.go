package namespace

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
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

const (
	adminsMapKey = "admins"

	globalNamespace = "global"
)

func getNamespace(ctx context.Context, namespace string) (namespaceE *entity.Namespace) {
	// 数据库查询
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Scan(&namespaceE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}

func isNamespaceOwner(userId int64, namespaceE *entity.Namespace) bool {
	return userId == namespaceE.OwnerId
}

func isNamespaceOwnerOrAdmin(ctx context.Context, userId int64, namespaceE *entity.Namespace) bool {
	// 判断 owner
	if isNamespaceOwner(userId, namespaceE) {
		return true
	}
	// 解析 setting json
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	// 判断 admin
	return settingJson.Get(adminsMapKey).Get(gconv.String(userId)).Valid()
}

func (s *sNamespace) IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) bool {
	// 参数合法性校验
	if userId == 0 || !legalNamespaceNameRe.MatchString(namespace) {
		return false
	}
	// 过程
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return false
	}
	return isNamespaceOwnerOrAdmin(ctx, userId, namespaceE)
}

func (s *sNamespace) IsNamespaceOwnerOrAdminOrOperator(ctx context.Context, namespace string, userId int64) bool {
	// 参数合法性校验
	if userId == 0 || !legalNamespaceNameRe.MatchString(namespace) {
		return false
	}
	// 过程
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return false
	}
	return isNamespaceOwnerOrAdmin(ctx, userId, namespaceE) ||
		service.User().CouldOpNamespace(ctx, gconv.Int64(userId))
}

func (s *sNamespace) IsGlobalNamespace(namespace string) bool {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return false
	}
	// 判断
	return namespace == globalNamespace
}

func (s *sNamespace) GetGlobalNamespace() string {
	return globalNamespace
}
