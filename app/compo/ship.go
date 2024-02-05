package compo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type ship struct {
	initialized bool
	x           float64
	y           float64
	direction   float64
	vx          float64
	vy          float64
	photons     int
	coolDown    int
	destroyed   bool
}

func (s *ship) Type() Type {
	return TypeShip
}

func (s *ship) Bounds() (x, y, radius float64) {
	return s.x, s.y, scale
}

type buffer struct {
	ship *ship
}

func (s *buffer) Bounds() (x, y, radius float64) {
	x, y, radius = s.ship.Bounds()
	return x, y, radius * 20
}

func (s *ship) Buffer() BoundsProvider {
	return &buffer{s}
}

func (s *ship) Destroy() {
	s.destroyed = true
}

func (s *ship) Update(screenWidth, screenHeight int, resChanged bool) (result []Compo) {

	if s.destroyed {
		return []Compo{Explosion(s.x, s.y, s.vx, s.vy)}
	}

	if !s.initialized {
		s.initialized = true
		s.x = float64(screenWidth / 2)
		s.y = float64(screenHeight / 2)
	}

	// handle left/right
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.direction -= 5
		if s.direction < 0 {
			s.direction += 360
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.direction += 5
		if s.direction > 360 {
			s.direction -= 360
		}
	}
	// handle thrust
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newVx := s.vx + math.Sin(s.direction*(math.Pi/180))*0.2
		newVy := s.vy - math.Cos(s.direction*(math.Pi/180))*0.2
		// limit max speed
		if math.Sqrt(math.Pow(newVx, 2)+math.Pow(newVy, 2)) < 10 {
			s.vx = newVx
			s.vy = newVy
		}
	}

	if s.coolDown > 0 {
		s.coolDown--
	}

	// add the ship
	result = append(result, s)

	// handle fire photon
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if s.coolDown == 0 && s.photons < 5 {
			s.coolDown = 10
			s.photons++
			p := Photon(
				s.x+math.Sin(s.direction*toRadians)*(scale*1.5),
				s.y-math.Cos(s.direction*toRadians)*(scale*1.5),
				math.Sin(s.direction*toRadians)*4+s.vx,
				-(math.Cos(s.direction*toRadians)*4 - s.vy),
				func() { s.photons-- },
			)
			result = append(result, p)
		}
	}

	// move the ship
	s.x, s.y = WrapScreen(s.x+s.vx, s.y+s.vy, screenWidth, screenHeight)

	return result
}

const toRadians = math.Pi / 180
const scale float64 = 20

func DrawShip(screen *ebiten.Image, x, y, direction, shipScale float64, width float32) (image.Point, image.Point) {

	wingAngle := float64(90 + 55)
	engineAngle := float64(90 + 44)
	engineDistance := shipScale * 0.60

	nose := image.Point{
		X: int(x + math.Sin(direction*toRadians)*shipScale),
		Y: int(y - math.Cos(direction*toRadians)*shipScale),
	}
	right := image.Point{
		X: int(x + math.Sin((direction+wingAngle)*toRadians)*shipScale),
		Y: int(y - math.Cos((direction+wingAngle)*toRadians)*shipScale),
	}
	rightInner := image.Point{
		X: int(x + math.Sin((direction+engineAngle)*toRadians)*engineDistance),
		Y: int(y - math.Cos((direction+engineAngle)*toRadians)*engineDistance),
	}
	left := image.Point{
		X: int(x + math.Sin((direction-wingAngle)*toRadians)*shipScale),
		Y: int(y - math.Cos((direction-wingAngle)*toRadians)*shipScale),
	}
	leftInner := image.Point{
		X: int(x + math.Sin((direction-engineAngle)*toRadians)*engineDistance),
		Y: int(y - math.Cos((direction-engineAngle)*toRadians)*engineDistance),
	}

	vector.StrokeLine(screen, float32(nose.X), float32(nose.Y), float32(right.X), float32(right.Y), width, color.White, true)
	vector.StrokeLine(screen, float32(nose.X), float32(nose.Y), float32(left.X), float32(left.Y), width, color.White, true)
	vector.StrokeLine(screen, float32(leftInner.X), float32(leftInner.Y), float32(rightInner.X), float32(rightInner.Y), width, color.White, true)
	return leftInner, rightInner

}

func (s *ship) Draw(screen *ebiten.Image) {

	width := float32(2)
	leftInner, rightInner := DrawShip(screen, s.x, s.y, s.direction, scale, width)

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		thrustAngle := float64(180) + float64(rand.Intn(10)-5)*2
		thrust := image.Point{
			X: int(s.x + math.Sin((s.direction+thrustAngle)*toRadians)*scale),
			Y: int(s.y - math.Cos((s.direction+thrustAngle)*toRadians)*scale),
		}
		thrustColor := color.RGBA{R: 255, G: 200, B: 100, A: 255}
		vector.StrokeLine(screen, float32(leftInner.X), float32(leftInner.Y), float32(thrust.X), float32(thrust.Y), width, thrustColor, true)
		vector.StrokeLine(screen, float32(thrust.X), float32(thrust.Y), float32(rightInner.X), float32(rightInner.Y), width, thrustColor, true)
	}

}

func Ship() Compo {
	return &ship{}
}
