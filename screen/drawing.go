package screen

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Takes coordinates based on Screen and Adjusts automatically for World (Screen x1,y1 are 0,0)
func (s *CustomScreen) DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions) {
	cameraMatrix := s.Camera.GetOffsetMatrix()
	op.GeoM.Concat(cameraMatrix)

	// op.GeoM.Concat(s.OffsetMatrix)
	s.Image.DrawImage(image, op)
}

func (s *CustomScreen) Fill(col color.Color) {
	s.Image.Fill(col)
}

func (s *CustomScreen) DrawLine(x1, y1, x2, y2 float64, col color.Color) {
	offx, offy := s.GetCamera().GetOffset()
	ebitenutil.DrawLine(s.Image, x1+offx, y1+offy, x2+offx, y2+offy, col)
}

func (s *CustomScreen) DrawRect(x, y, width, height float64, solid bool, clr color.Color) {
	offx, offy := s.GetCamera().GetOffset()
	if solid {
		ebitenutil.DrawRect(s.Image, x+offx, y+offy, width, height, clr)
	} else {
		x2 := x + width
		y2 := y + height
		s.DrawLine(x, y, x2, y, clr)
		s.DrawLine(x+1, y, x+1, y2, clr)
		s.DrawLine(x2, y, x2, y2, clr)
		s.DrawLine(x, y2-1, x2, y2-1, clr)
	}
}

func (s *CustomScreen) DebugPrint(text string) {
	s.DebugPrintAt(text, 0, 0)
}

func (s *CustomScreen) DebugPrintAt(text string, x, y int) {
	offx, offy := s.GetCamera().GetOffset()
	ebitenutil.DebugPrintAt(s.Image, text, x+int(offx), y+int(offy))
}
