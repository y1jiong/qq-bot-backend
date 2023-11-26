package consts

import (
	"github.com/gogf/gf/v2"
	"runtime"
)

const (
	ProjName = "qq-bot-backend"
	Version  = "v1.5.0"
)

var (
	GitTag      = ""
	GitCommit   = ""
	BuildTime   = ""
	Description = "Version: " + Version +
		"\nGo Version: " + runtime.Version() +
		"\nGoFrame Version: " + gf.VERSION +
		"\nGit Tag: " + GitTag +
		"\nGit Commit: " + GitCommit +
		"\nBuild Time: " + BuildTime
)
