// Package imagesplit provides utilities for splitting images into smaller tiles.
package imagesplit

// SplitOptions defines configurable options for image splitting operations.
type SplitOptions struct {
    // OutputDir is the directory where split images will be written. If empty,
    // the directory of the input image will be used.
    OutputDir string
    // FilePrefix is the prefix used when naming generated tiles. If empty,
    // the base name of the input image (without extension) will be used.
    FilePrefix string
    // Format determines the output image format. Supported values are "jpeg"
    // and "png" (case-insensitive). When left empty, the input image format is
    // used.
    Format string
    // Quality controls JPEG encoding quality (1-100). It is ignored for PNG
    // output. When set to 0, a default of 90 is used.
    Quality int
}

// GridSplit divides an input image into a grid defined by the provided number
// of rows and columns. It returns the list of generated file paths on success.
func GridSplit(inputPath string, rows, cols int, opts SplitOptions) ([]string, error) {
    return gridSplit(inputPath, rows, cols, opts)
}

// TileSplit divides an input image into tiles of the specified width and height
// (in pixels). It returns the list of generated file paths on success.
func TileSplit(inputPath string, tileWidth, tileHeight int, opts SplitOptions) ([]string, error) {
    return tileSplit(inputPath, tileWidth, tileHeight, opts)
}
