package internal

import (
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/monochromegane/go-gitignore"
)

func AnalyzeFile(path string) (*FileAnalysis, error) {
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

	symbols, index := BuildSymbols(tags, path, lang)

	fileInfo := FileInfo{
		Lang:    lang,
		LOC:     len(strings.Split(string(content), "\n")),
		Imports: imports,
		Symbols: symbols,
	}

	return &FileAnalysis{
		FileInfo: fileInfo,
		Index:    index,
	}, nil
}

func AnalyzeRepo(root string) (*AnalysisResult, error) {
	gitIgnore, err := gitignore.NewGitIgnore(".hermesignore", root)

	merged := &AnalysisResult{
		Files: make(map[string]FileInfo),
		Index: make(map[string]Location),
	}

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		//		matched := gitIgnore.Match(relPath, d.IsDir())

		if gitIgnore.Match(relPath, d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		result, err := AnalyzeFile(path)
		if err != nil {
			return err
		}

		MergeFileAnalysis(path, result, merged)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return merged, nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func MergeFileAnalysis(
	path string,
	file *FileAnalysis,
	result *AnalysisResult,
) {
	result.Files[path] = file.FileInfo

	for k, v := range file.Index {
		MergeIndexLocation(result.Index, k, v)
	}
}
