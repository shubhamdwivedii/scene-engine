package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	cam "github.com/shubhamdwivedii/scene-engine/camera"
	gop "github.com/shubhamdwivedii/scene-engine/gopher"
	scr "github.com/shubhamdwivedii/scene-engine/screen"
	vpt "github.com/shubhamdwivedii/scene-engine/viewport"
)

type Game struct{}

const (
	WORLD_W, WORLD_H = 360, 280
	VIEW_W, VIEW_H   = 320, 240
)

var gamescreen *scr.Screen
var viewport *vpt.Viewport
var camera *cam.Camera
var gopher *gop.Gopher

func init() {
	gopher = gop.New(WORLD_W/2, WORLD_H/2, 7)
	viewport = vpt.New(VIEW_W, VIEW_H, WORLD_W, WORLD_H)
	camera = cam.New(WORLD_W, WORLD_H, 120, 120, WORLD_W/3, WORLD_H/2)
	camera.FocusOn(gopher)
	gamescreen = scr.New(VIEW_W, VIEW_H, WORLD_W, WORLD_H, viewport, camera)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		gamescreen.Shake()
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		gamescreen.Viewport.MoveBy(-1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		viewport.MoveBy(1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		viewport.MoveBy(0, -1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		viewport.MoveBy(0, 1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		viewport.ZoomBy(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		viewport.ZoomBy(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		viewport.RoatateBy(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		viewport.Reset()
	}

	gopher.Update()
	// Update Camera After FocusEntity has been updated. (Or else you'll see jitter)
	camera.Update()
	gamescreen.Update()
	return nil
}

func (g *Game) Draw(renderScreen *ebiten.Image) {
	// Draw to game screen first
	gamescreen.Fill(color.RGBA{202, 244, 244, 0xff})
	camMatrix := camera.GetOffsetMatrix()
	gopher.Draw(gamescreen.Image, camMatrix)
	drawPlatforms(gamescreen.Image)
	camera.Draw(gamescreen.Image)
	gamescreen.Draw(renderScreen)
}

func drawPlatforms(screen *ebiten.Image) {
	pw := 40.0
	ph := 20.0
	gap := 50
	py := 160.0
	offx, offy := camera.GetOffset()
	for i := 0; i < 20; i++ {
		ebitenutil.DrawRect(screen, float64(i*gap)+offx, py+offy, pw, ph, color.RGBA{255, 0, 0, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	gamescreen.SetShakeIntensity(3.5)
	gamescreen.SetDebug(true, true)
	// viewport.AllowOutOfBounds = true
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
