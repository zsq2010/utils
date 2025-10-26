package imagesplit

import (
    "fmt"
    "image"
)

func gridSplit(inputPath string, rows, cols int, opts SplitOptions) ([]string, error) {
    if rows <= 0 {
        return nil, fmt.Errorf("rows must be greater than zero")
    }
    if cols <= 0 {
        return nil, fmt.Errorf("cols must be greater than zero")
    }

    ctx, err := prepareSplit(inputPath, opts)
    if err != nil {
        return nil, err
    }

    colWidths, err := distributeSize(ctx.bounds.Dx(), cols)
    if err != nil {
        return nil, err
    }
    rowHeights, err := distributeSize(ctx.bounds.Dy(), rows)
    if err != nil {
        return nil, err
    }

    result := make([]string, 0, rows*cols)

    y := ctx.bounds.Min.Y
    for r := 0; r < rows; r++ {
        h := rowHeights[r]
        x := ctx.bounds.Min.X
        for c := 0; c < cols; c++ {
            w := colWidths[c]
            rect := image.Rect(x, y, x+w, y+h)
            name := fmt.Sprintf("%s_row%d_col%d", ctx.options.prefix, r, c)
            output, err := saveTile(ctx.img, rect, ctx.options, name)
            if err != nil {
                return nil, err
            }
            result = append(result, output)
            x += w
        }
        y += h
    }

    return result, nil
}
