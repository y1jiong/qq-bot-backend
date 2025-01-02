package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/lukesampson/figlet/figletlib"
	"golang.org/x/term"
	"os"
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
			defer fmt.Println(consts.Description)
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				return
			}
			flfFont, err := figletlib.ReadFontFromBytes(fonts.SlantFontBytes)
			if err != nil {
				return
			}
			figletlib.PrintMsg(strings.ToUpper(consts.ProjName), flfFont, width, flfFont.Settings(), "left")
			return
		},
	}
)
