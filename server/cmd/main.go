package main

import (
	m "raspi-bot/messages"
	"raspi-bot/server/config"
	h "raspi-bot/server/handlers"
	. "raspi-bot/utils"
	"raspi-bot/ws"
)
const (
	motorAPin1 = "16"
	motorAPin2 = "18"
	motorAPinE = "22"
)


func main() {


	var bh *h.BaseHandler

	cfg := config.GetConfig()

	registerHandlers( bh )

	server := ws.NewWsServer( cfg )

	hub :=  ws.NewWsHub(m.NewJsonCodec() )

	hub.RunAsync()

	server.Run()


}

func registerHandlers( bh *h.BaseHandler){

	for _, entry := range wsHandlers{
		InjectFieldValue(  entry.handler, "BaseHandler", bh )
		ws.RegisterWsHandler( entry.opcode, entry.handler)
	}

}


/*

	numbPtr := flag.Int("pwm", 42, "pwm value")
	flag.Parse()
	log.Printf("pwm value: %d", *numbPtr )


	r := raspi.NewAdaptor()

	cleanUp := func(){
		r.DigitalWrite(motorAPinE,0 )
		r.DigitalWrite(motorAPin1,0 )
		r.DigitalWrite(motorAPin2,0 )
	}

	cleanUp()


	pwm := byte( *numbPtr )


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)


	pwmPin, err  := r.PWMPin( motorAPinE )

	if err != nil {
		log.Fatalf("error open PWM pin : \n%s", err )
	}


	if err := pwmPin.Enable( true ); err != nil {
		log.Fatalf("error open PWM pin : \n%s", err )
	}


	if err := r.DigitalWrite(motorAPin1,1 );err != nil {
		log.Fatalf("error write to pin 16: \n%s", err )
	}

	if err := r.DigitalWrite(motorAPin2,0 );err != nil {
		log.Fatalf("error write to pin 18: \n%s", err )
	}

	var stopMotor =false
    var motorStopped = false

	runMotor := func() {

		motorStopped = false

		for  !stopMotor  {

			pwm += 10
			log.Printf("pwm = %d", pwm )

			if pwm > 254 {
				break
			}

			pwmDuty := uint32(gobot.FromScale(float64( pwm ), 0, 255) * float64(r.PiBlasterPeriod))

			if err := pwmPin.SetDutyCycle( pwmDuty ); err != nil {
				log.Fatalf("error setting pwm duty: \n%s", err)
			}

			time.Sleep(time.Millisecond * 100)
		}

		if err := pwmPin.InvertPolarity( true ); err !=nil {
			log.Fatalf("error invert  pwm polarity : \n%s", err)
		}

		for  !stopMotor {
			pwm -= 10
			log.Printf("pwm = %d", pwm)

			if pwm < 100 {
				break
			}

			pwmDuty := uint32(gobot.FromScale(float64( pwm ), 0, 255) * float64(r.PiBlasterPeriod))

			if err := pwmPin.SetDutyCycle( pwmDuty ); err != nil {
				log.Fatalf("error setting pwm duty: \n%s", err)
			}

			time.Sleep(time.Millisecond * 100)
		}

		motorStopped = true
	}

	go runMotor()
<-signalChan

close( stopCh )

cleanUp()

r.Finalize()
*/