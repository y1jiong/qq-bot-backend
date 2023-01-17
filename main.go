package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	_ "he3-bot/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"he3-bot/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
