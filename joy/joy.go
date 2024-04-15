package joy

import (
	"fmt"
	"math"
	"rainbow/starter"
	"rainbow/wiz"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

type dim struct {
	s       sync.Mutex
	dimming uint
	d       time.Duration
	done    chan bool
	message chan<- []wiz.Option
}

func (d *dim) incrementDimming() uint {
	d.s.Lock()
	defer d.s.Unlock()
	d.dimming = uint(math.Min(float64(d.dimming)+5, 100))
	return d.dimming
}

func (d *dim) reduceDimming() uint {
	d.s.Lock()
	defer d.s.Unlock()
	d.dimming = uint(math.Max(float64(d.dimming)-5, 5))
	return d.dimming
}

func (d *dim) up() {
	d.changeState(d.incrementDimming)
}
func (d *dim) down() {
	d.changeState(d.reduceDimming)
}

func (d *dim) cancel() {
	d.s.Lock()
	defer d.s.Unlock()
	if d.done != nil {
		close(d.done)
		d.done = nil
	}
}

type colorManager struct {
	last      *wiz.Color
	lastPulse *wiz.Color
	sync.RWMutex
}

func (cm *colorManager) resetLast(data interface{}) {
	cm.Lock()
	defer cm.Unlock()
	cm.last = nil
}

func maxUint8(a, b uint8) uint8 {
	return uint8(math.Max(float64(a), float64(b)))
}

func (cm *colorManager) getColor(c wiz.Color) wiz.Color {
	cm.Lock()
	defer cm.Unlock()
	if cm.last == nil {
		cm.last = &c
		return c
	}
	nc := wiz.Color{
		Blue:  maxUint8(cm.last.Blue, c.Blue),
		Green: maxUint8(cm.last.Green, c.Green),
		Red:   maxUint8(cm.last.Red, c.Red),
	}
	cm.last = &nc
	return nc
}

func (d *dim) changeState(fn func() uint) {
	once := sync.OnceFunc(func() {
		d.s.Lock()
		d.done = make(chan bool)
		d.s.Unlock()
	})
	ticker := time.NewTicker(d.d)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			once()
			d.message <- []wiz.Option{wiz.WithDimming(fn())}
		case <-d.done:
			return
		}
	}
}

func Off() []wiz.Option {
	return []wiz.Option{
		wiz.WithState(false),
	}
}

func New(messages chan<- []wiz.Option) starter.Starter {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3")
	dimManager := &dim{
		d:       time.Millisecond * 50,
		message: messages,
	}
	colorManager := &colorManager{}
	work := func() {
		// setters

		stick.On(joystick.L1Press, func(data interface{}) {
			messages <- []wiz.Option{
				wiz.WithColor(
					colorManager.getColor(wiz.Color{
						Red: 255,
					}),
				),
			}

			fmt.Println("L1Press", data)
		})

		stick.On(joystick.L1Release, colorManager.resetLast)

		stick.On(joystick.TrianglePress, func(data interface{}) {
			messages <- []wiz.Option{
				wiz.WithColor(
					colorManager.getColor(wiz.Color{
						Green: 255,
					}),
				),
			}

			fmt.Println("triangle_press")
		})
		stick.On(joystick.TriangleRelease, colorManager.resetLast)

		stick.On(joystick.CirclePress, func(data interface{}) {
			messages <- []wiz.Option{
				wiz.WithColor(
					colorManager.getColor(wiz.Color{
						Blue: 255,
					}),
				),
			}
			fmt.Println("circle_press")
		})

		stick.On(joystick.CircleRelease, colorManager.resetLast)

		// pulse

		// buttons
		stick.On(joystick.SquarePress, func(data interface{}) {
			messages <- Off()

			fmt.Println("square_press")

		})
		stick.On(joystick.SquareRelease, func(data interface{}) {
			fmt.Println("square_release")
		})
		stick.On(joystick.XPress, func(data interface{}) {
			messages <- Off()
			fmt.Println("x_press")
		})
		stick.On(joystick.XRelease, func(data interface{}) {
			fmt.Println("x_release")
		})
		stick.On(joystick.StartPress, func(data interface{}) {
			fmt.Println("start_press")
		})
		stick.On(joystick.StartRelease, func(data interface{}) {
			fmt.Println("start_release")
		})
		stick.On(joystick.SelectPress, func(data interface{}) {
			fmt.Println("select_press")
		})
		stick.On(joystick.SelectRelease, func(data interface{}) {
			fmt.Println("select_release")
		})

		// joysticks
		stick.On(joystick.LeftX, func(data interface{}) {
			fmt.Println("left_x", data)
		})
		stick.On(joystick.LeftY, func(data interface{}) {
			fmt.Println("left_y", data)
			if value, ok := data.(int16); ok {
				fmt.Println(value)
				if value == 0 {
					go dimManager.cancel()
				} else if value > 0 {
					go dimManager.down()
				} else if value < 0 {
					go dimManager.up()
				}

			}
		})
		stick.On(joystick.RightX, func(data interface{}) {
			fmt.Println("right_x", data)
		})
		stick.On(joystick.RightY, func(data interface{}) {
			fmt.Println("right_y", data)
		})

		// triggers
		stick.On(joystick.R1Press, func(data interface{}) {
			messages <- []wiz.Option{
				wiz.WithSceneID(30),
			}
			fmt.Println("R1Press", data)
		})
		stick.On(joystick.R1Release, func(data interface{}) {
			fmt.Println("R1Release", data)
		})
		stick.On(joystick.R2Press, func(data interface{}) {
			messages <- Off()

			fmt.Println("R2Press", data)
		})
		stick.On(joystick.R2Release, func(data interface{}) {
			fmt.Println("R2Release", data)
		})
		stick.On(joystick.L2Press, func(data interface{}) {
			messages <- Off()
			fmt.Println("L2Press", data)
		})
		stick.On(joystick.L2Release, func(data interface{}) {
			fmt.Println("L2Release", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick},
		work,
	)
	return robot
}
