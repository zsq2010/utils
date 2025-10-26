package imagesplit

import (
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "strings"
    "testing"

    testdata "github.com/cto-new/imagesplit/imagesplit/testdata"
)

func createSampleImages(t *testing.T) (string, string) {
    t.Helper()
    dir := t.TempDir()

    pngPath := filepath.Join(dir, "gradient.png")
    if err := testdata.WriteGradientPNG(pngPath); err != nil {
        t.Fatalf("write gradient png: %v", err)
    }

    jpegPath := filepath.Join(dir, "blocks.jpg")
    if err := testdata.WriteBlocksJPEG(jpegPath); err != nil {
        t.Fatalf("write blocks jpeg: %v", err)
    }

    return pngPath, jpegPath
}

func TestGridSplitPNG(t *testing.T) {
    pngPath, _ := createSampleImages(t)
    outDir := t.TempDir()

    files, err := GridSplit(pngPath, 3, 4, SplitOptions{OutputDir: outDir})
    if err != nil {
        t.Fatalf("GridSplit returned error: %v", err)
    }
    if len(files) != 12 {
        t.Fatalf("expected 12 tiles, got %d", len(files))
    }

    expectedHeights := []int{4, 3, 3}
    expectedWidths := []int{3, 3, 2, 2}
    cols := 4

    for i, path := range files {
        if _, err := os.Stat(path); err != nil {
            t.Fatalf("expected file %s: %v", path, err)
        }

        f, err := os.Open(path)
        if err != nil {
            t.Fatalf("open tile: %v", err)
        }
        img, err := png.Decode(f)
        f.Close()
        if err != nil {
            t.Fatalf("decode tile png: %v", err)
        }
        bounds := img.Bounds()
        row := i / cols
        col := i % cols
        if bounds.Dx() != expectedWidths[col] || bounds.Dy() != expectedHeights[row] {
            t.Errorf("unexpected tile size for row %d col %d: got %dx%d", row, col, bounds.Dx(), bounds.Dy())
        }
    }
}

func TestGridSplitJPEGInput(t *testing.T) {
    _, jpegPath := createSampleImages(t)
    outDir := t.TempDir()

    files, err := GridSplit(jpegPath, 2, 3, SplitOptions{OutputDir: outDir})
    if err != nil {
        t.Fatalf("GridSplit returned error: %v", err)
    }
    if len(files) != 6 {
        t.Fatalf("expected 6 tiles, got %d", len(files))
    }

    for _, path := range files {
        if filepath.Ext(path) != ".jpg" {
            t.Errorf("expected jpg extension, got %s", filepath.Ext(path))
        }
        f, err := os.Open(path)
        if err != nil {
            t.Fatalf("open tile: %v", err)
        }
        if _, err := jpeg.Decode(f); err != nil {
            t.Fatalf("decode jpeg tile: %v", err)
        }
        f.Close()
    }
}

func TestTileSplitRemainders(t *testing.T) {
    pngPath, _ := createSampleImages(t)
    outDir := t.TempDir()

    files, err := TileSplit(pngPath, 4, 3, SplitOptions{OutputDir: outDir, FilePrefix: "tiles"})
    if err != nil {
        t.Fatalf("TileSplit returned error: %v", err)
    }
    if len(files) != 12 {
        t.Fatalf("expected 12 tiles, got %d", len(files))
    }

    last := files[len(files)-1]
    f, err := os.Open(last)
    if err != nil {
        t.Fatalf("open last tile: %v", err)
    }
    img, err := png.Decode(f)
    f.Close()
    if err != nil {
        t.Fatalf("decode last tile: %v", err)
    }
    bounds := img.Bounds()
    if bounds.Dx() != 2 || bounds.Dy() != 1 {
        t.Errorf("expected last tile size 2x1, got %dx%d", bounds.Dx(), bounds.Dy())
    }

    for i, path := range files {
        if !strings.Contains(filepath.Base(path), "tiles_tile_") {
            t.Errorf("unexpected filename for tile %d: %s", i, path)
        }
    }
}

func TestTileSplitJPEGOutput(t *testing.T) {
    pngPath, _ := createSampleImages(t)
    outDir := t.TempDir()

    files, err := TileSplit(pngPath, 6, 5, SplitOptions{OutputDir: outDir, Format: "jpeg", Quality: 75})
    if err != nil {
        t.Fatalf("TileSplit returned error: %v", err)
    }
    if len(files) != 4 {
        t.Fatalf("expected 4 tiles, got %d", len(files))
    }

    for _, path := range files {
        if filepath.Ext(path) != ".jpg" {
            t.Errorf("expected jpg extension, got %s", filepath.Ext(path))
        }
        f, err := os.Open(path)
        if err != nil {
            t.Fatalf("open tile: %v", err)
        }
        if _, err := jpeg.Decode(f); err != nil {
            t.Fatalf("decode jpeg tile: %v", err)
        }
        f.Close()
    }
}

func TestGridSplitInvalidParameters(t *testing.T) {
    if _, err := GridSplit("ignored.png", 0, 2, SplitOptions{}); err == nil {
        t.Fatalf("expected error for zero rows")
    }
    if _, err := GridSplit("ignored.png", 2, -1, SplitOptions{}); err == nil {
        t.Fatalf("expected error for negative cols")
    }
}

func TestTileSplitInvalidParameters(t *testing.T) {
    if _, err := TileSplit("ignored.png", 0, 5, SplitOptions{}); err == nil {
        t.Fatalf("expected error for zero width")
    }
    if _, err := TileSplit("ignored.png", 5, 0, SplitOptions{}); err == nil {
        t.Fatalf("expected error for zero height")
    }
}

func TestUnsupportedOutputFormat(t *testing.T) {
    pngPath, _ := createSampleImages(t)
    _, err := GridSplit(pngPath, 2, 2, SplitOptions{Format: "gif"})
    if err == nil {
        t.Fatalf("expected error for unsupported output format")
    }
}

func TestUnsupportedInputFormat(t *testing.T) {
    tmp := filepath.Join(t.TempDir(), "not_image.txt")
    if err := os.WriteFile(tmp, []byte("not an image"), 0o644); err != nil {
        t.Fatalf("write temp file: %v", err)
    }

    if _, err := GridSplit(tmp, 2, 2, SplitOptions{}); err == nil {
        t.Fatalf("expected error for invalid input format")
    }
}

func TestOutputDirectoryCreated(t *testing.T) {
    pngPath, _ := createSampleImages(t)
    base := t.TempDir()
    outDir := filepath.Join(base, "nested", "dir")

    files, err := TileSplit(pngPath, 5, 5, SplitOptions{OutputDir: outDir})
    if err != nil {
        t.Fatalf("TileSplit returned error: %v", err)
    }
    if len(files) == 0 {
        t.Fatalf("expected at least one output file")
    }

    if _, err := os.Stat(outDir); err != nil {
        t.Fatalf("expected output directory to be created: %v", err)
    }
}

func TestSplitDirectoryGridMode(t *testing.T) {
    inputDir := t.TempDir()
    if err := testdata.WriteGradientPNG(filepath.Join(inputDir, "gradient.png")); err != nil {
        t.Fatalf("write gradient png: %v", err)
    }
    if err := testdata.WriteBlocksJPEG(filepath.Join(inputDir, "blocks.jpg")); err != nil {
        t.Fatalf("write blocks jpeg: %v", err)
    }
    if err := os.WriteFile(filepath.Join(inputDir, "ignore.txt"), []byte("not an image"), 0o644); err != nil {
        t.Fatalf("write ignore file: %v", err)
    }

    outDir := t.TempDir()

    results, err := SplitDirectory(inputDir, outDir, DirectorySplitConfig{
        Mode: DirectorySplitModeGrid,
        Rows: 2,
        Cols: 2,
    })
    if err != nil {
        t.Fatalf("SplitDirectory returned error: %v", err)
    }
    if len(results) != 2 {
        t.Fatalf("expected 2 processed images, got %d", len(results))
    }

    for inputPath, files := range results {
        if len(files) != 4 {
            t.Fatalf("expected 4 tiles for %s, got %d", inputPath, len(files))
        }
        base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
        imageDir := filepath.Join(outDir, base)
        if stat, err := os.Stat(imageDir); err != nil || !stat.IsDir() {
            t.Fatalf("expected directory %s to exist", imageDir)
        }
        for _, file := range files {
            if filepath.Dir(file) != imageDir {
                t.Fatalf("expected file %s to be inside %s", file, imageDir)
            }
        }
    }
}

func TestSplitDirectoryTileMode(t *testing.T) {
    inputDir := t.TempDir()
    if err := testdata.WriteGradientPNG(filepath.Join(inputDir, "gradient.png")); err != nil {
        t.Fatalf("write gradient png: %v", err)
    }

    outDir := t.TempDir()

    results, err := SplitDirectory(inputDir, outDir, DirectorySplitConfig{
        Mode:       DirectorySplitModeTile,
        TileWidth:  4,
        TileHeight: 3,
        Options: SplitOptions{
            Format: "jpeg",
            Quality: 80,
        },
    })
    if err != nil {
        t.Fatalf("SplitDirectory returned error: %v", err)
    }
    if len(results) != 1 {
        t.Fatalf("expected 1 processed image, got %d", len(results))
    }
    for _, files := range results {
        for _, file := range files {
            if filepath.Ext(file) != ".jpg" {
                t.Fatalf("expected jpg output when forcing jpeg format, got %s", filepath.Ext(file))
            }
        }
    }
}

func TestSplitDirectoryHandlesDuplicateNames(t *testing.T) {
    inputDir := t.TempDir()
    if err := testdata.WriteGradientPNG(filepath.Join(inputDir, "photo.png")); err != nil {
        t.Fatalf("write gradient png: %v", err)
    }
    if err := testdata.WriteBlocksJPEG(filepath.Join(inputDir, "photo.jpg")); err != nil {
        t.Fatalf("write blocks jpeg: %v", err)
    }

    outDir := t.TempDir()

    results, err := SplitDirectory(inputDir, outDir, DirectorySplitConfig{
        Mode: DirectorySplitModeGrid,
        Rows: 2,
        Cols: 2,
    })
    if err != nil {
        t.Fatalf("SplitDirectory returned error: %v", err)
    }
    if len(results) != 2 {
        t.Fatalf("expected 2 processed images, got %d", len(results))
    }

    dirs := make(map[string]struct{})
    for _, files := range results {
        if len(files) == 0 {
            t.Fatalf("expected at least one output file per image")
        }
        dir := filepath.Base(filepath.Dir(files[0]))
        dirs[dir] = struct{}{}
    }
    if len(dirs) != 2 {
        t.Fatalf("expected unique directories per image, got %v", dirs)
    }
}
