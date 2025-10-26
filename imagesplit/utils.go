package imagesplit

import (
    "fmt"
    "image"
    "image/draw"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "strings"
)

type normalizedOptions struct {
    outputDir string
    prefix    string
    format    string
    extension string
    quality   int
}

type splitContext struct {
    img     image.Image
    bounds  image.Rectangle
    options normalizedOptions
}

func prepareSplit(inputPath string, opts SplitOptions) (*splitContext, error) {
    if inputPath == "" {
        return nil, fmt.Errorf("input path is required")
    }

    img, srcFormat, err := loadImage(inputPath)
    if err != nil {
        return nil, err
    }

    normalized, err := normalizeOptions(inputPath, opts, srcFormat)
    if err != nil {
        return nil, err
    }

    return &splitContext{
        img:     img,
        bounds:  img.Bounds(),
        options: normalized,
    }, nil
}

func loadImage(path string) (image.Image, string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, "", fmt.Errorf("open image: %w", err)
    }
    defer f.Close()

    img, format, err := image.Decode(f)
    if err != nil {
        return nil, "", fmt.Errorf("decode image: %w", err)
    }

    switch strings.ToLower(format) {
    case "jpeg", "jpg", "png":
        return img, format, nil
    default:
        return nil, "", fmt.Errorf("unsupported image format: %s", format)
    }
}

func normalizeOptions(inputPath string, opts SplitOptions, sourceFormat string) (normalizedOptions, error) {
    format := strings.TrimSpace(strings.ToLower(opts.Format))
    if format == "" {
        format = strings.ToLower(sourceFormat)
    }

    var extension string
    switch format {
    case "jpeg", "jpg":
        format = "jpeg"
        extension = "jpg"
    case "png":
        extension = "png"
    default:
        return normalizedOptions{}, fmt.Errorf("unsupported output format: %s", format)
    }

    quality := opts.Quality
    if quality <= 0 || quality > 100 {
        quality = 90
    }

    outputDir := opts.OutputDir
    if strings.TrimSpace(outputDir) == "" {
        outputDir = filepath.Dir(inputPath)
    }

    if err := os.MkdirAll(outputDir, 0o755); err != nil {
        return normalizedOptions{}, fmt.Errorf("create output directory: %w", err)
    }

    prefix := strings.TrimSpace(opts.FilePrefix)
    if prefix == "" {
        base := filepath.Base(inputPath)
        if ext := filepath.Ext(base); ext != "" {
            base = base[:len(base)-len(ext)]
        }
        if base == "" {
            base = "tile"
        }
        prefix = base
    }

    return normalizedOptions{
        outputDir: outputDir,
        prefix:    prefix,
        format:    format,
        extension: extension,
        quality:   quality,
    }, nil
}

func saveTile(img image.Image, rect image.Rectangle, opts normalizedOptions, name string) (string, error) {
    if rect.Dx() <= 0 || rect.Dy() <= 0 {
        return "", fmt.Errorf("invalid tile dimensions: %dx%d", rect.Dx(), rect.Dy())
    }

    tile := cropImage(img, rect)
    filename := fmt.Sprintf("%s.%s", name, opts.extension)
    outputPath := filepath.Join(opts.outputDir, filename)

    file, err := os.Create(outputPath)
    if err != nil {
        return "", fmt.Errorf("create output file: %w", err)
    }
    defer func() {
        file.Close()
        if err != nil {
            os.Remove(outputPath)
        }
    }()

    switch opts.format {
    case "png":
        err = png.Encode(file, tile)
    case "jpeg":
        err = jpeg.Encode(file, tile, &jpeg.Options{Quality: opts.quality})
    default:
        err = fmt.Errorf("unsupported output format: %s", opts.format)
    }

    if err != nil {
        return "", fmt.Errorf("encode image: %w", err)
    }

    return outputPath, nil
}

func cropImage(img image.Image, rect image.Rectangle) *image.RGBA {
    dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
    draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
    return dst
}

func distributeSize(total, parts int) ([]int, error) {
    if parts <= 0 {
        return nil, fmt.Errorf("parts must be positive")
    }
    base := total / parts
    remainder := total % parts

    sizes := make([]int, parts)
    for i := 0; i < parts; i++ {
        sizes[i] = base
        if i < remainder {
            sizes[i]++
        }
        if sizes[i] == 0 {
            return nil, fmt.Errorf("image dimension %d too small for %d segments", total, parts)
        }
    }
    return sizes, nil
}
