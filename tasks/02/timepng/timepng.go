package timepng

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"time"
)

const (
	defaultWidth  = 3
	defaultHeight = 5
	spaceWidth    = 1
)

// TimePNG записывает в `out` картинку в формате png с текущим временем
func TimePNG(out io.Writer, t time.Time, c color.Color, scale int) {
	img := buildTimeImage(t, c, scale)
	png.Encode(out, img)
}

// buildTimeImage создает новое изображение с временем `t`
func buildTimeImage(t time.Time, c color.Color, scale int) *image.RGBA {
	width, height := defaultWidth*scale*5+spaceWidth*scale*4, defaultHeight*scale

	timeStr := t.Format("15:04")

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	widthStep := defaultWidth * scale
	spaceStep := spaceWidth * scale

	for i, v := range timeStr {
		wStart, hStart, wEnd, hEnd := widthStep*i+spaceStep*i, 0, widthStep*(i+1)+spaceStep*i, height
		sub := img.SubImage(image.Rect(wStart, hStart, wEnd, hEnd)).(*image.RGBA)
		sub.Rect = image.Rect(0, 0, widthStep, height)
		fillWithMask(sub, nums[v], c, scale)
	}

	return img
}

// fillWithMask заполняет изображение `img` цветом `c` по маске `mask`. Маска `mask`
// должна иметь пропорциональные размеры `img` с учетом фактора `scale`
// NOTE: Так как это вспомогательная функция, можно считать, что mask имеет размер (3x5)
func fillWithMask(img *image.RGBA, mask []int, c color.Color, scale int) {
	width, height := defaultWidth, defaultHeight

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			maskValue := mask[x+y*width]
			scaleSquared := scale * scale
			for s := 0; s < scaleSquared; s++ {
				xExpr := x*scale + s%scale
				yExpr := y*scale + s/scale
				if maskValue != 0 {
					img.Set(xExpr, yExpr, c)
				}
			}
		}
	}
}

var nums = map[rune][]int{
	'0': {
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'1': {
		0, 1, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	},
	'2': {
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	'3': {
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	'4': {
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
	},
	'5': {
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	'6': {
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'7': {
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	'8': {
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'9': {
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	':': {
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0,
	},
}
