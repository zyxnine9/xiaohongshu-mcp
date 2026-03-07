# CLI 改造实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task.

**Goal:** 将现有的小红书 MCP 服务器改造为同时支持 CLI 工具和 MCP 服务器。

**Architecture:** 使用 cobra 框架构建 CLI，无参数运行启动 MCP 服务器，使用子命令执行 CLI 操作。复用现有 XiaohongshuService 的所有业务逻辑。

**Tech Stack:** Go, cobra, sirupsen/logrus (已有)

---

## Task 1: 添加 cobra 依赖

**Step 1: 添加 cobra 依赖**

```bash
go get github.com/spf13/cobra
```

**Step 2: 提交**

```bash
git add go.mod go.sum
git commit -m "feat: 添加 cobra CLI 框架依赖"
```

---

## Task 2: 创建 CLI 命令基础结构

**Files:**
- Create: `cmd/root.go`
- Modify: `main.go`

**Step 1: 创建 cmd/root.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var RootCmd = &cobra.Command{
	Use:   "xhs",
	Short: "小红书 CLI 工具",
	Long:  `小红书 MCP 服务器兼 CLI 工具`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "o", false, "输出 JSON 格式")
}
```

**Step 2: 修改 main.go**

```go
package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/cmd"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

func main() {
	// 检查是否有子命令参数
	if len(os.Args) > 1 {
		// 有子命令，运行 CLI
		cmd.Execute()
		return
	}

	// 无参数，运行 MCP 服务器
	runServer()
}

var (
	headless bool
	binPath  string
	port     string
)

func init() {
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&port, "port", ":18060", "端口")
	flag.Parse()
}

func runServer() {
	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)

	xiaohongshuService := NewXiaohongshuService()
	appServer := NewAppServer(xiaohongshuService)
	if err := appServer.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
```

**Step 3: 提交**

```bash
git add cmd/root.go main.go
git commit -m "feat: 添加 CLI 基础结构"
```

---

## Task 3: 创建 output 工具包

**Files:**
- Create: `pkg/output/output.go`

**Step 1: 创建输出工具**

```go
package output

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

type Output interface {
	Success(msg string)
	Error(err error)
	Print(data interface{})
}

// HumanOutput 人类可读输出
type HumanOutput struct{}

func (h *HumanOutput) Success(msg string) {
	color.New(color.FgGreen).Println("✓", msg)
}

func (h *HumanOutput) Error(err error) {
	color.New(color.FgRed).Println("✗", err.Error())
}

func (h *HumanOutput) Print(data interface{}) {
	fmt.Printf("%#v\n", data)
}

// JSONOutput JSON 输出
type JSONOutput struct{}

func (j *JSONOutput) Success(msg string) {
	j.Print(map[string]string{"status": "success", "message": msg})
}

func (j *JSONOutput) Error(err error) {
	j.Print(map[string]string{"status": "error", "error": err.Error()})
}

func (j *JSONOutput) Print(data interface{}) {
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonBytes))
}

// NewOutput 创建输出器
func NewOutput(isJSON bool) Output {
	if isJSON {
		return &JSONOutput{}
	}
	return &HumanOutput{}
}
```

**Step 2: 提交**

```bash
git add pkg/output/output.go
git commit -m "feat: 添加输出工具包"
```

---

## Task 4: 创建 login 命令

**Files:**
- Create: `cmd/login.go`

**Step 1: 创建 login 命令**

```go
package cmd

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
	RootCmd.AddCommand(loginCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/login.go
git commit -m "feat: 添加 login 命令"
```

---

## Task 5: 创建 status 命令

**Files:**
- Create: `cmd/status.go`

**Step 1: 创建 status 命令**

```go
package cmd

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
	RootCmd.AddCommand(statusCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/status.go
git commit -m "feat: 添加 status 命令"
```

---

## Task 6: 创建 publish 命令

**Files:**
- Create: `cmd/publish.go`

**Step 1: 创建 publish 命令**

```go
package cmd

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

	RootCmd.AddCommand(publishCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/publish.go
git commit -m "feat: 添加 publish 命令"
```

---

## Task 7: 创建 publish-video 命令

**Files:**
- Create: `cmd/publish_video.go`

**Step 1: 创建 publish-video 命令**

```go
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	videoPath string
)

var publishVideoCmd = &cobra.Command{
	Use:   "publish-video",
	Short: "发布视频笔记",
	Example: `  xhs publish-video -t "标题" -c "正文" -v video.mp4
  xhs publish-video -t "标题" -c "正文" -v video.mp4 --json`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		req := &PublishVideoRequest{
			Title:      title,
			Content:    content,
			Video:      videoPath,
			Tags:       tags,
			ScheduleAt: schedule,
		}

		resp, err := service.PublishVideo(ctx, req)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("发布成功")
		output.NewOutput(jsonOutput).Print(resp)
	},
}

func init() {
	publishVideoCmd.Flags().StringVarP(&title, "title", "t", "", "笔记标题 (必填)")
	publishVideoCmd.Flags().StringVarP(&content, "content", "c", "", "笔记正文 (必填)")
	publishVideoCmd.Flags().StringVarP(&videoPath, "video", "v", "", "视频文件路径 (必填)")
	publishVideoCmd.Flags().StringSliceVarP(&tags, "tag", "", []string{}, "标签，支持多个")
	publishVideoCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "定时发布时间 (ISO8601格式)")

	RootCmd.AddCommand(publishVideoCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/publish_video.go
git commit -m "feat: 添加 publish-video 命令"
```

---

## Task 8: 创建 search 命令

**Files:**
- Create: `cmd/search.go`

**Step 1: 创建 search 命令**

```go
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	keyword string
	limit   int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索笔记",
	Example: `  xhs search -k "美食"
  xhs search -k "旅游" -l 20`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.SearchFeeds(ctx, keyword)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if limit > 0 && len(result.Feeds) > limit {
			result.Feeds = result.Feeds[:limit]
		}

		output.NewOutput(jsonOutput).Success("搜索完成")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	searchCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "搜索关键词 (必填)")
	searchCmd.Flags().IntVarP(&limit, "limit", "l", 0, "返回结果数量限制")

	RootCmd.AddCommand(searchCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/search.go
git commit -m "feat: 添加 search 命令"
```

---

## Task 9: 创建 get-feed 命令

**Files:**
- Create: `cmd/get_feed.go`

**Step 1: 创建 get-feed 命令**

```go
package cmd

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

	RootCmd.AddCommand(getFeedCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/get_feed.go
git commit -m "feat: 添加 get-feed 命令"
```

---

## Task 10: 创建 get-user 命令

**Files:**
- Create: `cmd/get_user.go`

**Step 1: 创建 get-user 命令**

```go
package cmd

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

	RootCmd.AddCommand(getUserCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/get_user.go
git commit -m "feat: 添加 get-user 命令"
```

---

## Task 11: 创建 comment 命令

**Files:**
- Create: `cmd/comment.go`

**Step 1: 创建 comment 命令**

```go
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var commentContent string

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "发表评论",
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

	RootCmd.AddCommand(commentCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/comment.go
git commit -m "feat: 添加 comment 命令"
```

---

## Task 12: 创建 reply 命令

**Files:**
- Create: `cmd/reply.go`

**Step 1: 创建 reply 命令**

```go
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	commentID string
	userID    string
)

var replyCmd = &cobra.Command{
	Use:   "reply",
	Short: "回复评论",
	Example: `  xhs reply --feed-id "xxx" --comment-id "xxx" --user-id "xxx" -c "回复内容"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.ReplyCommentToFeed(ctx, feedID, xsecToken, commentID, userID, commentContent)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("回复成功")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	replyCmd.Flags().StringVar(&feedID, "feed-id", "", "笔记 ID (必填)")
	replyCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")
	replyCmd.Flags().StringVar(&commentID, "comment-id", "", "评论 ID (必填)")
	replyCmd.Flags().StringVar(&userID, "user-id", "", "被回复用户 ID (必填)")
	replyCmd.Flags().StringVarP(&commentContent, "content", "c", "", "回复内容 (必填)")

	RootCmd.AddCommand(replyCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/reply.go
git commit -m "feat: 添加 reply 命令"
```

---

## Task 13: 创建 like 命令

**Files:**
- Create: `cmd/like.go`

**Step 1: 创建 like 命令**

```go
package cmd

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
  xhs like --feed-id "xxx" --unlink`,
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

	RootCmd.AddCommand(likeCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/like.go
git commit -m "feat: 添加 like 命令"
```

---

## Task 14: 创建 favorite 命令

**Files:**
- Create: `cmd/favorite.go`

**Step 1: 创建 favorite 命令**

```go
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var unfavorite bool

var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "收藏/取消收藏笔记",
	Example: `  xhs favorite --feed-id "xxx"
  xhs favorite --feed-id "xxx" --unfavorite`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		var result interface{}
		var err error

		if unfavorite {
			result, err = service.UnfavoriteFeed(ctx, feedID, xsecToken)
		} else {
			result, err = service.FavoriteFeed(ctx, feedID, xsecToken)
		}

		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if unfavorite {
			output.NewOutput(jsonOutput).Success("取消收藏成功")
		} else {
			output.NewOutput(jsonOutput).Success("收藏成功")
		}
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	favoriteCmd.Flags().StringVar(&feedID, "feed-id", "", "笔记 ID (必填)")
	favoriteCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")
	favoriteCmd.Flags().BoolVar(&unfavorite, "unfavorite", false, "取消收藏")

	RootCmd.AddCommand(favoriteCmd)
}
```

**Step 2: 提交**

```bash
git add cmd/favorite.go
git commit -m "feat: 添加 favorite 命令"
```

---

## Task 15: 创建 server 命令

**Files:**
- Create: `cmd/server.go`

**Step 1: 创建 server 命令**

```go
package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

var (
	serverPort  string
	serverHeadless bool
	serverBin   string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 MCP 服务器",
	Example: `  xhs server
  xhs server --port 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.InitHeadless(serverHeadless)
		configs.SetBinPath(serverBin)

		xiaohongshuService := NewXiaohongshuService()
		appServer := NewAppServer(xiaohongshuService)
		if err := appServer.Start(serverPort); err != nil {
			logrus.Fatalf("failed to run server: %v", err)
		}
	},
}

func init() {
	serverCmd.Flags().StringVar(&serverPort, "port", ":18060", "服务端口")
	serverCmd.Flags().BoolVar(&serverHeadless, "headless", true, "是否无头模式")
	serverCmd.Flags().StringVar(&serverBin, "bin", "", "浏览器二进制路径")

	RootCmd.AddCommand(serverCmd)
}

// 为了避免循环导入，将这些函数声明移到单独文件
func NewXiaohongshuService() *XiaohongshuService {
	return newXiaohongshuServiceForCLI()
}

func newXiaohongshuServiceForCLI() *XiaohongshuService {
	return &XiaohongshuService{}
}
```

**Step 2: 创建 cmd/service_helper.go**

```go
package cmd

// 提供给 CLI 使用的服务包装
func newXiaohongshuServiceForCLI() *XiaohongshuService {
	return &XiaohongshuService{}
}
```

**Step 3: 提交**

```bash
git add cmd/server.go cmd/service_helper.go
git commit -m "feat: 添加 server 命令"
```

---

## Task 16: 格式化代码并验证构建

**Step 1: 格式化代码**

```bash
gofmt -w cmd/ pkg/output/
```

**Step 2: 构建测试**

```bash
go build -o xhs .
```

**Step 3: 测试 CLI**

```bash
./xhs --help
./xhs login --help
```

**Step 4: 提交**

```bash
git add .
git commit -m "chore: 格式化代码并验证构建"
```

---

**Plan complete and saved to `docs/plans/2026-03-07-cli-implementation.md`.**

两个执行选项：

1. **Subagent-Driven (本会话)** - 我为每个任务派遣新的子代理，任务间进行审查，快速迭代

2. **Parallel Session (单独会话)** - 在新会话中使用 executing-plans，分批执行并设置检查点

你选择哪种方式？
