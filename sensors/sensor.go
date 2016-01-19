package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"github.com/rjrobinson/go-powerplant/dto"
	"github.com/rjrobinson/go-powerplant/qutils"
)

var (
	freq     = flag.Uint("freq", 5, "update frequency in cycle/sex")
	max      = flag.Float64("max", 5., "max value generated")
	min      = flag.Float64("min", 1., "min value generated")
	nom      = (*max-*min)/2 + *min
	r        = rand.New(rand.NewSource(time.Now().UnixNano()))
	stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")
	url      = "ampq://guest@localhost:5672"
	value    = r.Float64()*(*max-*min) + *min
	name     = flag.String("name", "sensor", "name of the sensor")
)

func main() {

	flag.Parse()
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)
	sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)

	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish("", sensorQueue.Name, false, false, msg)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()

		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			TimeStamp: time.Now(),
		}

		buf.Reset()
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"",             //exchange,
			dataQueue.Name, // key,
			false,          // mandatory,
			false,          //immediate,
			msg)
		log.Printf("Reading sent. Value %v\n", value)
	}
}
func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64() * (maxStep - minStep)
}
