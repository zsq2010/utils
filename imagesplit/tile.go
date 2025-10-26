package imagesplit

import (
    "fmt"
    "image"
)

func tileSplit(inputPath string, tileWidth, tileHeight int, opts SplitOptions) ([]string, error) {
    if tileWidth <= 0 {
        return nil, fmt.Errorf("tileWidth must be greater than zero")
    }
    if tileHeight <= 0 {
        return nil, fmt.Errorf("tileHeight must be greater than zero")
    }

    ctx, err := prepareSplit(inputPath, opts)
    if err != nil {
        return nil, err
    }

    result := []string{}
    index := 0

    for y := ctx.bounds.Min.Y; y < ctx.bounds.Max.Y; y += tileHeight {
        h := tileHeight
        if y+h > ctx.bounds.Max.Y {
            h = ctx.bounds.Max.Y - y
        }
        if h <= 0 {
            break
        }

        for x := ctx.bounds.Min.X; x < ctx.bounds.Max.X; x += tileWidth {
            w := tileWidth
            if x+w > ctx.bounds.Max.X {
                w = ctx.bounds.Max.X - x
            }
            if w <= 0 {
                break
            }

            rect := image.Rect(x, y, x+w, y+h)
            name := fmt.Sprintf("%s_tile_%d", ctx.options.prefix, index)
            output, err := saveTile(ctx.img, rect, ctx.options, name)
            if err != nil {
                return nil, err
            }
            result = append(result, output)
            index++
        }
    }

    return result, nil
}
