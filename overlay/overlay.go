package overlay

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Overlay is like a Static Screen. (No Shake, No Move)
type Overlay interface {
	DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions)
	Render(screen *ebiten.Image)
}

// Can be used for Overlay, Effects or Transitions
type StaticScreen struct {
	Width       int
	Height      int
	Image       *ebiten.Image
	DrawOP      *ebiten.DrawImageOptions
	Debug       bool
	AutoScaling bool
}

func New(width, height int) Overlay {
	screenImg := ebiten.NewImage(width, height)

	// To Test
	screenImg.Fill(color.RGBA{64, 220, 14, 64})

	return &StaticScreen{
		Image:       screenImg,
		Width:       width,
		Height:      height,
		DrawOP:      &ebiten.DrawImageOptions{},
		AutoScaling: true,
	}
}

func (s *StaticScreen) Render(screen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()

	// Scaling Screen Image to Render Resolution
	if !s.AutoScaling {
		resX, resY := screen.Bounds().Dx(), screen.Bounds().Dy()
		if resX != s.Width || resY != s.Height {
			scaleX, scaleY := float64(resX)/float64(s.Width), float64(resY)/float64(s.Height)
			s.DrawOP.GeoM.Scale(scaleX, scaleY)
		}
	}

	screen.DrawImage(s.Image, s.DrawOP)
}

func (s *StaticScreen) Fill(col color.Color) {
	s.Image.Fill(col)
}

func (s *StaticScreen) DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions) {
	s.Image.DrawImage(image, op)
}

func (s *StaticScreen) DrawLine(x1, y1, x2, y2 float64, col color.Color) {
	ebitenutil.DrawLine(s.Image, x1, y1, x2, y2, col)
}

func (s *StaticScreen) DrawRect(x, y, width, height float64, solid bool, clr color.Color) {
	if solid {
		ebitenutil.DrawRect(s.Image, x, y, width, height, clr)
	} else {
		x2 := x + width
		y2 := y + height
		s.DrawLine(x, y, x2, y, clr)
		s.DrawLine(x+1, y, x+1, y2, clr)
		s.DrawLine(x2, y, x2, y2, clr)
		s.DrawLine(x, y2-1, x2, y2-1, clr)
	}
}

func (s *StaticScreen) DebugPrint(text string) {
	ebitenutil.DebugPrint(s.Image, text)
}

func (s *StaticScreen) DebugPrintAt(text string, x, y int) {
	ebitenutil.DebugPrintAt(s.Image, text, x, y)
}
