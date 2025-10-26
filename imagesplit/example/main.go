package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cto-new/imagesplit/imagesplit"
)

func main() {
	outputDir := filepath.Join("output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("create output directory: %v", err)
	}

	inputImage := filepath.Join(outputDir, "test.png")
	// if err := testdata.WriteGradientPNG(inputImage); err != nil {
	// 	log.Fatalf("prepare input image: %v", err)
	// }

	gridFiles, err := imagesplit.GridSplit(inputImage, 2, 3, imagesplit.SplitOptions{
		OutputDir:  outputDir,
		FilePrefix: "sample_grid",
	})
	if err != nil {
		log.Fatalf("grid split failed: %v", err)
	}
	fmt.Println("Grid split results:")
	for _, path := range gridFiles {
		fmt.Println(" -", path)
	}

	tileFiles, err := imagesplit.TileSplit(inputImage, 750, 100, imagesplit.SplitOptions{
		OutputDir:  outputDir,
		FilePrefix: "sample_tile",
		Format:     "jpeg",
		Quality:    80,
	})
	if err != nil {
		log.Fatalf("tile split failed: %v", err)
	}
	fmt.Println("\nTile split results:")
	for _, path := range tileFiles {
		fmt.Println(" -", path)
	}
}
