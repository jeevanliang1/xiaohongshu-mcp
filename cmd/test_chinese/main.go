package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("测试修复后的中文字体支持...")

	// 使用真正的中文字体
	fontPath := "/System/Library/Fonts/STHeiti Medium.ttc"
	testText := "生活不止眼前的苟且，还有诗和远方"

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
	dc := gg.NewContext(800, 600)

	// 绘制渐变背景
	gradientPattern := gg.NewLinearGradient(0, 0, 800, 600)
	gradientPattern.AddColorStop(0, color.RGBA{102, 126, 234, 255}) // 蓝紫色
	gradientPattern.AddColorStop(1, color.RGBA{118, 75, 162, 255})  // 紫色

	dc.SetFillStyle(gradientPattern)
	dc.DrawRectangle(0, 0, 800, 600)
	dc.Fill()

	// 设置字体和颜色
	dc.SetFontFace(face)
	dc.SetColor(color.RGBA{255, 255, 255, 255})

	// 测量文字尺寸
	textWidth, textHeight := dc.MeasureString(testText)

	// 计算文字位置（居中）
	textX := 400.0
	textY := 300.0

	// 绘制半透明背景框
	padding := 20.0
	boxWidth := textWidth + 2*padding
	boxHeight := textHeight + 2*padding
	boxX := textX - boxWidth/2
	boxY := textY - boxHeight/2

	// 绘制圆角矩形背景
	dc.SetColor(color.RGBA{0, 0, 0, 153}) // 黑色半透明
	dc.DrawRoundedRectangle(boxX, boxY, boxWidth, boxHeight, 15)
	dc.Fill()

	// 绘制文字
	dc.SetColor(color.RGBA{255, 255, 255, 255}) // 白色
	dc.DrawStringAnchored(testText, textX, textY, 0.5, 0.5)

	// 保存测试图片
	filename := "chinese_font_test.png"
	if err := dc.SavePNG(filename); err != nil {
		fmt.Printf("❌ 保存图片失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试图片已保存: %s\n", filename)

		// 检查文件大小
		if stat, err := os.Stat(filename); err == nil {
			fmt.Printf("📁 文件大小: %d bytes\n", stat.Size())
		}
	}

	fmt.Println("测试完成！")
}
