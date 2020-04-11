package synth

import (
	"math"
	"math/rand"
)

// OscParams contains parameters which control an oscillator.
type OscParams struct {
	Amplitude, Frequency float64
}

// Freq gets the frequency from a OscParam struct.
// These methods may seem redundant but they are used so that oscillators can easily fufil the Oscillator interface by
// just embedding the OscParams struct.
func (p *OscParams) Freq() float64 {
	return p.Frequency
}

// SetFreq sets the frequency for an OscParam struct.
// These methods may seem redundant but they are used so that oscillators can easily fufil the Oscillator interface by
// just embedding the OscParams struct.
func (p *OscParams) SetFreq(f float64) {
	p.Frequency = f
}

// Amp gets the amplitude from a OscParam struct.
// These methods may seem redundant but they are used so that oscillators can easily fufil the Oscillator interface by
// just embedding the OscParams struct.
func (p *OscParams) Amp() float64 {
	return p.Amplitude
}

// SetAmp sets the amplitude for an OscParam struct.
// These methods may seem redundant but they are used so that oscillators can easily fufil the Oscillator interface by
// just embedding the OscParams struct.
func (p *OscParams) SetAmp(a float64) {
	p.Amplitude = a
}

// Oscillator is any streamer that provides a repeating, oscillating signal.
type Oscillator interface {
	Streamer

	Amp() float64
	Freq() float64
	SetAmp(t float64)
	SetFreq(t float64)
}

// Sine is a sine wave.
type Sine struct {
	*OscParams
}

// Stream generates the required sample for a given point on a sine wave.
func (w *Sine) Stream(t float64) float64 {
	return w.Amp() * math.Sin(W(w.Freq())*t)
}

// NewSine returns a new sine wave.
func NewSine(amp, freq float64) *Sine {
	return &Sine{
		&OscParams{amp, freq},
	}
}

// Square is a square wave.
type Square struct {
	*OscParams
}

// Stream generates the required sample for a given point on a square wave.
func (w *Square) Stream(t float64) float64 {
	if w.Amp()*math.Sin(W(w.Freq())*t) > 0 {
		return 1
	}

	return -1
}

// NewSquare returns a new square wave.
func NewSquare(amp, freq float64) *Square {
	return &Square{
		&OscParams{amp, freq},
	}
}

// AnalogSquare is an analog square wave.
type AnalogSquare struct {
	*OscParams
	Iterations int
}

// Stream generates the required sample for a given point on an analog square wave.
func (w *AnalogSquare) Stream(t float64) float64 {
	output := 0.0

	for i := 1.0; i <= float64(w.Iterations); i++ {
		output += ((1 - math.Cos(math.Pi*i)) / (math.Pi * i)) * math.Sin(W(w.Freq())*t*i)
	}

	return output * (2.0 / math.Pi) * w.Amp()
}

// NewAnalogSquare returns a new analog square wave.
func NewAnalogSquare(amp, freq float64, iterations int) *AnalogSquare {
	return &AnalogSquare{
		&OscParams{amp, freq},
		iterations,
	}
}

// Triangle is a triange wave.
type Triangle struct {
	*OscParams
}

// Stream generates the required sample for a given point on a trainge wave.
func (w *Triangle) Stream(t float64) float64 {
	return math.Asin(w.Amp() * math.Sin(W(w.Freq())*t) * (2 / math.Pi))
}

// NewTriangle returns a new triangle wave.
func NewTriangle(amp, freq float64) *Triangle {
	return &Triangle{
		&OscParams{amp, freq},
	}
}

// Sawtooth is a sawtooth wave.
// This sawtooth wave is built using the mod function.
type Sawtooth struct {
	*OscParams
}

// Stream generates the required sample for a given point on a sawtooth wave.
func (w *Sawtooth) Stream(t float64) float64 {
	return (2.0 * w.Amp() / math.Pi) * (w.Freq()*math.Pi*math.Mod(t, 1.0/w.Freq()) - (math.Pi / 2.0))
}

// NewSawtooth returns a new sawtooth wave.
func NewSawtooth(amp, freq float64) *Sawtooth {
	return &Sawtooth{
		&OscParams{amp, freq},
	}
}

// AnalogSawtooth is an analog sawtooth wave.
// Instead of using mod function, it approximates the value using a summation of sine waves.
type AnalogSawtooth struct {
	*OscParams
	Iterations int
}

// Stream generates the required sample for a given point on an analog sawtooth wave.
func (w *AnalogSawtooth) Stream(t float64) float64 {
	output := 0.0

	for i := 1.0; i <= float64(w.Iterations); i++ {
		output += (math.Sin(i * W(w.Freq()) * t)) / i
	}

	return output * (2.0 / math.Pi) * w.Amp()
}

// NewAnalogSawtooth returns a new sawtooth wave.
func NewAnalogSawtooth(amp, freq float64, i int) *AnalogSawtooth {
	return &AnalogSawtooth{
		&OscParams{amp, freq},
		i,
	}
}

// Noise represents random noise.
type Noise struct {
	*OscParams
}

// Stream generates random samples.
func (w *Noise) Stream(t float64) float64 {
	return w.Amp() * (2.0*rand.Float64() - 1)
}

// NewNoise returns a new noise oscillator.
func NewNoise(amp float64) *Noise {
	return &Noise{
		&OscParams{amp, 0},
	}
}
