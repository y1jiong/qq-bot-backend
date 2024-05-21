package namespace

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	propertiesMapKey = "properties"

	propertyPublic = "public"
)

func (s *sNamespace) getNamespaceProperties(ctx context.Context, namespace string) (properties map[string]any) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	properties, _ = settingJson.Get(propertiesMapKey).Map()
	if properties == nil {
		properties = make(map[string]any)
	}
	return
}

func (s *sNamespace) IsNamespacePropertyPublic(ctx context.Context, namespace string) bool {
	properties := s.getNamespaceProperties(ctx, namespace)
	if v, ok := properties[propertyPublic].(bool); ok {
		return v
	}
	return false
}
