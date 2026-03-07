package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var commentContent string

var commentCmd = &cobra.Command{
	Use:     "comment",
	Short:   "发表评论",
	Example: `  xhs comment --feed-id "xxx" -c "评论内容"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.PostCommentToFeed(ctx, feedID, xsecToken, commentContent)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("评论成功")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	commentCmd.Flags().StringVar(&feedID, "feed-id", "", "笔记 ID (必填)")
	commentCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")
	commentCmd.Flags().StringVarP(&commentContent, "content", "c", "", "评论内容 (必填)")

	rootCmd.AddCommand(commentCmd)
}
