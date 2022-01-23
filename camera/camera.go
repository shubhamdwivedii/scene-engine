package camera

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
)

type FocusableEntity interface {
	GetPosition() (posX, posY float64)
}

type Camera struct {
	Position      f64.Vec2
	WorldView     f64.Vec2
	FocusView     f64.Vec2
	FocusCenter   f64.Vec2
	Debug         bool
	AutoFocus     bool
	FocusedEntity FocusableEntity
}

// worldWidth, worldHeight is the width/height of the World (including out-of-screen area)
// focusWidth, focusHeight is the width/height of focus area (used to check if x,y is out of focus)
// focusX, focusY is point of focus within the the WorldView (center of FocusView)
func New(worldWidth, worldHeight, focusWidth, focusHeight int, focusX, focusY float64) *Camera {
	return &Camera{
		Position:    f64.Vec2{float64(worldWidth) / 2, float64(worldHeight) / 2},
		WorldView:   f64.Vec2{float64(worldWidth), float64(worldHeight)},
		FocusView:   f64.Vec2{float64(focusWidth), float64(focusHeight)},
		FocusCenter: f64.Vec2{focusX, focusY},
	}
	// ORIGIN is (0,0), FocusedEntity is nil
}

func (c *Camera) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		c.MoveBy(-4, 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyL) {
		c.MoveBy(4, 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyI) {
		c.MoveBy(0, -4)
	}

	if ebiten.IsKeyPressed(ebiten.KeyK) {
		c.MoveBy(0, 4)
	}

	if c.AutoFocus && c.FocusedEntity != nil {
		xPos, yPos := c.FocusedEntity.GetPosition()
		c.Refocus(xPos, yPos)
	}

	return nil
}

// Only used for Debugging // Only draw to main render screen (no offsets)
func (c *Camera) Draw(screen *ebiten.Image) {
	if c.Debug {
		offx, offy := c.GetOffset()
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("Camera-X: %0.2f Camera-Y: %0.2f", c.Position[0], c.Position[1]),
			0, 32,
		)
		x1 := offx + c.FocusCenter[0] - c.FocusView[0]/2
		x2 := x1 + c.FocusView[0]
		y1 := offy + c.FocusCenter[1] - c.FocusView[1]/2
		y2 := y1 + c.FocusView[1]

		ebitenutil.DrawLine(screen, x1, y1, x2, y1, color.RGBA{0, 0, 255, 255})
		ebitenutil.DrawLine(screen, x1, y1, x1, y2, color.RGBA{0, 0, 255, 255})
		ebitenutil.DrawLine(screen, x2, y1, x2, y2, color.RGBA{0, 0, 255, 255})
		ebitenutil.DrawLine(screen, x1, y2, x2, y2, color.RGBA{0, 0, 255, 255})
	}
}

func (c *Camera) EnableDebug() {
	c.Debug = true
}

func (c *Camera) FocusOn(entity FocusableEntity) {
	c.FocusedEntity = entity
	c.AutoFocus = true
}

func (c *Camera) DisableAutoFocus() {
	c.AutoFocus = false
}

func (c *Camera) Refocus(x, y float64) {
	dx, dy := 0.0, 0.0
	lx := c.FocusCenter[0] - c.FocusView[0]/2
	if x < lx {
		dx = x - lx
	}
	rx := c.FocusCenter[0] + c.FocusView[0]/2
	if x > rx {
		dx = x - rx
	}

	ty := c.FocusCenter[1] - c.FocusView[1]/2
	if y < ty {
		dy = y - ty
	}

	by := c.FocusCenter[1] + c.FocusView[1]/2
	if y > by {
		dy = y - by
	}

	c.MoveBy(dx, dy)
}

// (w/2,h/2) is default position
func (c *Camera) MoveTo(x, y float64) {
	dx, dy := x-c.Position[0], y-c.Position[1]
	c.MoveBy(dx, dy)
}

func (c *Camera) MoveBy(dx, dy float64) {
	c.Position[0] += dx
	c.Position[1] += dy
	c.FocusCenter[0] += dx
	c.FocusCenter[1] += dy
}

// offset is (0,0) when camera position is (ww/2, wh/2)
func (c *Camera) GetOffset() (float64, float64) {
	// Right +ve, Left -ve, Up -ve, Down +ve
	dx, dy := c.Position[0]-c.WorldView[0]/2, c.Position[1]-c.WorldView[1]/2

	// if camera moves towards right ??? Why this works ?
	return -dx, -dy
}

// Concat this matrix for adjust for camera position
func (c *Camera) GetOffsetMatrix() ebiten.GeoM {
	matrix := ebiten.GeoM{}
	dx, dy := c.GetOffset()
	matrix.Translate(dx, dy)
	return matrix
}
