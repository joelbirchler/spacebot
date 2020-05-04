package main

import (
	"fmt"
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
		func() { gobot.Every(100*time.Millisecond, tick) },
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

	fmt.Println(x, y, z, altitude, pressure, temp)
}

func altPressTemp() (float32, float32, float32, error) {
	var err error

	apt := struct {
		a float32
		p float32
		t float32
	}{}

	for _, read := range []struct {
		val *float32
		fn  func() (float32, error)
	}{
		{val: &apt.a, fn: bmp280.Altitude},
		{val: &apt.p, fn: bmp280.Pressure},
		{val: &apt.t, fn: bmp280.Temperature},
	} {
		*read.val, err = read.fn()
		if err != nil {
			return 0.0, 0.0, 0.0, err
		}
	}

	return apt.a, apt.p, apt.t, nil
}
