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
	done    *chan bool
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
	close(*d.done)
	dc := make(chan bool)
	d.done = &dc
}

func (d *dim) changeState(fn func() uint) {
	ticker := time.NewTicker(d.d)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("aca")
			d.message <- []wiz.Option{wiz.WithDimming(fn())}
		case <-*d.done:
			fmt.Println("por aca")
			return
		}
	}
}

func New(messages chan<- []wiz.Option) starter.Starter {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3")
	d := make(chan bool)
	dimManager := &dim{
		d:       time.Millisecond * 50,
		done:    &d,
		message: messages,
	}
	work := func() {
		// buttons
		stick.On(joystick.SquarePress, func(data interface{}) {
			fmt.Println("square_press")
			messages <- []wiz.Option{
				wiz.WithColor(wiz.Color{
					Red: 255,
				}),
				wiz.WithDimming(10),
			}
		})
		stick.On(joystick.SquareRelease, func(data interface{}) {
			fmt.Println("square_release")
		})
		stick.On(joystick.TrianglePress, func(data interface{}) {
			messages <- []wiz.Option{
				wiz.WithColor(wiz.Color{
					Green: 255,
				}),
				wiz.WithDimming(100),
			}
			fmt.Println("triangle_press")
		})
		stick.On(joystick.TriangleRelease, func(data interface{}) {
			fmt.Println("triangle_release")
		})
		stick.On(joystick.CirclePress, func(data interface{}) {
			messages <- []wiz.Option{wiz.WithColor(wiz.Color{
				Blue: 255,
			})}
			fmt.Println("circle_press")
		})
		stick.On(joystick.CircleRelease, func(data interface{}) {
			fmt.Println("circle_release")
		})
		stick.On(joystick.XPress, func(data interface{}) {
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
			fmt.Println("R1Press", data)
		})
		stick.On(joystick.R1Release, func(data interface{}) {
			fmt.Println("R1Release", data)
		})
		stick.On(joystick.R2Press, func(data interface{}) {
			fmt.Println("R2Press", data)
		})
		stick.On(joystick.R2Release, func(data interface{}) {
			fmt.Println("R2Release", data)
		})
		stick.On(joystick.L1Press, func(data interface{}) {
			fmt.Println("L1Press", data)
		})
		stick.On(joystick.L1Release, func(data interface{}) {
			fmt.Println("L1Release", data)
		})
		stick.On(joystick.L2Press, func(data interface{}) {
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
