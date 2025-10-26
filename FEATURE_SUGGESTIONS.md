# 功能建议 (Feature Suggestions)

本文档列出了可以为 `imagesplit` 库增加的常用工具和功能建议。

## 📊 优先级分类

### 🔥 高优先级 (与核心功能强相关)

#### 1. 图片合并工具 (Merge/Stitch)

**功能描述**: 将分割的小图片重新合并成大图，是分割操作的逆向功能。

**建议 API**:
```go
// MergeGrid 按网格顺序合并图片
// tilePaths: 图片路径列表，按 row0_col0, row0_col1, ..., row1_col0 的顺序
// rows, cols: 网格的行列数
func MergeGrid(tilePaths []string, rows, cols int, outputPath string, opts MergeOptions) error

// MergeOptions 合并选项
type MergeOptions struct {
    Format  string // 输出格式 "png" 或 "jpeg"
    Quality int    // JPEG 质量 (1-100)
}
```

**使用场景**:
- 还原之前分割的图片
- 图片拼接/全景图制作
- 分布式图片处理后的结果合并

**示例**:
```go
tiles := []string{
    "output/image_row0_col0.png",
    "output/image_row0_col1.png",
    "output/image_row1_col0.png",
    "output/image_row1_col1.png",
}
err := imagesplit.MergeGrid(tiles, 2, 2, "merged.png", imagesplit.MergeOptions{
    Format: "png",
})
```

---

#### 2. 图片信息获取工具 (ImageInfo)

**功能描述**: 获取图片的基本信息，帮助用户在分割前了解图片属性，确定合适的分割参数。

**建议 API**:
```go
// GetImageInfo 获取图片的详细信息
func GetImageInfo(imagePath string) (*ImageInfo, error)

// ImageInfo 图片信息
type ImageInfo struct {
    Width      int    // 宽度（像素）
    Height     int    // 高度（像素）
    Format     string // 格式 "png", "jpeg" 等
    FileSize   int64  // 文件大小（字节）
    ColorModel string // 色彩模型 "RGBA", "NRGBA", "YCbCr" 等
}
```

**使用场景**:
- 分割前检查图片尺寸，计算合适的分割参数
- 批量处理时过滤不符合条件的图片
- 图片管理和分类

**示例**:
```go
info, err := imagesplit.GetImageInfo("input.jpg")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Size: %dx%d, Format: %s, FileSize: %d bytes\n", 
    info.Width, info.Height, info.Format, info.FileSize)

// 根据图片大小决定分割策略
if info.Width > 4000 && info.Height > 4000 {
    imagesplit.GridSplit("input.jpg", 4, 4, opts)
} else {
    imagesplit.GridSplit("input.jpg", 2, 2, opts)
}
```

---

#### 3. 分割预览工具 (PreviewSplit)

**功能描述**: 生成带网格线的预览图，显示分割效果，但不实际分割图片。

**建议 API**:
```go
// PreviewGridSplit 生成网格分割预览图
func PreviewGridSplit(inputPath string, rows, cols int, opts PreviewOptions) (string, error)

// PreviewTileSplit 生成固定尺寸分割预览图
func PreviewTileSplit(inputPath string, tileWidth, tileHeight int, opts PreviewOptions) (string, error)

// PreviewOptions 预览选项
type PreviewOptions struct {
    OutputPath  string // 输出路径
    LineColor   string // 网格线颜色 "red", "blue", "black" 等
    LineWidth   int    // 网格线宽度（像素）
    ShowNumbers bool   // 是否显示分块编号
}
```

**使用场景**:
- 在实际分割前预览效果，避免参数设置错误
- 调试分割参数
- 向用户展示分割方案

**示例**:
```go
preview, err := imagesplit.PreviewGridSplit("input.png", 3, 4, imagesplit.PreviewOptions{
    OutputPath:  "preview.png",
    LineColor:   "red",
    LineWidth:   2,
    ShowNumbers: true,
})
```

---

### 📦 中优先级 (增强易用性)

#### 4. 图片缩放工具 (Resize)

**功能描述**: 调整图片大小，可用于分割前的预处理。

**建议 API**:
```go
// Resize 调整图片到指定尺寸
func Resize(inputPath string, width, height int, opts ResizeOptions) (string, error)

// ResizeByRatio 按比例缩放图片
func ResizeByRatio(inputPath string, ratio float64, opts ResizeOptions) (string, error)

// ResizeToFit 缩放图片使其适合指定的最大尺寸（保持宽高比）
func ResizeToFit(inputPath string, maxWidth, maxHeight int, opts ResizeOptions) (string, error)

// ResizeOptions 缩放选项
type ResizeOptions struct {
    OutputPath      string // 输出路径
    Format          string // 输出格式
    Quality         int    // JPEG 质量
    Interpolation   string // 插值方法 "nearest", "bilinear", "bicubic"
    KeepAspectRatio bool   // 是否保持宽高比
}
```

**使用场景**:
- 分割前统一图片尺寸
- 生成不同分辨率版本
- 减小图片以提高处理速度

**示例**:
```go
// 缩放到 50%
resized, err := imagesplit.ResizeByRatio("large.jpg", 0.5, imagesplit.ResizeOptions{
    OutputPath: "medium.jpg",
    Format:     "jpeg",
    Quality:    85,
})

// 缩放后再分割
imagesplit.GridSplit(resized, 3, 3, opts)
```

---

#### 5. 格式转换工具 (Convert)

**功能描述**: 独立的格式转换功能，虽然现有代码支持输出格式转换，但独立出来更方便使用。

**建议 API**:
```go
// Convert 转换图片格式
func Convert(inputPath, outputPath string, format string, opts ConvertOptions) error

// BatchConvert 批量转换目录中的所有图片
func BatchConvert(inputDir, outputDir string, format string, opts ConvertOptions) ([]string, error)

// ConvertOptions 转换选项
type ConvertOptions struct {
    Quality       int  // JPEG 质量
    OverwriteExisting bool // 是否覆盖已存在的文件
}
```

**使用场景**:
- PNG 转 JPEG 以减小文件大小
- JPEG 转 PNG 以保留透明度支持
- 批量格式统一

**示例**:
```go
// 单个文件转换
err := imagesplit.Convert("input.png", "output.jpg", "jpeg", imagesplit.ConvertOptions{
    Quality: 90,
})

// 批量转换
results, err := imagesplit.BatchConvert("input-dir", "output-dir", "png", 
    imagesplit.ConvertOptions{})
```

---

#### 6. 图片裁剪工具 (Crop)

**功能描述**: 自定义区域裁剪，不同于批量分割，更灵活地裁剪单个区域。

**建议 API**:
```go
// Crop 裁剪指定区域
func Crop(inputPath string, x, y, width, height int, opts SplitOptions) (string, error)

// CropCenter 从中心裁剪指定尺寸
func CropCenter(inputPath string, width, height int, opts SplitOptions) (string, error)

// CropMultiple 裁剪多个指定区域
func CropMultiple(inputPath string, regions []CropRegion, opts SplitOptions) ([]string, error)

// CropRegion 裁剪区域定义
type CropRegion struct {
    X, Y          int
    Width, Height int
    Name          string // 可选的区域名称
}
```

**使用场景**:
- 提取图片的特定部分
- 去除边缘无用区域
- 批量裁剪特定位置（如去除水印）

**示例**:
```go
// 裁剪左上角 500x500 区域
cropped, err := imagesplit.Crop("input.jpg", 0, 0, 500, 500, imagesplit.SplitOptions{
    OutputDir: "output",
    Format:    "png",
})

// 居中裁剪
centered, err := imagesplit.CropCenter("input.jpg", 800, 600, opts)
```

---

#### 7. 智能分割工具 (SmartSplit)

**功能描述**: 根据目标自动计算最优分割参数。

**建议 API**:
```go
// SplitToTargetSize 分割到目标尺寸（每个块不超过指定大小）
func SplitToTargetSize(inputPath string, maxTileSize int, opts SplitOptions) ([]string, error)

// SplitToTargetCount 分割到目标数量（尽可能均匀分割成指定数量）
func SplitToTargetCount(inputPath string, targetCount int, opts SplitOptions) ([]string, error)

// SplitToTargetFileSize 分割使每个文件不超过指定大小
func SplitToTargetFileSize(inputPath string, maxFileSizeKB int, opts SplitOptions) ([]string, error)
```

**使用场景**:
- 自动优化分割参数
- 适配不同显示设备的最大分辨率
- 控制输出文件大小

**示例**:
```go
// 自动分割，每块不超过 512x512
tiles, err := imagesplit.SplitToTargetSize("large.jpg", 512, opts)

// 分割成约 16 块
tiles, err := imagesplit.SplitToTargetCount("input.png", 16, opts)
```

---

### 🎨 低优先级 (进阶功能)

#### 8. 图片旋转/翻转工具 (Rotate/Flip)

**功能描述**: 调整图片方向。

**建议 API**:
```go
// Rotate90 旋转 90 度
func Rotate90(inputPath string, clockwise bool, opts SplitOptions) (string, error)

// Rotate180 旋转 180 度
func Rotate180(inputPath string, opts SplitOptions) (string, error)

// FlipHorizontal 水平翻转
func FlipHorizontal(inputPath string, opts SplitOptions) (string, error)

// FlipVertical 垂直翻转
func FlipVertical(inputPath string, opts SplitOptions) (string, error)
```

**使用场景**:
- 修正图片方向
- 数据增强（机器学习）

---

#### 9. 缩略图生成 (Thumbnail)

**功能描述**: 快速生成缩略图用于预览。

**建议 API**:
```go
// GenerateThumbnail 生成缩略图（保持宽高比）
func GenerateThumbnail(inputPath string, maxWidth, maxHeight int, opts SplitOptions) (string, error)

// GenerateThumbnails 批量生成缩略图
func GenerateThumbnails(inputDir, outputDir string, maxWidth, maxHeight int, opts SplitOptions) ([]string, error)
```

---

#### 10. 并发批量处理

**功能描述**: 提高大批量处理的性能。

**建议 API**:
```go
// SplitDirectoryConcurrent 并发批量处理
func SplitDirectoryConcurrent(inputDir, outputDir string, cfg DirectorySplitConfig, workers int) (map[string][]string, error)
```

---

#### 11. 水印添加 (Watermark)

**功能描述**: 给图片添加文字或图片水印。

**建议 API**:
```go
// AddWatermark 添加图片水印
func AddWatermark(inputPath, watermarkPath string, position WatermarkPosition, opts WatermarkOptions) (string, error)

// AddTextWatermark 添加文字水印
func AddTextWatermark(inputPath, text string, opts TextWatermarkOptions) (string, error)

type WatermarkPosition string
const (
    WatermarkTopLeft     WatermarkPosition = "top-left"
    WatermarkTopRight    WatermarkPosition = "top-right"
    WatermarkBottomLeft  WatermarkPosition = "bottom-left"
    WatermarkBottomRight WatermarkPosition = "bottom-right"
    WatermarkCenter      WatermarkPosition = "center"
)
```

---

#### 12. 图片质量分析

**功能描述**: 分析图片质量指标。

**建议 API**:
```go
// AnalyzeQuality 分析图片质量
func AnalyzeQuality(imagePath string) (*QualityReport, error)

type QualityReport struct {
    Sharpness   float64 // 清晰度
    Brightness  float64 // 亮度
    Contrast    float64 // 对比度
    Compression float64 // 压缩率
}
```

---

#### 13. 边框添加

**功能描述**: 给图片添加边框或内边距。

**建议 API**:
```go
// AddBorder 添加边框
func AddBorder(inputPath string, borderWidth int, color string, opts SplitOptions) (string, error)

// AddPadding 添加内边距
func AddPadding(inputPath string, padding int, color string, opts SplitOptions) (string, error)
```

---

## 🎯 实现建议

### 第一阶段：核心功能补全
1. **MergeGrid** - 合并功能（最重要的缺失功能）
2. **GetImageInfo** - 信息获取（使用频率高）
3. **PreviewSplit** - 预览功能（避免错误操作）

### 第二阶段：实用工具
4. **Resize** - 缩放功能（常见需求）
5. **Convert** - 独立转换工具
6. **Crop** - 单区域裁剪

### 第三阶段：智能优化
7. **SmartSplit** - 智能分割
8. **SplitDirectoryConcurrent** - 并发处理

### 第四阶段：进阶功能
9. 根据用户反馈和实际需求添加其他功能

---

## 📝 注意事项

1. **保持 API 一致性**: 新功能应该遵循现有的命名和参数风格
2. **错误处理**: 继续保持完善的错误处理和验证
3. **测试覆盖**: 每个新功能都应该有对应的单元测试
4. **文档更新**: 及时更新 README.md 和代码注释
5. **向后兼容**: 确保新功能不会破坏现有 API

---

## 🤝 贡献

欢迎对这些建议提出反馈或直接贡献代码实现。如有其他功能建议，请提交 Issue 讨论。
