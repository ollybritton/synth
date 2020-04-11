package instruments

import "github.com/ollybritton/synth"

// Bell returns a basic bell-like instrument.
func Bell() *synth.Synth {
	return synth.NewSynth(
		func(amp, freq, t float64) float64 {
			output := 0.0

			output += 1.00 * synth.NewSine(amp, freq).Stream(t+synth.NewSine(0.0005, 2).Stream(t))
			output += 0.50 * synth.NewSine(amp, freq*2).Stream(t)
			output += 0.05 * synth.NewSine(amp, freq*3).Stream(t)

			return output
		},
		synth.NewADEnvelope(1, 0.1, 1.0),
		0.2,
	)
}
