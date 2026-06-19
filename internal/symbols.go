package internal

import (
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

func MapKind(kind string) string {
	switch kind {
	case "class":
		return "cls"
	case "function":
		return "fn"
	case "member":
		return "method"
	case "variable":
		return "var"
	case "struct":
		return "struct"
	case "const":
		return "const"
	case "package":
		return "package"
	default:
		return kind
	}
}

func BuildSymbols(tags []CTag, file string, lang string) ([]Symbol, map[string]Location) {
	symbols := []Symbol{}
	index := make(map[string]Location)

	for _, tag := range tags {
		name := tag.Name

		if tag.Scope != "" {
			name = tag.Scope + "." + tag.Name
		}

		symbols = append(symbols, Symbol{
			Name:   name,
			Type:   MapKind(tag.Kind),
			Line:   tag.Line,
			Public: IsPublicSymbol(tag, lang),
		})

		AddIndexLocation(index, name, file, tag.Line)
	}

	return symbols, index
}

func MergeIndexLocation(index map[string]Location, name string, location Location) {
	AddIndexLocation(index, name, location.File, location.Line)
}

func AddIndexLocation(index map[string]Location, name string, file string, line int) {
	location := Location{
		File: file,
		Line: line,
	}

	if _, exists := index[name]; !exists {
		index[name] = location
		return
	}

	index[qualifiedIndexName(file, name)] = location
}

func qualifiedIndexName(file string, name string) string {
	return filepath.ToSlash(file) + "#" + name
}

func IsPublicSymbol(tag CTag, lang string) bool {
	if tag.Access == "public" {
		return true
	}

	if tag.Access == "private" || tag.Access == "protected" {
		return false
	}

	name := tag.Name
	if name == "" || tag.Kind == "package" {
		return false
	}

	switch lang {
	case "go":
		return startsWithUpper(name)
	case "python":
		return !strings.HasPrefix(name, "_")
	case "rust":
		return strings.Contains(tag.Pattern, "pub ")
	case "javascript", "typescript":
		return strings.Contains(tag.Pattern, "export ")
	default:
		return false
	}
}

func startsWithUpper(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return r != utf8.RuneError && unicode.IsUpper(r)
}
