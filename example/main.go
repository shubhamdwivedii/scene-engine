package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	cam "github.com/shubhamdwivedii/scene-engine/camera"
	gop "github.com/shubhamdwivedii/scene-engine/gopher"
	ovr "github.com/shubhamdwivedii/scene-engine/overlay"
	scr "github.com/shubhamdwivedii/scene-engine/screen"
	vpt "github.com/shubhamdwivedii/scene-engine/viewport"
)

type Game struct{}

const (
	WORLD_W, WORLD_H = 360, 280
	VIEW_W, VIEW_H   = 320, 240
)

var gameScreen scr.Screen
var overlayScreen ovr.Overlay
var viewport *vpt.Viewport
var camera *cam.Camera
var gopher *gop.Gopher
var crateBox *ebiten.Image

func init() {
	var err error
	crateBox, _, err = ebitenutil.NewImageFromFile("./assets/cratebox.png")
	if err != nil {
		log.Fatal(err)
	}
	gopher = gop.New(WORLD_W/2, WORLD_H/2, 7)
	viewport = vpt.New(VIEW_W, VIEW_H, WORLD_W, WORLD_H)
	camera = cam.New(WORLD_W, WORLD_H, 120, 120, WORLD_W/3, WORLD_H/2)
	camera.FocusOn(gopher)
	gameScreen = scr.New(VIEW_W, VIEW_H, WORLD_W, WORLD_H, viewport, camera)
	overlayScreen = ovr.New(VIEW_W, VIEW_H)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		gameScreen.Shake()
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		viewport.MoveBy(-1, 0)
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
	gameScreen.Update()
	return nil
}

func (g *Game) Draw(renderScreen *ebiten.Image) {
	// Draw to game screen first
	gameScreen.Fill(color.RGBA{202, 244, 244, 0xff})
	gopher.Draw(gameScreen)

	drawPlatforms(gameScreen)
	gameScreen.Render(renderScreen)

	// Render Overlay Over the GameScreen
	crateBoxOP := &ebiten.DrawImageOptions{}
	// Transparency Doesn't work without this
	crateBoxOP.CompositeMode = ebiten.CompositeModeCopy
	crateBoxOP.ColorM.Scale(1, 1, 1, 0.25)
	overlayScreen.DrawImage(crateBox, crateBoxOP)
	overlayScreen.Render(renderScreen)
}

func drawPlatforms(screen scr.Screen) {
	pw := 40.0
	ph := 20.0
	gap := 50
	py := 160.0
	for i := 0; i < 20; i++ {
		screen.DrawRect(float64(i*gap), py, pw, ph, true, color.RGBA{255, 0, 0, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 1024, 768 // To Test Resolution Independent Scaling
	return VIEW_W, VIEW_H // Ideally Return Internal Resolution Here.
}

func main() {
	ebiten.SetWindowSize(640, 480)
	gameScreen.SetShakeIntensity(3.5)
	gameScreen.SetDebug(true)
	// viewport.AllowOutOfBounds = true
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
