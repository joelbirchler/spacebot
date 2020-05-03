package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	//log.Println("Hello!")

	pi := raspi.NewAdaptor()
	adxl := i2c.NewADXL345Driver(pi)
	bmp280 := i2c.NewBMP280Driver(pi, i2c.WithAddress(0x76))
	bmp280.Start()

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			x, y, z, _ := adxl.XYZ()
			a, _ := bmp280.Altitude()
			p, err := bmp280.Pressure()
			t, _ := bmp280.Temperature()

			if err != nil {
				log.Fatal(err)
			}
			// some functions would make this more sensable
			//log.Printf("x: %.7f | y: %.7f | z: %.7f \n", x, y, z)

			fmt.Print("\033[G\033[K")
			fmt.Printf("x: %.4f | y: %.4f | z: %.4f |+| a: %.4f | p: %.4f | t: %.4f", x, y, z, a, p, t)
		})
	}

	robot := gobot.NewRobot(
		"spacebot",
		[]gobot.Connection{pi},
		[]gobot.Device{adxl},
		work)

	robot.Start()
}
