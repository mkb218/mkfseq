package fseq

import "testing"
import "sndfile"

func TestAnalyze(t *testing.T) {
	var i sndfile.Info
	file, err := sndfile.Open("ok.aiff", sndfile.Read, &i)
	if err != nil {
		t.Fatal(err)
	}
	s := spectral_analyze(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Freqs) != Frames || len(s.Frames) != Frames {
		t.Error("unexpected size for data members")
	}
}
