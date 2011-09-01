package fseq

import "testing"
import "sndfile"


func TestAnalyze(t *testing.T) {
	var i sndfile.Info
	file, err := sndfile.Open("ok.aiff", sndfile.Read, &i)
	if err != nil {
		t.Fatal(err)
	}
	_ = Analyze(file)
	if err != nil {
		t.Fatal(err)
	}
}