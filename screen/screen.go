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
	"golang.org/x/image/math/f64"

	cam "github.com/shubhamdwivedii/scene-engine/camera"
	vpt "github.com/shubhamdwivedii/scene-engine/viewport"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Screen interface {
	Shake()
	SetShakeIntensity(intensity float64)
	SetDebug(debugOn bool)
	Update() error
	Render(screen *ebiten.Image)
	GetImage() (screenImage *ebiten.Image)
	GetViewport() (viewport *vpt.Viewport)
	GetCamera() (camera *cam.Camera)

	DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions)
	DrawLine(x1, y1, x2, y2 float64, col color.Color)
	DrawRect(x, y, width, height float64, fill bool, col color.Color)
	Fill(col color.Color)
	DebugPrint(text string)
}

type CustomScreen struct {
	ScreenWidth  int
	ScreenHeight int
	WorldWidth   int
	WorldHeight  int
	Offset       f64.Vec2
	OffsetMatrix ebiten.GeoM
	Image        *ebiten.Image
	Viewport     *vpt.Viewport
	Camera       *cam.Camera
	MaxIntensity float64
	Intensity    float64
	Duration     float64
	DrawOP       *ebiten.DrawImageOptions
	Debug        bool
	AutoScaling  bool // Automatically Scales To Target Screen Resolution on Render
}

// Width + Padding = World_Width
func New(screenWidth, screenHeight, worldWidth, worldHeight int, viewport *vpt.Viewport, camera *cam.Camera) Screen {
	screenImg := ebiten.NewImage(worldWidth, worldHeight)

	// Offsets are used to render relative to screenOrigin (instead of worldOrigin)
	offx, offy := float64(worldWidth-screenWidth)/2, float64(worldHeight-screenHeight)/2

	offsetMatrix := ebiten.GeoM{}
	offsetMatrix.Translate(offx, offy)

	return &CustomScreen{
		Image:        screenImg,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		WorldWidth:   worldWidth,
		WorldHeight:  worldHeight,
		OffsetMatrix: offsetMatrix,
		Offset:       f64.Vec2{offx, offy},
		Viewport:     viewport,
		Camera:       camera,
		MaxIntensity: 10.0,
		Intensity:    1.0,
		Duration:     1.0,
		DrawOP:       &ebiten.DrawImageOptions{},
		AutoScaling:  true,
	}
}

func (s *CustomScreen) SetDebug(debugOn bool) {
	s.Debug = debugOn
	s.Camera.Debug = debugOn
}

func (s *CustomScreen) Shake() {
	s.Intensity = 0.0
}

// 10.0 = Very Intense, 1.0  = Non Existent
func (s *CustomScreen) SetShakeIntensity(maxIntensity float64) {
	s.MaxIntensity = maxIntensity
	s.Duration = maxIntensity / 10.0
}

func (s *CustomScreen) Update() error {
	s.Intensity += 1 / 60.0 // 60 FPS fixed.
	return nil
}

// func (s *CustomScreen) AdjustForOffset(x, y float64) (float64, float64) {
// 	return x + s.Offset[0], y + s.Offset[1]
// }

// Draws CustomScreen to RenderScreen
func (s *CustomScreen) Render(screen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()

	// Adjusting for Offset so screen is centered by ScreenOrigin (instaead of WorldOrigin)
	invertedOffset := s.OffsetMatrix
	invertedOffset.Invert()
	s.DrawOP.GeoM.Concat(invertedOffset)

	if s.Intensity < 1 {
		lerped := gfx.Lerp(s.Duration, 0, s.Intensity)
		amplitude := s.MaxIntensity * lerped
		dx := amplitude * (2*rand.Float64() - 1)
		dy := amplitude * (2*rand.Float64() - 1)
		s.DrawOP.GeoM.Translate(-dx, -dy)
	}

	transformMatrix := s.Viewport.RenderMatrix()
	s.DrawOP.GeoM.Concat(transformMatrix)

	if s.Debug {
		// Debug stuff to render on game scene screen
		s.drawViewportArea()
		s.drawCameraFocusArea()
	}

	// Scaling Screen Image to Render Resolution
	if s.AutoScaling {
		resX, resY := screen.Bounds().Dx(), screen.Bounds().Dy()
		if resX != s.ScreenWidth || resY != s.ScreenHeight {
			scaleX, scaleY := float64(resX)/float64(s.ScreenWidth), float64(resY)/float64(s.ScreenHeight)
			s.DrawOP.GeoM.Scale(scaleX, scaleY)
		}
	}

	// Render Screen Image to Real Render Screen
	screen.DrawImage(s.Image, s.DrawOP)

	if s.Debug {
		// Print debug content on real render screen
		worldX, worldY := s.Viewport.ScreenToWorld(ebiten.CursorPosition())

		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()),
		)

		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
				s.Viewport.String(),
				worldX, worldY),
			0, 240-32,
		)
	}
}

func (s *CustomScreen) GetImage() *ebiten.Image {
	return s.Image
}

func (s *CustomScreen) GetViewport() *vpt.Viewport {
	return s.Viewport
}

func (s *CustomScreen) GetCamera() *cam.Camera {
	return s.Camera
}

func (s *CustomScreen) drawViewportArea() {
	x1, y1 := s.Offset[0], s.Offset[1]
	x2, y2 := float64(s.ScreenWidth)+x1, float64(s.ScreenHeight)+y1
	s.DrawLine(x1, y1, x2, y1, color.RGBA{255, 0, 0, 255})
	s.DrawLine(x1+1, y1, x1+1, y2, color.RGBA{255, 0, 0, 255})
	s.DrawLine(x2, y1, x2, y2, color.RGBA{255, 0, 0, 255})
	s.DrawLine(x1, y2-1, x2, y2-1, color.RGBA{255, 0, 0, 255})
}

func (s *CustomScreen) drawCameraFocusArea() {
	// offx, offy := s.Camera.GetOffset()
	cPosition := s.Camera.Position
	cFocusCenter := s.Camera.FocusCenter
	cFocusView := s.Camera.FocusView
	s.DebugPrintAt(fmt.Sprintf("Camera-X: %0.2f Camera-Y: %0.2f", cPosition[0], cPosition[1]), 0, 32)

	x1 := cFocusCenter[0] - cFocusView[0]/2
	x2 := x1 + cFocusView[0]
	y1 := cFocusCenter[1] - cFocusView[1]/2
	y2 := y1 + cFocusView[1]

	// Camera Offset is adjusted in CustomScreen.DebugPrintAt()

	s.DrawLine(x1, y1, x2, y1, color.RGBA{0, 0, 255, 255})
	s.DrawLine(x1+1, y1, x1+1, y2, color.RGBA{0, 0, 255, 255})
	s.DrawLine(x2, y1, x2, y2, color.RGBA{0, 0, 255, 255})
	s.DrawLine(x1, y2-1, x2, y2-1, color.RGBA{0, 0, 255, 255})
}
