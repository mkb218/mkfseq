package fseq

//import "fmt"
import "sndfile"
//import "os"
import "math"
import "sort"

type OperatorFrame struct {
	Amp, Freq float64
}

type Operator []OperatorFrame

type LoopModeT int

const (
	OneWay LoopModeT = 0
	Round  LoopModeT = 1
)

type PitchModeT int

const (
	FseqPitch PitchModeT = 0
	FreePitch PitchModeT = 1
)

type Fseq struct {
	Pitches                           []float64
	Voiced, Unvoiced                  [8]Operator
	Title                             string
	LoopStart, LoopEnd                int
	LoopMode                          LoopModeT
	SpeedAdjust, VelocitySens         int
	PitchMode                         PitchModeT
	NoteAssign, PitchTuning, SeqDelay int
}

func newFseq(length int) (f *Fseq) {
	f = new(Fseq)
	f.Pitches = make([]float64, 0, length)
	for i := range f.Voiced {
		f.Voiced[i] = make(Operator, length)
		f.Unvoiced[i] = make(Operator, length)
	}
	f.Title = "Untitled"
	f.LoopEnd = 511
	f.LoopMode = OneWay
	f.SpeedAdjust = 26
	f.PitchMode = FseqPitch
	f.NoteAssign = 54
	f.PitchTuning = 63
	return
}

func CreateFseq(af *sndfile.File, length int, fftbins int) (f *Fseq) {
	f = newFseq(length)
	f.fdetect(af, length, fftbins)
	f.pdetect(af, 55.0, 880.0, length, fftbins)
	return
}

const FormantDetectBandwidth = 7
const VoicedEnergyRatio = 0.15
const UnvoicedEnergyRatio = 0.25
const VoicedOps = 8
const UnvoicedOps = 8
const ImportHighestFormantFreq float64 = 10000.0
const FormantDetectDisallowNeighbors = 7
const SampleRate = 44100.0
const QuickCombInterval = 20

func clamp(signal, min, max float64) float64 {
	if signal > max {
		return max
	} else if signal < min {
		return min
	}

	return signal
}

func (f *Fseq) pdetect(af *sndfile.File, lower, upper float64, length, fftbins int) {
	windowWidth := SampleRate / lower
	t := math.Trunc(windowWidth)
	for index := 0; index < length; index++ {
		var width uint
		if int(math.Trunc((windowWidth-t)*10.0)) < 5 {
			width = uint(t)
		} else {
			width = uint(math.Ceil(windowWidth))
		}
		full := make([]float64, int32(width*2)*af.Format.Channels)

		_, err := af.Seek(af.Format.Frames * int64(index) / int64(length), sndfile.Set)
		if err != nil {
			panic(err)
		}

		read, err := af.ReadFrames(full)
		if err != nil {
			panic(err)
		}
		
		if sumread := read ; read != int64(width*2) {
			_, err = af.Seek(0, sndfile.Set)

			if err != nil {
				panic(err)
			}

			read, err = af.ReadFrames(full[(int32(sumread)*af.Format.Channels):])
			
			if err != nil {
				panic(err)
			}

			if read + sumread != int64(width*2) {
				width = uint((read + sumread) / 2)
			}
			
		}
		samps := mixdown(full, int(af.Format.Channels))
		bestComb := uint(0)
		bestPower := float64(-999999)

		combLow := width
		combHigh := uint(math.Ceil(SampleRate / upper))
		for comb := combLow; comb >= combHigh; comb -= QuickCombInterval {
			var power float64
			for w := uint(0); w < comb; w++ {
				power += samps[w] * samps[w+comb]
			}

			if power > bestPower {
				bestPower = power
				bestComb = comb
			}
		}

		f.Pitches = append(f.Pitches, SampleRate/float64(bestComb))
	}
}

func (f *Fseq) fdetect(af *sndfile.File, length, fftbins int) {
	s := spectral_analyze(af, length, fftbins)
	window := make([]float64, FormantDetectBandwidth)
	for i, _ := range window {
		window[i] = math.Cos((float64(i) / float64(FormantDetectBandwidth+1)) * math.Pi * 0.5)
	}

	for i, frame := range s.Frames {
		vpowers := make([]float64, len(frame))
		upowers := make([]float64, len(frame))
		vfreqs := make([]float64, len(frame))
		ufreqs := make([]float64, len(frame))
		vuratios := make([]float64, len(frame))

		for j, _ := range frame {
			unwindowed := make([]float64, len(window)*2+1)
			windowed := make([]float64, len(window)*2-1)
			var power float64
			var freqSum float64

			for w := 1 - len(window); w < len(window); w++ {
				band := w + 1

				setIndex := w + len(window) - 1

				if band < 0 || len(frame) <= band {
					windowed[setIndex] = 0
					unwindowed[setIndex] = 0
					continue
				} else {
					bandFreq := s.Freqs[band]
					unwindowed[setIndex] = frame[band]
					y := w
					if y < 0 {
						y *= -1
					}
					thisPower := frame[band] * window[y]
					windowed[setIndex] = thisPower
					power += thisPower
					freqSum += bandFreq * thisPower
				}
			}

			var maxEnergyRatio float64
			for w, _ := range windowed {
				if unwindowed[w]/power > maxEnergyRatio {
					maxEnergyRatio = unwindowed[w] / power
				}
			}

			ufreqs[j] = freqSum / power
			vfreqs[j] = ufreqs[j]

			vuratio := clamp((maxEnergyRatio-UnvoicedEnergyRatio)/(VoicedEnergyRatio-UnvoicedEnergyRatio), 0.0, 1.0)
			vuratios[j] = vuratio
			vpowers[j] = power * vuratio
			upowers[i] = power * (1.0 - vuratio)
		}
		pickStoreFormants(f, i, VoicedOps, vpowers, vfreqs, upowers, ufreqs, s.Freqs, vuratios)
	}
}

func pickStoreFormants(f *Fseq, i, count int, vpowers, vfreqs, upowers, ufreqs, bandFreqs, vuratios []float64) {
	okToPick := make([]bool, len(vpowers))
	for j := 0; j < len(vpowers); j++ {
		okToPick[j] = bandFreqs[j] < ImportHighestFormantFreq
	}

	bestIndexes := make([]int, count)
	for v, _ := range bestIndexes {
		bestIndex := -1
		for k, _ := range vpowers {
			if !okToPick[k] {
				continue
			}
			if bestIndex == -1 || vpowers[k] > vpowers[bestIndex] {
				bestIndex = k
			}
		}
		if bestIndex == -1 {
			bestIndex = len(vpowers) - 1 // probably bad
		}
		bestIndexes[v] = bestIndex

		disallow := FormantDetectDisallowNeighbors
		for j := disallow * -1; j < disallow; j++ {
			disIndex := bestIndex + 1
			if disIndex < 0 || len(okToPick) <= disIndex {
				continue
			}
			okToPick[disIndex] = false
		}
	}
	sort.Ints(bestIndexes)

	for v := 0; v < VoicedOps; v++ {
		b:= bestIndexes[v]
		vp := vpowers[b]
		vf := vfreqs[b]
		up := upowers[b]
		uf := ufreqs[b]
		fv := &(f.Voiced[v])
		(*fv)[i] = OperatorFrame{vp,vf}
		fu := &(f.Unvoiced[v])
		(*fu)[i] = OperatorFrame{up,uf}
	}
}
