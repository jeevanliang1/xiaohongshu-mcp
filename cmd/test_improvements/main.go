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

	fmt.Println("测试文字渲染改进功能...")

	// 创建图片生成器
	imageGenerator := NewImageGenerator("assets")

	// 测试用例
	testCases := []struct {
		name     string
		text     string
		fontSize int
	}{
		{
			name:     "自动换行测试",
			text:     "这是一段很长的文字，用来测试自动换行功能是否正常工作，当文字宽度超过图片宽度时应该自动换行",
			fontSize: 48,
		},
		{
			name:     "手动换行测试",
			text:     "第一行文字\n第二行文字\n第三行文字",
			fontSize: 48,
		},
		{
			name:     "emoji测试",
			text:     "生活不止眼前的苟且 😊\n还有诗和远方 🌟\n加油！💪",
			fontSize: 48,
		},
		{
			name:     "混合测试",
			text:     "标题：美好生活 🏠\n\n这是一段包含emoji和换行的长文字，用来测试所有功能是否都能正常工作，包括自动换行、手动换行和emoji渲染。\n\n结尾：谢谢！🙏",
			fontSize: 36,
		},
	}

	for i, testCase := range testCases {
		fmt.Printf("\n=== 测试 %d: %s ===\n", i+1, testCase.name)
		fmt.Printf("文字内容: %s\n", testCase.text)

		// 测试生成封面图片
		req := &CoverImageRequest{
			Text:      testCase.text,
			Width:     1080,
			Height:    1440,
			FontSize:  testCase.fontSize,
			TextColor: "#FFFFFF",
			Style:     "gradient",
		}

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
