package internal

import (
	"sort"
	"strings"
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"
	c_ts "github.com/tree-sitter/tree-sitter-c/bindings/go"
	cpp_ts "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	go_ts "github.com/tree-sitter/tree-sitter-go/bindings/go"
	java_ts "github.com/tree-sitter/tree-sitter-java/bindings/go"
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
		"c": {
			Language: c_ts.Language,
			Imports:  extractCImports,
		},
		"cpp": {
			Language: cpp_ts.Language,
			Imports:  extractCImports,
		},
		"go": {
			Language: go_ts.Language,
			Imports:  extractGoImports,
		},
		"java": {
			Language: java_ts.Language,
			Imports:  extractJavaImports,
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
	if lang == "bazel" {
		return uniqueSorted(extractBazelImports(content))
	}
	if lang == "terraform" {
		return uniqueSorted(extractTerraformImports(content))
	}

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

func extractJavaImports(root *sitter.Node, content []byte) []string {
	imports := []string{}

	walk(root, func(node *sitter.Node) {
		if node.Kind() != "import_declaration" {
			return
		}

		text := strings.TrimSpace(node.Utf8Text(content))
		text = strings.TrimPrefix(text, "import")
		text = strings.TrimSpace(text)
		text = strings.TrimPrefix(text, "static")
		text = strings.TrimSuffix(text, ";")
		text = strings.TrimSpace(text)
		if text != "" {
			imports = append(imports, text)
		}
	})

	return imports
}

func extractCImports(root *sitter.Node, content []byte) []string {
	imports := []string{}

	walk(root, func(node *sitter.Node) {
		if node.Kind() != "preproc_include" {
			return
		}

		path := node.ChildByFieldName("path")
		if path != nil {
			imports = append(imports, cleanStringLiteral(path.Utf8Text(content)))
		}
	})

	return imports
}

func extractBazelImports(content []byte) []string {
	imports := []string{}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "load(") {
			continue
		}

		line = strings.TrimPrefix(line, "load(")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		imports = append(imports, cleanStringLiteral(firstBazelLoadArg(line)))
	}

	return imports
}

func extractTerraformImports(content []byte) []string {
	imports := []string{}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "source") {
			continue
		}

		rest := strings.TrimSpace(strings.TrimPrefix(line, "source"))
		if !strings.HasPrefix(rest, "=") {
			continue
		}

		rest = strings.TrimSpace(strings.TrimPrefix(rest, "="))
		if len(rest) < 2 {
			continue
		}

		quote := rest[0]
		if quote != '"' && quote != '\'' {
			continue
		}

		end := strings.IndexByte(rest[1:], quote)
		if end < 0 {
			continue
		}

		value := rest[1 : end+1]
		if value != "" {
			imports = append(imports, value)
		}
	}

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

func firstBazelLoadArg(value string) string {
	if value == "" {
		return ""
	}

	quote := value[0]
	if quote != '"' && quote != '\'' {
		return ""
	}

	end := strings.IndexByte(value[1:], quote)
	if end < 0 {
		return ""
	}

	return value[:end+2]
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
	if quote == '<' {
		return strings.TrimSuffix(strings.TrimPrefix(value, "<"), ">")
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
