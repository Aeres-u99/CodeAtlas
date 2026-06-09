package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Aeres-u99/hermes/v2/internal"
	"os"
)

func main() {
	var inputFile string
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "%s - The Code Map you will Ever need!", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "Have Fun!")
	}
	flag.StringVar(&inputFile, "input", "code.py", "Code to Parse")
	flag.Parse()
	result, err := internal.AnalyzeFile(inputFile)
	if err != nil {
		panic(err)
	}

	output := internal.BuildOutput(
		inputFile,
		result,
	)

	data, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
