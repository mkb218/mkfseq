package fseq

import "os"
import "fmt"
import "testing"
import "sndfile"

func TestAnalyze(t *testing.T) {
	var i sndfile.Info
	file, err := sndfile.Open("ok.aiff", sndfile.Read, &i)
	if err != nil {
		t.Fatal(err)
	}
	s := spectral_analyze(file, 512, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Freqs) != 511 || len(s.Frames) != 512 {
		t.Errorf("unexpected size for data members freqs %d frames %d", len(s.Freqs), len(s.Frames))
	}
	f, err := os.Create("spectrum.csv")
	if err != nil {
		return
	}
	fmt.Fprint(f, "freq")
	for j := 0; j < len(s.Frames[0]); j++ {
		fmt.Fprintf(f, ",%d", j)
	}
	fmt.Fprint(f, "\n")
	
	for i := 0; i < len(s.Freqs); i++ {
		fmt.Fprintf(f, "%f", s.Freqs[i])
		for j := 0; j < len(s.Frames[i]); j++ {
			fmt.Fprintf(f, ",%f", s.Frames[i][j])
		}
		fmt.Fprint(f, "\n")
	}
}
