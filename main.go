package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	adxl345Address = 0x53
	bmp280Address  = 0x76
)

var (
	pi      *raspi.Adaptor
	adxl345 *i2c.ADXL345Driver
	bmp280  *i2c.BMP280Driver
)

func init() {
	pi = raspi.NewAdaptor()
	adxl345 = i2c.NewADXL345Driver(pi, i2c.WithAddress(adxl345Address))
	bmp280 = i2c.NewBMP280Driver(pi, i2c.WithAddress(bmp280Address))
}

func main() {
	log.Println("Hello!")

	robot := gobot.NewRobot(
		"spacebot",
		[]gobot.Connection{pi},
		[]gobot.Device{adxl345, bmp280},
		func() { gobot.Every(60*time.Second, tick) },
	)

	robot.Start()
}

func tick() {
	x, y, z, adxlErr := adxl345.XYZ()
	altitude, pressure, temp, bmpErr := altPressTemp()

	if adxlErr != nil {
		log.Println("adxl read error:", adxlErr)
	}

	if bmpErr != nil {
		log.Println("bmp read error:", bmpErr)
	}

	log.Println(x, y, z, altitude, pressure, temp)
}

func altPressTemp() (float32, float32, float32, error) {
	// FIXME: Clean this up...
	a, err := bmp280.Altitude()
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}

	p, err := bmp280.Pressure()
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}

	t, err := bmp280.Temperature()
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}

	return a, p, t, nil
}
