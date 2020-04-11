# Synth
Synth is a basic synthesizer libary, written in Go. At the moment is can combine oscillators and evnelopes in order to create basic synths.

```go
// Create a new synth capable of polyphony
s := synth.NewPolySynth(instruments.Harmonica())

// Add it to the global mixer so that it can be heard
synth.GlobalMixer.Add(s)

// Play an A major chord for 1 second
s.TriggerAttackRelease(1 * time.Second, []float64{440, 550, 660})
```

It's stil under development and a lot of stuff is probably wrong, but I hope to add support for it in the future, so I wouldn't recommend using this library.

## Examples
A basic program which accepts midi input and uses it to control a synth can be found in the `cmd/synth` directory.