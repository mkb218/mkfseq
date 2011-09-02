package fseq

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
	if len(s.Freqs) != 512 || len(s.Frames) != 512 {
		t.Error("unexpected size for data members")
	}
}
