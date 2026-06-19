package internal

import (
	"sort"
	"strings"
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"
	go_ts "github.com/tree-sitter/tree-sitter-go/bindings/go"
	js_ts "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	py_ts "github.com/tree-sitter/tree-sitter-python/bindings/go"
	rs_ts "github.com/tree-sitter/tree-sitter-rust/bindings/go"
	ts_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

type LanguageDefinition struct {
	Language func() unsafe.Pointer
	Imports  func(*sitter.Node, []byte) []string
}

func GetLanguage(lang string) *LanguageDefinition {
	registry := map[string]LanguageDefinition{
		"go": {
			Language: go_ts.Language,
			Imports:  extractGoImports,
		},
		"javascript": {
			Language: js_ts.Language,
			Imports:  extractJavaScriptImports,
		},
		"python": {
			Language: py_ts.Language,
			Imports:  extractPythonImports,
		},
		"rust": {
			Language: rs_ts.Language,
			Imports:  extractRustImports,
		},
		"typescript": {
			Language: ts_ts.LanguageTypescript,
			Imports:  extractJavaScriptImports,
		},
	}

	definition, ok := registry[lang]
	if !ok {
		return nil
	}

	return &definition
}

func ExtractImports(content []byte, lang string) []string {
	definition := GetLanguage(lang)
	if definition == nil {
		return []string{}
	}

	parser := sitter.NewParser()
	parser.SetLanguage(sitter.NewLanguage(definition.Language()))

	tree := parser.Parse(content, nil)
	if tree == nil {
		return []string{}
	}

	return uniqueSorted(definition.Imports(tree.RootNode(), content))
}

func extractPythonImports(root *sitter.Node, content []byte) []string {
	imports := []string{}

	walk(root, func(node *sitter.Node) {
		switch node.Kind() {
		case "import_statement":
			imports = append(imports, collectNodeText(node, content, "dotted_name")...)
		case "import_from_statement":
			module := firstNodeText(node, content, "dotted_name", "relative_import")
			if module != "" {
				imports = append(imports, module)
			}
		}
	})

	return imports
}

func extractGoImports(root *sitter.Node, content []byte) []string {
	imports := []string{}

	walk(root, func(node *sitter.Node) {
		if node.Kind() != "import_spec" {
			return
		}

		path := node.ChildByFieldName("path")
		if path == nil {
			path = firstChild(node, "interpreted_string_literal", "raw_string_literal")
		}
		if path != nil {
			imports = append(imports, cleanStringLiteral(path.Utf8Text(content)))
		}
	})

	return imports
}

func extractJavaScriptImports(root *sitter.Node, content []byte) []string {
	imports := []string{}

	walk(root, func(node *sitter.Node) {
		switch node.Kind() {
		case "import_statement", "export_statement":
			source := node.ChildByFieldName("source")
			if source != nil {
				imports = append(imports, cleanStringLiteral(source.Utf8Text(content)))
			}
		case "call_expression":
			function := node.ChildByFieldName("function")
			if function == nil || function.Utf8Text(content) != "import" {
				return
			}

			arguments := node.ChildByFieldName("arguments")
			if arguments == nil {
				return
			}

			source := firstChild(arguments, "string")
			if source != nil {
				imports = append(imports, cleanStringLiteral(source.Utf8Text(content)))
			}
		}
	})

	return imports
}

func extractRustImports(root *sitter.Node, content []byte) []string {
	imports := []string{}

	walk(root, func(node *sitter.Node) {
		if node.Kind() != "use_declaration" {
			return
		}

		text := strings.TrimSpace(node.Utf8Text(content))
		text = strings.TrimPrefix(text, "use")
		text = strings.TrimSuffix(text, ";")
		text = strings.TrimSpace(text)
		if text != "" {
			imports = append(imports, text)
		}
	})

	return imports
}

func walk(node *sitter.Node, visit func(*sitter.Node)) {
	if node == nil {
		return
	}

	visit(node)

	for i := uint(0); i < node.ChildCount(); i++ {
		walk(node.Child(i), visit)
	}
}

func collectNodeText(node *sitter.Node, content []byte, kind string) []string {
	values := []string{}

	walk(node, func(child *sitter.Node) {
		if child.Kind() == kind {
			values = append(values, child.Utf8Text(content))
		}
	})

	return values
}

func firstNodeText(node *sitter.Node, content []byte, kinds ...string) string {
	child := firstChild(node, kinds...)
	if child == nil {
		return ""
	}

	return child.Utf8Text(content)
}

func firstChild(node *sitter.Node, kinds ...string) *sitter.Node {
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		for _, kind := range kinds {
			if child.Kind() == kind {
				return child
			}
		}
	}

	return nil
}

func cleanStringLiteral(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 2 {
		return value
	}

	quote := value[0]
	if quote == '"' || quote == '\'' || quote == '`' {
		return strings.Trim(value, string(quote))
	}

	return value
}

func uniqueSorted(values []string) []string {
	seen := make(map[string]bool)
	unique := []string{}

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}

		seen[value] = true
		unique = append(unique, value)
	}

	sort.Strings(unique)
	return unique
}
