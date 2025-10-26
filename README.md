# imagesplit

`imagesplit` 是一个用于 Go 语言的图片分割库，支持按网格（行×列）和固定尺寸两种方式对 PNG / JPEG 图片进行分割。库提供简洁的 API，可以作为独立依赖被其他项目调用。

## 功能特性

- ✅ 支持 PNG (`.png`) 和 JPEG (`.jpg`, `.jpeg`) 格式
- ✅ 网格分割：按照指定的行列数自动生成小图块
- ✅ 固定尺寸分割：按照固定的宽高切割，自动处理边缘剩余区域
- ✅ 灵活的输出配置：输出目录、文件前缀、图片格式、JPEG 质量
- ✅ 完善的错误处理：格式不支持、参数错误、输出目录创建失败等

## 安装

```bash
go get github.com/zsq2010/utils
```

## 快速上手

### 单张图片分割

```go
package main

import (
    "fmt"
    "log"

    "github.com/zsq2010/utils/imagesplit"
)

func main() {
    files, err := imagesplit.GridSplit("input.png", 3, 4, imagesplit.SplitOptions{
        OutputDir:  "output",
        FilePrefix: "sample",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(files)
}
```

### 批量处理目录

```go
cfg := imagesplit.DirectorySplitConfig{
    Mode: imagesplit.DirectorySplitModeGrid,
    Rows: 3,
    Cols: 3,
    Options: imagesplit.SplitOptions{Format: "png"},
}
results, err := imagesplit.SplitDirectory("input-dir", "output-dir", cfg)
if err != nil {
    log.Fatal(err)
}
for src, generated := range results {
    fmt.Println("source:", src)
    for _, path := range generated {
        fmt.Println("  -", path)
    }
}
```

## API 文档

```go
func GridSplit(inputPath string, rows, cols int, opts imagesplit.SplitOptions) ([]string, error)
```
- `rows` / `cols`: 网格行列数，必须大于 0。
- `SplitOptions`
  - `OutputDir`: 输出目录（为空时使用原图所在目录）。
  - `FilePrefix`: 输出文件前缀（为空时使用原图文件名）。
  - `Format`: 输出格式（`"png"`、`"jpeg"`，为空使用原图格式）。
  - `Quality`: JPEG 质量，范围 1-100（默认 90）。
- 返回值为生成的文件路径列表。

```go
func TileSplit(inputPath string, tileWidth, tileHeight int, opts imagesplit.SplitOptions) ([]string, error)
```
- `tileWidth` / `tileHeight`: 图块宽高，必须大于 0。
- 其余参数与 `GridSplit` 一致。

```go
func SplitDirectory(inputDir, outputDir string, cfg imagesplit.DirectorySplitConfig) (map[string][]string, error)
```
- `inputDir`: 输入图片所在目录。
- `outputDir`: 输出根目录，每张图片会在该目录下创建一个以图片名命名的子目录。
- `cfg.Mode`: 分割模式（`imagesplit.DirectorySplitModeGrid` / `imagesplit.DirectorySplitModeTile`）。
- `cfg.Rows`, `cfg.Cols`: 网格模式的行列数。
- `cfg.TileWidth`, `cfg.TileHeight`: 固定尺寸模式的宽高。
- `cfg.Options`: 其它分割选项（输出格式、JPEG 质量等），`OutputDir` 会被自动覆盖为图片专属子目录。

## 命名规则

- 网格分割：`{prefix}_row{i}_col{j}.{ext}` → 例如：`image_row0_col2.png`
- 固定尺寸：`{prefix}_tile_{index}.{ext}` → 例如：`image_tile_5.jpg`

## 示例

仓库包含一个完整的示例程序，位于 `imagesplit/example/main.go`：

```bash
cd imagesplit/example
go run .
```

示例会自动生成一张示例 PNG 图片，并在 `imagesplit/example/output` 目录下演示网格分割、固定尺寸分割以及目录批量处理，同时给出 JPEG 输出示例。

## 测试图片

项目提供 `imagesplit/testdata` 辅助包，可动态生成内置的测试图片：

```go
import testdata "github.com/zsq2010/utils/imagesplit/testdata"

data, _ := testdata.GradientPNG()    // 获取示例 PNG 图片字节流
testdata.WriteBlocksJPEG("blocks.jpg") // 生成示例 JPEG 文件
```

这些工具可用于测试、示例或文档代码中，避免手动维护二进制测试资源。

## 测试

项目提供了覆盖主要功能和边界情况的单元测试，包括：

- 网格分割及无法整除时的边缘尺寸验证
- 固定尺寸分割及剩余像素处理
- PNG / JPEG 输入输出
- 错误处理与参数校验
- 输出目录自动创建

运行测试：

```bash
cd imagesplit
go test ./...
```

## 许可证

该项目基于 MIT License 发布，欢迎自由使用与贡献。
