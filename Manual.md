# Manual

v1.0.0

{required} [optional]

## List

| Command                              | Description                        | Comment                                          |
| ------------------------------------ | ---------------------------------- | ------------------------------------------------ |
| /list join {list_name} {key} [value] | 将 key[:value] 添加到 list_name 中 | 需要 list_name 对应的 namespace admin 及以上权限 |
| /list leave {list_name} {key}        | 将 key 从 list_name 中移除         | 同上                                             |
| /list query {list_name}              | 查询 list_name                     | 同上                                             |
| /list add {list_name} {namespace}    | 在 namespace 下新增 list_name      | 同上                                             |
| /list rm {list_name}                 | 删除 list_name                     | 同上                                             |

## Group

| Command                                   | Description                                                  | Comment                                        |
| ----------------------------------------- | ------------------------------------------------------------ | ---------------------------------------------- |
| /group bind {namespace}                   | 将当前 group 绑定到 namespace 中（会重置当前 group 的所有配置） | 需要 group admin 和 namespace admin 及以上权限 |
| /group unbind                             | 解除当前 group 的绑定                                        | 同上                                           |
| /group query                              | 查询当前 group 的配置                                        | 需要 namespace admin 及以上权限                |
| /group approval enable mc                 | 入群审批流程启用 mc 正版用户名验证                           | 需要 group admin 和 namespace admin 及以上权限 |
| /group approval enable regexp             | 入群审批流程启用正则表达式                                   | 同上                                           |
| /group approval enable whitelist          | 入群审批流程启用白名单                                       | 同上                                           |
| /group approval enable blacklist          | 入群审批流程启用黑名单                                       | 同上                                           |
| /group approval set regexp {regexp}       | 指定入群审批流程的正则表达式（若有子表达式，则会采用第一个子表达式） | 同上                                           |
| /group approval add whitelist {list_name} | 新增入群审批流程白名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                           |
| /group approval add blacklist {list_name} | 新增入群审批流程黑名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                           |
| /group approval disable mc                | 入群审批流程禁用 mc 正版用户名验证                           | 同上                                           |
| /group approval disable regexp            | 入群审批流程禁用正则表达式                                   | 同上                                           |
| /group approval disable whitelist         | 入群审批流程禁用白名单                                       | 同上                                           |
| /group approval disable blacklist         | 入群审批流程禁用黑名单                                       | 同上                                           |
| /group approval rm whitelist {list_name}  | 移除入群审批流程白名单 list_name                             | 同上                                           |
| /group approval rm blacklist {list_name}  | 移除入群审批流程黑名单 list_name                             | 同上                                           |

## User

| Command                           | Description                            | Comment                   |
| --------------------------------- | -------------------------------------- | ------------------------- |
| /user join {namespace} {user_id}  | 将 user_id 添加到 namespace admin 名单 | 需要 namespace owner 权限 |
| /user leave {namespace} {user_id} | 将 user_id 移除到 namespace admin 名单 | 同上                      |

## Namespace

| Command                      | Description              | Comment                           |
| ---------------------------- | ------------------------ | --------------------------------- |
| /namespace add {namespace}   | 新建 namespace           | 需要系统授予的操作 namespace 权限 |
| /namespace rm {namespace}    | 删除 namespace           | 同上                              |
| /namespace query             | 查询自己所有的 namespace | 同上                              |
| /namespace {namespace}       | 查询 namespace 配置      | 需要 namespace admin 及以上权限   |
| /namespace {namespace} reset | 重置 namespace 的 admin  | 需要 namespace owner 权限         |

## Model

| Command            | Description | Comment        |
| ------------------ | ----------- | -------------- |
| /model set {model} | 设置机型    | 需要受系统信任 |

## Token

| Command                   | Description                    | Comment        |
| ------------------------- | ------------------------------ | -------------- |
| /token add {name} {token} | 添加可让 user 接入本系统的令牌 | 需要受系统信任 |
| /token rm {name}          | 删除令牌                       | 需要受系统信任 |
| /token query              | 查询自己所有的令牌             | 需要受系统信任 |

## System

| Command               | Description           | Comment               |
| --------------------- | --------------------- | --------------------- |
| No permission to view | No permission to view | No permission to view |

