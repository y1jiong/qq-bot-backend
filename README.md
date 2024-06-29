# qq-bot-backend

支持 OneBot11 协议 CQ 码消息上报格式

适配 [go-cqhttp](https://docs.go-cqhttp.org) [LLOneBot](https://github.com/LLOneBot/LLOneBot)

## 项目规范

1. 如果在 `defer` 语句中或过程中（例如 `defer wg.Wait()`）需要修改返回值，尽量使用命名返回值，以避免在代码维护和理解时出现混淆。

2. 参数合法性校验在最后进行