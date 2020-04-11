package synth

// Streamer represents anything that produces audio at a given point.
// The method Stream returns a value in the range -1 to 1 for a given time `t`.
type Streamer interface {
	Stream(t float64) float64
}

// StreamerFunc is a streamer made of a function.
type StreamerFunc func(t float64) float64

// Stream returns the corresponding samples from a streamer function.
func (f StreamerFunc) Stream(t float64) float64 {
	return f(t)
}

// Mixer is a streamer which combines the streams from lots of different streamers.
type Mixer struct {
	streamers []Streamer
}

// Stream streams the combination of several streamers.
func (m *Mixer) Stream(t float64) float64 {
	sum := 0.0

	for _, streamer := range m.streamers {
		sum += streamer.Stream(t)
	}

	return sum
}

// Add adds one or more streamers to the mixer.
func (m *Mixer) Add(streamers ...Streamer) {
	m.streamers = append(m.streamers, streamers...)
}

// NewMixer returns a new mixer for a set of streamers.
func NewMixer(streamers ...Streamer) *Mixer {
	return &Mixer{streamers: streamers}
}
