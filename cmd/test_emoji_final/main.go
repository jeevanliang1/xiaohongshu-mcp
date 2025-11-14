package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("测试最终emoji渲染...")

	// 使用Apple Color Emoji字体
	fontPath := "/System/Library/Fonts/Apple Color Emoji.ttc"
	testText := "美好生活 😊 加油 💪 星星 🌟"

	if _, err := os.Stat(fontPath); err != nil {
		fmt.Printf("❌ 字体文件不存在: %s\n", fontPath)
		return
	}

	face, err := gg.LoadFontFace(fontPath, 48)
	if err != nil {
		fmt.Printf("❌ 字体加载失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 字体加载成功: %s\n", fontPath)

	// 创建测试图片
	dc := gg.NewContext(600, 200)

	// 绘制背景
	dc.SetColor(color.RGBA{240, 240, 240, 255})
	dc.DrawRectangle(0, 0, 600, 200)
	dc.Fill()

	// 设置字体和颜色
	dc.SetFontFace(face)
	dc.SetColor(color.RGBA{0, 0, 0, 255})

	// 测量文字
	textWidth, textHeight := dc.MeasureString(testText)
	fmt.Printf("文字尺寸: %.1f x %.1f\n", textWidth, textHeight)

	// 绘制文字
	dc.DrawStringAnchored(testText, 300, 100, 0.5, 0.5)

	// 保存图片
	filename := "emoji_final_test.png"
	if err := dc.SavePNG(filename); err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
	} else {
		fmt.Printf("✅ 保存成功: %s\n", filename)
	}

	fmt.Println("测试完成！")
}
