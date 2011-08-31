package mkfseq

import "math"
import "fmt"
import "os"
import "github.com/runningwild/go-fftw"

const Frames = 512
const FftBins = 1024
const SpectrumBands = FftBins / 2 - 1

type SpecFrame []float64
type SpectralAnalysis struct {
    Frames []SpecFrame
    Freqs []float64
}

func Analyze(af AudioFile) (s *SpectralAnalysis) {
    s = new(SpectralAnalysis)
    s.Frames = make([]SpecFrame, Frames)
    
    var maxPower float64
    for f := uint(0); f < Frames; f++ {
        samps, err := af.GetSamplesAt(FftBins, f * af.Length() / Frames)
        if err != nil {
            panic(err)
        }

        c := fftw.Alloc1d(len(samps))
		o := fftw.Alloc1d(len(samps))
        for i := 0; i < FftBins; i++ {
            hann := (float64(1) - math.Cos(float64(i)/float64(FftBins))*math.Pi*2.0)
            hamming := hann*0.92 + 0.08;
            c[i] = complex(samps[i] * hamming, 0)
			defer func() {
				if x := recover(); x != nil {
					fmt.Fprintf(os.Stderr, "panicked at i = %d\n", i)
				}
			}()
        }
        
        p := fftw.PlanDft1d(c, o, fftw.Forward, 0)
		p.Execute()
        s.Frames[f] = make(SpecFrame, SpectrumBands)
        for i := 0; i < SpectrumBands; i++ {
            pow := math.Sqrt( math.Pow(real(c[i+1]), 2) + math.Pow(imag(c[i+1]), 2) )
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

