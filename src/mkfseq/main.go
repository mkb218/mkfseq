package main

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
	s := fseq.Analyze(file)
	fmt.Printf("%v\n", len(s.Freqs))
	if err != nil {
		panic(err)
	}
}