package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlctrez/asteroids/app"
	"log"
)

func main() {
	if err := ebiten.RunGame(app.New()); err != nil {
		log.Fatal(err)
	}
}
