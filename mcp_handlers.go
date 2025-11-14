package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/volcengine/volc-sdk-golang/service/visual"
)

// MCP 工具处理函数

// getLocalIP 获取本机IP地址
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// convertToHTTPURL 将本地文件路径转换为HTTP URL
func convertToHTTPURL(localPath string, port string) string {
	// 提取文件名
	filename := filepath.Base(localPath)
	// 获取本机IP
	ip := getLocalIP()

	// 处理端口参数，提取端口号
	var portOnly string
	if strings.Contains(port, ":") {
		// 如果包含冒号，提取端口号部分
		portOnly = port[strings.LastIndex(port, ":"):]
	} else {
		// 如果不包含冒号，直接使用
		portOnly = port
	}

	// 构建HTTP URL
	return fmt.Sprintf("http://%s%s/images/%s", ip, portOnly, filename)
}

// handleCheckLoginStatus 处理检查登录状态
func (s *AppServer) handleCheckLoginStatus(ctx context.Context) *MCPToolResult {
	logrus.Info("MCP: 检查登录状态")

	status, err := s.xiaohongshuService.CheckLoginStatus(ctx)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "检查登录状态失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	resultText := fmt.Sprintf("登录状态检查成功: %+v", status)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: resultText,
		}},
	}
}

// handleGetLoginQrcode 处理获取登录二维码请求。
// 返回二维码图片的 Base64 编码和超时时间，供前端展示扫码登录。
func (s *AppServer) handleGetLoginQrcode(ctx context.Context) *MCPToolResult {
	logrus.Info("MCP: 获取登录扫码图片")

	result, err := s.xiaohongshuService.GetLoginQrcode(ctx)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{Type: "text", Text: "获取登录扫码图片失败: " + err.Error()}},
			IsError: true,
		}
	}

	if result.IsLoggedIn {
		return &MCPToolResult{
			Content: []MCPContent{{Type: "text", Text: "你当前已处于登录状态"}},
		}
	}

	now := time.Now()
	deadline := func() string {
		d, err := time.ParseDuration(result.Timeout)
		if err != nil {
			return now.Format("2006-01-02 15:04:05")
		}
		return now.Add(d).Format("2006-01-02 15:04:05")
	}()

	// 已登录：文本 + 图片
	contents := []MCPContent{
		{Type: "text", Text: "请用小红书 App 在 " + deadline + " 前扫码登录 👇"},
		{
			Type:     "image",
			MimeType: "image/png",
			Data:     strings.TrimPrefix(result.Img, "data:image/png;base64,"),
		},
	}
	return &MCPToolResult{Content: contents}
}

// handlePublishContent 处理发布内容
func (s *AppServer) handlePublishContent(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 发布内容")

	// 解析参数
	title, _ := args["title"].(string)
	content, _ := args["content"].(string)
	imagePathsInterface, _ := args["images"].([]interface{})
	tagsInterface, _ := args["tags"].([]interface{})

	var imagePaths []string
	for _, path := range imagePathsInterface {
		if pathStr, ok := path.(string); ok {
			imagePaths = append(imagePaths, pathStr)
		}
	}

	var tags []string
	for _, tag := range tagsInterface {
		if tagStr, ok := tag.(string); ok {
			tags = append(tags, tagStr)
		}
	}

	logrus.Infof("MCP: 发布内容 - 标题: %s, 图片数量: %d, 标签数量: %d", title, len(imagePaths), len(tags))

	// 构建发布请求
	req := &PublishRequest{
		Title:   title,
		Content: content,
		Images:  imagePaths,
		Tags:    tags,
	}

	// 执行发布
	result, err := s.xiaohongshuService.PublishContent(ctx, req)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发布失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	resultText := fmt.Sprintf("内容发布成功: %+v", result)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: resultText,
		}},
	}
}

// handlePublishVideo 处理发布视频内容（仅本地单个视频文件）
func (s *AppServer) handlePublishVideo(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 发布视频内容（本地）")

	title, _ := args["title"].(string)
	content, _ := args["content"].(string)
	videoPath, _ := args["video"].(string)
	tagsInterface, _ := args["tags"].([]interface{})

	var tags []string
	for _, tag := range tagsInterface {
		if tagStr, ok := tag.(string); ok {
			tags = append(tags, tagStr)
		}
	}

	if videoPath == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发布失败: 缺少本地视频文件路径",
			}},
			IsError: true,
		}
	}

	logrus.Infof("MCP: 发布视频 - 标题: %s, 标签数量: %d", title, len(tags))

	// 构建发布请求
	req := &PublishVideoRequest{
		Title:   title,
		Content: content,
		Video:   videoPath,
		Tags:    tags,
	}

	// 执行发布
	result, err := s.xiaohongshuService.PublishVideo(ctx, req)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发布失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	resultText := fmt.Sprintf("视频发布成功: %+v", result)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: resultText,
		}},
	}
}

// handleListFeeds 处理获取Feeds列表
func (s *AppServer) handleListFeeds(ctx context.Context) *MCPToolResult {
	logrus.Info("MCP: 获取Feeds列表")

	result, err := s.xiaohongshuService.ListFeeds(ctx)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取Feeds列表失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 格式化输出，转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("获取Feeds列表成功，但序列化失败: %v", err),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}

// handleSearchFeeds 处理搜索Feeds
func (s *AppServer) handleSearchFeeds(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 搜索Feeds")

	// 解析参数
	keyword, ok := args["keyword"].(string)
	if !ok || keyword == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "搜索Feeds失败: 缺少关键词参数",
			}},
			IsError: true,
		}
	}

	logrus.Infof("MCP: 搜索Feeds - 关键词: %s", keyword)

	result, err := s.xiaohongshuService.SearchFeeds(ctx, keyword)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "搜索Feeds失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 格式化输出，转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("搜索Feeds成功，但序列化失败: %v", err),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}

// handleGetFeedDetail 处理获取Feed详情
func (s *AppServer) handleGetFeedDetail(ctx context.Context, args map[string]any) *MCPToolResult {
	logrus.Info("MCP: 获取Feed详情")

	// 解析参数
	feedID, ok := args["feed_id"].(string)
	if !ok || feedID == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取Feed详情失败: 缺少feed_id参数",
			}},
			IsError: true,
		}
	}

	xsecToken, ok := args["xsec_token"].(string)
	if !ok || xsecToken == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取Feed详情失败: 缺少xsec_token参数",
			}},
			IsError: true,
		}
	}

	logrus.Infof("MCP: 获取Feed详情 - Feed ID: %s", feedID)

	result, err := s.xiaohongshuService.GetFeedDetail(ctx, feedID, xsecToken)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取Feed详情失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 格式化输出，转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("获取Feed详情成功，但序列化失败: %v", err),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}

// handleUserProfile 获取用户主页
func (s *AppServer) handleUserProfile(ctx context.Context, args map[string]any) *MCPToolResult {
	logrus.Info("MCP: 获取用户主页")

	// 解析参数
	userID, ok := args["user_id"].(string)
	if !ok || userID == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取用户主页失败: 缺少user_id参数",
			}},
			IsError: true,
		}
	}

	xsecToken, ok := args["xsec_token"].(string)
	if !ok || xsecToken == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取用户主页失败: 缺少xsec_token参数",
			}},
			IsError: true,
		}
	}

	logrus.Infof("MCP: 获取用户主页 - User ID: %s", userID)

	result, err := s.xiaohongshuService.UserProfile(ctx, userID, xsecToken)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取用户主页失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 格式化输出，转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("获取用户主页，但序列化失败: %v", err),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}

// handleDownloadImages 处理下载并保存图片
func (s *AppServer) handleDownloadImages(ctx context.Context, args map[string]any) *MCPToolResult {
	logrus.Info("MCP: 下载并保存图片")

	// 参数解析：支持 images 数组，或 image 单个字符串
	var images []string

	if imgs, ok := args["images"].([]interface{}); ok {
		for _, v := range imgs {
			if str, ok := v.(string); ok && str != "" {
				images = append(images, str)
			}
		}
	}
	if single, ok := args["image"].(string); ok && single != "" {
		images = append(images, single)
	}

	if len(images) == 0 {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "下载失败: 需要提供 images 数组或 image 字符串"}}, IsError: true}
	}

	// 解析保存目录参数
	saveDir, _ := args["save_dir"].(string)

	// 执行下载
	res, err := s.xiaohongshuService.DownloadImages(ctx, images, saveDir)
	if err != nil {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "下载失败: " + err.Error()}}, IsError: true}
	}

	// 返回JSON结果
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "下载成功，但结果序列化失败: " + err.Error()}}, IsError: true}
	}
	return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: string(data)}}}
}

// handleLikeFeed 处理点赞/取消点赞
func (s *AppServer) handleLikeFeed(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	feedID, ok := args["feed_id"].(string)
	if !ok || feedID == "" {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "操作失败: 缺少feed_id参数"}}, IsError: true}
	}
	xsecToken, ok := args["xsec_token"].(string)
	if !ok || xsecToken == "" {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "操作失败: 缺少xsec_token参数"}}, IsError: true}
	}
	unlike, _ := args["unlike"].(bool)

	var res *ActionResult
	var err error

	if unlike {
		res, err = s.xiaohongshuService.UnlikeFeed(ctx, feedID, xsecToken)
	} else {
		res, err = s.xiaohongshuService.LikeFeed(ctx, feedID, xsecToken)
	}

	if err != nil {
		action := "点赞"
		if unlike {
			action = "取消点赞"
		}
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: action + "失败: " + err.Error()}}, IsError: true}
	}

	action := "点赞"
	if unlike {
		action = "取消点赞"
	}
	return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: fmt.Sprintf("%s成功 - Feed ID: %s", action, res.FeedID)}}}
}

// handleFavoriteFeed 处理收藏/取消收藏
func (s *AppServer) handleFavoriteFeed(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	feedID, ok := args["feed_id"].(string)
	if !ok || feedID == "" {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "操作失败: 缺少feed_id参数"}}, IsError: true}
	}
	xsecToken, ok := args["xsec_token"].(string)
	if !ok || xsecToken == "" {
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: "操作失败: 缺少xsec_token参数"}}, IsError: true}
	}
	unfavorite, _ := args["unfavorite"].(bool)

	var res *ActionResult
	var err error

	if unfavorite {
		res, err = s.xiaohongshuService.UnfavoriteFeed(ctx, feedID, xsecToken)
	} else {
		res, err = s.xiaohongshuService.FavoriteFeed(ctx, feedID, xsecToken)
	}

	if err != nil {
		action := "收藏"
		if unfavorite {
			action = "取消收藏"
		}
		return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: action + "失败: " + err.Error()}}, IsError: true}
	}

	action := "收藏"
	if unfavorite {
		action = "取消收藏"
	}
	return &MCPToolResult{Content: []MCPContent{{Type: "text", Text: fmt.Sprintf("%s成功 - Feed ID: %s", action, res.FeedID)}}}
}

// handlePostComment 处理发表评论到Feed
func (s *AppServer) handlePostComment(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 发表评论到Feed")

	// 解析参数
	feedID, ok := args["feed_id"].(string)
	if !ok || feedID == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发表评论失败: 缺少feed_id参数",
			}},
			IsError: true,
		}
	}

	xsecToken, ok := args["xsec_token"].(string)
	if !ok || xsecToken == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发表评论失败: 缺少xsec_token参数",
			}},
			IsError: true,
		}
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发表评论失败: 缺少content参数",
			}},
			IsError: true,
		}
	}

	logrus.Infof("MCP: 发表评论 - Feed ID: %s, 内容长度: %d", feedID, len(content))

	// 发表评论
	result, err := s.xiaohongshuService.PostCommentToFeed(ctx, feedID, xsecToken, content)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发表评论失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 返回成功结果，只包含feed_id
	resultText := fmt.Sprintf("评论发表成功 - Feed ID: %s", result.FeedID)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: resultText,
		}},
	}
}

// handleTextToImage 处理文生图请求
func (s *AppServer) handleTextToImage(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 文生图请求")

	// 解析参数
	prompt, ok := args["prompt"].(string)
	if !ok || prompt == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: 缺少prompt参数",
			}},
			IsError: true,
		}
	}

	// 解析可选参数
	width := 512 // 默认值
	if w, ok := args["width"].(int); ok {
		width = w
	} else if w, ok := args["width"].(float64); ok {
		width = int(w)
	}

	height := 512 // 默认值
	if h, ok := args["height"].(int); ok {
		height = h
	} else if h, ok := args["height"].(float64); ok {
		height = int(h)
	}

	// 参数范围检查
	if width < 256 || width > 768 {
		width = 512
	}
	if height < 256 || height > 768 {
		height = 512
	}

	logrus.Infof("MCP: 文生图 - 提示词: %s, 尺寸: %dx%d", prompt, width, height)

	// 设置访问密钥
	visual.DefaultInstance.Client.SetAccessKey("AKLTMTU5ZGZjMWRkYzBjNGNhMWFlYTBiNzU0MmFhNWM3NjA")
	visual.DefaultInstance.Client.SetSecretKey("WVRKaE9UWTNaRFUyTm1Oa05EaGtPV0k0WldJMU9HWm1ZMlpsTlRBelpXSQ==")
	// 设置区域为即梦服务区域
	visual.DefaultInstance.SetRegion("cn-north-1")

	// 构建请求参数
	req := map[string]interface{}{
		"req_key": "jimeng_high_aes_general_v21_L",
		"prompt":  prompt,
		"width":   width,
		"height":  height,
	}

	// 调用文生图API
	resp, statusCode, err := visual.DefaultInstance.CVProcess(req)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 检查HTTP状态码
	if statusCode != 200 {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("文生图失败: HTTP状态码 %d", statusCode),
			}},
			IsError: true,
		}
	}

	// resp已经是map[string]interface{}类型，直接使用
	result := resp

	// 检查响应状态
	code, ok := result["code"].(float64)
	if !ok {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: 响应格式错误，缺少code字段",
			}},
			IsError: true,
		}
	}

	if code != 10000 {
		message, _ := result["message"].(string)
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("文生图失败: code=%v, message=%s", code, message),
			}},
			IsError: true,
		}
	}

	// 提取图片URL
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: 响应格式错误，缺少data字段",
			}},
			IsError: true,
		}
	}

	binaryDataInterface, ok := data["binary_data_base64"].([]interface{})
	if !ok {
		// 打印data的数据类型
		logrus.Errorf("文生图失败: 响应格式错误，缺少binary_data_base64字段")
		logrus.Errorf("data类型: %T", data)

		// 打印所有key
		var keys []string
		for key := range data {
			keys = append(keys, key)
		}
		logrus.Errorf("data包含的key: %v", keys)

		// 打印binary_data_base64字段的具体内容
		logrus.Errorf("data[\"binary_data_base64\"]类型: %T, 值: %+v", data["binary_data_base64"], data["binary_data_base64"])

		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: 响应格式错误，缺少binary_data_base64字段",
			}},
			IsError: true,
		}
	}

	// 创建存储目录
	saveDir := "generated_images"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: 无法创建存储目录 - " + err.Error(),
			}},
			IsError: true,
		}
	}

	var savedImagePaths []string

	for i, base64Data := range binaryDataInterface {
		base64Str, ok := base64Data.(string)
		if !ok {
			logrus.Errorf("跳过无效的Base64数据，索引: %d", i)
			continue
		}

		// 解码Base64数据
		imageData, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			logrus.Errorf("Base64解码失败，索引: %d, 错误: %v", i, err)
			continue
		}

		// 生成唯一文件名
		filename := generateUniqueFileNameGlobal("generated_image", "jpg")
		filePath := filepath.Join(saveDir, filename)

		// 保存图片到本地
		if err := os.WriteFile(filePath, imageData, 0644); err != nil {
			logrus.Errorf("保存图片失败，路径: %s, 错误: %v", filePath, err)
			continue
		}

		// 将路径转换为绝对路径，确保其他工具可以正确读取
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			logrus.Errorf("转换绝对路径失败，路径: %s, 错误: %v", filePath, err)
			// 如果转换失败，仍然使用相对路径
			savedImagePaths = append(savedImagePaths, filePath)
			logrus.Infof("图片已保存: %s (使用相对路径)", filePath)
		} else {
			savedImagePaths = append(savedImagePaths, absPath)
			logrus.Infof("图片已保存: %s (绝对路径: %s)", filePath, absPath)
		}
	}

	if len(savedImagePaths) == 0 {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图失败: 未成功保存任何图片",
			}},
			IsError: true,
		}
	}

	// 将本地文件路径转换为HTTP URL
	var httpURLs []string
	// 获取服务器端口，默认使用18060
	port := ":18060"
	if s.httpServer != nil && s.httpServer.Addr != "" {
		port = s.httpServer.Addr
	}

	for _, localPath := range savedImagePaths {
		httpURL := convertToHTTPURL(localPath, port)
		httpURLs = append(httpURLs, httpURL)
		logrus.Infof("图片HTTP访问地址: %s", httpURL)
	}

	// 构建成功响应
	response := &TextToImageResponse{
		Success:    true,
		ImageURLs:  httpURLs,
		ImagePaths: savedImagePaths, // 添加本地文件路径
		Message:    fmt.Sprintf("成功生成并保存 %d 张图片，可通过HTTP访问", len(httpURLs)),
	}

	// 转换为JSON字符串
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "文生图成功，但结果序列化失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}

// handleImageToImage 处理图生图请求（包含创建任务和轮询状态）
func (s *AppServer) handleImageToImage(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 图生图请求")

	// 解析参数
	prompt, ok := args["prompt"].(string)
	if !ok || prompt == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 缺少prompt参数",
			}},
			IsError: true,
		}
	}

	imagePath, ok := args["image_path"].(string)
	if !ok || imagePath == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 缺少image_path参数",
			}},
			IsError: true,
		}
	}

	// 解析可选参数
	width := 512 // 默认值
	if w, ok := args["width"].(int); ok {
		width = w
	} else if w, ok := args["width"].(float64); ok {
		width = int(w)
	}

	height := 512 // 默认值
	if h, ok := args["height"].(int); ok {
		height = h
	} else if h, ok := args["height"].(float64); ok {
		height = int(h)
	}

	strength := 0.8 // 默认值
	if s, ok := args["strength"].(float64); ok {
		strength = s
	}

	// 参数范围检查
	if width < 256 || width > 768 {
		width = 512
	}
	if height < 256 || height > 768 {
		height = 512
	}
	if strength < 0.1 || strength > 1.0 {
		strength = 0.8
	}

	logrus.Infof("MCP: 图生图 - 提示词: %s, 图片路径: %s, 尺寸: %dx%d, 强度: %.2f", prompt, imagePath, width, height, strength)

	// 检查图片文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 图片文件不存在 - " + imagePath,
			}},
			IsError: true,
		}
	}

	// 读取图片文件并转换为base64
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 无法读取图片文件 - " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 将图片数据转换为base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// 设置访问密钥
	visual.DefaultInstance.Client.SetAccessKey("AKLTMTU5ZGZjMWRkYzBjNGNhMWFlYTBiNzU0MmFhNWM3NjA")
	visual.DefaultInstance.Client.SetSecretKey("WVRKaE9UWTNaRFUyTm1Oa05EaGtPV0k0WldJMU9HWm1ZMlpsTlRBelpXSQ==")
	// 设置区域为即梦服务区域
	visual.DefaultInstance.SetRegion("cn-north-1")

	// 步骤1: 创建图生图任务
	logrus.Info("创建图生图任务...")
	createReq := map[string]interface{}{
		"req_key":            "jimeng_i2i_v30", // 图生图的req_key
		"prompt":             prompt,
		"width":              width,
		"height":             height,
		"strength":           strength,
		"binary_data_base64": []string{imageBase64}, // 上传的参考图片base64
	}

	// 调用创建图生图任务API
	createResp, statusCode, err := visual.DefaultInstance.CVSync2AsyncSubmitTask(createReq)
	if err != nil {
		logrus.Errorf("创建图生图任务API调用失败: %v", err)
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 创建任务失败 - " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 检查HTTP状态码
	if statusCode != 200 {
		logrus.Errorf("创建图生图任务失败: HTTP状态码 %d", statusCode)
		logrus.Errorf("创建任务响应: %+v", createResp)
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("图生图失败: 创建任务失败，HTTP状态码 %d", statusCode),
			}},
			IsError: true,
		}
	}

	// 检查响应状态
	createResult := createResp
	code, ok := createResult["code"].(float64)
	if !ok {
		logrus.Errorf("创建图生图任务响应格式错误，缺少code字段")
		logrus.Errorf("创建任务完整响应: %+v", createResult)
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 创建任务响应格式错误，缺少code字段",
			}},
			IsError: true,
		}
	}

	if code != 10000 {
		message, _ := createResult["message"].(string)
		logrus.Errorf("创建图生图任务失败: code=%v, message=%s", code, message)
		logrus.Errorf("创建任务完整响应: %+v", createResult)
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("图生图失败: 创建任务失败，code=%v, message=%s", code, message),
			}},
			IsError: true,
		}
	}

	// 提取任务数据
	createData, ok := createResult["data"].(map[string]interface{})
	if !ok {
		logrus.Errorf("创建图生图任务响应格式错误，缺少data字段")
		logrus.Errorf("创建任务完整响应: %+v", createResult)
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 创建任务响应格式错误，缺少data字段",
			}},
			IsError: true,
		}
	}

	// 检查是否直接返回了binary_data_base64（同步返回）
	if binaryDataInterface, ok := createData["binary_data_base64"].([]interface{}); ok && len(binaryDataInterface) > 0 {
		logrus.Info("图生图任务同步完成，直接返回了图片数据，开始处理...")

		// 创建存储目录
		saveDir := "generated_images"
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			return &MCPToolResult{
				Content: []MCPContent{{
					Type: "text",
					Text: "图生图失败: 无法创建存储目录 - " + err.Error(),
				}},
				IsError: true,
			}
		}

		var savedImagePaths []string

		for i, base64Data := range binaryDataInterface {
			base64Str, ok := base64Data.(string)
			if !ok {
				logrus.Errorf("跳过无效的Base64数据，索引: %d", i)
				continue
			}

			// 解码Base64数据
			imageData, err := base64.StdEncoding.DecodeString(base64Str)
			if err != nil {
				logrus.Errorf("Base64解码失败，索引: %d, 错误: %v", i, err)
				continue
			}

			// 生成唯一文件名
			filename := generateUniqueFileNameGlobal("img2img", "jpg")
			filePath := filepath.Join(saveDir, filename)

			// 保存图片到本地
			if err := os.WriteFile(filePath, imageData, 0644); err != nil {
				logrus.Errorf("保存图片失败，路径: %s, 错误: %v", filePath, err)
				continue
			}

			savedImagePaths = append(savedImagePaths, filePath)
			logrus.Infof("图片已保存: %s", filePath)
		}

		if len(savedImagePaths) == 0 {
			return &MCPToolResult{
				Content: []MCPContent{{
					Type: "text",
					Text: "图生图失败: 未成功保存任何图片",
				}},
				IsError: true,
			}
		}

		// 将本地文件路径转换为HTTP URL
		var httpURLs []string
		// 获取服务器端口，默认使用18060
		port := ":18060"
		if s.httpServer != nil && s.httpServer.Addr != "" {
			port = s.httpServer.Addr
		}

		for _, localPath := range savedImagePaths {
			httpURL := convertToHTTPURL(localPath, port)
			httpURLs = append(httpURLs, httpURL)
			logrus.Infof("图片HTTP访问地址: %s", httpURL)
		}

		// 构建成功响应
		response := &ImageToImageResponse{
			Success:   true,
			ImageURLs: httpURLs,
			Message:   fmt.Sprintf("图生图完成，成功生成并保存 %d 张图片，可通过HTTP访问", len(httpURLs)),
		}

		// 转换为JSON字符串
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return &MCPToolResult{
				Content: []MCPContent{{
					Type: "text",
					Text: "图生图成功，但结果序列化失败: " + err.Error(),
				}},
				IsError: true,
			}
		}

		logrus.Infof("图生图同步完成，共生成 %d 张图片", len(httpURLs))
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: string(jsonData),
			}},
		}
	}

	// 如果没有直接返回图片数据，则提取任务ID进行异步处理
	taskID, ok := createData["task_id"].(string)
	if !ok {
		logrus.Errorf("图生图失败: 创建任务响应格式错误，缺少task_id字段")
		logrus.Errorf("data类型: %T", createData)
		var keys []string
		for key := range createData {
			keys = append(keys, key)
		}
		logrus.Errorf("data包含的key: %v", keys)
		logrus.Errorf("data[\"task_id\"]类型: %T, 值: %+v", createData["task_id"], createData["task_id"])

		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "图生图失败: 创建任务响应格式错误，缺少task_id字段",
			}},
			IsError: true,
		}
	}

	logrus.Infof("图生图任务创建成功，任务ID: %s，开始轮询状态...", taskID)

	// 步骤2: 轮询任务状态直到完成
	maxRetries := 60                 // 最多轮询60次
	retryInterval := 5 * time.Second // 每次间隔5秒
	maxWaitTime := 5 * time.Minute   // 最大等待5分钟

	startTime := time.Now()
	for i := 0; i < maxRetries; i++ {
		// 检查是否超时
		if time.Since(startTime) > maxWaitTime {
			return &MCPToolResult{
				Content: []MCPContent{{
					Type: "text",
					Text: fmt.Sprintf("图生图失败: 任务超时，等待时间超过 %v", maxWaitTime),
				}},
				IsError: true,
			}
		}

		logrus.Infof("查询任务状态，第 %d/%d 次，任务ID: %s", i+1, maxRetries, taskID)

		// 构建查询任务状态的请求参数
		statusReq := map[string]interface{}{
			"req_key": "jimeng_i2i_v30", // 图生图的req_key
			"task_id": taskID,
		}

		// 调用查询图生图任务状态API
		statusResp, statusCode, err := visual.DefaultInstance.CVSync2AsyncGetResult(statusReq)
		if err != nil {
			logrus.Errorf("查询任务状态API调用失败: %v", err)
			time.Sleep(retryInterval)
			continue
		}

		// 检查HTTP状态码
		if statusCode != 200 {
			logrus.Errorf("查询任务状态失败: HTTP状态码 %d", statusCode)
			logrus.Errorf("查询状态响应: %+v", statusResp)
			time.Sleep(retryInterval)
			continue
		}

		// 检查响应状态
		statusResult := statusResp
		code, ok := statusResult["code"].(float64)
		if !ok {
			logrus.Errorf("查询任务状态响应格式错误，缺少code字段")
			logrus.Errorf("查询状态完整响应: %+v", statusResult)
			time.Sleep(retryInterval)
			continue
		}

		if code != 10000 {
			message, _ := statusResult["message"].(string)
			logrus.Errorf("查询任务状态失败: code=%v, message=%s", code, message)
			logrus.Errorf("查询状态完整响应: %+v", statusResult)
			time.Sleep(retryInterval)
			continue
		}

		// 提取任务状态和结果
		statusData, ok := statusResult["data"].(map[string]interface{})
		if !ok {
			logrus.Errorf("查询任务状态响应格式错误，缺少data字段")
			logrus.Errorf("查询状态完整响应: %+v", statusResult)
			time.Sleep(retryInterval)
			continue
		}

		status, ok := statusData["status"].(string)
		if !ok {
			logrus.Errorf("查询任务状态响应格式错误，缺少status字段")
			logrus.Errorf("查询状态data字段: %+v", statusData)
			time.Sleep(retryInterval)
			continue
		}

		logrus.Infof("任务状态: %s", status)

		// 如果任务完成，处理生成的图片
		if status == "completed" {
			logrus.Infof("图生图任务完成，开始处理生成的图片，任务ID: %s", taskID)
			binaryDataInterface, ok := statusData["binary_data_base64"].([]interface{})
			if !ok || len(binaryDataInterface) == 0 {
				logrus.Errorf("任务完成但未返回图片数据，任务ID: %s", taskID)
				logrus.Errorf("任务完成时的data字段: %+v", statusData)
				return &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "图生图失败: 任务完成但未返回图片数据",
					}},
					IsError: true,
				}
			}

			// 创建存储目录
			saveDir := "generated_images"
			if err := os.MkdirAll(saveDir, 0755); err != nil {
				return &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "图生图失败: 无法创建存储目录 - " + err.Error(),
					}},
					IsError: true,
				}
			}

			var savedImagePaths []string
			for i, base64Data := range binaryDataInterface {
				base64Str, ok := base64Data.(string)
				if !ok {
					logrus.Errorf("跳过无效的Base64数据，索引: %d", i)
					continue
				}

				// 解码Base64数据
				imageData, err := base64.StdEncoding.DecodeString(base64Str)
				if err != nil {
					logrus.Errorf("Base64解码失败，索引: %d, 错误: %v", i, err)
					continue
				}

				// 生成唯一文件名
				filename := generateUniqueFileNameGlobal("img2img", "jpg")
				filePath := filepath.Join(saveDir, filename)

				// 保存图片到本地
				if err := os.WriteFile(filePath, imageData, 0644); err != nil {
					logrus.Errorf("保存图片失败，路径: %s, 错误: %v", filePath, err)
					continue
				}

				savedImagePaths = append(savedImagePaths, filePath)
				logrus.Infof("图片已保存: %s", filePath)
			}

			if len(savedImagePaths) == 0 {
				logrus.Errorf("图生图任务完成但未成功保存任何图片，任务ID: %s", taskID)
				logrus.Errorf("图片数据数量: %d", len(binaryDataInterface))
				return &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "图生图失败: 未成功保存任何图片",
					}},
					IsError: true,
				}
			}

			// 将本地文件路径转换为HTTP URL
			var httpURLs []string
			// 获取服务器端口，默认使用18060
			port := ":18060"
			if s.httpServer != nil && s.httpServer.Addr != "" {
				port = s.httpServer.Addr
			}

			for _, localPath := range savedImagePaths {
				httpURL := convertToHTTPURL(localPath, port)
				httpURLs = append(httpURLs, httpURL)
				logrus.Infof("图片HTTP访问地址: %s", httpURL)
			}

			// 构建成功响应
			response := &ImageToImageResponse{
				Success:   true,
				ImageURLs: httpURLs,
				Message:   fmt.Sprintf("图生图完成，成功生成并保存 %d 张图片，可通过HTTP访问", len(httpURLs)),
			}

			// 转换为JSON字符串
			jsonData, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				return &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: "图生图成功，但结果序列化失败: " + err.Error(),
					}},
					IsError: true,
				}
			}

			logrus.Infof("图生图任务完成，共生成 %d 张图片", len(httpURLs))
			return &MCPToolResult{
				Content: []MCPContent{{
					Type: "text",
					Text: string(jsonData),
				}},
			}
		}

		// 如果任务失败
		if status == "failed" {
			logrus.Errorf("图生图任务执行失败，任务ID: %s", taskID)
			logrus.Errorf("任务失败时的完整响应: %+v", statusResult)

			// 尝试获取失败原因
			if errorMsg, ok := statusData["error"].(string); ok && errorMsg != "" {
				logrus.Errorf("任务失败原因: %s", errorMsg)
				return &MCPToolResult{
					Content: []MCPContent{{
						Type: "text",
						Text: fmt.Sprintf("图生图失败: 任务执行失败 - %s", errorMsg),
					}},
					IsError: true,
				}
			}

			return &MCPToolResult{
				Content: []MCPContent{{
					Type: "text",
					Text: "图生图失败: 任务执行失败",
				}},
				IsError: true,
			}
		}

		// 如果任务还在处理中，等待后继续轮询
		if status == "pending" || status == "processing" {
			logrus.Infof("任务状态: %s，等待 %v 后继续查询...", status, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		// 未知状态，等待后继续轮询
		logrus.Warnf("未知任务状态: %s，等待 %v 后继续查询...", status, retryInterval)
		time.Sleep(retryInterval)
	}

	// 轮询超时
	logrus.Errorf("图生图任务轮询超时，任务ID: %s", taskID)
	logrus.Errorf("轮询统计: 已尝试 %d 次，总耗时 %v，最大等待时间 %v", maxRetries, time.Since(startTime), maxWaitTime)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: fmt.Sprintf("图生图失败: 任务轮询超时，已尝试 %d 次，总耗时 %v", maxRetries, time.Since(startTime)),
		}},
		IsError: true,
	}
}

// handleGenerateCoverImage 处理生成封面图片请求
func (s *AppServer) handleGenerateCoverImage(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 生成封面图片")

	// 解析参数
	text, ok := args["text"].(string)
	if !ok || text == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "生成封面图片失败: 缺少text参数",
			}},
			IsError: true,
		}
	}

	// 解析可选参数
	width := 1080 // 默认值
	if w, ok := args["width"].(int); ok {
		width = w
	} else if w, ok := args["width"].(float64); ok {
		width = int(w)
	}

	height := 1440 // 默认值
	if h, ok := args["height"].(int); ok {
		height = h
	} else if h, ok := args["height"].(float64); ok {
		height = int(h)
	}

	fontSize := 48 // 默认值
	logrus.Infof("MCP: 原始font_size参数: %v (类型: %T)", args["font_size"], args["font_size"])

	// 尝试解析为int类型
	if fs, ok := args["font_size"].(int); ok {
		fontSize = fs
		logrus.Infof("MCP: 成功解析font_size (int): %d", fontSize)
	} else if fs, ok := args["font_size"].(float64); ok {
		fontSize = int(fs)
		logrus.Infof("MCP: 成功解析font_size (float64): %d", fontSize)
	} else {
		logrus.Warnf("MCP: font_size参数解析失败，使用默认值: %d", fontSize)
	}

	textColor, _ := args["text_color"].(string)
	bgColor, _ := args["bg_color"].(string)
	style, _ := args["style"].(string)
	backgroundImage, _ := args["background_image"].(string)
	outputPath, _ := args["output_path"].(string)

	// 解析文字垂直偏移值
	textOffsetY := 0 // 默认值
	if offset, ok := args["text_offset_y"].(int); ok {
		textOffsetY = offset
	} else if offset, ok := args["text_offset_y"].(float64); ok {
		textOffsetY = int(offset)
	}

	logrus.Infof("MCP: 生成封面图片 - 文字: %s, 尺寸: %dx%d, 字体大小: %d, 背景图: %s, 文字偏移: %d", text, width, height, fontSize, backgroundImage, textOffsetY)

	// 创建图片生成器
	imageGenerator := NewImageGenerator("assets")

	// 构建请求
	req := &CoverImageRequest{
		Text:            text,
		Width:           width,
		Height:          height,
		FontSize:        fontSize,
		TextColor:       textColor,
		BgColor:         bgColor,
		Style:           style,
		BackgroundImage: backgroundImage,
		TextOffsetY:     textOffsetY,
		OutputPath:      outputPath,
	}

	// 生成封面图片
	result, err := imageGenerator.GenerateCoverImage(req)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "生成封面图片失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 将本地文件路径转换为HTTP URL
	if result.ImagePath != "" {
		// 获取服务器端口，默认使用18060
		port := ":18060"
		if s.httpServer != nil && s.httpServer.Addr != "" {
			port = s.httpServer.Addr
		}
		result.ImageURL = convertToHTTPURL(result.ImagePath, port)
		logrus.Infof("封面图片HTTP访问地址: %s", result.ImageURL)
	}

	// 转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "生成封面图片成功，但结果序列化失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}
