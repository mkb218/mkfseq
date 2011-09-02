package fseq

import (
	"os"
	"rand"
	"math"
	)

const syxsize = 4096
const AMP_STEPS_PER_6DB = 8.5

func (f *Fseq) WriteToSyx(filename, title string) (err os.Error) {
	framelen := len(f.Pitches)
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	buf := make([]byte, 0, 4096)
	buf = append(buf, 0xf0,0x43,0x00,0x5e, 0x64, 0x10, 0x60, 0, 0)	
	// title 8 bytes
	t := []byte(title)
	for ; len(t) < 8 ; t = append(t, ' ') {}
	buf = append(buf, t...)
	
	buf = append(buf, make([]byte, 8)...)
	
	buf = write14bit(buf, int16(f.LoopStart))
	buf = write14bit(buf, int16(f.LoopEnd))

	buf = append(buf, byte( f.LoopMode ))
	buf = append(buf, byte( f.SpeedAdjust ))
	buf = append(buf, byte( f.VelocitySens ))
	buf = append(buf, byte( f.PitchMode ))
	buf = append(buf, byte( f.NoteAssign ))
	buf = append(buf, byte( f.PitchTuning ))
	buf = append(buf, byte( f.SeqDelay ))
	buf = append(buf, byte( framelen / 128 - 1 ))
	buf = append(buf, 0, 0)

	switch framelen {
	case 128:
		fallthrough
	case 256:
		fallthrough
	case 384:
		fallthrough
	case 512:
		buf = write14bit(buf, int16(framelen))
	default:
		err = os.NewError("unsupported frame size "+string(framelen))
		return
	}

	var skiplast bool
	for frame:=0; frame<framelen; {	// Every frame:
		var fq int

		// Write the pitch first
		buf := write14bit( buf, int16(freqToSyx(f.Pitches[frame])) )

		freqs := make([]byte,16)
		// Voiced: hi freq bytes, low freq bytes, then level.
		for op:=0; op<8; op++  {
			fq = freqToSyx(f.Voiced[op][frame].Freq)
			freqs[op] = byte(fq>>7) & 0x7f
			freqs[op+8] = byte( fq & 0x7f )
		}
		
		buf = append(buf, freqs...)
		
		for op:=0; op<8; op++ {
			buf = append(buf, byte(ampToSyx(f.Voiced[op][frame].Amp)))
		}
		
		// Unvoiced: hi freq bytes, low freq bytes, then level.
		for op:=0; op<8; op++  {
			fq = freqToSyx(f.Unvoiced[op][frame].Freq)
			freqs[op] = byte(fq>>7) & 0x7f
			freqs[op+8] = byte( fq & 0x7f )
		}
		
		buf = append(buf, freqs...)
		
		for op:=0; op<8; op++ {
			buf = append(buf, byte(ampToSyx(f.Unvoiced[op][frame].Amp)))
		}
		
		// Depending on the number of frames being saved, we may skip over some frames in our FormantSequence.
		// FormantSequences always have 512 frames, but .syx format can have a reduced number.
		// TODO: Interpolation between frames?
		switch framelen {
			case 128:	frame += 4
			case 256:	frame += 2
			case 384:	
				frame += 1
				if skiplast {
					frame++
				}
				skiplast = !skiplast
			 // TODO: interpolate
			case 512:	frame += 1
		}
	}
	

	_, err = file.Write(buf)
	if err != nil {
		return
	}
	
	err = file.Close()
	return
}

func dither(d float64) int {
	dec := d - math.Floor(d)
	if rand.Float64() < dec {
		 return int(math.Ceil( d )) 
	}
	return int(math.Floor( d ))
}

func freqToSyx(freq float64) (syx int) {
	diff := 14.302 - math.Log2(float64(freq))
	syx = dither( float64(0x3fff) - diff * 512)
	return 
}

func ampToSyx(amp float64) (syx int) {
	if amp < 0.000001 { 
		amp = 0.000001
	}
	
	logdrops := math.Log2( 1.0 / amp )
	syx = int(logdrops * AMP_STEPS_PER_6DB)
	if syx > 0x7f {
		syx = 0x7f
	}
	if syx < 0 {
		syx = 0
	}
	
	return
}

func write14bit(buf []byte, v int16) []byte {
	buf = append(buf, byte(v >> 7) & 0x7f)
	buf = append(buf, (byte(v) & 0x7f))
	return buf
}