package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	title    string
	content  string
	images   []string
	tags     []string
	schedule string
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "发布图文笔记",
	Example: `  xhs publish -t "标题" -c "正文" -i img1.jpg -i img2.jpg
  xhs publish -t "标题" -c "正文" -i img.jpg --json`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		req := &PublishRequest{
			Title:      title,
			Content:    content,
			Images:     images,
			Tags:       tags,
			ScheduleAt: schedule,
		}

		resp, err := service.PublishContent(ctx, req)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("发布成功")
		output.NewOutput(jsonOutput).Print(resp)
	},
}

func init() {
	publishCmd.Flags().StringVarP(&title, "title", "t", "", "笔记标题 (必填)")
	publishCmd.Flags().StringVarP(&content, "content", "c", "", "笔记正文 (必填)")
	publishCmd.Flags().StringSliceVarP(&images, "image", "i", []string{}, "图片路径，支持多个 (必填)")
	publishCmd.Flags().StringSliceVarP(&tags, "tag", "", []string{}, "标签，支持多个")
	publishCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "定时发布时间 (ISO8601格式)")

	rootCmd.AddCommand(publishCmd)
}
