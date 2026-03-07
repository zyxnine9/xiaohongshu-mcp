package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	feedID    string
	xsecToken string
)

var getFeedCmd = &cobra.Command{
	Use:   "get-feed",
	Short: "获取笔记详情",
	Example: `  xhs get-feed --id "xxx"
  xhs get-feed --id "xxx" --token "xxx"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.GetFeedDetail(ctx, feedID, xsecToken, false)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("获取成功")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	getFeedCmd.Flags().StringVar(&feedID, "id", "", "笔记 ID (必填)")
	getFeedCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")

	rootCmd.AddCommand(getFeedCmd)
}
