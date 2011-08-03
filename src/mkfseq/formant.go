package mkfseq

type OperatorFrame struct {
    Amp, Freq float64
}

type Operator []OperatorFrame

type LoopModeT int

const (
    OneWay LoopModeT = 0
    Round LoopModeT = 1
)

type PitchModeT int

const (
    FseqPitch PitchModeT = 0
    FreePitch PitchModeT = 1
)

type Fseq struct {
    Pitches Operator
    Voiced, Unvoiced [8]Operator
    Title string
    LoopStart, LoopEnd int
    LoopMode LoopModeT
    SpeedAdjust, VelocitySens int
    PitchMode PitchModeT
    NoteAssign, PitchTuning, SeqDelay int
}

func NewFseq (f *Fseq) {
    f = new(Fseq)
    f.Pitches = make(Operator, 0)
    for i := range f.Voiced {
        f.Voiced[i] = make(Operator, 0)
        f.Unvoiced[i] = make(Operator, 0)
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