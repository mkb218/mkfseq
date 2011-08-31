package main

import "os"
import "mkfseq"
import "fmt"

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	af, err := mkfseq.ParseAIFF(file)
	fmt.Printf("%v\n", af)
	if err != nil {
		panic(err)
	}
	f := mkfseq.CreateFseq(af)	
}