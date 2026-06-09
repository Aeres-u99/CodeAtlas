package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Aeres-u99/hermes/v2/internal"
	sitter "github.com/tree-sitter/go-tree-sitter"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	// go_ts "github.com/tree-sitter/tree-sitter-go/bindings/go"
	// js_ts "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	py_ts "github.com/tree-sitter/tree-sitter-python/bindings/go"
	// rs_ts "github.com/tree-sitter/tree-sitter-rust/bindings/go"
	// ts_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
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
	content, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	output := internal.Output{
		Version:   1,
		Generated: time.Now().UTC().Format(time.RFC3339),
		Files:     make(map[string]internal.FileInfo),
		Index:     make(map[string]internal.Location),
	}
	lang := detectLanguage(inputFile)
	loc := len(strings.Split(string(content), "\n"))
	output.Files[inputFile] = internal.FileInfo{
		Lang:    lang,
		LOC:     loc,
		Imports: []string{},
		Symbols: []internal.Symbol{},
	}
	parser := sitter.NewParser()
	tsLang := sitter.NewLanguage(py_ts.Language())
	parser.SetLanguage(tsLang)
	tree := parser.Parse(content, nil)
	root := tree.RootNode()
	imports := []string{}
	for i := uint(0); i < root.ChildCount(); i++ {
		child := root.Child(i)
		if child.Kind() != "import_statement" {
			continue
		}

		for j := uint(0); j < child.ChildCount(); j++ {
			grandchild := child.Child(j)

			if grandchild.Kind() == "dotted_name" {
				imports = append(imports, grandchild.Utf8Text(content))
			}

		}
	}
	cmd := exec.Command(
		"ctags",
		"--output-format=json",
		"--fields=+n",
		inputFile,
	)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	lines := strings.Split(stdout.String(), "\n")
	var tags []internal.CTag
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var tag internal.CTag

		if err := json.Unmarshal([]byte(line), &tag); err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	symbols := []internal.Symbol{}
	for _, tag := range tags {
		name := tag.Name
		if tag.Scope != "" {
			name = tag.Scope + "." + tag.Name
		}

		symbols = append(symbols, internal.Symbol{
			Name: name,
			Type: mapKind(tag.Kind),
			Line: tag.Line,
		})

		output.Index[name] = internal.Location{
			File: inputFile,
			Line: tag.Line,
		}
	}
	fileInfo := output.Files[inputFile]
	fileInfo.Symbols = symbols
	fileInfo.Imports = imports
	output.Files[inputFile] = fileInfo
	data, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func detectLanguage(path string) string {
	switch filepath.Ext(path) {
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".rs":
		return "rust"
	case ".c":
		return "c"
	case ".bazel":
		return "bazel"
	case ".cpp":
		return "cpp"
	case ".java":
		return "java"
	case ".lua":
		return "lua"
	default:
		return "Unknown"
	}
}

func mapKind(kind string) string {
	switch kind {
	case "class":
		return "cls"
	case "function":
		return "fn"
	case "member":
		return "method"
	case "variable":
		return "var"
	default:
		return kind
	}
}
