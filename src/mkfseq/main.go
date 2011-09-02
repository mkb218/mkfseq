package main

import "flag"
import "os"
import "fseq"
import "sndfile"
import "fmt"

func main() {
	var i sndfile.Info
	file, err := sndfile.Open(os.Args[1], sndfile.Read, &i)
	if err != nil {
		panic(err)
	}
	length := flag.Int("length", 512, "fseq length")
	fftbins := flag.Int("fftbins", 1024, "FFT bins")
	flag.Parse()
	f := fseq.CreateFseq(file, *length, *fftbins)
	fmt.Printf("%v\n", f)
	if err != nil {
		panic(err)
	}
}