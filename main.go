package main

import (
	"log"
	"os"
	"time"

	"github.com/thejerf/suture"
)

type datapoint struct {
	time        time.Time
	temperature float64
	wattHours   int64
}

var interval = time.Minute
var debug = os.Getenv("EDSDEBUG") != ""

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	edsURL := os.Getenv("EDSURL")
	dbURL := os.Getenv("INFLUXURL")

	results := make(chan datapoint)
	pSrv := &poster{
		url: dbURL,
		in:  results,
	}
	rSrv := &reader{
		url:  edsURL,
		out:  results,
		intv: interval,
	}

	elapsed := time.Now().UnixNano() % interval.Nanoseconds()
	sleep := time.Duration(interval.Nanoseconds() - elapsed)
	if debug {
		log.Println("Waiting", sleep, "to get in step")
	}
	time.Sleep(sleep)

	srv := suture.NewSimple("main")
	srv.Add(pSrv)
	srv.Add(rSrv)
	srv.Serve()
}
