package internal

import (
	"log/slog"
	"path/filepath"
)

func DetectLanguage(path string) string {
	slog.Debug("Checking language")
	switch filepath.Base(path) {
	case "BUILD", "BUILD.bazel", "MODULE.bazel", "WORKSPACE", "WORKSPACE.bazel":
		return "bazel"
	case "Makefile":
		return "makefile"
	}

	switch filepath.Ext(path) {
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".js", ".jsx":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".c":
		return "c"
	case ".bazel", ".bzl":
		return "bazel"
	case ".cc", ".cpp", ".cxx", ".h", ".hh", ".hpp", ".hxx":
		return "cpp"
	case ".tf", ".tfvars":
		return "terraform"
	case ".mk":
		return "makefile"
	case ".md", ".markdown":
		return "markdown"
	case ".lua":
		return "lua"
	default:
		return "Unknown"
	}
}
