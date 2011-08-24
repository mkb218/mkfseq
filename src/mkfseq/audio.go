package audio

import "os"
import "fmt"
import "io"
import "bufio"
import binary "encoding/binary"

type AudioFile interface {
	GetSamplesAt(s, at uint) ([]float64, os.Error)
	Length() uint
}

type AIFF []float64

func (a AIFF) Length() uint {
	return uint(len(a))
}

func (a AIFF) GetSamplesAt(s, at uint) (f []float64, e os.Error) {
	e = nil
	f = (a)[s:s+at-1]
	return
}

func ParseAIFF(in io.ReadSeeker) (a AIFF, error os.Error) {
	var channels, bits uint
	for chunk := 0; chunk < 100; chunk++ {
		readbuf := make([]byte, 4)
		read, err := in.Read(readbuf)
		if read != 4 {
			break
		}
		if err != nil {
			panic(err)
		}
		var chunklen uint32
		err = binary.Read(in, binary.BigEndian, chunklen)
		var commset bool
		switch string(readbuf) {
		case "FORM":
			aiffCode := make([]byte, 4)
			read, err = in.Read(aiffCode)
			if read != 4 || err != nil {
				error = err
				return
			}
			if string(aiffCode) != "AIFF" && string(aiffCode) != "AIFC" {
				_, err = in.Seek(-4, 1)
				if err != nil {
					error = err
					return
				}
				var filesize uint32
				err = binary.Read(in, binary.BigEndian, filesize)
				read, err = in.Read(aiffCode)
				if read != 4 || err != nil {
					error = err
					return
				}
				if string(aiffCode) != "AIFF" && string(aiffCode) != "AIFC" {
					error = os.NewError("FORM does not precede AIFF: " + string(aiffCode))
					return
				}
			}
			break
		case "COMM":
			if chunklen != 18 {
				error = os.NewError("COMM chunk is not size 18: size=" + string(chunklen))
				return
			}
			var samplecount uint32
			err = binary.Read(in, binary.BigEndian, channels)
			if err != nil {
				error = err
				return
			}
			err = binary.Read(in, binary.BigEndian, samplecount)
			if err != nil {
				error = err
				return
			}
			a = make([]float64, samplecount)
			err = binary.Read(in, binary.BigEndian, bits)
			if err != nil {
				error = err
				return
			}
			// 44100 is assumed, ignore 80-bit float
			_, err = in.Seek(10, 1)
			if err != nil {
				error = err
				return
			}
			commset = true
		case "SSND":
			// sound data woot
			if !commset {
				error = os.NewError("Sound data appeared before COMM chunk. Can't parse this file.")
				return
			}
			// Block-aligned data: Some AIFF sound data will be left- (and maybe right-) padded with sound data
			// that should not be read.
			var offset uint32
			err = binary.Read(in, binary.BigEndian, offset)
			if err != nil {
				error = err
				return
			}
			var blockSize uint32
			err = binary.Read(in, binary.BigEndian, blockSize)
			if err != nil {
				error = err
				return
			}
			seekbase, err := in.Seek(int64(offset), 1)
			if err != nil {
				error = err
				return
			}

			bufr, err := bufio.NewReaderSize(in, int(bits/8*1024))
			if err != nil {
				panic(err)
			}
			// read in sample data while converting to mono
			var maxsampf float64 = float64(int64(0x7ffffffff) >> (32 - bits))
			for index := 0; index < len(a); index++ {
				var sampf float64
				for c := uint(0); c < channels; c++ {
					var sampi int
					switch bits {
					case 8:
						{
							var tmp int8
							binary.Read(bufr, binary.BigEndian, tmp)
							sampi = int(tmp)
						}
					case 16:
						{
							var tmp int16
							binary.Read(bufr, binary.BigEndian, tmp)
							sampi = int(tmp)
						}
					case 32:
						binary.Read(bufr, binary.BigEndian, sampi)
					case 24:
						{
							var tmp1 int8
							var tmp2 int16
							binary.Read(bufr, binary.BigEndian, tmp1)
							binary.Read(bufr, binary.BigEndian, tmp2)
							sampi = int(int16(tmp1)<<16 | tmp2)
						}
					default:
						panic("wtf bit depth")
					}
					sampf += float64(sampi) / maxsampf
				}
				a[index] = sampf
			}
			_, err = in.Seek(seekbase+int64(uint(len(a))*channels*(bits>>3)), 0)
			if err != nil {
				panic(err)
			}

		default:
			fmt.Fprintf(os.Stderr, "wtf chunk type! %s\n", readbuf)
			fallthrough
		// Things we're not going to parse:
		case "FVER":
			fallthrough // format version
		case "MARK":
			fallthrough // markers
		case "INST":
			fallthrough // instrument data
		case "MIDI":
			fallthrough // midi data
		case "AESD":
			fallthrough // recording-related
		case "APPL":
			fallthrough // application-specific data
		case "COMT":
			fallthrough // comments
		case "NAME":
			fallthrough // name
		case "AUTH":
			fallthrough // author
		case "(c) ":
			fallthrough // copyright
		case "ANNO": // annotation
			// skip over these chunks
			_, err = in.Seek(int64(chunklen), 1)
			if err != nil {
				panic(err)
			}
		}

	}
	return
}
