package screen

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/font"

	cam "github.com/shubhamdwivedii/scene-engine/camera"
	vpt "github.com/shubhamdwivedii/scene-engine/viewport"
)

const (
	AUTO_PADDING = 20
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
	DebugPrintAt(text string, x, y int)
	DrawText(text string, fnt font.Face, x, y float64, clr color.Color)
}

type CustomScreen struct {
	ScreenWidth  int
	ScreenHeight int
	WorldWidth   int
	WorldHeight  int
	// Offset            f64.Vec2
	// OffsetMatrix      ebiten.GeoM
	Image             *ebiten.Image
	Viewport          *vpt.Viewport
	Camera            *cam.Camera
	MaxShakeIntensity float64
	ShakeIntensity    float64
	ShakeDuration     float64
	DrawOP            *ebiten.DrawImageOptions
	Debug             bool
	AutoScaling       bool // Automatically Scales To Target Screen Resolution on Render
	AutoPadding       bool
	StaticViewport    bool
	StaticCamera      bool
}

type ScreenOptions struct {
	AutoScaling   bool
	FixedViewport bool
}

// Width + Padding = World_Width
func New(screenWidth, screenHeight, worldWidth, worldHeight int, viewport *vpt.Viewport, camera *cam.Camera) (Screen, error) {
	autoPadding := false
	if screenWidth == worldWidth && screenHeight == worldHeight {
		if viewport != nil {
			return nil, errors.New("viewport should be nil when screen-size equals world-size")
		}
		autoPadding = true
		worldWidth += AUTO_PADDING * 2
		worldHeight += AUTO_PADDING * 2
	} else {
		if viewport == nil {
			return nil, errors.New("viewport cannot be nil if screen-size and world-size are different")
		}
	}
	screenImg := ebiten.NewImage(worldWidth, worldHeight)

	// Offsets are used to render relative to screenOrigin (instead of worldOrigin)
	// offx, offy := float64(worldWidth-screenWidth)/2, float64(worldHeight-screenHeight)/2

	// offsetMatrix := ebiten.GeoM{}
	// offsetMatrix.Translate(offx, offy)

	return &CustomScreen{
		Image:             screenImg,
		ScreenWidth:       screenWidth,
		ScreenHeight:      screenHeight,
		WorldWidth:        worldWidth,
		WorldHeight:       worldHeight,
		Viewport:          viewport,
		Camera:            camera,
		MaxShakeIntensity: 10.0,
		ShakeIntensity:    1.0,
		ShakeDuration:     1.0,
		DrawOP:            &ebiten.DrawImageOptions{},
		AutoScaling:       true,
		StaticViewport:    viewport == nil,
		StaticCamera:      camera == nil,
		AutoPadding:       autoPadding,
	}, nil
}

func (s *CustomScreen) SetDebug(debugOn bool) {
	s.Debug = debugOn
	if s.Camera != nil {
		s.Camera.Debug = debugOn
	}
}

func (s *CustomScreen) Shake() {
	s.ShakeIntensity = 0.0
}

// 10.0 = Very Intense, 1.0  = Non Existent
func (s *CustomScreen) SetShakeIntensity(maxIntensity float64) {
	s.MaxShakeIntensity = maxIntensity
	s.ShakeDuration = maxIntensity / 10.0
}

func (s *CustomScreen) Update() error {
	s.ShakeIntensity += 1 / 60.0 // 60 FPS fixed.
	return nil
}

// func (s *CustomScreen) AdjustForOffset(x, y float64) (float64, float64) {
// 	return x + s.Offset[0], y + s.Offset[1]
// }

// Draws CustomScreen to RenderScreen
func (s *CustomScreen) Render(screen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()

	if s.ShakeIntensity < 1 {
		lerped := gfx.Lerp(s.ShakeDuration, 0, s.ShakeIntensity)
		amplitude := s.MaxShakeIntensity * lerped
		dx := amplitude * (2*rand.Float64() - 1)
		dy := amplitude * (2*rand.Float64() - 1)
		s.DrawOP.GeoM.Translate(-dx, -dy)
	}

	if s.AutoPadding && s.Viewport == nil {
		// Need To Render CustomScreen slightly off left/top (on RenderScreen) to adjust for AutoPadding
		s.DrawOP.GeoM.Translate(-AUTO_PADDING, -AUTO_PADDING)

	} else {
		transformMatrix := s.Viewport.RenderMatrix()
		s.DrawOP.GeoM.Concat(transformMatrix)
	}

	if s.Debug {
		// Debug stuff to render on game scene screen
		if s.Viewport != nil && !s.AutoPadding {
			s.drawViewportArea()
		} else {
			s.drawFixedViewArea()
		}

		if s.Camera != nil {
			s.drawCameraFocusArea()
		}
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
		// worldX, worldY := s.Viewport.ScreenToWorld(ebiten.CursorPosition())

		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()),
		)

		// ebitenutil.DebugPrintAt(
		// 	screen,
		// 	fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
		// 		s.Viewport.String(),
		// 		worldX, worldY),
		// 	0, 240-32,
		// )
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
	offx, offy := 0.0, 0.0

	if s.Camera != nil {
		// CameraOffset is already added once in s.DrawLine
		offx, offy = s.Camera.GetOffsets()
		// Adding CameraOffsets twice would make this move with camera (will appear static on screen)
	}

	// Current Position of the Viewport RED
	x1, y1 := s.Viewport.Position[0]-offx, s.Viewport.Position[1]-offy
	x2, y2 := float64(s.ScreenWidth)+x1, float64(s.ScreenHeight)+y1
	s.DrawLine(x1, y1, x2, y1, color.RGBA{255, 0, 0, 255})
	s.DrawLine(x1+1, y1, x1+1, y2, color.RGBA{255, 0, 0, 255})
	s.DrawLine(x2, y1, x2, y2, color.RGBA{255, 0, 0, 255})
	s.DrawLine(x1, y2-1, x2, y2-1, color.RGBA{255, 0, 0, 255})

	// Inital Position of the Viewport (for reference) GREEN
	x1, y1 = s.Viewport.InitialPosition[0]-offx, s.Viewport.InitialPosition[1]-offy
	x2, y2 = float64(s.ScreenWidth)+x1, float64(s.ScreenHeight)+y1
	s.DrawLine(x1, y1, x2, y1, color.RGBA{0, 255, 0, 255})
	s.DrawLine(x1+1, y1, x1+1, y2, color.RGBA{0, 255, 0, 255})
	s.DrawLine(x2, y1, x2, y2, color.RGBA{0, 255, 0, 255})
	s.DrawLine(x1, y2-1, x2, y2-1, color.RGBA{0, 255, 0, 255})
}

func (s *CustomScreen) drawFixedViewArea() {
	offx, offy := 0.0, 0.0

	if s.Camera != nil {
		offx, offy = s.Camera.GetOffsets()
		// Adding CameraOffsets twice would make this move with camera (will appear static on screen)
	}

	x1, y1 := -offx, -offy // why -ve works ???
	x2, y2 := float64(s.ScreenWidth)+x1, float64(s.ScreenHeight)+y1
	s.DrawLine(x1, y1, x2, y1, color.RGBA{0, 64, 135, 255})
	s.DrawLine(x1+1, y1, x1+1, y2, color.RGBA{0, 64, 135, 255})
	s.DrawLine(x2, y1, x2, y2, color.RGBA{0, 64, 135, 255})
	s.DrawLine(x1, y2-1, x2, y2-1, color.RGBA{0, 64, 135, 255})

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

// Includes Padding-Offset if Viewport is Nil
// Includes Camera Padding if Camera is not-Nil
func (s *CustomScreen) GetOffsets() (dx, dy float64) {
	if s.Camera != nil {
		dx, dy = s.Camera.GetOffsets()
	}
	if s.AutoPadding && s.Viewport == nil {
		dx += AUTO_PADDING
		dy += AUTO_PADDING
	}
	return
}

func (s *CustomScreen) GetOffsetMatrix() (offsetMatrix ebiten.GeoM) {
	if s.Camera != nil {
		offsetMatrix = s.Camera.GetOffsetMatrix()
	}
	if s.AutoPadding && s.Viewport == nil {
		offsetMatrix.Translate(AUTO_PADDING, AUTO_PADDING)
	}
	return
}
