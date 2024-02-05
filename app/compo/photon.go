package compo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type photon struct {
	x              float64
	y              float64
	vx             float64
	vy             float64
	rotation       float64
	rotationChange float64
	distance       int
	done           PhotonDone
	destroyed      bool
}

func (p *photon) Type() Type {
	return TypePhoton
}

func (p *photon) Bounds() (x, y, radius float64) {
	return p.x, p.y, scale * photonScale
}

func (p *photon) Destroy() {
	p.destroyed = true
}

const photonScale float64 = 0.2

func (p *photon) Update(screenWidth, screenHeight int, resChanged bool) []Compo {
	p.x, p.y = WrapScreen(p.x+p.vx, p.y+p.vy, screenWidth, screenHeight)
	p.rotation += p.rotationChange
	p.distance--
	if p.destroyed || p.distance < 0 {
		p.done()
		return nil
	}
	return []Compo{p}
}

func (p *photon) Draw(screen *ebiten.Image) {
	one := image.Point{
		X: int(p.x + math.Sin(p.rotation*toRadians)*scale*photonScale),
		Y: int(p.y - math.Cos(p.rotation*toRadians)*scale*photonScale),
	}
	two := image.Point{
		X: int(p.x + math.Sin((p.rotation+120)*toRadians)*scale*photonScale),
		Y: int(p.y - math.Cos((p.rotation+120)*toRadians)*scale*photonScale),
	}
	three := image.Point{
		X: int(p.x + math.Sin((p.rotation-120)*toRadians)*scale*photonScale),
		Y: int(p.y - math.Cos((p.rotation-120)*toRadians)*scale*photonScale),
	}
	photonLineWidth := float32(2)
	vector.StrokeLine(screen, float32(one.X), float32(one.Y), float32(two.X), float32(two.Y), photonLineWidth, color.White, true)
	vector.StrokeLine(screen, float32(two.X), float32(two.Y), float32(three.X), float32(three.Y), photonLineWidth, color.White, true)
	vector.StrokeLine(screen, float32(three.X), float32(three.Y), float32(one.X), float32(one.Y), photonLineWidth, color.White, true)
}

type PhotonDone func()

func Photon(x, y, vx, vy float64, done PhotonDone) Compo {
	rotationChange := RandomNegate(rand.Float64()*5 + 5)
	return &photon{
		x:              x,
		y:              y,
		vx:             vx,
		vy:             vy,
		distance:       100,
		rotationChange: rotationChange,
		done:           done,
	}
}
