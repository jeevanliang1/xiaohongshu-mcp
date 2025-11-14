package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)

	fmt.Println("测试字体大小参数...")

	// 创建图片生成器
	imageGenerator := NewImageGenerator("assets")

	// 测试不同的字体大小
	fontSizes := []int{48, 100, 200}
	testText := "字体大小测试"

	for _, fontSize := range fontSizes {
		fmt.Printf("\n=== 测试字体大小: %d ===\n", fontSize)

		// 测试生成封面图片
		req := &CoverImageRequest{
			Text:      testText,
			Width:     400,
			Height:    200,
			FontSize:  fontSize,
			TextColor: "#000000",
			Style:     "solid",
			BgColor:   "#F0F0F0",
		}

		fmt.Printf("请求参数: FontSize=%d\n", req.FontSize)

		result, err := imageGenerator.GenerateCoverImage(req)
		if err != nil {
			log.Printf("❌ 生成封面图片失败: %v", err)
			continue
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

	fmt.Println("\n测试完成！")
}
