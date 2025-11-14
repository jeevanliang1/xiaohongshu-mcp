package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("创建字体大小对比图...")

	// 测试不同的字体大小
	fontSizes := []int{24, 48, 72, 96}
	fontPath := "/System/Library/Fonts/STHeiti Medium.ttc"
	testText := "字体大小测试"

	if _, err := os.Stat(fontPath); err != nil {
		fmt.Printf("❌ 字体文件不存在: %s\n", fontPath)
		return
	}

	// 创建一个大的画布来显示所有字体大小
	canvasWidth := 800
	canvasHeight := 600
	dc := gg.NewContext(canvasWidth, canvasHeight)

	// 绘制白色背景
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.DrawRectangle(0, 0, float64(canvasWidth), float64(canvasHeight))
	dc.Fill()

	// 绘制标题
	titleFace, err := gg.LoadFontFace(fontPath, 32)
	if err != nil {
		fmt.Printf("❌ 标题字体加载失败: %v\n", err)
		return
	}
	dc.SetFontFace(titleFace)
	dc.SetColor(color.RGBA{0, 0, 0, 255})
	dc.DrawStringAnchored("字体大小对比测试", float64(canvasWidth)/2, 50, 0.5, 0.5)

	// 绘制不同字体大小的文字
	yPos := 120.0
	for _, fontSize := range fontSizes {
		fmt.Printf("绘制字体大小: %d\n", fontSize)

		face, err := gg.LoadFontFace(fontPath, float64(fontSize))
		if err != nil {
			fmt.Printf("  ❌ 字体加载失败: %v\n", err)
			continue
		}

		// 设置字体和颜色
		dc.SetFontFace(face)
		dc.SetColor(color.RGBA{0, 0, 0, 255})

		// 测量文字尺寸
		textWidth, textHeight := dc.MeasureString(testText)
		fmt.Printf("  📏 文字尺寸: %.1f x %.1f\n", textWidth, textHeight)

		// 绘制标签
		labelFace, _ := gg.LoadFontFace(fontPath, 20)
		dc.SetFontFace(labelFace)
		dc.SetColor(color.RGBA{100, 100, 100, 255})
		dc.DrawStringAnchored(fmt.Sprintf("字体大小: %d", fontSize), 100, yPos, 0, 0.5)

		// 绘制测试文字
		dc.SetFontFace(face)
		dc.SetColor(color.RGBA{0, 0, 0, 255})
		dc.DrawStringAnchored(testText, 400, yPos, 0.5, 0.5)

		yPos += textHeight + 40
	}

	// 保存对比图
	filename := "font_size_comparison.png"
	if err := dc.SavePNG(filename); err != nil {
		fmt.Printf("❌ 保存图片失败: %v\n", err)
	} else {
		fmt.Printf("✅ 字体大小对比图已保存: %s\n", filename)

		// 检查文件大小
		if stat, err := os.Stat(filename); err == nil {
			fmt.Printf("📁 文件大小: %d bytes\n", stat.Size())
		}
	}

	fmt.Println("测试完成！")
}
