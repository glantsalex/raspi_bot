package raspi

import (
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

const(
	pinTrigger = "38"
	pinEcho ="40"

	maxDistance = 220
	timeOut = maxDistance * 60

	LOW  = 0
	HIGH = 1
)

type Sonar struct {
	rasp *raspi.Adaptor
	out chan float64
}

func NewSonar() *Sonar {

	return &Sonar{
		rasp: raspi.NewAdaptor(),
		out: make( chan  float64),
	}
}

func ( s* Sonar) Run() chan float64{
	go s.loop()
	return s.out
}

func ( s* Sonar) loop(){

	for{
		distance:= s.getDistance()
		s.out <- distance
		time.Sleep( time.Second )
	}
}

func  ( s* Sonar) getDistance() float64 {

	var pingTime int64 = 0

	distance := 0.0

	s.rasp.DigitalWrite( pinTrigger, HIGH )
	time.Sleep( time.Microsecond * 10 )
	s.rasp.DigitalWrite( pinTrigger, LOW )

	pingTime = int64(s.pulseIn( pinEcho, HIGH, timeOut ))

	distance = float64( pingTime ) * 340.0 / 2.0 / 10000.0

	return distance
}


func ( s* Sonar) pulseIn( pin string, level int, timeout int ) int {

	var tn,t0,t1  time.Time
	var micros int


	valFromSonarNotEqLevel := func() bool {
		val,_ := s.rasp.DigitalRead( pin )
		return val != level
	}

	valFromSonarEqLevel := func() bool {
		val,_ := s.rasp.DigitalRead( pin )
		return val == level
	}

	t0 	= time.Now()

	for valFromSonarNotEqLevel() {

		tn := time.Now()

		if tn.Second()  > t0.Second() {
			micros = 1e6
		}else{
			micros = 0
		}

		micros+=  ( tn.Nanosecond() - t0.Nanosecond()) /1e3

		if micros > timeout {
			log.Printf("TIMEOUT occurred first loop. micros: %d, timout: %d ", micros,timeout)
			return 0
		}
	}

	t1 = time.Now()

	for valFromSonarEqLevel() {

		tn = time.Now()

		if tn.Second() > t0.Second() {
			micros = 1e6
		}else{
			micros = 0
		}

		micros +=  ( tn.Nanosecond() - t0.Nanosecond() ) / 1e3

		if micros > timeout {
			log.Print("TIMEOUT occurred second loop" )
			return 0
		}
	}

	if tn.Second() > t1.Second() {
		micros = 1e6
	}else{
		micros = 0
	}

	micros +=  (tn.Nanosecond() - t1.Nanosecond() )/1e3

	//a:=tn.Nanosecond()/1e3;b:=t1.Nanosecond()/1e3
	//log.Printf("tn: %d, t1: %d, diff: %d micros: %d ", a,b, a-b ,  micros )

	return micros
}