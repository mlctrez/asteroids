package compo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
)

type explosion struct {
	x    float64
	y    float64
	vx   float64
	vy   float64
	size float64
}

func (e *explosion) Update(screenWidth, screenHeight int, resChanged bool) []Compo {
	e.x, e.y = WrapScreen(e.x+e.vx, e.y+e.vy, screenWidth, screenHeight)
	e.size += 1
	if e.size > 15 {
		return []Compo{}
	}
	return []Compo{e}
}

func (e *explosion) Draw(screen *ebiten.Image) {
	for i := 0; i < 360; i = i + 7 {
		particle := image.Point{
			X: int(e.x + math.Sin(float64(i)*toRadians)*e.size*float64(i%2+1)*float64(i%3+1)),
			Y: int(e.y - math.Cos(float64(i)*toRadians)*e.size*float64(i%2+1)*float64(i%3+1)),
		}
		col := RandomColor(100)
		for x := -1; x <= 1; x = x + 1 {
			for y := -1; y <= 1; y = y + 1 {
				screen.Set(particle.X+x, particle.Y+y, col)
			}
		}
	}
}

func (e *explosion) Type() Type {
	return TypeNone
}

func (e *explosion) Bounds() (x, y, radius float64) {
	return e.x, e.y, 0.0
}

func (e *explosion) Destroy() {
}

func Explosion(x, y, vx, vy float64) Compo {
	return &explosion{x: x, y: y, vx: vx, vy: vy, size: 1}
}
