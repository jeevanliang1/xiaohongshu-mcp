package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("测试左对齐文字渲染...")

	// 测试华文黑体字体
	fontPath := "/System/Library/Fonts/STHeiti Medium.ttc"

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
	dc := gg.NewContext(600, 300)

	// 绘制背景
	dc.SetColor(color.RGBA{240, 240, 240, 255})
	dc.DrawRectangle(0, 0, 600, 300)
	dc.Fill()

	// 设置字体和颜色
	dc.SetFontFace(face)
	dc.SetColor(color.RGBA{0, 0, 0, 255})

	// 绘制文字（左对齐）
	lines := []string{"第一行文字", "第二行文字", "第三行文字"}
	lineHeight := 60.0
	startY := 50.0
	textStartX := 50.0 // 左边距50像素

	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		dc.DrawStringAnchored(line, textStartX, y, 0, 0.5) // 左对齐
	}

	// 保存图片
	filename := "left_align_test.png"
	if err := dc.SavePNG(filename); err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
	} else {
		fmt.Printf("✅ 保存成功: %s\n", filename)
	}

	fmt.Println("测试完成！")
}
