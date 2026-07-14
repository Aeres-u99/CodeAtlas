package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/monochromegane/go-gitignore"
)

func AnalyzeFile(path string) (*FileAnalysis, error) {
	slog.Debug("Analyzing File", "file", path)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lang := DetectLanguage(path)
	slog.Debug("Detected Language", "language", lang)

	imports := ExtractImports(content, lang)
	slog.Debug("Found Imports: ", "imports", imports)
	slog.Debug("Found Tags")
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
	slog.Debug("Found symbols", "symbols", symbols)

	return &FileAnalysis{
		FileInfo: fileInfo,
		Index:    index,
	}, nil
}

func AnalyzeRepo(root string) (*AnalysisResult, error) {
	ignoreFile := ".codeatlasignore"
	_, err := os.Stat(ignoreFile)
	if os.IsNotExist(err) {
		slog.Debug("Creating the default .codeatlasignore file", "file", ignoreFile)
		if err := os.WriteFile(ignoreFile, []byte(DefaultCodeAtlasIgnore), 0644); err != nil {
			slog.Debug("Problem Creating CodeAtlasIgnore file", "file", ".codeatlasignore")
			return nil, fmt.Errorf("failed to create %s: %w", ignoreFile, err)
		}
		slog.Info("Generated Default CodeAtlasIgnore File", "file", ignoreFile)
	}
	gitIgnore, err := gitignore.NewGitIgnore(".codeatlasignore", root)

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
	slog.Debug("Checking for directory: ", "isdir", info)
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
	slog.Debug("Merging Analysis")
	result.Files[path] = file.FileInfo
	for k, v := range file.Index {
		MergeIndexLocation(result.Index, k, v)
	}
}
