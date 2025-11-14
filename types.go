package main

// HTTP API 响应类型

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

// MCP 相关类型（用于内部转换）

// MCPToolResult MCP 工具结果（内部使用）
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent MCP 内容（内部使用）
type MCPContent struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// FeedDetailRequest Feed详情请求
type FeedDetailRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
}

// FeedDetailResponse Feed详情响应
type FeedDetailResponse struct {
	FeedID string `json:"feed_id"`
	Data   any    `json:"data"`
}

// PostCommentRequest 发表评论请求
type PostCommentRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// PostCommentResponse 发表评论响应
type PostCommentResponse struct {
	FeedID  string `json:"feed_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserProfileRequest 用户主页请求
type UserProfileRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
}

// ActionResult 通用动作响应（点赞/收藏等）
type ActionResult struct {
	FeedID  string `json:"feed_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TextToImageRequest 文生图请求
type TextToImageRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// TextToImageResponse 文生图响应
type TextToImageResponse struct {
	Success    bool     `json:"success"`
	ImageURLs  []string `json:"image_urls"`
	ImagePaths []string `json:"image_paths"` // 本地文件路径
	Message    string   `json:"message,omitempty"`
	Error      string   `json:"error,omitempty"`
}

// ImageToImageRequest 图生图请求
type ImageToImageRequest struct {
	Prompt    string  `json:"prompt" binding:"required"`
	ImagePath string  `json:"image_path" binding:"required"`
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	Strength  float64 `json:"strength,omitempty"`
}

// ImageToImageResponse 图生图响应
type ImageToImageResponse struct {
	Success   bool     `json:"success"`
	ImageURLs []string `json:"image_urls"`
	Message   string   `json:"message,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// CoverImageRequest 封面图片生成请求
type CoverImageRequest struct {
	Text            string `json:"text" binding:"required"`    // 要显示的文字
	Width           int    `json:"width,omitempty"`            // 图片宽度，默认800
	Height          int    `json:"height,omitempty"`           // 图片高度，默认600
	FontSize        int    `json:"font_size,omitempty"`        // 字体大小，默认48
	TextColor       string `json:"text_color,omitempty"`       // 文字颜色，默认白色
	BgColor         string `json:"bg_color,omitempty"`         // 背景颜色，默认随机渐变
	Style           string `json:"style,omitempty"`            // 样式：gradient, solid, pattern
	BackgroundImage string `json:"background_image,omitempty"` // 背景图片路径，如果设置则使用背景图替代纯色或渐变
	TextOffsetY     int    `json:"text_offset_y,omitempty"`    // 文字垂直偏移值（像素），默认0（居中），正值向下，负值向上
	OutputPath      string `json:"output_path,omitempty"`      // 输出路径，默认自动生成
}

// CoverImageResponse 封面图片生成响应
type CoverImageResponse struct {
	Success   bool   `json:"success"`
	ImagePath string `json:"image_path,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}
