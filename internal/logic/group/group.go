package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
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
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	var err error
	if gEntity == nil {
		// 初始化 group 对象
		gEntity = &entity.Group{
			GroupId:     groupId,
			Namespace:   namespace,
			SettingJson: "{}",
		}
		// 数据库插入
		_, err = dao.Group.Ctx(ctx).
			Data(gEntity).
			OmitEmpty().
			Insert()
	} else {
		if gEntity.Namespace != "" {
			service.Bot().SendPlainMsg(ctx,
				"当前 group("+gconv.String(groupId)+") 已经绑定了 namespace("+gEntity.Namespace+")")
			return
		}
		// 重置 setting
		gEntity = &entity.Group{
			Namespace:   namespace,
			SettingJson: "{}",
		}
		// 数据库更新
		_, err = dao.Group.Ctx(ctx).
			Where(dao.Group.Columns().GroupId, groupId).
			Data(gEntity).
			OmitEmpty().
			Update()
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已绑定当前 group("+gconv.String(groupId)+") 到 namespace("+namespace+")")
}

func (s *sGroup) Unbind(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	if gEntity.Namespace == "" {
		return
	}
	// 数据库更新
	_, err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().Namespace, "").
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已解除 group("+gconv.String(groupId)+") 的 namespace 绑定")
}

func (s *sGroup) QueryGroup(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 回执
	msg := dao.Group.Columns().Namespace + ": " + gEntity.Namespace + "\n" +
		dao.Group.Columns().SettingJson + ": " + gEntity.SettingJson + "\n" +
		dao.Group.Columns().UpdatedAt + ": " + gEntity.UpdatedAt.String()
	service.Bot().SendPlainMsg(ctx, msg)
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
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	approvalProcessMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
	// 处理 args
	if len(args) > 0 {
		switch processName {
		case consts.RegexpCmd:
			if service.Module().IsIncludeCqCode(args[0]) {
				// 包含 CQ Code 时发送表情 gun
				service.Bot().SendMsg(ctx, "[CQ:face,id=288]")
				return
			}
			// 解码被 CQ Code 转义的字符
			args[0] = service.Module().DecodeCqCode(args[0])
			// 处理正则表达式
			_, err = regexp.Compile(args[0])
			if err != nil {
				service.Bot().SendPlainMsg(ctx, "输入的正则表达式无法通过编译")
				return
			}
			settingJson.Set(approvalRegexpKey, args[0])
		case consts.WhitelistCmd:
			// 处理白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			whitelists := settingJson.Get(approvalWhitelistMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(approvalWhitelistMapKey, whitelists)
		case consts.BlacklistCmd:
			// 处理黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			blacklists := settingJson.Get(approvalBlacklistMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(approvalBlacklistMapKey, blacklists)
		}
	} else {
		// 添加 processName
		approvalProcessMap[processName] = nil
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
		service.Bot().SendPlainMsg(ctx,
			"已添加 group("+gconv.String(groupId)+") 入群审批流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已添加 group("+gconv.String(groupId)+") 入群审批流程 "+processName)
	}
}

func (s *sGroup) RemoveApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	approvalProcessMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
	// 处理 args
	if len(args) > 0 {
		switch processName {
		case consts.WhitelistCmd:
			// 处理白名单
			whitelists := settingJson.Get(approvalWhitelistMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(whitelists, args[0])
			settingJson.Set(approvalWhitelistMapKey, whitelists)
		case consts.BlacklistCmd:
			// 处理黑名单
			blacklists := settingJson.Get(approvalBlacklistMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(blacklists, args[0])
			settingJson.Set(approvalBlacklistMapKey, blacklists)
		}
	} else {
		// 删除 processName
		if _, ok := approvalProcessMap[processName]; !ok {
			service.Bot().SendPlainMsg(ctx, processName+" 不存在")
			return
		}
		delete(approvalProcessMap, processName)
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
		service.Bot().SendPlainMsg(ctx,
			"已移除 group("+gconv.String(groupId)+") 入群审批流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已移除 group("+gconv.String(groupId)+") 入群审批流程 "+processName)
	}
}

func (s *sGroup) GetWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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

func (s *sGroup) GetBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists = settingJson.Get(approvalBlacklistMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetRegexp(ctx context.Context, groupId int64) (re string) {
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
	re = settingJson.Get(approvalRegexpKey).MustString()
	return
}
