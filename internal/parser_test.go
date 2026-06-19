package internal

import (
	"reflect"
	"testing"
)

func TestExtractImportsByLanguage(t *testing.T) {
	tests := []struct {
		name    string
		lang    string
		content string
		want    []string
	}{
		{
			name: "go",
			lang: "go",
			content: `package main

import (
	"fmt"
	alias "net/http"
)
`,
			want: []string{"fmt", "net/http"},
		},
		{
			name: "javascript",
			lang: "javascript",
			content: `import fs from "fs";
export { thing } from "./thing.js";
const lazy = import("lazy-module");
`,
			want: []string{"./thing.js", "fs", "lazy-module"},
		},
		{
			name: "typescript",
			lang: "typescript",
			content: `import type { User } from "./types";
import { make } from "@app/core";
`,
			want: []string{"./types", "@app/core"},
		},
		{
			name: "rust",
			lang: "rust",
			content: `use std::fmt;
use crate::config::{self, Config};
`,
			want: []string{"crate::config::{self, Config}", "std::fmt"},
		},
		{
			name: "python",
			lang: "python",
			content: `import os
import json as js
from pathlib import Path
`,
			want: []string{"json", "os", "pathlib"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractImports([]byte(tt.content), tt.lang)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ExtractImports() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestGetLanguageUnknown(t *testing.T) {
	if GetLanguage("unknown") != nil {
		t.Fatal("GetLanguage() returned a definition for an unknown language")
	}
}
