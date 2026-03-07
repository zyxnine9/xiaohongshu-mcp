package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var unlike bool

var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "点赞/取消点赞笔记",
	Example: `  xhs like --feed-id "xxx"
  xhs like --feed-id "xxx" --unlike`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		var result interface{}
		var err error

		if unlike {
			result, err = service.UnlikeFeed(ctx, feedID, xsecToken)
		} else {
			result, err = service.LikeFeed(ctx, feedID, xsecToken)
		}

		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if unlike {
			output.NewOutput(jsonOutput).Success("取消点赞成功")
		} else {
			output.NewOutput(jsonOutput).Success("点赞成功")
		}
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	likeCmd.Flags().StringVar(&feedID, "feed-id", "", "笔记 ID (必填)")
	likeCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")
	likeCmd.Flags().BoolVar(&unlike, "unlike", false, "取消点赞")

	rootCmd.AddCommand(likeCmd)
}
