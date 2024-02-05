package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mlctrez/asteroids/app/compo"
	"github.com/mlctrez/asteroids/app/res"
	"golang.org/x/image/font"
	"image/color"
	"math/rand"
	"os"
)

func New() *Game {
	if os.Getenv("GOARCH") != "wasm" {
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		//ebiten.SetFullscreen(true)
		ebiten.SetWindowSize(1280, 1024)
	}
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	return &Game{
		smallFont:    res.CachedFontFace("DMMono-Regular.ttf", 24),
		largeFont:    res.CachedFontFace("DMMono-Regular.ttf", 128),
		lives:        0,
		splashScreen: true,
	}
}

type Game struct {
	components       []compo.Compo
	fieldInitialized bool
	smallFont        font.Face
	largeFont        font.Face
	score            uint16
	level            uint16
	lives            uint8
	splashScreen     bool
}

func (g *Game) Update() (err error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.lives = 0
		g.splashScreen = true
		g.score = 0
		g.level = 1
		g.components = []compo.Compo{}
		g.fieldInitialized = false
	}

	if g.lives == 0 && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.lives = 3
		g.score = 0
		g.level = 1
		g.components = []compo.Compo{compo.Ship()}
		g.fieldInitialized = false
		g.splashScreen = false
	}

	var newComponents []compo.Compo
	for _, component := range g.components {
		newComponents = append(newComponents, component.Update(screenWidth, screenHeight, resChanged)...)
	}

	var photons []compo.Compo
	for _, component := range newComponents {
		if component.Type() == compo.TypePhoton {
			photons = append(photons, component)
		}
	}

	var ship compo.Compo
	for _, component := range newComponents {
		if component.Type() == compo.TypeShip {
			ship = component
		}
	}

	var asteroids []compo.Compo
	for _, component := range newComponents {
		if component.Type() == compo.TypeAsteroid {
			asteroids = append(asteroids, component)
		}
	}

	for _, asteroid := range asteroids {
		for _, photon := range photons {
			if compo.Overlaps(photon, asteroid) {
				g.addScore(asteroid.(compo.PointsProvider).Points())
				asteroid.Destroy()
				photon.Destroy()
				break
			}
		}
	}

	for _, asteroid := range asteroids {
		if compo.Overlaps(ship, asteroid) {
			g.addScore(asteroid.(compo.PointsProvider).Points())
			asteroid.Destroy()
			ship.Destroy()
			g.lives--
			if g.lives == 0 {
				fmt.Println("Game Over")
			}
			break
		}
	}

	if g.lives > 0 && ship == nil && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		newComponents = append(newComponents, compo.Ship())
	}

	if len(asteroids) == 0 && g.splashScreen {
		for i := 0; i < 8; i++ {
			rx := rand.Float64() * float64(screenWidth)
			ry := rand.Float64() * float64(screenHeight)
			asteroid := compo.Asteroid(rx, ry, 3)
			newComponents = append(newComponents, asteroid)
		}
	}

	if len(asteroids) == 0 && ship != nil {
		g.level++

		shipBuffer := ship.(compo.BufferProvider)
		for i := 0; i < 8; i++ {
			for {
				rx := rand.Float64() * float64(screenWidth)
				ry := rand.Float64() * float64(screenHeight)
				asteroid := compo.Asteroid(rx, ry, 3)
				if !compo.Overlaps(shipBuffer.Buffer(), asteroid) {
					asteroids = append(asteroids, asteroid)
					break
				}
			}
		}
		newComponents = append(newComponents, asteroids...)
	}

	g.components = newComponents
	return err
}

const extraLifeEvery = 10000

func (g *Game) addScore(points uint16) {
	extraBefore := g.score / extraLifeEvery
	g.score += points
	extraAfter := g.score / extraLifeEvery
	if extraAfter > extraBefore {
		// TODO: play a high-pitched beeping sound each extra life
		g.lives++
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
	for _, component := range g.components {
		component.Draw(screen)
	}

	if !g.splashScreen {
		scoreText := fmt.Sprintf("%02d", g.score)
		for len(scoreText) < 6 {
			scoreText = " " + scoreText
		}

		text.Draw(screen, scoreText, g.smallFont, 20, 40, color.White)
		for i := 0; i < int(g.lives); i++ {
			compo.DrawShip(screen, float64(i*12)+70, 60, 0, 10, 1)
		}
	}

	if g.lives == 0 {
		if g.splashScreen {
			text.Draw(screen, "ASTEROIDS", g.largeFont, screenWidth/2-350, screenHeight/2, color.White)
		} else {
			text.Draw(screen, "GAME OVER", g.largeFont, screenWidth/2-350, screenHeight/2, color.White)
		}
		text.Draw(screen, "PRESS P TO PLAY", g.smallFont, screenWidth/2-110, screenHeight/2+80, color.White)
		text.Draw(screen, "W=THRUST A=LEFT D=RIGHT SPACE=FIRE", g.smallFont, screenWidth/2-250, screenHeight/2+120, color.White)
		text.Draw(screen, "S=SPAWN R=RESET", g.smallFont, screenWidth/2-110, screenHeight/2+150, color.White)
	}
}

var screenWidth, screenHeight int
var resChanged bool

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	resChanged = outsideWidth != screenWidth || outsideHeight != screenHeight
	screenWidth, screenHeight = outsideWidth, outsideHeight
	return screenWidth, screenHeight
}
