package mkfseq

//#cgo CFLAGS: -I/opt/local/include 
//#cgo LDFLAGS: -L/opt/local/lib -ldjbfft
//#include <fftc8.h>
import "C"

import "math"
import "unsafe"

const Frames = 512
const FftBins = 1024
const SpectrumBands = FftBins / 2 - 1

type SpecFrame []float64
type SpectralAnalysis struct {
    Frames []SpecFrame
    Freqs []float64
}

func sampsToComplex(s []float64) []C.complex8 {
    c := make([]C.complex8, len(s))
    for i, r := range s {
        c[i].re = C.real8(r)
    }
    return c
}

func Analyze(af AudioFile) (s *SpectralAnalysis) {
    s = new(SpectralAnalysis)
    s.Frames = make([]SpecFrame, Frames)
    
    var maxPower float64
    for f := uint(0); f < Frames; f++ {
        samps, err := af.GetSamplesAt(1024, f / Frames * af.Length())
        if err != nil {
            panic(err)
        }
        c := sampsToComplex(samps)
        for i := 0; i < FftBins; i++ {
            hann := (float64(1) - math.Cos(float64(i)/float64(FftBins))*math.Pi*2.0)
            hamming := hann*0.92 + 0.08;
            c[i].re *= C.real8(hamming)
        }
        
        C.fftc8_1024((*C.complex8)(unsafe.Pointer(&c[0])))
        s.Frames[f] = make(SpecFrame, SpectrumBands)
        for i := 0; i < SpectrumBands; i++ {
            pow := math.Sqrt( math.Pow(float64(c[i+1].re), 2) + math.Pow(float64(c[i+1].im), 2) )
            s.Frames[f][i] = pow
            if pow > maxPower {
                maxPower = pow
            }
        }
    }
    
    for f := 0; f < Frames; f++ {
        for i := 0; i < SpectrumBands; i++ {
            s.Frames[f][i] /= maxPower
        }
    }
    
    s.Freqs = make([]float64, Frames)
    for i := 0; i < Frames; i++ {
        s.Freqs[i] = float64(((float64(44100) * 0.5) / (Frames+float64(1))) * float64(i+1))
    }

    return
}

