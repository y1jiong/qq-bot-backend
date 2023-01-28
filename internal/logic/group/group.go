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
	// approval
	approvalProcessMapKey    = "approvalProcess"
	approvalRegexpKey        = "approvalRegexp"
	approvalWhitelistsMapKey = "approvalWhitelists"
	approvalBlacklistsMapKey = "approvalBlacklists"
	// keyword
	keywordProcessMapKey    = "keywordProcess"
	keywordWhitelistsMapKey = "keywordWhitelists"
	keywordBlacklistsMapKey = "keywordBlacklists"
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
	if len(args) > 0 {
		// 处理 args
		switch processName {
		case consts.WhitelistCmd:
			// 处理白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			whitelists := settingJson.Get(approvalWhitelistsMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(approvalWhitelistsMapKey, whitelists)
		case consts.BlacklistCmd:
			// 处理黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			blacklists := settingJson.Get(approvalBlacklistsMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(approvalBlacklistsMapKey, blacklists)
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
		}
	} else {
		// 添加 processName
		processMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
		processMap[processName] = nil
		settingJson.Set(approvalProcessMapKey, processMap)
	}
	// 保存数据
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
			"已添加 group("+gconv.String(groupId)+") 入群审批流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已启用 group("+gconv.String(groupId)+") 入群审批流程 "+processName)
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
	if len(args) > 0 {
		// 处理 args
		switch processName {
		case consts.WhitelistCmd:
			// 处理白名单
			whitelists := settingJson.Get(approvalWhitelistsMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(whitelists, args[0])
			settingJson.Set(approvalWhitelistsMapKey, whitelists)
		case consts.BlacklistCmd:
			// 处理黑名单
			blacklists := settingJson.Get(approvalBlacklistsMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(blacklists, args[0])
			settingJson.Set(approvalBlacklistsMapKey, blacklists)
		}
	} else {
		// 删除 processName
		processMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
		if _, ok := processMap[processName]; !ok {
			service.Bot().SendPlainMsg(ctx, processName+" 不存在")
			return
		}
		delete(processMap, processName)
		settingJson.Set(approvalProcessMapKey, processMap)
	}
	// 保存数据
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
		service.Bot().SendPlainMsg(ctx, "已禁用 group("+gconv.String(groupId)+") 入群审批流程 "+processName)
	}
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

func (s *sGroup) GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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
	whitelists = settingJson.Get(approvalWhitelistsMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists = settingJson.Get(approvalBlacklistsMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetApprovalRegexp(ctx context.Context, groupId int64) (exp string) {
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
	exp = settingJson.Get(approvalRegexpKey).MustString()
	return
}

func (s *sGroup) AddKeywordProcess(ctx context.Context, groupId int64, processName string, args ...string) {
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
	if len(args) > 0 {
		// 处理 args
		switch processName {
		case consts.BlacklistCmd:
			// 处理黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			blacklists := settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(keywordBlacklistsMapKey, blacklists)
		case consts.WhitelistCmd:
			// 处理白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			whitelists := settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(keywordWhitelistsMapKey, whitelists)
		}
	} else {
		// 添加 processName
		processMap := settingJson.Get(keywordProcessMapKey).MustMap(make(map[string]any))
		processMap[processName] = nil
		settingJson.Set(keywordProcessMapKey, processMap)
	}
	// 保存数据
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
			"已添加 group("+gconv.String(groupId)+") 关键词检查流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已启用 group("+gconv.String(groupId)+") 关键词检查流程 "+processName)
	}
}

func (s *sGroup) RemoveKeywordProcess(ctx context.Context, groupId int64, processName string, args ...string) {
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
	if len(args) > 0 {
		// 处理 args
		switch processName {
		case consts.BlacklistCmd:
			// 处理黑名单
			blacklists := settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(blacklists, args[0])
			settingJson.Set(keywordBlacklistsMapKey, blacklists)
		case consts.WhitelistCmd:
			// 处理白名单
			whitelists := settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(whitelists, args[0])
			settingJson.Set(keywordWhitelistsMapKey, whitelists)
		}
	} else {
		// 删除 processName
		processMap := settingJson.Get(keywordProcessMapKey).MustMap(make(map[string]any))
		if _, ok := processMap[processName]; !ok {
			service.Bot().SendPlainMsg(ctx, processName+" 不存在")
			return
		}
		delete(processMap, processName)
		settingJson.Set(keywordProcessMapKey, processMap)
	}
	// 保存数据
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
			"已移除 group("+gconv.String(groupId)+") 关键词检查流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已禁用 group("+gconv.String(groupId)+") 关键词检查流程 "+processName)
	}
}

func (s *sGroup) GetKeywordProcess(ctx context.Context, groupId int64) (process map[string]any) {
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
	process = settingJson.Get(keywordProcessMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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
	whitelists = settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists = settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
	return
}
