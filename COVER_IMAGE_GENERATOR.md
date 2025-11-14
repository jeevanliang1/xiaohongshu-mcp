# 封面图片生成器 MCP 工具

## 功能概述

新增的 `generate_cover_image` MCP 工具可以根据文字内容自动生成美观的封面图片，支持多种背景样式和自定义参数，非常适合用于小红书封面、海报等场景。

## 功能特性

- 🎨 **多种背景样式**：支持渐变、纯色、图案三种背景样式
- 🎯 **智能文字布局**：自动居中显示，带半透明背景框
- 🌈 **丰富色彩搭配**：内置10种美观的渐变色彩组合
- ⚙️ **高度可定制**：支持自定义尺寸、字体大小、颜色等参数
- 📱 **跨平台字体**：自动适配 macOS、Linux、Windows 系统字体
- 📝 **智能换行**：支持自动换行和手动换行符（\n）
- 😊 **Emoji支持**：完美渲染emoji表情符号
- 📏 **智能布局**：自动100像素padding和左对齐

## 使用方法

### MCP 工具调用

```json
{
  "name": "generate_cover_image",
  "arguments": {
    "text": "生活不止眼前的苟且，还有诗和远方",
    "width": 800,
    "height": 600,
    "font_size": 48,
    "text_color": "#FFFFFF",
    "style": "gradient"
  }
}
```

### 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `text` | string | ✅ | - | 要显示在封面上的文字内容 |
| `width` | int | ❌ | 1080 | 图片宽度（像素） |
| `height` | int | ❌ | 1440 | 图片高度（像素） |
| `font_size` | int | ❌ | 48 | 字体大小 |
| `text_color` | string | ❌ | "#FFFFFF" | 文字颜色（十六进制格式） |
| `bg_color` | string | ❌ | 随机 | 背景颜色（十六进制格式，仅solid样式有效） |
| `style` | string | ❌ | "gradient" | 背景样式：gradient（渐变）、solid（纯色）、pattern（图案） |
| `output_path` | string | ❌ | 自动生成 | 输出文件路径 |

### 背景样式说明

1. **gradient（渐变）**：使用随机选择的渐变色彩组合，效果最佳
2. **solid（纯色）**：使用纯色背景，可通过 `bg_color` 参数自定义
3. **pattern（图案）**：渐变背景 + 装饰性圆形图案

### 内置渐变色彩

工具内置了10种精心挑选的渐变色彩组合：
- 蓝紫色渐变
- 粉红色渐变  
- 蓝色渐变
- 绿色渐变
- 橙粉色渐变
- 青粉色渐变
- 粉紫色渐变
- 橙黄色渐变
- 紫色渐变
- 粉色渐变

## 使用示例

### 基础用法
```json
{
  "name": "generate_cover_image",
  "arguments": {
    "text": "今天也要加油呀！💪"
  }
}
```

### 多行文字和emoji
```json
{
  "name": "generate_cover_image",
  "arguments": {
    "text": "美好生活 🏠\n\n生活不止眼前的苟且\n还有诗和远方 🌟\n\n加油！💪",
    "font_size": 48
  }
}
```

### 自定义样式
```json
{
  "name": "generate_cover_image", 
  "arguments": {
    "text": "美食分享",
    "width": 1080,
    "height": 1440,
    "font_size": 60,
    "text_color": "#FF6B6B",
    "style": "pattern"
  }
}
```

### 纯色背景
```json
{
  "name": "generate_cover_image",
  "arguments": {
    "text": "简约风格",
    "style": "solid",
    "bg_color": "#4ECDC4",
    "text_color": "#FFFFFF"
  }
}
```

## 输出结果

工具会返回包含以下信息的 JSON 响应：

```json
{
  "success": true,
  "image_path": "generated_images/cover_image_20251026_165317.png",
  "image_url": "http://192.168.0.113:18060/images/cover_image_20251026_165317.png",
  "message": "封面图片生成成功"
}
```

- `image_path`：本地文件路径
- `image_url`：HTTP 访问地址（可用于在线预览）
- `message`：操作结果消息

## 技术实现

- **图形库**：使用 `github.com/fogleman/gg` 进行图像绘制
- **字体支持**：自动检测系统字体，支持中英文显示
- **颜色解析**：支持十六进制颜色格式
- **文件管理**：自动创建输出目录，生成唯一文件名

## 注意事项

1. ✅ **中文字体支持**：已完美支持中文显示，使用真正的中文字体（华文黑体）
2. ✅ **自动换行**：文字会自动换行，不会被裁减，支持手动换行符（\n）
3. ✅ **Emoji支持**：完美支持emoji表情符号渲染
4. ✅ **智能布局**：自动添加100像素padding，文字左对齐，确保文字不会被裁减
5. 生成的图片保存在 `generated_images/` 目录下
6. 默认尺寸为1080×1440（竖屏比例），适合小红书等社交媒体平台
7. 字体加载优先级：Apple Color Emoji > STHeiti Medium > STHeiti Light > PingFang > Hiragino Sans GB > 其他字体
8. ✅ **Emoji渲染**：优先使用Apple Color Emoji字体，完美支持emoji表情符号

## 更新日志

- **v1.5.2** (2025-10-26)：✅ 修改文字对齐方式为左对齐，提升阅读体验
- **v1.5.1** (2025-10-26)：✅ 修复emoji渲染问题，优化字体加载优先级，优先使用Apple Color Emoji
- **v1.5.0** (2025-10-26)：✅ 重大更新！新增自动换行、emoji支持、智能布局功能
- **v1.4.0** (2025-10-26)：✅ 修复字体大小参数不生效的问题，避免字体缓存问题
- **v1.3.0** (2025-10-26)：✅ 更新默认尺寸为1080×1440，更适合小红书等社交媒体平台
- **v1.2.0** (2025-10-26)：✅ 彻底修复中文字体显示问题，使用华文黑体等真正的中文字体
- **v1.1.0** (2025-10-26)：✅ 修复中文字体支持，完美支持中文显示
- **v1.0.0** (2025-10-26)：初始版本，支持基础封面图片生成功能
