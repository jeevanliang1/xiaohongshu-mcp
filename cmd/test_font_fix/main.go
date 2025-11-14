package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("测试字体修复...")

	// 测试华文黑体字体
	fontPath := "/System/Library/Fonts/STHeiti Medium.ttc"
	testText := "测试中文文字 😊 emoji 💪"

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
	filename := "font_fix_test.png"
	if err := dc.SavePNG(filename); err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
	} else {
		fmt.Printf("✅ 保存成功: %s\n", filename)
	}

	fmt.Println("测试完成！")
}
