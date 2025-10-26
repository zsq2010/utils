package imagesplit

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// DirectorySplitMode indicates how images should be split when processing a directory.
type DirectorySplitMode string

const (
    // DirectorySplitModeGrid splits each image using a grid (rows x columns).
    DirectorySplitModeGrid DirectorySplitMode = "grid"
    // DirectorySplitModeTile splits each image into fixed-size tiles.
    DirectorySplitModeTile DirectorySplitMode = "tile"
)

// DirectorySplitConfig configures how images in a directory are processed.
type DirectorySplitConfig struct {
    Mode       DirectorySplitMode
    Rows       int
    Cols       int
    TileWidth  int
    TileHeight int
    Options    SplitOptions
}

// SplitDirectory walks through the input directory, splitting every supported image
// using the provided configuration. Each image gets its own subdirectory inside
// outputDir, named after the image file (duplicate names receive numeric/format suffixes).
// The function returns a map keyed by the input image path containing the generated file paths.
func SplitDirectory(inputDir, outputDir string, cfg DirectorySplitConfig) (map[string][]string, error) {
    if strings.TrimSpace(inputDir) == "" {
        return nil, fmt.Errorf("input directory is required")
    }
    if strings.TrimSpace(outputDir) == "" {
        return nil, fmt.Errorf("output directory is required")
    }

    if err := validateDirectoryConfig(cfg); err != nil {
        return nil, err
    }

    info, err := os.Stat(inputDir)
    if err != nil {
        return nil, fmt.Errorf("stat input directory: %w", err)
    }
    if !info.IsDir() {
        return nil, fmt.Errorf("input path is not a directory: %s", inputDir)
    }

    if err := os.MkdirAll(outputDir, 0o755); err != nil {
        return nil, fmt.Errorf("create output directory: %w", err)
    }

    entries, err := os.ReadDir(inputDir)
    if err != nil {
        return nil, fmt.Errorf("read input directory: %w", err)
    }

    results := make(map[string][]string)
    usedDirs := make(map[string]struct{})

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        name := entry.Name()
        ext := strings.ToLower(filepath.Ext(name))
        if ext == "" {
            continue
        }
        switch ext {
        case ".png", ".jpg", ".jpeg":
        default:
            continue
        }

        inputPath := filepath.Join(inputDir, name)
        base := strings.TrimSuffix(name, ext)
        if base == "" {
            base = strings.TrimPrefix(name, ".")
            if base == "" {
                base = "image"
            }
        }
        dirName := base
        if _, exists := usedDirs[dirName]; exists {
            altBase := dirName
            if suffix := strings.TrimPrefix(ext, "."); suffix != "" {
                altBase = fmt.Sprintf("%s_%s", dirName, suffix)
            }
            candidate := altBase
            counter := 1
            for {
                if _, taken := usedDirs[candidate]; !taken {
                    dirName = candidate
                    break
                }
                counter++
                candidate = fmt.Sprintf("%s_%d", altBase, counter)
            }
        }
        usedDirs[dirName] = struct{}{}
        subOutput := filepath.Join(outputDir, dirName)

        if err := os.RemoveAll(subOutput); err != nil {
            return nil, fmt.Errorf("remove existing output directory: %w", err)
        }

        opts := cfg.Options
        opts.OutputDir = subOutput

        var generated []string
        switch cfg.Mode {
        case DirectorySplitModeGrid:
            generated, err = GridSplit(inputPath, cfg.Rows, cfg.Cols, opts)
        case DirectorySplitModeTile:
            generated, err = TileSplit(inputPath, cfg.TileWidth, cfg.TileHeight, opts)
        default:
            err = fmt.Errorf("unsupported directory split mode: %s", cfg.Mode)
        }
        if err != nil {
            return nil, fmt.Errorf("split image %s: %w", inputPath, err)
        }
        results[inputPath] = generated
    }

    return results, nil
}

func validateDirectoryConfig(cfg DirectorySplitConfig) error {
    switch cfg.Mode {
    case DirectorySplitModeGrid:
        if cfg.Rows <= 0 {
            return fmt.Errorf("rows must be greater than zero for grid mode")
        }
        if cfg.Cols <= 0 {
            return fmt.Errorf("cols must be greater than zero for grid mode")
        }
    case DirectorySplitModeTile:
        if cfg.TileWidth <= 0 {
            return fmt.Errorf("tileWidth must be greater than zero for tile mode")
        }
        if cfg.TileHeight <= 0 {
            return fmt.Errorf("tileHeight must be greater than zero for tile mode")
        }
    default:
        return fmt.Errorf("unsupported directory split mode: %s", cfg.Mode)
    }
    return nil
}
