package screen

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/peterhellberg/gfx"

	cam "github.com/shubhamdwivedii/scene-engine/camera"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Screen struct {
	ScreenWidth  int
	ScreenHeight int
	WorldWidth   int
	WorldHeight  int
	OffsetX      float64
	OffsetY      float64
	Image        *ebiten.Image
	Camera       *cam.Camera
	MaxIntensity float64
	Intensity    float64
	Duration     float64
	DrawOP       *ebiten.DrawImageOptions
	Debug        bool
}

// Width + Padding = World_Width
func New(screenWidth, screenHeight, worldWidth, worldHeight int, camera *cam.Camera) *Screen {
	screenImg := ebiten.NewImage(worldWidth, worldHeight)

	// Offsets are used to render relative to screenOrigin (instead of worldOrigin)
	offx, offy := float64(worldWidth-screenWidth)/2, float64(worldHeight-screenHeight)/2

	return &Screen{
		Image:        screenImg,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		WorldWidth:   worldWidth,
		WorldHeight:  worldHeight,
		OffsetX:      offx,
		OffsetY:      offy,
		Camera:       camera,
		MaxIntensity: 10.0,
		Intensity:    1.0,
		Duration:     1.0,
		DrawOP:       &ebiten.DrawImageOptions{},
	}
}

func (s *Screen) EnableDebug() {
	s.Debug = true
}

func (s *Screen) Shake() {
	s.Intensity = 0.0
}

// 10.0 = Very Intense, 1.0  = Non Existent
func (s *Screen) SetShakeIntensity(maxIntensity float64) {
	s.MaxIntensity = maxIntensity
	s.Duration = maxIntensity / 10.0
}

func (s *Screen) Update() error {
	s.Intensity += 1 / 60.0 // 60 FPS fixed.
	return nil
}

func (s *Screen) AdjustForOffset(x, y float64) (float64, float64) {
	return x + s.OffsetX, y + s.OffsetY
}

func (s *Screen) Draw(screen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()

	// Adjusting for Offset so screen is centered by ScreenOrigin (instaead of WorldOrigin)
	s.DrawOP.GeoM.Translate(-s.OffsetX, -s.OffsetY)

	if s.Intensity < 1 {
		lerped := gfx.Lerp(s.Duration, 0, s.Intensity)
		amplitude := s.MaxIntensity * lerped
		dx := amplitude * (2*rand.Float64() - 1)
		dy := amplitude * (2*rand.Float64() - 1)
		s.DrawOP.GeoM.Translate(-dx, -dy)
	}

	transformMatrix := s.Camera.RenderMatrix()
	s.DrawOP.GeoM.Concat(transformMatrix)

	if s.Debug {
		// Debug stuff to render on game scene screen
		x1, y1 := s.OffsetX, s.OffsetY
		x2, y2 := float64(s.ScreenWidth)+x1, float64(s.ScreenHeight)+y1
		ebitenutil.DrawLine(s.Image, x1, y1, x1, y2, color.RGBA{255, 0, 0, 255})
		ebitenutil.DrawLine(s.Image, x1, y1, x2, y1, color.RGBA{255, 0, 0, 255})
		ebitenutil.DrawLine(s.Image, x2, y1, x2, y2, color.RGBA{255, 0, 0, 255})
		ebitenutil.DrawLine(s.Image, x1, y2, x2, y2, color.RGBA{255, 0, 0, 255})
	}

	// Render Screen Image to Real Render Screen
	screen.DrawImage(s.Image, s.DrawOP)

	if s.Debug {
		// Print debug content on real render screen
		worldX, worldY := s.Camera.ScreenToWorld(ebiten.CursorPosition())

		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()),
		)

		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
				s.Camera.String(),
				worldX, worldY),
			0, 240-32,
		)
	}
}
