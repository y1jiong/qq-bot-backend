package consts

import (
	"runtime"

	"github.com/gogf/gf/v2"
)

const (
	ProjName = "qq-bot-backend"
	Version  = "1.9.4"
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
