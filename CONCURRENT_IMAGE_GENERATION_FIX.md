# 图片生成并发问题修复

## 问题描述

原来的图片生成系统存在以下并发问题：

1. **时间戳精度不够**：使用 `time.Now().Format("20060102_150405")` 只精确到秒，多个并发请求可能在同一秒内生成，导致文件名冲突
2. **没有线程安全保护**：`ImageGenerator` 没有互斥锁保护
3. **文件名生成逻辑分散**：在多个地方都有文件名生成逻辑，没有统一管理

## 解决方案

### 1. 实现线程安全的唯一文件名生成机制

在 `image_generator.go` 中添加了：

```go
// ImageGenerator 图片生成器
type ImageGenerator struct {
    assetsDir string
    mutex     sync.Mutex // 互斥锁，确保并发安全
}

// generateUniqueFileNameGlobal 全局唯一文件名生成函数，供其他模块使用
func generateUniqueFileNameGlobal(prefix, extension string) string {
    // 使用纳秒时间戳确保高精度
    nanos := time.Now().UnixNano()
    
    // 生成8字节的随机数
    randomBytes := make([]byte, 8)
    if _, err := rand.Read(randomBytes); err != nil {
        // 如果crypto/rand失败，回退到math/rand
        mathrand.Seed(time.Now().UnixNano())
        for i := range randomBytes {
            randomBytes[i] = byte(mathrand.Intn(256))
        }
    }
    randomStr := hex.EncodeToString(randomBytes)
    
    // 组合：前缀_纳秒时间戳_随机字符串.扩展名
    return fmt.Sprintf("%s_%d_%s.%s", prefix, nanos, randomStr, extension)
}
```

### 2. 更新所有图片生成相关的文件名生成逻辑

修改了以下文件中的文件名生成逻辑：

- `image_generator.go` - 封面图片生成
- `mcp_handlers.go` - 文生图和图生图功能

**原来的文件名格式：**
```
cover_image_20250126_150405.png
generated_image_20250126_150405_1.jpg
img2img_20250126_150405_1.jpg
```

**新的文件名格式：**
```
cover_image_1761487032795583000_62118b8058fa99d5.png
generated_image_1761487032795853000_bc3482c0bb5cb2c5.jpg
img2img_1761487032795555000_80f3d73224b25f92.jpg
```

### 3. 文件名格式说明

新的文件名格式：`{前缀}_{纳秒时间戳}_{16位随机数}.{扩展名}`

- **前缀**：标识图片类型（cover_image, generated_image, img2img）
- **纳秒时间戳**：确保高精度时间戳，避免并发冲突
- **16位随机数**：额外的随机性保证，即使在同一纳秒内也能保证唯一性
- **扩展名**：文件格式（png, jpg）

## 测试结果

通过并发测试验证了修复效果：

- 测试了10个并发请求，每个请求生成3个文件名
- 总共生成30个文件名，全部唯一
- 无任何文件名冲突

```
✅ 所有文件名都是唯一的，并发测试通过！
```

## 优势

1. **高并发支持**：纳秒级时间戳 + 随机数确保即使在极高并发下也不会产生文件名冲突
2. **线程安全**：使用互斥锁保护关键操作
3. **统一管理**：所有文件名生成逻辑统一使用 `generateUniqueFileNameGlobal` 函数
4. **向后兼容**：不影响现有功能，只是改进了文件名生成机制
5. **可扩展性**：新的文件名生成机制可以轻松扩展到其他需要唯一文件名的场景

## 使用示例

```go
// 生成封面图片文件名
filename := generateUniqueFileNameGlobal("cover_image", "png")
// 结果: cover_image_1761487032795583000_62118b8058fa99d5.png

// 生成文生图文件名
filename := generateUniqueFileNameGlobal("generated_image", "jpg")
// 结果: generated_image_1761487032795853000_bc3482c0bb5cb2c5.jpg

// 生成图生图文件名
filename := generateUniqueFileNameGlobal("img2img", "jpg")
// 结果: img2img_1761487032795555000_80f3d73224b25f92.jpg
```

现在你的图片生成系统可以安全地处理多个并发请求，不会再出现文件名冲突的问题！
