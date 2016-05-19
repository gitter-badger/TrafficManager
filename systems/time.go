package systems

import (
	"image/color"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

const (
	SpeedOne   = 1
	SpeedTwo   = 2
	SpeedThree = 15

	SpeedOneButton   = "speed1"
	SpeedTwoButton   = "speed2"
	SpeedThreeButton = "speed3"

	clockSize    float64 = 24
	clockPadding float32 = 4
	clockZIndex  float32 = 1000
)

type TimeComponent struct {
	Time  time.Time
	Speed float32
}

type clock struct {
	ecs.BasicEntity
	TimeComponent
	common.RenderComponent
	common.SpaceComponent
}

type TimeSystem struct {
	clock      clock
	clockCache string

	robotoFont common.Font
}

func (*TimeSystem) Remove(ecs.BasicEntity) {}

func (t *TimeSystem) New(w *ecs.World) {
	// Set default values
	t.clock.Time = time.Now()
	t.clock.Speed = SpeedOne

	// Register buttons
	engo.Input.RegisterButton(SpeedOneButton, engo.NumOne, engo.One)
	engo.Input.RegisterButton(SpeedTwoButton, engo.NumTwo, engo.Two)
	engo.Input.RegisterButton(SpeedThreeButton, engo.NumThree, engo.Three)

	// Load the preloaded font
	t.robotoFont = common.Font{
		URL:  "fonts/Roboto-Regular.ttf",
		FG:   color.Black,
		Size: clockSize,
	}
	if err := t.robotoFont.CreatePreloaded(); err != nil {
		panic(err)
	}

	// Create graphical representation of the clock
	t.clock.BasicEntity = ecs.NewBasic()
	t.clock.RenderComponent = common.RenderComponent{
		Drawable: t.robotoFont.Render(t.clock.Time.Format("15:04")),
		Color:    color.Black,
	}
	t.clock.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: engo.CanvasWidth() - t.clock.RenderComponent.Drawable.Width() - clockPadding,
			Y: clockPadding,
		},
		Width:  t.clock.RenderComponent.Drawable.Width() + 2*clockPadding,
		Height: t.clock.RenderComponent.Drawable.Height() + 2*clockPadding,
	}
	t.clock.SetZIndex(clockZIndex)
	t.clock.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&t.clock.BasicEntity, &t.clock.RenderComponent, &t.clock.SpaceComponent)
		case *CommuterSystem:
			sys.SetClock(&t.clock.BasicEntity, &t.clock.TimeComponent, &t.clock.SpaceComponent)
		}
	}
}

func (t *TimeSystem) Update(dt float32) {
	// Update the visual clock
	t.clock.Time = t.clock.Time.Add(time.Duration(float32(time.Minute) * dt * t.clock.Speed))
	if timeString := t.clock.Time.Format("15:04"); timeString != t.clockCache {
		t.clock.Drawable.Close()
		t.clock.Drawable = t.robotoFont.Render(timeString)
		t.clockCache = timeString
	}
	t.clock.Position.X = engo.CanvasWidth() - t.clock.Width

	// Watch for speed changes
	if engo.Input.Button(SpeedOneButton).Down() {
		t.clock.Speed = SpeedOne
	} else if engo.Input.Button(SpeedTwoButton).Down() {
		t.clock.Speed = SpeedTwo
	} else if engo.Input.Button(SpeedThreeButton).Down() {
		t.clock.Speed = SpeedThree
	}
}