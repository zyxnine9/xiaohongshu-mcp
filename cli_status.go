package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "检查登录状态",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		status, err := service.CheckLoginStatus(ctx)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if status.IsLoggedIn {
			output.NewOutput(jsonOutput).Success("已登录: " + status.Username)
		} else {
			output.NewOutput(jsonOutput).Success("未登录，请先执行 login 命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
