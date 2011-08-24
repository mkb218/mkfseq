package main

import "os"
import "mkfseq"
import "fmt"

func main() {
	file, err := os.Open(os.Args[0])
	if err != nil {
		panic(err)
	}
	af, err := mkfseq.ParseAIFF(file)
	f := mkfseq.CreateFseq(af)	
	fmt.Printf("%v\n", f)
}