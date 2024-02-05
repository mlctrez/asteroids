package compo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

type asteroid struct {
	x           float64
	y           float64
	vx          float64
	vy          float64
	rotation    float64
	rotationDir float64
	color       color.Color
	// 1=small, 2=medium, 3=large
	size           int
	destroyed      bool
	vertMultiplier []float64
}

func (a *asteroid) Type() Type {
	return TypeAsteroid
}

func (a *asteroid) Bounds() (x, y, radius float64) {
	return a.x, a.y, scale * float64(a.size)
}

func (a *asteroid) Destroy() {
	a.destroyed = true
}

func (a *asteroid) Points() uint16 {
	switch a.size {
	case 1:
		return 100
	case 2:
		return 50
	case 3:
		return 20
	default:
		return 0
	}
}

func (a *asteroid) Update(screenWidth, screenHeight int, resChanged bool) []Compo {

	if a.destroyed {
		result := []Compo{Explosion(a.x, a.y, a.vx, a.vy)}
		if a.size > 1 {
			result = append(result, Asteroid(a.x, a.y, a.size-1), Asteroid(a.x, a.y, a.size-1))
		}
		return result
	}

	a.x, a.y = WrapScreen(a.x+a.vx, a.y+a.vy, screenWidth, screenHeight)
	a.rotation += a.rotationDir

	return []Compo{a}
}

func (a *asteroid) Draw(screen *ebiten.Image) {

	lineWidth := float32(a.size)

	asteroidScale := scale * float64(a.size)
	_ = asteroidScale
	num := len(a.vertMultiplier)
	for i := 0; i < num; i++ {
		to := i + 1
		if to == num {
			to = 0
		}
		var factor float64 = float64(360 / num)
		r1 := (factor*float64(i) + a.rotation) * toRadians
		r2 := (factor*float64(to) + a.rotation) * toRadians
		fx := a.x + math.Cos(r1)*asteroidScale*a.vertMultiplier[i]
		fy := a.y - math.Sin(r1)*asteroidScale*a.vertMultiplier[i]
		tx := a.x + math.Cos(r2)*asteroidScale*a.vertMultiplier[to]
		ty := a.y - math.Sin(r2)*asteroidScale*a.vertMultiplier[to]
		vector.StrokeLine(screen, float32(fx), float32(fy), float32(tx), float32(ty), lineWidth, a.color, true)
	}
}

func Asteroid(x, y float64, size int) Compo {
	speedMultiplier := float64(4-size) / 2
	vx := RandomNegate(rand.Float64()*speedMultiplier + 0.5)
	vy := RandomNegate(rand.Float64()*speedMultiplier + 0.5)
	ast := &asteroid{x: x, y: y, vx: vx, vy: vy, size: size}
	ast.color = RandomColor(180)

	var randomPart float64 = .2
	vertTotal := 16
	for i := 0; i < vertTotal; i++ {
		val := rand.Float64()*randomPart + (1 - randomPart)
		ast.vertMultiplier = append(ast.vertMultiplier, val)
	}
	ast.vertMultiplier[0] *= 0.5
	ast.vertMultiplier[rand.Intn(vertTotal/2)] *= 0.5

	ast.rotation = rand.Float64() * 360
	ast.rotationDir = RandomNegate((rand.Float64() * 0.5) + 0.5)
	return ast
}
