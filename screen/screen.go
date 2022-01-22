package screen

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/peterhellberg/gfx"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Screen struct {
	Image        *ebiten.Image
	Padding      int
	MaxIntensity float64
	Intensity    float64
	Duration     float64
	DrawOP       *ebiten.DrawImageOptions
}

func New(width, height, padding int) *Screen {
	screenImg := ebiten.NewImage(width+(padding*2), height+(padding*2))
	return &Screen{
		Image:        screenImg,
		Padding:      padding,
		MaxIntensity: 10.0,
		Intensity:    0.0,
		Duration:     1.0,
		DrawOP:       &ebiten.DrawImageOptions{},
	}
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

func (s *Screen) AdjustPadding(values ...float64) []float64 {
	for value := range values {
		value += s.Padding
	}
	return values
}

func (s *Screen) Draw(screen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()
	s.DrawOP.GeoM.Translate(-float64(s.Padding), -float64(s.Padding))
	if s.Intensity < 1 {
		lerp := gfx.Lerp(s.Duration, 0, s.Intensity)
		fmt.Println("Lerp", lerp)
		amplitude := s.MaxIntensity * lerp
		dx := amplitude * (2*rand.Float64() - 1)
		dy := amplitude * (2*rand.Float64() - 1)
		s.DrawOP.GeoM.Translate(-dx, -dy)
	}
	screen.DrawImage(s.Image, s.DrawOP)
}
