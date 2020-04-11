package instruments

import "github.com/ollybritton/synth"

// Harmonica returns a basic harmonica-like instrument.
// Let's be real here, it sounds nothing like a harmonica.
func Harmonica() *synth.Synth {
	return synth.NewSynth(
		func(amp, freq, t float64) float64 {
			output := 0.0

			output += 1.00 * synth.NewAnalogSquare(amp, freq, 30).Stream(t+synth.NewSine(0.001, 5).Stream(t))
			output += 0.50 * synth.NewAnalogSquare(amp, freq*2, 30).Stream(t)
			output += 0.05 * synth.NewNoise(amp).Stream(t)

			return output
		},
		synth.NewADSREnvelope(1, 0.95, 0.05, 0.1, 0.2),
		0.2,
	)
}
