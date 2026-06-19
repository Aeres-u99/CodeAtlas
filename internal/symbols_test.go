package internal

import "testing"

func TestBuildSymbolsMarksPublicSymbols(t *testing.T) {
	tags := []CTag{
		{Name: "Exported", Kind: "function", Line: 3},
		{Name: "private", Kind: "function", Line: 7},
	}

	symbols, _ := BuildSymbols(tags, "demo.go", "go")

	if !symbols[0].Public {
		t.Fatal("expected exported Go symbol to be public")
	}

	if symbols[1].Public {
		t.Fatal("expected lowercase Go symbol to be private")
	}
}

func TestMergeIndexLocationKeepsDuplicateSymbols(t *testing.T) {
	index := make(map[string]Location)

	MergeIndexLocation(index, "Serve", Location{File: "a.go", Line: 10})
	MergeIndexLocation(index, "Serve", Location{File: "pkg/b.go", Line: 20})

	if got := index["Serve"]; got.File != "a.go" || got.Line != 10 {
		t.Fatalf("bare index overwritten: %#v", got)
	}

	got, ok := index["pkg/b.go#Serve"]
	if !ok {
		t.Fatal("duplicate symbol did not get a qualified index key")
	}

	if got.File != "pkg/b.go" || got.Line != 20 {
		t.Fatalf("qualified index = %#v, want pkg/b.go:20", got)
	}
}
