package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
)

// ImageGenerator 图片生成器
type ImageGenerator struct {
	assetsDir string
	mutex     sync.Mutex // 互斥锁，确保并发安全
}

// NewImageGenerator 创建图片生成器
func NewImageGenerator(assetsDir string) *ImageGenerator {
	return &ImageGenerator{
		assetsDir: assetsDir,
	}
}

// generateUniqueFileName 生成唯一的文件名，确保并发安全
func (ig *ImageGenerator) generateUniqueFileName(prefix, extension string) string {
	ig.mutex.Lock()
	defer ig.mutex.Unlock()

	return generateUniqueFileNameGlobal(prefix, extension)
}

// generateUniqueFileNameGlobal 全局唯一文件名生成函数，供其他模块使用
func generateUniqueFileNameGlobal(prefix, extension string) string {
	// 使用纳秒时间戳确保高精度
	nanos := time.Now().UnixNano()

	// 生成8字节的随机数
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果crypto/rand失败，回退到math/rand
		mathrand.Seed(time.Now().UnixNano())
		for i := range randomBytes {
			randomBytes[i] = byte(mathrand.Intn(256))
		}
	}
	randomStr := hex.EncodeToString(randomBytes)

	// 组合：前缀_纳秒时间戳_随机字符串.扩展名
	return fmt.Sprintf("%s_%d_%s.%s", prefix, nanos, randomStr, extension)
}

// GenerateCoverImage 生成封面图片
func (ig *ImageGenerator) GenerateCoverImage(req *CoverImageRequest) (*CoverImageResponse, error) {
	// 如果设置了背景图，使用背景图的尺寸并等比例缩放
	if req.BackgroundImage != "" {
		// 读取背景图片尺寸
		bgWidth, bgHeight, err := ig.getImageDimensions(req.BackgroundImage)
		if err != nil {
			return nil, fmt.Errorf("读取背景图片失败: %v", err)
		}

		// 等比例缩放，以1080为最大值
		maxDimension := 1080
		scale := 1.0
		if bgWidth > maxDimension || bgHeight > maxDimension {
			scaleWidth := float64(maxDimension) / float64(bgWidth)
			scaleHeight := float64(maxDimension) / float64(bgHeight)
			// 选择较小的缩放比例以确保两个维度都不超过1080
			scale = math.Min(scaleWidth, scaleHeight)
		}

		req.Width = int(float64(bgWidth) * scale)
		req.Height = int(float64(bgHeight) * scale)

		logrus.Infof("背景图尺寸: %dx%d, 缩放后尺寸: %dx%d (缩放比例: %.2f)", bgWidth, bgHeight, req.Width, req.Height, scale)
	} else {
		// 设置默认值（仅当没有背景图时）
		if req.Width == 0 {
			req.Width = 1080
		}
		if req.Height == 0 {
			req.Height = 1440
		}
	}

	// 设置其他默认值
	if req.FontSize == 0 {
		req.FontSize = 48
	}
	if req.TextColor == "" {
		req.TextColor = "#FFFFFF"
	}
	if req.Style == "" {
		req.Style = "gradient"
	}

	// 创建画布
	dc := gg.NewContext(req.Width, req.Height)

	// 绘制背景
	if err := ig.drawBackground(dc, req); err != nil {
		return nil, fmt.Errorf("绘制背景失败: %v", err)
	}

	// 绘制文字
	if err := ig.drawText(dc, req); err != nil {
		return nil, fmt.Errorf("绘制文字失败: %v", err)
	}

	// 生成输出路径
	if req.OutputPath == "" {
		filename := ig.generateUniqueFileName("cover_image", "png")
		req.OutputPath = filepath.Join("generated_images", filename)
	}

	// 确保输出目录存在
	outputDir := filepath.Dir(req.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 保存图片
	if err := dc.SavePNG(req.OutputPath); err != nil {
		return nil, fmt.Errorf("保存图片失败: %v", err)
	}

	// 将路径转换为绝对路径，确保其他工具可以正确读取
	absPath, err := filepath.Abs(req.OutputPath)
	if err != nil {
		logrus.Errorf("转换绝对路径失败，路径: %s, 错误: %v", req.OutputPath, err)
		// 如果转换失败，仍然使用相对路径
		logrus.Infof("封面图片生成成功: %s (使用相对路径)", req.OutputPath)
		return &CoverImageResponse{
			Success:   true,
			ImagePath: req.OutputPath,
			Message:   "封面图片生成成功",
		}, nil
	}

	logrus.Infof("封面图片生成成功: %s (绝对路径: %s)", req.OutputPath, absPath)

	return &CoverImageResponse{
		Success:   true,
		ImagePath: absPath, // 返回绝对路径
		Message:   "封面图片生成成功",
	}, nil
}

// drawBackground 绘制背景
func (ig *ImageGenerator) drawBackground(dc *gg.Context, req *CoverImageRequest) error {
	// 如果设置了背景图，使用背景图
	if req.BackgroundImage != "" {
		return ig.drawImageBackground(dc, req)
	}

	// 否则使用原来的背景样式
	switch req.Style {
	case "gradient":
		return ig.drawGradientBackground(dc, req)
	case "solid":
		return ig.drawSolidBackground(dc, req)
	case "pattern":
		return ig.drawPatternBackground(dc, req)
	default:
		return ig.drawGradientBackground(dc, req)
	}
}

// drawGradientBackground 绘制渐变背景
func (ig *ImageGenerator) drawGradientBackground(dc *gg.Context, req *CoverImageRequest) error {
	// 随机生成渐变颜色
	mathrand.Seed(time.Now().UnixNano())

	// 预定义一些美观的渐变色彩组合
	gradients := [][]string{
		{"#667eea", "#764ba2"}, // 蓝紫色
		{"#f093fb", "#f5576c"}, // 粉红色
		{"#4facfe", "#00f2fe"}, // 蓝色
		{"#43e97b", "#38f9d7"}, // 绿色
		{"#fa709a", "#fee140"}, // 橙粉色
		{"#a8edea", "#fed6e3"}, // 青粉色
		{"#ff9a9e", "#fecfef"}, // 粉紫色
		{"#ffecd2", "#fcb69f"}, // 橙黄色
		{"#a18cd1", "#fbc2eb"}, // 紫色
		{"#fad0c4", "#ffd1ff"}, // 粉色
	}

	// 随机选择一个渐变
	gradient := gradients[mathrand.Intn(len(gradients))]

	// 解析颜色
	color1, err := parseColor(gradient[0])
	if err != nil {
		return err
	}
	color2, err := parseColor(gradient[1])
	if err != nil {
		return err
	}

	// 创建线性渐变
	gradientPattern := gg.NewLinearGradient(0, 0, float64(req.Width), float64(req.Height))
	gradientPattern.AddColorStop(0, color1)
	gradientPattern.AddColorStop(1, color2)

	// 设置渐变并填充
	dc.SetFillStyle(gradientPattern)
	dc.DrawRectangle(0, 0, float64(req.Width), float64(req.Height))
	dc.Fill()

	return nil
}

// drawSolidBackground 绘制纯色背景
func (ig *ImageGenerator) drawSolidBackground(dc *gg.Context, req *CoverImageRequest) error {
	var bgColor color.Color
	var err error

	if req.BgColor != "" {
		bgColor, err = parseColor(req.BgColor)
		if err != nil {
			return err
		}
	} else {
		// 随机选择一个纯色
		colors := []string{
			"#667eea", "#f093fb", "#4facfe", "#43e97b",
			"#fa709a", "#a8edea", "#ff9a9e", "#ffecd2",
		}
		bgColor, err = parseColor(colors[mathrand.Intn(len(colors))])
		if err != nil {
			return err
		}
	}

	dc.SetColor(bgColor)
	dc.DrawRectangle(0, 0, float64(req.Width), float64(req.Height))
	dc.Fill()

	return nil
}

// drawPatternBackground 绘制图案背景
func (ig *ImageGenerator) drawPatternBackground(dc *gg.Context, req *CoverImageRequest) error {
	// 先绘制基础渐变背景
	if err := ig.drawGradientBackground(dc, req); err != nil {
		return err
	}

	// 添加一些装饰性图案
	dc.SetColor(color.RGBA{255, 255, 255, 30}) // 半透明白色

	// 绘制圆形图案
	for i := 0; i < 20; i++ {
		x := float64(mathrand.Intn(req.Width))
		y := float64(mathrand.Intn(req.Height))
		radius := float64(mathrand.Intn(50) + 10)

		dc.DrawCircle(x, y, radius)
		dc.Fill()
	}

	return nil
}

// normalizeImagePath 规范化图片路径，将相对路径转换为绝对路径
func (ig *ImageGenerator) normalizeImagePath(imagePath string) (string, error) {
	// 如果路径已经是绝对路径，直接返回
	if filepath.IsAbs(imagePath) {
		return imagePath, nil
	}
	// 将相对路径转换为绝对路径
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return imagePath, fmt.Errorf("无法转换路径为绝对路径: %v", err)
	}
	return absPath, nil
}

// getImageDimensions 获取图片的宽高
func (ig *ImageGenerator) getImageDimensions(imagePath string) (int, int, error) {
	// 规范化路径
	normalizedPath, err := ig.normalizeImagePath(imagePath)
	if err != nil {
		return 0, 0, err
	}

	// 打开文件
	file, err := os.Open(normalizedPath)
	if err != nil {
		return 0, 0, fmt.Errorf("无法打开图片文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("无法解码图片: %v", err)
	}

	return img.Width, img.Height, nil
}

// drawImageBackground 绘制背景图片
func (ig *ImageGenerator) drawImageBackground(dc *gg.Context, req *CoverImageRequest) error {
	// 规范化路径
	normalizedPath, err := ig.normalizeImagePath(req.BackgroundImage)
	if err != nil {
		return fmt.Errorf("路径规范化失败: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(normalizedPath); os.IsNotExist(err) {
		return fmt.Errorf("背景图片文件不存在: %s (规范化后: %s)", req.BackgroundImage, normalizedPath)
	}

	// 加载背景图片
	backgroundImg, err := gg.LoadImage(normalizedPath)
	if err != nil {
		return fmt.Errorf("加载背景图片失败: %v (路径: %s)", err, normalizedPath)
	}

	// 获取原始图片尺寸
	imgWidth := float64(backgroundImg.Bounds().Dx())
	imgHeight := float64(backgroundImg.Bounds().Dy())

	// 计算缩放比例以填充整个画布
	scaleX := float64(req.Width) / imgWidth
	scaleY := float64(req.Height) / imgHeight

	// 应用缩放变换
	dc.Push()
	dc.Scale(scaleX, scaleY)

	// 绘制背景图
	dc.DrawImage(backgroundImg, 0, 0)

	// 恢复变换
	dc.Pop()

	return nil
}

// drawText 绘制文字
func (ig *ImageGenerator) drawText(dc *gg.Context, req *CoverImageRequest) error {
	// 解析文字颜色
	textColor, err := parseColor(req.TextColor)
	if err != nil {
		return err
	}

	// 设置字体大小
	logrus.Infof("drawText: 设置字体大小 %d", req.FontSize)

	// 尝试加载支持emoji的字体
	fontFace := ig.loadEmojiSupportFont(req.FontSize)
	if fontFace != nil {
		dc.SetFontFace(fontFace)
		logrus.Infof("drawText: 加载emoji支持字体成功，大小: %d", req.FontSize)
	} else {
		// 回退到原来的方法
		if fontFace := ig.getDefaultFontFace(req.FontSize); fontFace != nil {
			dc.SetFontFace(fontFace)
			logrus.Infof("drawText: 回退方法字体设置成功")
		} else {
			logrus.Warn("drawText: 字体设置失败，使用默认字体")
		}
	}
	dc.SetColor(textColor)

	// 设置padding
	padding := 50.0
	availableWidth := float64(req.Width) - 2*padding

	// 处理换行符和自动换行
	lines := ig.wrapText(dc, req.Text, availableWidth)

	// 计算总文字高度
	lineHeight := float64(req.FontSize) * 1.2 // 行高为字体大小的1.2倍
	totalTextHeight := float64(len(lines)) * lineHeight

	// 计算文字区域位置（居中，然后应用垂直偏移）
	textAreaX := float64(req.Width) / 2
	textAreaY := float64(req.Height)/2 + float64(req.TextOffsetY)

	// 绘制半透明背景框
	boxWidth := availableWidth
	boxHeight := totalTextHeight + 40 // 上下各加20的padding
	boxX := textAreaX - boxWidth/2
	boxY := textAreaY - boxHeight/2

	// 确保背景框不会超出100像素的padding边界
	if boxX < padding {
		boxX = padding
		boxWidth = float64(req.Width) - 2*padding
	}
	if boxY < padding {
		boxY = padding
	}
	if boxX+boxWidth > float64(req.Width)-padding {
		boxWidth = float64(req.Width) - 2*padding
		boxX = padding
	}
	if boxY+boxHeight > float64(req.Height)-padding {
		boxHeight = float64(req.Height) - 2*padding
		boxY = padding
	}

	// 绘制圆角矩形背景
	dc.SetColor(color.RGBA{0, 0, 0, 153}) // 黑色半透明
	dc.DrawRoundedRectangle(boxX, boxY, boxWidth, boxHeight, 15)
	dc.Fill()

	// 绘制文字行（左对齐）
	dc.SetColor(textColor)
	startY := textAreaY - totalTextHeight/2 + lineHeight/2
	textStartX := padding + 20 // 左边距：padding + 20像素
	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		dc.DrawStringAnchored(line, textStartX, y, 0, 0.5) // 左对齐：0, 0.5
	}

	return nil
}

// wrapText 处理文字换行，支持手动换行符和自动换行
func (ig *ImageGenerator) wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string

	// 首先按换行符分割
	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		if paragraph == "" {
			// 空行
			lines = append(lines, "")
			continue
		}

		// 对每个段落进行自动换行
		// 对于中文，按字符分割；对于英文，按单词分割
		words := ig.splitText(paragraph)
		if len(words) == 0 {
			continue
		}

		var currentLine strings.Builder
		currentLine.WriteString(words[0])

		for i := 1; i < len(words); i++ {
			word := words[i]
			separator := " "
			// 如果前一个字符是中文或当前字符是中文，不需要空格
			if len(currentLine.String()) > 0 && len(word) > 0 {
				lastChar := currentLine.String()[len(currentLine.String())-1:]
				firstChar := word[0:1]
				if ig.isChinese(lastChar) || ig.isChinese(firstChar) {
					separator = ""
				}
			}

			testLine := currentLine.String() + separator + word

			// 测量当前行的宽度
			width, _ := dc.MeasureString(testLine)

			if width <= maxWidth {
				// 可以添加这个词
				currentLine.WriteString(separator)
				currentLine.WriteString(word)
			} else {
				// 需要换行
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
			}
		}

		// 添加最后一行
		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		}
	}

	// 如果没有内容，返回一个空行
	if len(lines) == 0 {
		lines = append(lines, "")
	}

	return lines
}

// splitText 智能分割文本，支持中英文混合
func (ig *ImageGenerator) splitText(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range text {
		if ig.isChinese(string(r)) {
			// 中文字符，每个字符作为一个词
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			words = append(words, string(r))
		} else if r == ' ' || r == '\t' {
			// 空格或制表符，结束当前词
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		} else {
			// 英文字符，添加到当前词
			current.WriteRune(r)
		}
	}

	// 添加最后一个词
	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// isChinese 判断字符是否为中文
func (ig *ImageGenerator) isChinese(char string) bool {
	if len(char) == 0 {
		return false
	}
	r := []rune(char)[0]
	return r >= 0x4e00 && r <= 0x9fff
}

// loadEmojiSupportFont 加载支持emoji的字体
func (ig *ImageGenerator) loadEmojiSupportFont(size int) font.Face {
	// 优先使用支持中文和emoji的字体（按优先级排序）
	emojiFonts := []string{
		// macOS 系统字体，优先支持中文的字体
		"/System/Library/Fonts/STHeiti Medium.ttc",    // 华文黑体-中（支持中文和emoji）
		"/System/Library/Fonts/STHeiti Light.ttc",     // 华文黑体-细（支持中文和emoji）
		"/System/Library/Fonts/PingFang.ttc",          // 苹方（支持中文和emoji）
		"/System/Library/Fonts/Hiragino Sans GB.ttc",  // 冬青黑体简体中文（支持中文和emoji）
		"/System/Library/Fonts/Helvetica.ttc",         // Helvetica（支持emoji）
		"/System/Library/Fonts/Arial.ttf",             // Arial（支持emoji）
		"/System/Library/Fonts/Apple Color Emoji.ttc", // Apple Color Emoji（仅支持emoji）

		// Linux 字体
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
		"/usr/share/fonts/truetype/noto/NotoColorEmoji.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",

		// Windows 字体
		"C:/Windows/Fonts/msyh.ttc",     // 微软雅黑（支持中文和emoji）
		"C:/Windows/Fonts/arial.ttf",    // Arial（支持emoji）
		"C:/Windows/Fonts/seguiemj.ttf", // Segoe UI Emoji
	}

	for _, fontPath := range emojiFonts {
		if _, err := os.Stat(fontPath); err == nil {
			if fontFace, err := gg.LoadFontFace(fontPath, float64(size)); err == nil {
				logrus.Infof("成功加载emoji支持字体: %s, 大小: %d", fontPath, size)
				return fontFace
			} else {
				logrus.Debugf("emoji字体加载失败: %s, 错误: %v", fontPath, err)
			}
		}
	}

	logrus.Warn("无法加载任何emoji支持字体")
	return nil
}

// getDefaultFontFace 获取默认字体
func (ig *ImageGenerator) getDefaultFontFace(size int) font.Face {
	logrus.Infof("getDefaultFontFace: 请求字体大小 %d", size)

	// 优先使用华文黑体，这是最可靠的中文字体
	primaryFontPath := "/System/Library/Fonts/STHeiti Medium.ttc"

	if _, err := os.Stat(primaryFontPath); err == nil {
		if fontFace, err := gg.LoadFontFace(primaryFontPath, float64(size)); err == nil {
			logrus.Infof("成功加载主字体: %s, 大小: %d", primaryFontPath, size)
			return fontFace
		} else {
			logrus.Warnf("主字体加载失败: %s, 错误: %v", primaryFontPath, err)
		}
	}

	// 备用字体列表
	fontPaths := []string{
		"/System/Library/Fonts/STHeiti Light.ttc",    // 华文黑体-细
		"/System/Library/Fonts/PingFang.ttc",         // 苹方（如果支持）
		"/System/Library/Fonts/Hiragino Sans GB.ttc", // 冬青黑体简体中文（如果支持）
		"/System/Library/Fonts/Geneva.ttf",           // 备用字体
		"/System/Library/Fonts/NewYork.ttf",          // 备用字体
		"/System/Library/Fonts/SFCompact.ttf",        // 备用字体
		"/System/Library/Fonts/Symbol.ttf",           // 备用字体
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
		"C:/Windows/Fonts/msyh.ttc",   // 微软雅黑
		"C:/Windows/Fonts/simsun.ttc", // 宋体
		"C:/Windows/Fonts/arial.ttf",  // Arial
		"assets/fonts/NotoSansCJK-Regular.ttc",
		"assets/fonts/SourceHanSansCN-Regular.otf",
	}

	for _, fontPath := range fontPaths {
		if _, err := os.Stat(fontPath); err == nil {
			if fontFace, err := gg.LoadFontFace(fontPath, float64(size)); err == nil {
				logrus.Infof("成功加载备用字体: %s, 大小: %d", fontPath, size)
				return fontFace
			} else {
				logrus.Debugf("备用字体加载失败: %s, 错误: %v", fontPath, err)
			}
		}
	}

	// 如果都加载失败，使用默认字体
	logrus.Warn("无法加载任何系统字体，使用默认字体（可能不支持中文）")
	return nil
}

// parseColor 解析颜色字符串
func parseColor(colorStr string) (color.Color, error) {
	// 移除 # 前缀
	if strings.HasPrefix(colorStr, "#") {
		colorStr = colorStr[1:]
	}

	// 确保是6位十六进制
	if len(colorStr) != 6 {
		return nil, fmt.Errorf("无效的颜色格式: %s", colorStr)
	}

	// 解析RGB值
	var r, g, b uint8
	if _, err := fmt.Sscanf(colorStr, "%02x%02x%02x", &r, &g, &b); err != nil {
		return nil, fmt.Errorf("解析颜色失败: %v", err)
	}

	return color.RGBA{r, g, b, 255}, nil
}
