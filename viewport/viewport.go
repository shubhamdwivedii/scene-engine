package viewport

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

type Viewport struct {
	WorldView        f64.Vec2
	Dimensions       f64.Vec2
	Position         f64.Vec2 // TopLeft (Not Middle)
	WorldCenter      f64.Vec2
	ZoomFactor       int
	Rotation         int
	AllowOutOfBounds bool
}

// Viewport should have same dimenstions as Viewable Screen
func New(screenWidth, screenHeight, worldWidth, worldHeight int) *Viewport {
	return &Viewport{
		Dimensions:  f64.Vec2{float64(screenWidth), float64(screenHeight)},
		WorldView:   f64.Vec2{float64(worldWidth), float64(worldHeight)},
		WorldCenter: f64.Vec2{float64(worldWidth) / 2, float64(worldHeight) / 2},
	}
	// rest is zero valued
}

func (v *Viewport) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		v.Position, v.Rotation, v.ZoomFactor,
	)
}

func (v *Viewport) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		v.Dimensions[0] * 0.5,
		v.Dimensions[1] * 0.5,
	}
}

func (v *Viewport) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-v.Position[0], -v.Position[1])
	// Scaling & Rotation is done around center of Screen/Image
	m.Translate(-v.viewportCenter()[0], -v.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(v.ZoomFactor)),
		math.Pow(1.01, float64(v.ZoomFactor)),
	)
	m.Rotate(float64(v.Rotation) * 2 * math.Pi / 360)
	m.Translate(v.viewportCenter()[0], v.viewportCenter()[1])
	return m
}

// Checks if Viewport is OutOfBounds (Outside WorldView)
// Retuns dx, dy to adjust In-Bound
func (v *Viewport) OutOfBounds() (dx float64, dy float64) {
	x1, y1 := v.Position[0], v.Position[1]
	x2, y2 := x1+v.Dimensions[0], y1+v.Dimensions[1]

	padX := (v.WorldView[0] - v.Dimensions[0]) / 2
	padY := (v.WorldView[1] - v.Dimensions[1]) / 2

	if x1 < 0-padX {
		dx = x1*(-1) - padX
	}

	if x2 > v.Dimensions[0]+padX {
		dx = (v.Dimensions[0] + padX) - x2
	}

	if y1 < 0-padY {
		dy = y1*(-1) - padY
	}

	if y2 > v.Dimensions[1]+padY {
		dy = (v.Dimensions[1] + padY) - y2
	}

	return dx, dy
}

func (v *Viewport) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: v.worldMatrix(),
	})
}

func (v *Viewport) RenderMatrix() ebiten.GeoM {
	return v.worldMatrix()
}

// Converts Screen Coordinates to World Coordinates
func (v *Viewport) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := v.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// when scaling its possible that matrix is not invertible
		return math.NaN(), math.NaN()
	}
}

func (v *Viewport) Reset() {
	v.Position[0] = 0
	v.Position[1] = 0
	v.Rotation = 0
	v.ZoomFactor = 0
}

// (0,0) is default origin
func (v *Viewport) MoveTo(x, y float64) {
	v.Position[0] = x
	v.Position[1] = y

	if !v.AllowOutOfBounds {
		dx, dy := v.OutOfBounds()
		fmt.Println("dx, dy", dx, dy)
		v.Position[0] += dx
		v.Position[1] += dy
	}
}

func (v *Viewport) MoveBy(dx, dy float64) {
	v.Position[0] += dx
	v.Position[1] += dy

	if !v.AllowOutOfBounds {
		odx, ody := v.OutOfBounds()
		fmt.Println("dx, dy", odx, ody)
		v.Position[0] += odx
		v.Position[1] += ody
	}
}

func (v *Viewport) ZoomBy(dz int) {
	if v.ZoomFactor > -2400 && v.ZoomFactor < 2400 {
		v.ZoomFactor += dz
	}
}

// default z = 0
func (v *Viewport) SetZoom(z int) {
	v.ZoomFactor = z
}

// default r = 0
func (v *Viewport) SetRotation(r int) {
	v.Rotation = r
}

func (v *Viewport) RoatateBy(dr int) {
	v.Rotation += dr
}
