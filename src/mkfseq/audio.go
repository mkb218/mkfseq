package mkfseq

import "os"

type AudioFile interface {
    GetSamplesAt(s, at uint) ([]float64, os.Error)
    Length() uint
}