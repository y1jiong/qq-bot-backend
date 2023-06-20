package consts

import (
	"github.com/gogf/gf/v2"
	"runtime"
)

const (
	ProjName = "qq-bot-backend"
	Version  = "v1.3.4"
)

var (
	GitCommit   = ""
	BuildTime   = ""
	Description = "Version: " + Version +
		"\nGo Version: " + runtime.Version() +
		"\nGoFrame Version: " + gf.VERSION +
		"\nGit Commit: " + GitCommit +
		"\nBuild Time: " + BuildTime
)
