package raspi

import (
	"time"
	"github.com/stianeikeland/go-rpio"
)

const (

	triggerPin = 20
	echoPin = 21

	soundSpeed = 331.3 + 0.606*21
	pulseDelay = 30000 * time.Nanosecond
)

// HCSR04 is a ultrasonic distance measurer
type HCSR04 struct {
	echo    rpio.Pin
	trigger rpio.Pin
	out     chan float64
}

// NewHCSR04 creates a new HCSR04 sensor
func NewHCSR04() *HCSR04 {

	echo := rpio.Pin(echoPin )
	echo.Input()

	trigger := rpio.Pin( triggerPin )
	trigger.Output()
	trigger.Low()

	return &HCSR04{echo, trigger, make( chan  float64)}
}

func ( sensor HCSR04) Run() chan float64{
	go sensor.loop()
	return sensor.out
}

func ( sensor HCSR04) loop(){

	for{
		distance:= sensor.Measure()
		sensor.out <- distance
		time.Sleep( time.Millisecond * 500 )
	}
}

// Measure returns distance in cm
func (sensor HCSR04) Measure() (value float64) {
	timeoutTime := time.Now()
	sensor.trigger.High()
	time.Sleep(pulseDelay)
	sensor.trigger.Low()

	for {
		if sensor.echo.Read() == rpio.High {
			break
		} else if time.Since(timeoutTime) >= 2*time.Second {
			return -1
		}
	}

	startTime := time.Now()

	for {
		if sensor.echo.Read() == rpio.Low {
			break
		} else if time.Since(timeoutTime) >= 2*time.Second {
			return -1
		}
	}

	duration := time.Since(startTime)
	return float64(duration.Nanoseconds()) / 10000000 * (soundSpeed / 2)
}
