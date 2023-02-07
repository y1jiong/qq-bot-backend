# Manual

v1.2

{required} [optional]

## List

| Command                              | Description                                                  | Comment                                          |
| ------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------ |
| /list add {list_name} {namespace}    | 在 namespace 下新增 list_name                                | 需要 list_name 对应的 namespace admin 及以上权限 |
| /list join {list_name} {key} [value] | 将 key[:value] 添加到 list_name 中（key value 可作为双因子认证使用）（key value 包含**空格**请用`%20`转义替换，包含**%**请用`%25`转义替换） | 同上                                             |
| /list leave {list_name} {key}        | 将 key 从 list_name 中移除（key 包含**空格**请用`%20`转义替换，包含**%**请用`%25`转义替换） | 同上                                             |
| /list query {list_name}              | 查询 list_name                                               | 同上                                             |
| /list set {list_name} {json}         | 用 json 覆盖 list_name 的数据                                | 同上                                             |
| /list reset {list_name}              | 重置 list_name 的数据                                        | 同上                                             |
| /list rm {list_name}                 | 删除 list_name（删除后原 list_name 不可使用）                | 同上                                             |

## Group

| Command                 | Description                                                  | Comment                                        |
| ----------------------- | ------------------------------------------------------------ | ---------------------------------------------- |
| /group bind {namespace} | 将当前 group 绑定到 namespace 中（会重置当前 group 的所有配置） | 需要 group admin 和 namespace admin 及以上权限 |
| /group unbind           | 解除当前 group 的绑定                                        | 同上                                           |
| /group query            | 查询当前 group 的配置                                        | 需要 namespace admin 及以上权限                |


## Group Approval

| Command                                   | Description                                                  | Comment                                        |
| ----------------------------------------- | ------------------------------------------------------------ | ---------------------------------------------- |
| /group approval enable mc                 | 入群审批流程启用 mc 正版用户名验证（将使用正版 UUID 作为双因子认证的输入） | 需要 group admin 和 namespace admin 及以上权限 |
| /group approval enable regexp             | 入群审批流程启用正则表达式（将使用匹配结果作为双因子认证的输入） | 同上                                           |
| /group approval enable whitelist          | 入群审批流程启用白名单                                       | 同上                                           |
| /group approval enable blacklist          | 入群审批流程启用黑名单                                       | 同上                                           |
| /group approval set regexp {regexp}       | 指定入群审批流程的正则表达式（若有子表达式，则会使用第一个子表达式的匹配结果） | 同上                                           |
| /group approval add whitelist {list_name} | 新增入群审批流程白名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                           |
| /group approval add blacklist {list_name} | 新增入群审批流程黑名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                           |
| /group approval rm whitelist {list_name}  | 移除入群审批流程白名单 list_name                             | 同上                                           |
| /group approval rm blacklist {list_name}  | 移除入群审批流程黑名单 list_name                             | 同上                                           |
| /group approval disable mc                | 入群审批流程禁用 mc 正版用户名验证                           | 同上                                           |
| /group approval disable regexp            | 入群审批流程禁用正则表达式                                   | 同上                                           |
| /group approval disable whitelist         | 入群审批流程禁用白名单                                       | 同上                                           |
| /group approval disable blacklist         | 入群审批流程禁用黑名单                                       | 同上                                           |

## Group Keyword

| Command                                  | Description                                                  | Comment                                        |
| ---------------------------------------- | ------------------------------------------------------------ | ---------------------------------------------- |
| /group keyword enable blacklist          | 关键词检查流程启用黑名单                                     | 需要 group admin 和 namespace admin 及以上权限 |
| /group keyword enable whitelist          | 关键词检查流程启用白名单                                     | 同上                                           |
| /group keyword add blacklist {list_name} | 新增关键词检查流程黑名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                           |
| /group keyword add whitelist {list_name} | 新增关键词检查流程白名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                           |
| /group keyword set reply {list_name}     | 设置关键词回复列表 list_name                                 | 同上                                           |
| /group keyword rm blacklist {list_name}  | 移除关键词检查流程黑名单 list_name                           | 同上                                           |
| /group keyword rm whitelist {list_name}  | 移除关键词检查流程白名单 list_name                           | 同上                                           |
| /group keyword rm reply                  | 移除关键词回复列表                                           | 同上                                           |
| /group keyword disable blacklist         | 关键词检查流程禁用黑名单                                     | 同上                                           |
| /group keyword disable whitelist         | 关键词检查流程禁用白名单                                     | 同上                                           |

## Group Log

| Command                          | Description            | Comment                                        |
| -------------------------------- | ---------------------- | ---------------------------------------------- |
| /group log leave set {list_name} | 设置离群记录 list_name | 需要 group admin 和 namespace admin 及以上权限 |
| /group log leave rm              | 移除离群记录 list_name | 同上                                           |

## Group Export

| Command                          | Description                    | Comment                                        |
| -------------------------------- | ------------------------------ | ---------------------------------------------- |
| /group export member {list_name} | 导出 group member 到 list_name | 需要 group admin 和 namespace admin 及以上权限 |

## User

| Command                           | Description                            | Comment                   |
| --------------------------------- | -------------------------------------- | ------------------------- |
| /user join {namespace} {user_id}  | 将 user_id 添加到 namespace admin 名单 | 需要 namespace owner 权限 |
| /user leave {namespace} {user_id} | 将 user_id 从 namespace admin 名单移除 | 同上                      |

## Namespace

| Command                      | Description              | Comment                           |
| ---------------------------- | ------------------------ | --------------------------------- |
| /namespace add {namespace}   | 新建 namespace           | 需要系统授予的操作 namespace 权限 |
| /namespace rm {namespace}    | 删除 namespace           | 同上                              |
| /namespace query             | 查询自己所有的 namespace | 同上                              |
| /namespace {namespace}       | 查询 namespace 配置      | 需要 namespace admin 及以上权限   |
| /namespace {namespace} reset | 重置 namespace 的 admin  | 需要 namespace owner 权限         |

## Extra

| Command                   | Description                    | Comment                       |
| ------------------------- | ------------------------------ | ----------------------------- |
| /raw {message}            | 获取 message 的原始信息        | 需要系统授予的获取 raw 的权限 |
| /model set {model}        | 设置机型                       | 需要受系统信任                |
| /token add {name} {token} | 添加可让 user 接入本系统的令牌 | 需要受系统信任                |
| /token rm {name}          | 删除令牌                       | 需要受系统信任                |
| /token query              | 查询自己所有的令牌             | 需要受系统信任                |

## System

| Command               | Description           | Comment               |
| --------------------- | --------------------- | --------------------- |
| No permission to view | No permission to view | No permission to view |

