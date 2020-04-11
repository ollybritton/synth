package synth

import (
	"log"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/hajimehoshi/oto"
)

const sr = beep.SampleRate(44100)

var t = 0.0
var tmu = &sync.Mutex{}

var (
	// NoiseFunc is the global noise function.
	NoiseFunc = makeNoise

	// GlobalMixer is the global mixer.
	GlobalMixer = &Mixer{}
)

// sampleToPCM converts a float in the range -1 to +1 to two bytes of that sample encoded in PCM format.
func sampleToPCM(val float64) (buf []byte) {
	intVal := int16(val * (1<<15 - 1))
	low := byte(intVal)
	high := byte(intVal >> 8)

	return []byte{low, high}
}

func makeNoise(t float64) float64 {
	return GlobalMixer.Stream(t)
}

// init initialises the audio handler and allows the playback of audio by the library's functions
func init() {
	context, err := oto.NewContext(int(sr), 1, 2, sr.N(time.Second/10))
	if err != nil {
		log.Fatalf("could not create audio context: %v", err)
	}

	player := context.NewPlayer()

	go func() {
		defer context.Close()

		log.Println("audio handler started")

		tmu.Lock()
		t = 0.0
		tmu.Unlock()

		for {
			val := NoiseFunc(t)
			player.Write(sampleToPCM(val))

			tmu.Lock()
			t += sr.D(1).Seconds()
			tmu.Unlock()
		}
	}()
}

// getGlobalTime gets the value of t atomically.
func getGlobalTime() float64 {
	tmu.Lock()
	defer tmu.Unlock()

	return t
}
