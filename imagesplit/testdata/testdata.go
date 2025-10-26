// Package testdata provides helper routines for generating sample images used
// in tests and examples.
package testdata

import (
    "bytes"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
)

// GradientPNG returns a small PNG image with a simple color gradient.
func GradientPNG() ([]byte, error) {
    img := gradientImage()
    buf := &bytes.Buffer{}
    if err := png.Encode(buf, img); err != nil {
        return nil, fmt.Errorf("encode gradient png: %w", err)
    }
    return buf.Bytes(), nil
}

// BlocksJPEG returns a JPEG image containing a smooth color transition.
func BlocksJPEG() ([]byte, error) {
    img := blocksImage()
    buf := &bytes.Buffer{}
    if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
        return nil, fmt.Errorf("encode blocks jpeg: %w", err)
    }
    return buf.Bytes(), nil
}

// WriteGradientPNG writes the gradient PNG image to the provided path.
func WriteGradientPNG(path string) error {
    data, err := GradientPNG()
    if err != nil {
        return err
    }
    return writeFile(path, data)
}

// WriteBlocksJPEG writes the blocks JPEG image to the provided path.
func WriteBlocksJPEG(path string) error {
    data, err := BlocksJPEG()
    if err != nil {
        return err
    }
    return writeFile(path, data)
}

func writeFile(path string, data []byte) error {
    dir := filepath.Dir(path)
    if dir != "." && dir != "" {
        if err := os.MkdirAll(dir, 0o755); err != nil {
            return fmt.Errorf("create directory: %w", err)
        }
    }
    if err := os.WriteFile(path, data, 0o644); err != nil {
        return fmt.Errorf("write file: %w", err)
    }
    return nil
}

func gradientImage() *image.RGBA {
    width, height := 10, 10
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r := uint8((x * 25) % 256)
            g := uint8((y * 25) % 256)
            b := uint8(((x + y) * 12) % 256)
            img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
        }
    }
    return img
}

func blocksImage() *image.RGBA {
    width, height := 12, 8
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r := uint8(float64(x) / float64(width-1) * 255)
            g := uint8(float64(y) / float64(height-1) * 255)
            img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: 180, A: 255})
        }
    }
    return img
}
