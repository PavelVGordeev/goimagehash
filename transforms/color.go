package transforms

import (
	"image"
	"math"
)

const (
	SPACES = 6
)
const (
	HChannel = iota
	SChannel
	IChannel
	YChannel
	CbChannel
	CrChannel
)

func min(x, y, z float64) float64 {
	if x <= y && x <= z {
		return x
	} else if y <= x && y <= z {
		return y
	} else {
		return z
	}
}

func GetColorSpaces(colorImg image.Image) [][][]float64 {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pixels := make([][][]float64, SPACES)
	for i := range pixels {
		pixels[i] = make([][]float64, h)
		for j := range pixels[i] {
			pixels[i][j] = make([]float64, w)
		}
	}
	for i := range pixels[0] {
		for j := range pixels[i] {
			color := colorImg.At(j, i)
			r, g, b, _ := color.RGBA()
			rn, gn, bn := float64(r)/0xffff, float64(g)/0xffff, float64(b)/0xffff
			minChannel := min(rn, gn, bn)
			pixels[SChannel][i][j] = 1 - 3*minChannel/(rn+gn+bn)
			pixels[IChannel][i][j] = 0.33333333333 * (rn + gn + bn)
			theta := math.Acos(0.5 * (2*rn - gn - bn) / math.Sqrt(rn*rn+bn*bn+gn*gn-rn*gn-rn*bn-gn*bn))
			pixels[HChannel][i][j] = 2*math.Pi - theta
			if bn <= gn {
				pixels[HChannel][i][j] = theta
			}
			pixels[YChannel][i][j] = 16 + 65.481*rn + 128.553*gn + 24.966*bn
			pixels[CbChannel][i][j] = 128 - 37.797*rn - 74.203*gn + 112*bn
			pixels[CrChannel][i][j] = 128 + 112*rn - 93.786*gn - 18.214*bn
		}
	}

	return pixels
}

func GetMomentum(pixels [][][]float64) []float64 {
	zmoments := make([]float64, SPACES)
	for space := 0; space < SPACES; space++ {
		//go func(space int) {
		for i := 0; i < len(pixels[space]); i++ {
			for j := 0; j < len(pixels[space][i]); j++ {
				zmoments[space] += pixels[space][i][j]
			}
		}
		//}(space)
	}
	return zmoments
}
