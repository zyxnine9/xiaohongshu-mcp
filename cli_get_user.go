package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var userID string

var getUserCmd = &cobra.Command{
	Use:   "get-user",
	Short: "获取用户主页",
	Example: `  xhs get-user --user-id "xxx"
  xhs get-user --user-id "xxx" --token "xxx"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.UserProfile(ctx, userID, xsecToken)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("获取成功")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	getUserCmd.Flags().StringVar(&userID, "user-id", "", "用户 ID (必填)")
	getUserCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")

	rootCmd.AddCommand(getUserCmd)
}
