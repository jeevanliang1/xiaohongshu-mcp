package main

import (
	"context"
	"encoding/base64"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
)

// MCP 工具参数结构体定义

// PublishContentArgs 发布内容的参数
type PublishContentArgs struct {
	Title   string   `json:"title" jsonschema:"内容标题（小红书限制：最多20个中文字或英文单词）"`
	Content string   `json:"content" jsonschema:"正文内容，不包含以#开头的标签内容，所有话题标签都用tags参数来生成和提供即可"`
	Images  []string `json:"images" jsonschema:"图片路径列表（至少需要1张图片）。支持两种方式：1. HTTP/HTTPS图片链接（自动下载）；2. 本地图片绝对路径（推荐，如:/Users/user/image.jpg）"`
	Tags    []string `json:"tags,omitempty" jsonschema:"话题标签列表（可选参数），如 [美食, 旅行, 生活]"`
}

// PublishVideoArgs 发布视频的参数（仅支持本地单个视频文件）
type PublishVideoArgs struct {
	Title   string   `json:"title" jsonschema:"内容标题（小红书限制：最多20个中文字或英文单词）"`
	Content string   `json:"content" jsonschema:"正文内容，不包含以#开头的标签内容，所有话题标签都用tags参数来生成和提供即可"`
	Video   string   `json:"video" jsonschema:"本地视频绝对路径（仅支持单个视频文件，如:/Users/user/video.mp4）"`
	Tags    []string `json:"tags,omitempty" jsonschema:"话题标签列表（可选参数），如 [美食, 旅行, 生活]"`
}

// SearchFeedsArgs 搜索内容的参数
type SearchFeedsArgs struct {
	Keyword string `json:"keyword" jsonschema:"搜索关键词"`
}

// FeedDetailArgs 获取Feed详情的参数
type FeedDetailArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"小红书笔记ID，从Feed列表获取"`
	XsecToken string `json:"xsec_token" jsonschema:"访问令牌，从Feed列表的xsecToken字段获取"`
}

// UserProfileArgs 获取用户主页的参数
type UserProfileArgs struct {
	UserID    string `json:"user_id" jsonschema:"小红书用户ID，从Feed列表获取"`
	XsecToken string `json:"xsec_token" jsonschema:"访问令牌，从Feed列表的xsecToken字段获取"`
}

// PostCommentArgs 发表评论的参数
type PostCommentArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"小红书笔记ID，从Feed列表获取"`
	XsecToken string `json:"xsec_token" jsonschema:"访问令牌，从Feed列表的xsecToken字段获取"`
	Content   string `json:"content" jsonschema:"评论内容"`
}

// LikeFeedArgs 点赞参数
type LikeFeedArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"小红书笔记ID，从Feed列表获取"`
	XsecToken string `json:"xsec_token" jsonschema:"访问令牌，从Feed列表的xsecToken字段获取"`
	Unlike    bool   `json:"unlike,omitempty" jsonschema:"是否取消点赞，true为取消点赞，false或未设置则为点赞"`
}

// FavoriteFeedArgs 收藏参数
type FavoriteFeedArgs struct {
	FeedID     string `json:"feed_id" jsonschema:"小红书笔记ID，从Feed列表获取"`
	XsecToken  string `json:"xsec_token" jsonschema:"访问令牌，从Feed列表的xsecToken字段获取"`
	Unfavorite bool   `json:"unfavorite,omitempty" jsonschema:"是否取消收藏，true为取消收藏，false或未设置则为收藏"`
}

// DownloadImagesArgs 下载图片参数
type DownloadImagesArgs struct {
	Images  []string `json:"images" jsonschema:"要下载的图片URL列表，支持http/https；也可传本地路径将直接返回"`
	SaveDir string   `json:"save_dir,omitempty" jsonschema:"可选参数，指定保存图片的文件夹路径。如果不传，则使用默认目录 image_file"`
}

// TextToImageArgs 文生图参数
type TextToImageArgs struct {
	Prompt string `json:"prompt" jsonschema:"用于生成图像的提示词，中英文均可输入"`
	Width  int    `json:"width,omitempty" jsonschema:"生成图像的宽度，默认值：512，取值范围：[256, 768]"`
	Height int    `json:"height,omitempty" jsonschema:"生成图像的高度，默认值：512，取值范围：[256, 768]"`
}

// ImageToImageArgs 图生图参数
type ImageToImageArgs struct {
	Prompt    string  `json:"prompt" jsonschema:"用于生成图像的提示词，中英文均可输入"`
	ImagePath string  `json:"image_path" jsonschema:"参考图片的本地文件路径（绝对路径），如：/Users/user/image.jpg"`
	Width     int     `json:"width,omitempty" jsonschema:"生成图像的宽度，默认值：512，取值范围：[256, 768]"`
	Height    int     `json:"height,omitempty" jsonschema:"生成图像的高度，默认值：512，取值范围：[256, 768]"`
	Strength  float64 `json:"strength,omitempty" jsonschema:"控制参考图片的影响强度，默认值：0.8，取值范围：[0.1, 1.0]，值越大参考图片影响越强"`
}

// GenerateCoverImageArgs 生成封面图片参数
type GenerateCoverImageArgs struct {
	Text            string `json:"text" jsonschema:"要显示在封面上的文字内容"`
	Width           int    `json:"width,omitempty" jsonschema:"图片宽度，默认值：1080"`
	Height          int    `json:"height,omitempty" jsonschema:"图片高度，默认值：1440"`
	FontSize        int    `json:"font_size,omitempty" jsonschema:"字体大小，默认值：48"`
	TextColor       string `json:"text_color,omitempty" jsonschema:"文字颜色，十六进制格式，如：#FFFFFF，默认白色"`
	BgColor         string `json:"bg_color,omitempty" jsonschema:"背景颜色，十六进制格式，如：#667eea，默认随机渐变"`
	Style           string `json:"style,omitempty" jsonschema:"背景样式：gradient（渐变，默认）、solid（纯色）、pattern（图案）"`
	BackgroundImage string `json:"background_image,omitempty" jsonschema:"背景图片路径，如果设置则使用背景图替代纯色或渐变，图片宽高将等比例缩放（最大1080）"`
	TextOffsetY     int    `json:"text_offset_y,omitempty" jsonschema:"文字垂直偏移值（像素），默认值：0（居中），正值向下偏移，负值向上偏移"`
	OutputPath      string `json:"output_path,omitempty" jsonschema:"输出文件路径，如不指定则自动生成"`
}

// InitMCPServer 初始化 MCP Server
func InitMCPServer(appServer *AppServer) *mcp.Server {
	// 创建 MCP Server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "xiaohongshu-mcp",
			Version: "2.0.0",
		},
		nil,
	)

	// 注册所有工具
	registerTools(server, appServer)

	logrus.Info("MCP Server initialized with official SDK")

	return server
}

// registerTools 注册所有 MCP 工具
func registerTools(server *mcp.Server, appServer *AppServer) {
	// 工具 1: 检查登录状态
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "check_login_status",
			Description: "检查小红书登录状态",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				result = appServer.handleCheckLoginStatus(ctx)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "检查登录状态失败: " + err.Error(),
					}},
					IsError: true,
				}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 2: 获取登录二维码
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_login_qrcode",
			Description: "获取登录二维码（返回 Base64 图片和超时时间）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleGetLoginQrcode(ctx)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 3: 发布内容
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_content",
			Description: "发布小红书图文内容",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args PublishContentArgs) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				// 转换参数格式到现有的 handler
				argsMap := map[string]interface{}{
					"title":   args.Title,
					"content": args.Content,
					"images":  convertStringsToInterfaces(args.Images),
					"tags":    convertStringsToInterfaces(args.Tags),
				}
				result = appServer.handlePublishContent(ctx, argsMap)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "发布内容失败: " + err.Error(),
					}},
					IsError: true,
				}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 4: 获取Feed列表
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "list_feeds",
			Description: "获取用户发布的内容列表",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleListFeeds(ctx)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 5: 搜索内容
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "search_feeds",
			Description: "搜索小红书内容（需要已登录）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args SearchFeedsArgs) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				argsMap := map[string]interface{}{
					"keyword": args.Keyword,
				}
				result = appServer.handleSearchFeeds(ctx, argsMap)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "搜索内容失败: " + err.Error(),
					}},
					IsError: true,
				}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 6: 获取Feed详情
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_feed_detail",
			Description: "获取小红书笔记详情，返回笔记内容、图片、作者信息、互动数据（点赞/收藏/分享数）及评论列表",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args FeedDetailArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
			}
			result := appServer.handleGetFeedDetail(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 7: 获取用户主页
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "user_profile",
			Description: "获取小红书用户主页，返回用户基本信息，关注、粉丝、获赞量及其笔记内容",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args UserProfileArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"user_id":    args.UserID,
				"xsec_token": args.XsecToken,
			}
			result := appServer.handleUserProfile(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 8: 发表评论
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "post_comment_to_feed",
			Description: "发表评论到小红书笔记",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args PostCommentArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
				"content":    args.Content,
			}
			result := appServer.handlePostComment(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 9: 发布视频（仅本地文件）
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_with_video",
			Description: "发布小红书视频内容（仅支持本地单个视频文件）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args PublishVideoArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"title":   args.Title,
				"content": args.Content,
				"video":   args.Video,
				"tags":    convertStringsToInterfaces(args.Tags),
			}
			result := appServer.handlePublishVideo(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 10: 点赞笔记
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "like_feed",
			Description: "为指定笔记点赞或取消点赞（如已点赞将跳过点赞，如未点赞将跳过取消点赞）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args LikeFeedArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
				"unlike":     args.Unlike,
			}
			result := appServer.handleLikeFeed(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 11: 收藏笔记
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "favorite_feed",
			Description: "收藏指定笔记或取消收藏（如已收藏将跳过收藏，如未收藏将跳过取消收藏）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args FavoriteFeedArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
				"unfavorite": args.Unfavorite,
			}
			result := appServer.handleFavoriteFeed(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 12: 下载并保存图片
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "download_images",
			Description: "下载并保存图片到本地（支持URL或本地路径混合输入，返回保存路径）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args DownloadImagesArgs) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				argsMap := map[string]any{
					"images":   convertStringsToInterfaces(args.Images),
					"save_dir": args.SaveDir,
				}
				result = appServer.handleDownloadImages(ctx, argsMap)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "下载失败: " + err.Error()}}, IsError: true}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 13: 文生图
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "text_to_image",
			Description: "使用即梦AI生成图片，根据文本提示词生成高质量图像",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args TextToImageArgs) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				argsMap := map[string]interface{}{
					"prompt": args.Prompt,
					"width":  args.Width,
					"height": args.Height,
				}
				result = appServer.handleTextToImage(ctx, argsMap)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "文生图失败: " + err.Error(),
					}},
					IsError: true,
				}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 14: 图生图
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "image_to_image",
			Description: "使用即梦AI进行图生图，基于参考图片和文本提示词生成新的图像（自动处理异步任务）",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args ImageToImageArgs) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				argsMap := map[string]interface{}{
					"prompt":     args.Prompt,
					"image_path": args.ImagePath,
					"width":      args.Width,
					"height":     args.Height,
					"strength":   args.Strength,
				}
				result = appServer.handleImageToImage(ctx, argsMap)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "图生图失败: " + err.Error(),
					}},
					IsError: true,
				}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	// 工具 15: 生成封面图片
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "generate_cover_image",
			Description: "生成带文字的封面图片，支持多种背景样式和自定义参数，适用于小红书封面、海报等场景",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, args GenerateCoverImageArgs) (*mcp.CallToolResult, any, error) {
			var result *MCPToolResult
			err := panicRecoveryWrapper(func() error {
				logrus.Infof("MCP工具调用: font_size=%d", args.FontSize)
				argsMap := map[string]interface{}{
					"text":             args.Text,
					"width":            args.Width,
					"height":           args.Height,
					"font_size":        args.FontSize,
					"text_color":       args.TextColor,
					"bg_color":         args.BgColor,
					"style":            args.Style,
					"background_image": args.BackgroundImage,
					"text_offset_y":    args.TextOffsetY,
					"output_path":      args.OutputPath,
				}
				logrus.Infof("MCP参数映射: %+v", argsMap)
				result = appServer.handleGenerateCoverImage(ctx, argsMap)
				return nil
			})
			if err != nil {
				result = &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "生成封面图片失败: " + err.Error(),
					}},
					IsError: true,
				}
			}
			return convertToMCPResult(result), nil, nil
		},
	)

	logrus.Infof("Registered %d MCP tools", 15)
}

// convertToMCPResult 将自定义的 MCPToolResult 转换为官方 SDK 的格式
func convertToMCPResult(result *MCPToolResult) *mcp.CallToolResult {
	var contents []mcp.Content
	for _, c := range result.Content {
		switch c.Type {
		case "text":
			contents = append(contents, &mcp.TextContent{Text: c.Text})
		case "image":
			// 解码 base64 字符串为 []byte
			imageData, err := base64.StdEncoding.DecodeString(c.Data)
			if err != nil {
				logrus.WithError(err).Error("Failed to decode base64 image data")
				// 如果解码失败，添加错误文本
				contents = append(contents, &mcp.TextContent{
					Text: "图片数据解码失败: " + err.Error(),
				})
			} else {
				contents = append(contents, &mcp.ImageContent{
					Data:     imageData,
					MIMEType: c.MimeType,
				})
			}
		}
	}

	return &mcp.CallToolResult{
		Content: contents,
		IsError: result.IsError,
	}
}

// convertStringsToInterfaces 辅助函数：将 []string 转换为 []interface{}
func convertStringsToInterfaces(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
