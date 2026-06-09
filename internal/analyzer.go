package internal

import (
	"os"
	"strings"
)

func AnalyzeFile(path string) (*AnalysisResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lang := DetectLanguage(path)

	imports := ExtractImports(content, lang)

	tags, err := GetTags(path)
	if err != nil {
		return nil, err
	}

	symbols, index := BuildSymbols(tags, path)

	fileInfo := FileInfo{
		Lang:    lang,
		LOC:     len(strings.Split(string(content), "\n")),
		Imports: imports,
		Symbols: symbols,
	}

	return &AnalysisResult{
		FileInfo: fileInfo,
		Index:    index,
	}, nil
}
