# qq-bot-backend

支持 OneBot11 协议 CQ 码消息上报格式

适配 [go-cqhttp](https://docs.go-cqhttp.org) [LLOneBot](https://github.com/LLOneBot/LLOneBot) and so on...

## 项目规范

1. 如果在 `defer` 语句中或过程中（例如 `defer wg.Wait()`）修改返回值，需要使用命名返回值

2. 命令参数合法性校验在最后进行

3. 不要用 `time.Sleep()`，要用 `select` `<-ctx.Done()` 和 `<-time.After()` 实现超时控制