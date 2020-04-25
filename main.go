package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)


func main() {
	log.Println("Hello!")

	pi := raspi.NewAdaptor()
	adxl := i2c.NewADXL345Driver(pi)

	work := func() {
		gobot.Every(100* time.Millisecond, func() {
			x, y, z, _ := adxl.XYZ()
			log.Printf("x: %.7f4 | y: %.7f | z: %.7f \n", x, y, z)
		})
	}

	robot := gobot.NewRobot(
		"spacebot",
		[]gobot.Connection{pi},
		[]gobot.Device{adxl},
		work)

	robot.Start()
}