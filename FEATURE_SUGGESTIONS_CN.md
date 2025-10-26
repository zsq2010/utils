# 常用工具建议 (快速参考)

## 🚀 最推荐实现的功能 (Top 5)

### 1. 🔄 图片合并 (MergeGrid)
```go
// 将分割的图片重新合并
imagesplit.MergeGrid(tiles, 2, 2, "merged.png", opts)
```
**优势**: 完善分割/合并闭环，实现可逆操作

---

### 2. ℹ️ 图片信息获取 (GetImageInfo)
```go
// 快速获取图片尺寸、格式等信息
info, err := imagesplit.GetImageInfo("input.jpg")
fmt.Printf("%dx%d, %s, %d bytes\n", info.Width, info.Height, info.Format, info.FileSize)
```
**优势**: 分割前了解图片属性，确定合适参数

---

### 3. 👁️ 分割预览 (PreviewSplit)
```go
// 生成带网格线的预览图，不实际分割
imagesplit.PreviewGridSplit("input.png", 3, 4, previewOpts)
```
**优势**: 避免错误分割，节省时间

---

### 4. 📐 图片缩放 (Resize)
```go
// 按比例缩放
imagesplit.ResizeByRatio("large.jpg", 0.5, opts)
// 适应最大尺寸
imagesplit.ResizeToFit("input.png", 800, 600, opts)
```
**优势**: 分割前统一尺寸，提高处理效率

---

### 5. ✂️ 自定义裁剪 (Crop)
```go
// 裁剪指定区域
imagesplit.Crop("input.jpg", 0, 0, 500, 500, opts)
// 居中裁剪
imagesplit.CropCenter("input.jpg", 800, 600, opts)
```
**优势**: 灵活裁剪单个区域，应用广泛

---

## 📦 其他实用功能

### 6. 🔄 格式转换 (Convert)
```go
imagesplit.Convert("input.png", "output.jpg", "jpeg", opts)
imagesplit.BatchConvert("input-dir", "output-dir", "png", opts)
```

### 7. 🎯 智能分割 (SmartSplit)
```go
// 自动计算参数，每块不超过 512x512
imagesplit.SplitToTargetSize("large.jpg", 512, opts)
// 自动分割成约 16 块
imagesplit.SplitToTargetCount("input.png", 16, opts)
```

### 8. 🔃 旋转翻转 (Rotate/Flip)
```go
imagesplit.Rotate90("input.jpg", true, opts)  // 顺时针旋转
imagesplit.FlipHorizontal("input.jpg", opts)  // 水平翻转
```

### 9. 🖼️ 缩略图生成 (Thumbnail)
```go
imagesplit.GenerateThumbnail("large.jpg", 200, 200, opts)
```

### 10. ⚡ 并发处理 (Concurrent)
```go
// 使用 4 个 worker 并发处理
imagesplit.SplitDirectoryConcurrent(inputDir, outputDir, cfg, 4)
```

---

## 💡 实现优先级

**P0 - 核心必备**:
- ✅ MergeGrid (合并)
- ✅ GetImageInfo (信息)
- ✅ PreviewSplit (预览)

**P1 - 常用功能**:
- ⭐ Resize (缩放)
- ⭐ Convert (转换)
- ⭐ Crop (裁剪)

**P2 - 增强功能**:
- SmartSplit (智能分割)
- Rotate/Flip (旋转翻转)
- Concurrent (并发)

**P3 - 进阶功能**:
- Watermark (水印)
- Thumbnail (缩略图)
- Quality Analysis (质量分析)

---

## 📖 详细说明

完整的功能说明、API 设计和使用示例请查看 [FEATURE_SUGGESTIONS.md](./FEATURE_SUGGESTIONS.md)

---

## 🎯 快速决策指南

**如果你的用户经常需要**:
- 还原/拼接图片 → 实现 **MergeGrid**
- 调试分割参数 → 实现 **PreviewSplit** + **GetImageInfo**
- 预处理图片 → 实现 **Resize** + **Crop**
- 提高处理速度 → 实现 **Concurrent** 版本
- 灵活的图片操作 → 实现 **Convert** + **Rotate**

**建议的实现顺序**: 
1. GetImageInfo (最简单，立即提升体验)
2. MergeGrid (最有价值，补齐核心功能)
3. PreviewSplit (中等难度，避免错误操作)
4. Resize (常用需求)
5. 根据用户反馈继续添加
