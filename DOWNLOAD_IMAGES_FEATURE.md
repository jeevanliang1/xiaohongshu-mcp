# 下载图片工具功能说明

## 概述

`download_images` 工具用于下载并保存图片到本地，支持自定义保存目录。

## 功能特性

- ✅ 支持 HTTP/HTTPS 图片URL下载
- ✅ 支持本地图片路径（直接返回路径）
- ✅ 支持混合输入（URL和本地路径）
- ✅ 自动创建保存目录
- ✅ 自定义保存目录（可选参数）
- ✅ 默认保存到 `image_file` 目录
- ✅ 防重复下载（相同URL生成相同文件名）
- ✅ 图片格式验证
- ✅ 错误处理和恢复机制

## MCP 工具参数

### 工具名称
`download_images`

### 参数说明

```json
{
  "images": ["string"],     // 必需：图片URL或本地路径列表
  "save_dir": "string"      // 可选：自定义保存目录路径
}
```

#### 参数详情

- **images** (必需): 字符串数组
  - 支持 HTTP/HTTPS 图片URL
  - 支持本地文件绝对路径
  - 可以混合使用URL和本地路径

- **save_dir** (可选): 字符串
  - 指定图片保存的目录路径
  - 如果不提供，默认使用 `image_file` 目录
  - 目录不存在时会自动创建
  - 支持相对路径和绝对路径

## 使用示例

### 示例1: 使用默认目录
```json
{
  "name": "download_images",
  "arguments": {
    "images": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.png"
    ]
  }
}
```
- 图片将保存到 `./image_file/` 目录

### 示例2: 使用自定义目录
```json
{
  "name": "download_images",
  "arguments": {
    "images": [
      "https://example.com/image1.jpg",
      "/Users/user/Pictures/local.png"
    ],
    "save_dir": "/Users/user/Downloads/my_images"
  }
}
```
- URL图片下载到 `/Users/user/Downloads/my_images/`
- 本地图片路径直接返回

### 示例3: 混合输入
```json
{
  "name": "download_images",
  "arguments": {
    "images": [
      "https://example.com/remote.jpg",
      "/Users/user/Pictures/local.png",
      "https://another-site.com/image.gif"
    ],
    "save_dir": "custom_folder"
  }
}
```

## 返回结果

### 成功响应
```json
{
  "saved_paths": [
    "/path/to/saved/image1.jpg",
    "/path/to/saved/image2.png"
  ],
  "count": 2
}
```

### 错误响应
```json
{
  "content": [
    {
      "type": "text",
      "text": "下载失败: 具体错误信息"
    }
  ],
  "isError": true
}
```

## 目录创建逻辑

1. **默认目录**: 如果不指定 `save_dir`，使用 `image_file` 目录
2. **自动创建**: 目录不存在时自动创建（权限 0755）
3. **路径处理**: 支持相对路径和绝对路径
4. **错误处理**: 目录创建失败时返回详细错误信息

## 文件命名规则

下载的图片文件使用以下命名规则：
- 格式: `img_{hash}_{timestamp}.{extension}`
- hash: URL的SHA256哈希值前16位
- timestamp: 下载时间戳
- extension: 检测到的图片格式扩展名

## 错误处理

- ✅ 无效URL格式
- ✅ 网络下载失败
- ✅ 非图片文件
- ✅ 目录创建失败
- ✅ 文件保存失败
- ✅ 参数验证错误

## 注意事项

1. **权限**: 确保对目标目录有写入权限
2. **磁盘空间**: 确保有足够的磁盘空间存储图片
3. **网络**: URL图片需要网络连接
4. **格式**: 只支持常见图片格式（jpg, png, gif, webp等）
5. **重复**: 相同URL会生成相同文件名，避免重复下载

## 技术实现

- 使用现有的 `downloader` 包
- 支持 panic 恢复机制
- 集成到 MCP 工具系统
- 遵循项目代码规范
