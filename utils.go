package synth

import "math"

// W converts a hertz value (e.g. 440) to an angular frequency value.
// https://en.wikipedia.org/wiki/Angular_frequency
func W(freq float64) float64 {
	return 2 * math.Pi * freq
}
