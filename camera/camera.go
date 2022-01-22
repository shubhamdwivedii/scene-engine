package camera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ScreenWidth  float64
	ScreenHeight float64
	ViewPort     f64.Vec2
	Position     f64.Vec2
	ZoomFactor   int
	Rotation     int
}

func New(screenWidth, screenHeight float64) *Camera {
	return &Camera{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		ViewPort:     f64.Vec2{screenWidth, screenHeight},
	}
	// rest is zero valued
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// Scaling & Rotation is done around center of Screen/Image
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) RenderMatrix() ebiten.GeoM {
	return c.worldMatrix()
}

// Converts Screen Coordinates to World Coordinates
func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// when scaling its possible that matrix is not invertible
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}

// (0,0) is default origin
func (c *Camera) MoveTo(x, y float64) {
	c.Position[0] = x
	c.Position[1] = y
}

func (c *Camera) MoveBy(dx, dy float64) {
	c.Position[0] += dx
	c.Position[1] += dy
}

func (c *Camera) ZoomBy(dz int) {
	if c.ZoomFactor > -2400 && c.ZoomFactor < 2400 {
		c.ZoomFactor += dz
	}
}

// default z = 0
func (c *Camera) SetZoom(z int) {
	c.ZoomFactor = z
}

// default r = 0
func (c *Camera) SetRotation(r int) {
	c.Rotation = r
}

func (c *Camera) RoatateBy(dr int) {
	c.Rotation += dr
}
