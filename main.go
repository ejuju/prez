package main

import (
	"fmt"
	"os"

	"github.com/ejuju/prez/pkg/prez"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Arguments required: input and output file paths")
		os.Exit(1)
	}
	fpathIn := os.Args[1]
	fpathOut := os.Args[2]

	doc, err := prez.ParseFile(fpathIn)
	if err != nil {
		fmt.Println("Parse file:", err)
		os.Exit(1)
	}

	f, err := os.OpenFile(fpathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Open output file:", err)
		os.Exit(1)
	}
	defer f.Close()

	err = prez.WritePDF(f, doc)
	if err != nil {
		fmt.Println("Export PDF:", err)
		os.Exit(1)
	}
}
