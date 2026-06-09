package internal

import (
	"time"
)

func BuildOutput(
	path string,
	result *AnalysisResult,
) Output {
	output := Output{
		Version:   1,
		Generated: time.Now().UTC().Format(time.RFC3339),
		Files:     make(map[string]FileInfo),
		Index:     make(map[string]Location),
	}

	output.Files[path] = result.FileInfo

	for k, v := range result.Index {
		output.Index[k] = v
	}

	return output
}
