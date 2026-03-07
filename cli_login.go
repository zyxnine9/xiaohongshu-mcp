package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "获取登录二维码",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		qrcode, err := service.GetLoginQrcode(ctx)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if qrcode.IsLoggedIn {
			output.NewOutput(jsonOutput).Success("已登录")
			return
		}

		output.NewOutput(jsonOutput).Success("二维码获取成功，请扫码登录")
		if qrcode.Img != "" {
			// 输出 base64 图片
			output.NewOutput(jsonOutput).Print(map[string]string{
				"qrcode": qrcode.Img,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
