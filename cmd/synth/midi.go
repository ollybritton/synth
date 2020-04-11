package main

import (
	"log"
	"math"

	"github.com/rakyll/portmidi"
)

func midiToFreq(num int64) float64 {
	return math.Pow(2, float64(num-69)/12.0) * 440
}

// OpenMidiInput opens a midi input channel. It should only be called once because it initializes the portmidi library.
func OpenMidiInput() *portmidi.Stream {
	portmidi.Initialize()

	in, err := portmidi.NewInputStream(portmidi.DefaultInputDeviceID(), 1024)
	if err != nil {
		log.Fatal(err)
	}

	return in
}
