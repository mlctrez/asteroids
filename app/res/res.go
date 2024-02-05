package res

import (
	"embed"
	_ "embed"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	_ "image/png"
	"log"
	"sync"
)

//go:embed ttf/*
var resources embed.FS

func Font(name string) *opentype.Font {
	open, err := resources.ReadFile("ttf/" + name)
	if err != nil {
		log.Fatal(err)
	}
	parse, err := opentype.Parse(open)
	if err != nil {
		log.Fatal(err)
	}
	return parse
}

var fontCache = make(map[string]font.Face)
var fontCacheMutex = &sync.Mutex{}

func CachedFontFace(name string, size float64) font.Face {
	fontCacheMutex.Lock()
	defer fontCacheMutex.Unlock()
	key := fmt.Sprintf("%s/%0.2f", name, size)
	if face, ok := fontCache[key]; ok {
		return face
	}
	opts := &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingVertical}
	face, err := opentype.NewFace(Font(name), opts)
	if err != nil {
		log.Fatal(err)
	}
	fontCache[key] = face
	return face
}
