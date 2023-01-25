package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
)

type sGroup struct{}

func New() *sGroup {
	return &sGroup{}
}

func init() {
	service.RegisterGroup(New())
}

const (
	approvalProcessMapKey   = "approvalProcess"
	approvalRegexpKey       = "approvalRegexp"
	approvalWhitelistMapKey = "approvalWhitelists"
	approvalBlacklistMapKey = "approvalBlacklists"
	regexpCmd               = "regexp"
	whitelistCmd            = "whitelist"
	blacklistCmd            = "blacklist"
)

func getGroup(ctx context.Context, groupId int64) (gEntity *entity.Group) {
	// 数据库查询
	err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Scan(&gEntity)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return
}

func (s *sGroup) BindNamespace(ctx context.Context, groupId int64, namespace string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) ||
		!service.Namespace().IsNamespaceOwnerOrAdmin(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据库计数
	n, err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if n > 0 {
		// 数据库更新
		_, err = dao.Group.Ctx(ctx).
			Where(dao.Group.Columns().GroupId, groupId).
			Data(dao.Group.Columns().Namespace, namespace).
			Update()
	} else {
		// 数据库插入
		gEntity := entity.Group{
			GroupId:     groupId,
			Namespace:   namespace,
			SettingJson: "{}",
		}
		_, err = dao.Group.Ctx(ctx).
			Data(gEntity).
			OmitEmpty().
			Insert()
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已绑定当前 group("+gconv.String(groupId)+") 到 namespace("+namespace+")")
}

func (s *sGroup) Unbind(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !s.IsGroupBindNamespaceOwnerOrAdmin(ctx, groupId, service.Bot().GetUserId(ctx)) {
		return
	}
	// 过程
	n, err := dao.Group.Ctx(ctx).Where(dao.Group.Columns().GroupId, groupId).Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if n < 1 {
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().Namespace, "").
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendMsg(ctx, "已解除 group("+gconv.String(groupId)+") 的 namespace 绑定")
}

func (s *sGroup) QueryGroup(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 过程
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		service.Bot().SendMsg(ctx, "没有任何数据")
		return
	}
	// 回执
	msg := dao.Group.Columns().Namespace + ": " + gEntity.Namespace + "\n" +
		dao.Group.Columns().SettingJson + ": " + gEntity.SettingJson + "\n" +
		dao.Group.Columns().UpdatedAt + ": " + gEntity.UpdatedAt.String()
	service.Bot().SendMsg(ctx, msg)
}

func (s *sGroup) IsGroupBindNamespaceOwnerOrAdmin(ctx context.Context, groupId, userId int64) (yes bool) {
	// 参数合法性校验
	if groupId < 1 || userId < 1 {
		return
	}
	// 获取 group 绑定的 namespace
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 过程
	return service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, userId)
}

func (s *sGroup) GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	process = settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) AddApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	approvalProcessMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
	// 添加 processName
	approvalProcessMap[processName] = nil
	// 处理 args
	if len(args) > 0 {
		switch processName {
		case regexpCmd:
			// 处理正则表达式
			_, err = regexp.Compile(args[0])
			if err != nil {
				service.Bot().SendMsg(ctx, "输入的正则表达式无法通过编译")
				return
			}
			settingJson.Set(approvalRegexpKey, args[0])
		case whitelistCmd:
			// 处理白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			whitelists := settingJson.Get(approvalWhitelistMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(approvalWhitelistMapKey, whitelists)
		case blacklistCmd:
			// 处理黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			blacklists := settingJson.Get(approvalBlacklistMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(approvalBlacklistMapKey, blacklists)
		}
	}
	// 保存数据
	settingJson.Set(approvalProcessMapKey, approvalProcessMap)
	approvalProcessBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(approvalProcessBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if len(args) > 0 {
		service.Bot().SendMsg(ctx,
			"已添加 group("+gconv.String(groupId)+") 入群审批流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendMsg(ctx, "已添加 group("+gconv.String(groupId)+") 入群审批流程 "+processName)
	}
}

func (s *sGroup) RemoveApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	approvalProcessMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
	// 删除 processName
	if _, ok := approvalProcessMap[processName]; ok {
		delete(approvalProcessMap, processName)
	} else {
		service.Bot().SendMsg(ctx, processName+" 不存在")
		return
	}
	// 处理 args
	if len(args) > 0 {
		switch processName {
		case whitelistCmd:
			// 处理白名单
			whitelists := settingJson.Get(approvalWhitelistMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; ok {
				delete(whitelists, args[0])
			} else {
				service.Bot().SendMsg(ctx, args[0]+" 不存在")
				return
			}
			settingJson.Set(approvalWhitelistMapKey, whitelists)
		case blacklistCmd:
			// 处理黑名单
			blacklists := settingJson.Get(approvalBlacklistMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; ok {
				delete(blacklists, args[0])
			} else {
				service.Bot().SendMsg(ctx, args[0]+" 不存在")
				return
			}
			settingJson.Set(approvalBlacklistMapKey, blacklists)
		}
	}
	// 保存数据
	settingJson.Set(approvalProcessMapKey, approvalProcessMap)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if len(args) > 0 {
		service.Bot().SendMsg(ctx,
			"已移除 group("+gconv.String(groupId)+") 入群审批流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendMsg(ctx, "已移除 group("+gconv.String(groupId)+") 入群审批流程 "+processName)
	}
}

func (s *sGroup) GetWhitelist(ctx context.Context, groupId int64) (whitelists map[string]any) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	whitelists = settingJson.Get(approvalWhitelistMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetBlacklist(ctx context.Context, groupId int64) (blacklists map[string]any) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	blacklists = settingJson.Get(blacklistCmd).MustMap(make(map[string]any))
	return
}
