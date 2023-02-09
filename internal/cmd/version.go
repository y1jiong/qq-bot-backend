package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcmd"
	"qq-bot-backend/internal/consts"
)

var (
	Version = gcmd.Command{
		Name:          "version",
		Brief:         "show version information of current binary",
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Println(consts.Description)
			return
		},
	}
)
