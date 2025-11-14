package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("测试修复后的功能...")

	// 测试华文黑体字体
	fontPath := "/System/Library/Fonts/STHeiti Medium.ttc"

	if _, err := os.Stat(fontPath); err != nil {
		fmt.Printf("❌ 字体文件不存在: %s\n", fontPath)
		return
	}

	// 测试不同的文字内容
	testCases := []struct {
		name string
		text string
	}{
		{
			name: "中文测试",
			text: "生活不止眼前的苟且，还有诗和远方",
		},
		{
			name: "emoji测试",
			text: "美好生活 😊 加油 💪",
		},
		{
			name: "混合测试",
			text: "Hello 世界 🌍\n美好生活 🏠",
		},
	}

	for i, testCase := range testCases {
		fmt.Printf("\n=== 测试 %d: %s ===\n", i+1, testCase.name)
		fmt.Printf("文字内容: %s\n", testCase.text)

		face, err := gg.LoadFontFace(fontPath, 48)
		if err != nil {
			fmt.Printf("❌ 字体加载失败: %v\n", err)
			continue
		}

		// 创建测试图片
		dc := gg.NewContext(400, 200)

		// 绘制背景
		dc.SetColor(color.RGBA{240, 240, 240, 255})
		dc.DrawRectangle(0, 0, 400, 200)
		dc.Fill()

		// 设置字体和颜色
		dc.SetFontFace(face)
		dc.SetColor(color.RGBA{0, 0, 0, 255})

		// 测量文字
		textWidth, textHeight := dc.MeasureString(testCase.text)
		fmt.Printf("文字尺寸: %.1f x %.1f\n", textWidth, textHeight)

		// 绘制文字
		dc.DrawStringAnchored(testCase.text, 200, 100, 0.5, 0.5)

		// 保存图片
		filename := fmt.Sprintf("test_fix_%d.png", i+1)
		if err := dc.SavePNG(filename); err != nil {
			fmt.Printf("❌ 保存失败: %v\n", err)
		} else {
			fmt.Printf("✅ 保存成功: %s\n", filename)
		}
	}

	fmt.Println("\n测试完成！")
}
