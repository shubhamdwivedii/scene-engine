package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	scr "github.com/shubhamdwivedii/scene-engine/screen"
)

type Game struct{}

var gopher *ebiten.Image
var gamescreen *scr.Screen

func init() {
	var err error
	gopher, _, err = ebitenutil.NewImageFromFile("./assets/gopher-title.png")
	if err != nil {
		log.Fatal(err)
	}
	gamescreen = scr.New(320, 240, 20)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		gamescreen.Shake()
	}
	gamescreen.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw to game screen first
	gamescreen.Image.Fill(color.RGBA{202, 244, 244, 0xff})
	OP := &ebiten.DrawImageOptions{}
	OP.GeoM.Translate(160-32+20, 120-32+20) // -GopherRadius+Padding
	gamescreen.Image.DrawImage(gopher, OP)
	gamescreen.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	gamescreen.SetShakeIntensity(3.5)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
