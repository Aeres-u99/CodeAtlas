package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Aeres-u99/hermes/v2/internal"
	"log/slog"
	"os"
)

func main() {
	var inputFile string
	var outputFile string
	var logLevel string

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nHermes - The Code Map You'll Ever Need\n\n")

		fmt.Fprintf(w, "Environment Variables:\n")
		fmt.Fprintf(w, "  HERMES_CTAGS   Path to a custom Universal Ctags binary\n")
		fmt.Fprintf(w, "  LOG_LEVEL      debug | info | warn | error\n\n")

		fmt.Fprintf(w, "Options:\n")
		flag.PrintDefaults()

		fmt.Fprintf(w, "\nHappy Hacking!\n")
		fmt.Fprintf(w, "Made with ♥  by KeiTachikawa\n")
	}

	flag.StringVar(&logLevel, "log-level", "info", "Log Level for debugging and Testing")
	flag.StringVar(&inputFile, "input", "code.py", "Code to Parse")
	flag.StringVar(&outputFile, "output", "hermes.json", "Code to Parse")
	slog.Debug("Found Flags",
		"logLevel", logLevel,
		"input", inputFile,
		"output", outputFile,
	)
	flag.Parse()
	if logLevel == "info" {
		// Verify if it has been set in Environment Variable
		if env := os.Getenv("HERMES_LOG_LEVEL"); env != "" {
			logLevel = env
		}
	}
	internal.InitLogger(logLevel)

	var result *internal.AnalysisResult
	var err error
	slog.Info("Initializing ...")
	if internal.IsDir(inputFile) {
		slog.Debug(
			"Checking Folder",
			"folder", inputFile,
		)
		result, err = internal.AnalyzeRepo(inputFile)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Debug(
			"Checking Folder",
		)
		fileResult, err := internal.AnalyzeFile(inputFile)
		if err != nil {
			panic(err)
		}
		slog.Debug(
			"Performing Analysis",
			"file", inputFile,
		)
		result = &internal.AnalysisResult{
			Files: map[string]internal.FileInfo{
				inputFile: fileResult.FileInfo,
			},
			Index: fileResult.Index,
		}
	}
	output := internal.BuildOutput(result)
	data, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(outputFile, data, 0644)
	slog.Info("Wrote Output", "file", outputFile)
}
