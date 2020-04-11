package synth

import (
	"fmt"
	"sync"
	"time"

	"github.com/mitchellh/copystructure"
)

// Synth defines an Synth.
type Synth struct {
	streamFunc func(amp, freq, t float64) float64
	Env        Envelope

	m    *sync.Mutex
	freq float64
	amp  float64

	finished bool
}

// Stream returns the correct sample for a given point in time `t`.
func (s *Synth) Stream(t float64) float64 {
	s.m.Lock()
	defer s.m.Unlock()

	if s.Env.GetAmplitude(t) <= 0 && s.Env.Started() {
		s.finished = true
	}

	return s.amp * s.streamFunc(s.Env.GetAmplitude(t), s.freq, t)
}

// TriggerAttack triggers the attack phase of the Synth's envelope.
func (s *Synth) TriggerAttack(freq float64) {
	s.SetFreq(freq)
	s.Env.Attack(getGlobalTime())
}

// TriggerRelease triggers the release phase of the Synth's envelope.
func (s *Synth) TriggerRelease() {
	s.Env.Release(getGlobalTime())
}

// TriggerAttackRelease triggers the attack phase of an Synth's envelope, followed by the release phase after the
// specified period of time.
func (s *Synth) TriggerAttackRelease(freq float64, t time.Duration) {
	s.TriggerAttack(freq)
	time.Sleep(t)
	s.TriggerRelease()
}

// SetFreq sets the frequency of the streamer.
func (s *Synth) SetFreq(freq float64) {
	s.m.Lock()
	s.freq = freq
	s.m.Unlock()
}

// SetAmp sets the amplitude of the streamer.
func (s *Synth) SetAmp(amp float64) {
	s.m.Lock()
	s.amp = amp
	s.m.Unlock()
}

// Finished returns true if the synth has finished playing the current tone.
func (s *Synth) Finished() bool {
	s.m.Lock()
	defer s.m.Unlock()

	return s.finished
}

// NewSynth returns a new synth struct with initialised values. The amp value is the maximum amp value.
func NewSynth(streamFunc func(amp, freq, t float64) float64, env Envelope, amp float64) *Synth {
	return &Synth{
		streamFunc: streamFunc,
		Env:        env,
		amp:        amp,

		m: &sync.Mutex{},
	}
}

// PolySynth defines a synth capable of polyphony.
type PolySynth struct {
	base  *Synth
	m     *sync.Mutex
	notes map[float64]*note
}

type note struct {
	synth *Synth

	on  float64
	off float64
}

// NewPolySynth takes an existing synth and makes it capable of multiple voices.
func NewPolySynth(synth *Synth) *PolySynth {
	return &PolySynth{
		base:  synth,
		m:     &sync.Mutex{},
		notes: map[float64]*note{},
	}
}

// Stream returns the correct sample for a given point in time `t`.
func (ps *PolySynth) Stream(t float64) float64 {
	sum := 0.0

	ps.m.Lock()
	defer ps.m.Unlock()

	for freq, note := range ps.notes {
		synth := note.synth
		val := synth.Stream(t)

		// TODO: remove from synth map if empty
		if val == 0 && note.synth.Finished() && note.off > note.on {
			delete(ps.notes, freq)
		}

		sum += val
		// fmt.Println(val)
	}

	return sum
}

// addSynth adds a synth to the internal synth map, copied from the base synth. It also returns a copy of the synth made.
func (ps *PolySynth) addSynth(freq float64) *note {
	copied, err := copystructure.Copy(ps.base)
	if err != nil {
		panic(err)
	}

	s := copied.(*Synth)
	s.m = &sync.Mutex{}
	s.streamFunc = ps.base.streamFunc
	s.SetFreq(freq)
	s.SetAmp(ps.base.amp)

	ps.m.Lock()
	defer ps.m.Unlock()

	n := &note{synth: s}

	ps.notes[freq] = n
	fmt.Println(s)
	return n
}

// getSynth attempts to get a synth for a given frequency
func (ps *PolySynth) getSynth(freq float64) (*Synth, bool) {
	ps.m.Lock()
	defer ps.m.Unlock()

	return ps.notes[freq].synth, ps.notes[freq] != nil
}

// TriggerAttack triggers the attack phase of the Synth's envelope.
func (ps *PolySynth) TriggerAttack(freq []float64) {
	for _, f := range freq {
		note := ps.addSynth(f)
		note.synth.Env.Attack(getGlobalTime())
		note.on = getGlobalTime()
	}
}

// TriggerRelease triggers the release phase of the Synth's envelope.
func (ps *PolySynth) TriggerRelease(freq []float64) {
	for _, f := range freq {
		for _, note := range ps.notes {
			if note.synth.freq == f {
				note.synth.Env.Release(getGlobalTime())
				note.off = getGlobalTime()
			}

		}
	}
}

// TriggerAttackRelease triggers the attack phase of an Synth's envelope, followed by the release phase after the
// specified period of time.
func (ps *PolySynth) TriggerAttackRelease(t time.Duration, freq []float64) {
	ps.TriggerAttack(freq)
	time.Sleep(t)
	ps.TriggerRelease(freq)
}
