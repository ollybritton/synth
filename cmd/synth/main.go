package main

import (
	"fmt"

	"github.com/ollybritton/synth"
	"github.com/ollybritton/synth/instruments"
)

func main() {
	// env := synth.NewADSREnvelope(0.2, 0.8, 0.05, 0.03, 0.05)
	// env := synth.NewASREnvelope(0.2, 0.05, 0.6)

	i := instruments.Bell()
	i.SetAmp(0.05)
	s := synth.NewPolySynth(i)

	synth.GlobalMixer.Add(s)

	in := OpenMidiInput()
	defer in.Close()

	for event := range in.Listen() {
		switch event.Status {
		case 144:
			fmt.Println("")
			fmt.Println("midi on", event.Data1)
			s.TriggerAttack([]float64{midiToFreq(event.Data1)})

		case 128:
			fmt.Println("midi off", event.Data1)
			s.TriggerRelease([]float64{midiToFreq(event.Data1)})
		}
	}

	select {}

}
