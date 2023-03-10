package transforms

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockImage struct {
	bitmap [][]MockColor
}

func NewMockImage(rgb [][][3]uint32) MockImage {
	bitmap := make([][]MockColor, len(rgb))
	for i := range rgb {
		bitmap[i] = make([]MockColor, len(rgb[i]))
		for j := range rgb[i] {
			bitmap[i][j] = NewMockColor(rgb[i][j])
		}
	}
	return MockImage{bitmap}
}

type MockColor struct {
	r uint32
	g uint32
	b uint32
}

func NewMockColor(rgb [3]uint32) MockColor {
	return MockColor{rgb[0], rgb[1], rgb[2]}
}

type MockColorModel struct {
}

func (m MockColor) RGBA() (uint32, uint32, uint32, uint32) {
	return m.r, m.g, m.b, 0
}

func (m MockColorModel) Convert(_ color.Color) color.Color {
	return MockColor{}
}

func (m MockImage) At(x int, y int) color.Color {
	return m.bitmap[y][x]
}

func (m MockImage) Bounds() image.Rectangle {
	return image.Rectangle{Min: image.Point{
		X: 0,
		Y: 0,
	}, Max: image.Point{
		X: 1,
		Y: 1,
	}}
}

func (m MockImage) ColorModel() color.Model {
	return &MockColorModel{}
}

func TestGetColorSpaces(t *testing.T) {
	tests := []struct {
		name   string
		bitmap [][][3]uint32
		want   []float64
	}{
		{
			name: "red",
			bitmap: [][][3]uint32{{[3]uint32{
				0xffff}}},
			want: []float64{0, 1.0, 0.33333333, 81.481, 90.203, 240},
		},
		{
			name: "mixed",
			bitmap: [][][3]uint32{{[3]uint32{
				0x1234,
				0x5678,
				0x9012}}},
			want: []float64{3.616314546, 0.7804579289, 0.323887998, 78.1284303044, 163.28026813, 94.03500962},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetColorSpaces(NewMockImage(tt.bitmap))
			for i := range got {
				assert.InDelta(t, tt.want[i], got[i][0][0], EPS)
			}

		})
	}

}

func TestGetMomentum(t *testing.T) {
	tests := []struct {
		name   string
		pixels [][][]float64
		want   []float64
	}{
		{
			name: "onepixels",
			pixels: [][][]float64{
				{{3.616314546}},
				{{0.7804579289}},
				{{0.323887998}},
				{{78.1284303044}},
				{{163.28026813}},
				{{94.03500962}}},
			want: []float64{3.616314546, 0.7804579289, 0.323887998, 78.1284303044, 163.28026813, 94.03500962},
		},
		//{
		//	name:   "multiple pixels",
		//	pixels: [][][]float64{{[]float64{3.616314546, 0.7804579289, 0.323887998, 78.1284303044, 163.28026813, 94.03500962}}},
		//	want:   []float64{1, 2, 3, 4, 5, 6},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMomentum(tt.pixels)
			assert.InDeltaSlice(t, tt.want, got, EPS)
		})
	}

}
func TestMinofthree(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "first",
			args: args{
				x: -1.0,
				y: 0,
				z: 2.0,
			},
			want: -1.0,
		},
		{
			name: "second",
			args: args{
				x: 4.0,
				y: 3.0,
				z: 5.0,
			},
			want: 3.0,
		},
		{
			name: "third",
			args: args{
				x: 5.0,
				y: 4.0,
				z: 3.0,
			},
			want: 3.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, min(tt.args.x, tt.args.y, tt.args.z), "min(%v, %v, %v)", tt.args.x, tt.args.y, tt.args.z)
		})
	}
}
