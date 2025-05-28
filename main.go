package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "qq-bot-backend/internal/logic"

	"context"
	"qq-bot-backend/internal/cmd"
)

func main() {
	if err := cmd.Main.AddCommand(
		&cmd.Version,
		&cmd.Install,
		&cmd.Uninstall,
	); err != nil {
		panic(err)
	}
	cmd.Main.Run(context.Background())
}
