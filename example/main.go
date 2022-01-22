package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	cam "github.com/shubhamdwivedii/scene-engine/camera"
	scr "github.com/shubhamdwivedii/scene-engine/screen"
)

type Game struct{}

var gopher *ebiten.Image
var gamescreen *scr.Screen
var camera *cam.Camera

func init() {
	var err error
	gopher, _, err = ebitenutil.NewImageFromFile("./assets/gopher-title.png")
	if err != nil {
		log.Fatal(err)
	}
	camera = cam.New(320, 240)
	gamescreen = scr.New(320, 240, 360, 280, camera)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		gamescreen.Shake()
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		gamescreen.Camera.MoveBy(-1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		camera.MoveBy(1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		camera.MoveBy(0, -1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		camera.MoveBy(0, 1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		camera.ZoomBy(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		camera.ZoomBy(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		camera.RoatateBy(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		camera.Reset()
	}

	gamescreen.Update()
	return nil
}

func (g *Game) Draw(renderScreen *ebiten.Image) {
	// Draw to game screen first
	gamescreen.Fill(color.RGBA{202, 244, 244, 0xff})
	OP := &ebiten.DrawImageOptions{}
	OP.GeoM.Translate(160-32, 120-32)
	gamescreen.DrawImage(gopher, OP) // Offset Adjusted Automatically
	OP.GeoM.Reset()
	OP.GeoM.Translate(-32, -32)
	gamescreen.DrawImage(gopher, OP)
	gamescreen.Draw(renderScreen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	gamescreen.SetShakeIntensity(3.5)
	gamescreen.EnableDebug()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
