package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "qq-bot-backend/internal/logic"

	"context"
	"qq-bot-backend/internal/cmd"
)

func main() {
	err := cmd.Main.AddCommand(&cmd.Install, &cmd.Uninstall, &cmd.Version)
	if err != nil {
		panic(err)
	}
	cmd.Main.Run(context.Background())
}
