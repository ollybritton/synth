package synth

import "fmt"

// Envelope is the interface representing an envelope. Given a time, relative to the start of the envelope,
// it will return an amplitude in the range 0 to 1.
type Envelope interface {
	GetAmplitude(t float64) float64

	Attack(t float64)
	Release(t float64)

	Finished() bool
	Started() bool
}

// ADEnvelope is an envelope with only an attack and decay phase. There is no sustain so all uses of this
// envelope last the same amount of time.
// Since there is no concept of 'release', calling .Release() does nothing.
type ADEnvelope struct {
	Amplitude      float64
	AttackDuration float64
	DecayDuration  float64

	attackTime float64
	started    bool
}

// GetAmplitude returns the given amplitude for a time, relative to the start of the envelope.
func (env *ADEnvelope) GetAmplitude(t float64) float64 {
	if !env.started {
		// No way of telling if attackTime actualy = 0 or it's just been initialised with 0
		return 0
	}

	current := t - env.attackTime

	// Attack
	if current <= env.AttackDuration {
		return (current / env.AttackDuration) * env.Amplitude
	}

	// Decay
	if current > env.AttackDuration && current <= env.AttackDuration+env.DecayDuration {
		return env.Amplitude - (env.Amplitude * ((current - env.AttackDuration) / env.DecayDuration))
	}

	// After decay
	return 0
}

// Attack triggers the start of the attack phase.
func (env *ADEnvelope) Attack(t float64) {
	fmt.Println("attack called", t)
	env.started = true
	env.attackTime = t
}

// Release does nothing as this is an attack-decay envelope.
func (env *ADEnvelope) Release(t float64) {}

// Finished returns true if the envelope has finished.
func (env *ADEnvelope) Finished() bool {
	return getGlobalTime() > (env.attackTime + env.AttackDuration + env.DecayDuration)
}

// Started returns true if the envelope has started.
func (env *ADEnvelope) Started() bool {
	return env.started
}

// NewADEnvelope returns a new, initialised ADEnvelope.
func NewADEnvelope(amp, attackDuration, decayDuration float64) *ADEnvelope {
	return &ADEnvelope{
		Amplitude:      amp,
		AttackDuration: attackDuration,
		DecayDuration:  decayDuration,

		attackTime: 0,
		started:    false,
	}
}

// ASREnvelope is an envelope with an attack, sustain and release phase. When activated, the envelope
// builds up to a constant amplitude, remains at the amplitude while the key is pressed, and then falls back
// down to 0 once the key is released.
type ASREnvelope struct {
	AttackDuration  float64
	ReleaseDuration float64
	Amplitude       float64

	attackTime  float64
	releaseTime float64
	isOn        bool

	started bool
}

// GetAmplitude returns the given amplitude for a time, relative to the start of the envelope.
func (env *ASREnvelope) GetAmplitude(t float64) float64 {
	// No way of telling if attackTime actualy = 0 or it's just been initialised with 0
	if !env.started {
		return 0
	}

	current := t - env.attackTime

	// Attack
	if current <= env.AttackDuration {
		return (current / env.AttackDuration) * env.Amplitude
	}

	// Sustain
	if current > env.AttackDuration && env.isOn {
		return env.Amplitude
	}

	// Release
	current = t - env.releaseTime
	if current <= env.ReleaseDuration {
		return env.Amplitude - (env.Amplitude * ((current) / env.ReleaseDuration))
	}

	// After full release
	return 0
}

// Attack triggers the start of the attack phase.
func (env *ASREnvelope) Attack(t float64) {

	env.attackTime = t
	env.isOn = true
	env.started = true
}

// Release triggers the start of the release phase and the end of the sustain phase.
func (env *ASREnvelope) Release(t float64) {
	env.releaseTime = t
	env.isOn = false
}

// Finished returns true if the envelope has finished.
func (env *ASREnvelope) Finished() bool {
	return getGlobalTime() > (env.releaseTime + env.ReleaseDuration)
}

// Started returns true if the envelope has started.
func (env *ASREnvelope) Started() bool {
	return env.started
}

// NewASREnvelope returns a new attack-sustain-release envelope.
func NewASREnvelope(amplitude, attackDuration, releaseDuration float64) *ASREnvelope {
	return &ASREnvelope{
		AttackDuration:  attackDuration,
		ReleaseDuration: releaseDuration,
		Amplitude:       amplitude,
	}
}

// ADSREnvelope is an envelope with an attack, decay, sustain and release. It is the most common type of envelope.
type ADSREnvelope struct {
	AttackDuration  float64
	DecayDuration   float64
	ReleaseDuration float64

	AttackAmplitude  float64
	SustainAmplitude float64

	attackTime  float64
	releaseTime float64
	isOn        bool

	finished bool
	started  bool
}

// GetAmplitude returns the given amplitude for a time, relative to the start of the envelope.
func (env *ADSREnvelope) GetAmplitude(t float64) float64 {
	// No way of telling if attackTime actualy = 0 or it's just been initialised with 0
	if !env.started {
		return 0
	}

	current := t - env.attackTime

	// Attack
	if current <= env.AttackDuration {
		return (current / env.AttackDuration) * env.AttackAmplitude
	}

	// Decay
	if current > env.AttackDuration && current <= (env.AttackDuration+env.DecayDuration) {
		return env.AttackAmplitude - ((env.AttackAmplitude - env.SustainAmplitude) * ((current - env.AttackDuration) / env.DecayDuration))
	}

	// Sustain
	if current > (env.AttackDuration+env.DecayDuration) && env.isOn {
		return env.SustainAmplitude
	}

	// Release
	current = t - env.releaseTime
	if current <= env.ReleaseDuration {
		return env.SustainAmplitude - (env.SustainAmplitude * ((current) / env.ReleaseDuration))
	}

	// After full release
	// env.finished = true
	return 0
}

// Attack triggers the start of the attack phase.
func (env *ADSREnvelope) Attack(t float64) {
	env.started = true
	// env.finished = false

	env.attackTime = t
	env.isOn = true
}

// Release triggers the start of the release phase and the end of the sustain phase.
func (env *ADSREnvelope) Release(t float64) {
	env.releaseTime = t
	env.isOn = false
}

// Finished returns true if the envelope has finished.
func (env *ADSREnvelope) Finished() bool {
	return getGlobalTime() > (env.releaseTime + env.ReleaseDuration)
}

// Started returns true if the envelope has started.
func (env *ADSREnvelope) Started() bool {
	return env.started
}

// NewADSREnvelope returns a new attack-decay-sustain-release envelope.
func NewADSREnvelope(attackAmplitude, sustainAmplitude, attackDuration, decayDuration, releaseDuration float64) *ADSREnvelope {
	return &ADSREnvelope{
		AttackDuration:   attackDuration,
		DecayDuration:    decayDuration,
		ReleaseDuration:  releaseDuration,
		AttackAmplitude:  attackAmplitude,
		SustainAmplitude: sustainAmplitude,
	}
}
