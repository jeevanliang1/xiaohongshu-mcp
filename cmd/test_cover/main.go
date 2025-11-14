package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xpzouying/xiaohongshu-mcp"
)

func main() {
	// 创建图片生成器
	imageGenerator := xiaohongshu.NewImageGenerator("assets")

	// 测试生成封面图片
	req := &xiaohongshu.CoverImageRequest{
		Text:      "生活不止眼前的苟且，还有诗和远方",
		Width:     800,
		Height:    600,
		FontSize:  48,
		TextColor: "#FFFFFF",
		Style:     "gradient",
	}

	fmt.Println("开始测试中文字体支持...")

	result, err := imageGenerator.GenerateCoverImage(req)
	if err != nil {
		log.Fatalf("❌ 生成封面图片失败: %v", err)
	}

	fmt.Printf("✅ 封面图片生成成功！\n")
	fmt.Printf("图片路径: %s\n", result.ImagePath)
	fmt.Printf("消息: %s\n", result.Message)

	// 检查文件是否存在
	if _, err := os.Stat(result.ImagePath); err == nil {
		fmt.Printf("✅ 图片文件已成功创建\n")

		// 获取文件大小
		if stat, err := os.Stat(result.ImagePath); err == nil {
			fmt.Printf("文件大小: %d bytes\n", stat.Size())
		}
	} else {
		fmt.Printf("❌ 图片文件创建失败: %v\n", err)
	}
}
