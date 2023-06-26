package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/lukesampson/figlet/figletlib"
	"github.com/nsf/termbox-go"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/resource/fonts"
	"strings"
)

var (
	Version = gcmd.Command{
		Name:          "version",
		Brief:         "show version information of current binary",
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			flfFont, err := figletlib.ReadFontFromBytes(fonts.SlantFontBytes)
			if err != nil {
				return
			}
			if err = termbox.Init(); err != nil {
				return
			}
			width, _ := termbox.Size()
			termbox.Close()
			figletlib.PrintMsg(strings.ToUpper(consts.ProjName), flfFont, width, flfFont.Settings(), "left")
			fmt.Println(consts.Description)
			return
		},
	}
)
