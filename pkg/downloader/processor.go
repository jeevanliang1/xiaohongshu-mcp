package downloader

import (
	"fmt"

	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

// ImageProcessor 图片处理器
type ImageProcessor struct {
	downloader *ImageDownloader
}

// NewImageProcessor 创建图片处理器
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		downloader: NewImageDownloader(configs.GetImagesPath()),
	}
}

// NewImageProcessorWithDir 创建使用指定目录的图片处理器
func NewImageProcessorWithDir(saveDir string) *ImageProcessor {
	return &ImageProcessor{
		downloader: NewImageDownloader(saveDir),
	}
}

// ProcessImages 处理图片列表，返回本地文件路径
// 支持三种输入格式：
// 1. URL格式 (http/https开头) - 自动下载到本地
// 2. 飞书 file_token - 使用 access_token 从飞书下载
// 3. 本地文件路径 - 直接使用
func (p *ImageProcessor) ProcessImages(images []string) ([]string, error) {
	return p.ProcessImagesWithFeishuToken(images, "")
}

// ProcessImagesWithFeishuToken 处理图片列表，支持飞书 file_token
// accessToken: 飞书访问令牌，仅在检测到 file_token 时使用
func (p *ImageProcessor) ProcessImagesWithFeishuToken(images []string, accessToken string) ([]string, error) {
	var localPaths []string
	var urlsToDownload []string
	var feishuTokens []string

	// 分离URL、飞书file_token和本地路径
	for _, image := range images {
		if IsImageURL(image) {
			urlsToDownload = append(urlsToDownload, image)
		} else if IsFeishuFileToken(image) {
			feishuTokens = append(feishuTokens, image)
		} else {
			// 本地路径直接添加
			localPaths = append(localPaths, image)
		}
	}

	// 批量下载URL图片
	if len(urlsToDownload) > 0 {
		downloadedPaths, err := p.downloader.DownloadImages(urlsToDownload)
		if err != nil {
			return nil, fmt.Errorf("failed to download images: %w", err)
		}
		localPaths = append(localPaths, downloadedPaths...)
	}

	// 批量下载飞书图片
	if len(feishuTokens) > 0 {
		if accessToken == "" {
			return nil, fmt.Errorf("access_token is required for feishu file_token download")
		}
		downloadedPaths, err := p.downloader.DownloadFeishuImages(feishuTokens, accessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to download feishu images: %w", err)
		}
		localPaths = append(localPaths, downloadedPaths...)
	}

	if len(localPaths) == 0 {
		return nil, fmt.Errorf("no valid images found")
	}

	return localPaths, nil
}
