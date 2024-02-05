package compo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"math/rand"
)

type Type int

const (
	TypeNone Type = iota
	TypeShip
	TypeAsteroid
	TypePhoton
)

type Compo interface {
	Update(screenWidth, screenHeight int, resChanged bool) []Compo
	Draw(screen *ebiten.Image)
	Type() Type
	Bounds() (x, y, radius float64)
	Destroy()
}

func WrapScreen(x, y float64, screenWidth, screenHeight int) (xo, yo float64) {
	xo = x
	yo = y
	dxF := float64(screenWidth)
	dyF := float64(screenHeight)
	if xo > dxF {
		xo -= dxF
	}
	if xo < 0 {
		xo += dxF
	}
	if yo > dyF {
		yo -= dyF
	}
	if yo < 0 {
		yo += dyF
	}
	return xo, yo
}

type PointsProvider interface {
	Points() uint16
}

type BoundsProvider interface {
	Bounds() (x, y, radius float64)
}

type BufferProvider interface {
	Buffer() BoundsProvider
}

func Overlaps(one, two BoundsProvider) bool {
	if one == nil || two == nil {
		return false
	}
	x1, y1, r1 := one.Bounds()
	x2, y2, r2 := two.Bounds()
	return math.Sqrt(math.Pow(x1-x2, 2)+math.Pow(y1-y2, 2)) < r1+r2
}

func RandomColor(lower uint8) color.RGBA {
	return color.RGBA{R: RandomUint8(lower), G: RandomUint8(lower), B: RandomUint8(lower), A: 0xff}
}

func RandomUint8(lower uint8) uint8 {
	return uint8(rand.Intn(int(255-lower))) + lower
}

func RandomNegate(input float64) float64 {
	if rand.Intn(2) == 0 {
		return -input
	}
	return input
}
