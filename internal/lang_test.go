package internal

import "testing"

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{path: "main.go", want: "go"},
		{path: "app.js", want: "javascript"},
		{path: "component.ts", want: "typescript"},
		{path: "lib.rs", want: "rust"},
		{path: "Main.java", want: "java"},
		{path: "native.c", want: "c"},
		{path: "native.cpp", want: "cpp"},
		{path: "native.hpp", want: "cpp"},
		{path: "Makefile", want: "makefile"},
		{path: "rules.mk", want: "makefile"},
		{path: "README.md", want: "markdown"},
		{path: "docs.markdown", want: "markdown"},
		{path: "BUILD", want: "bazel"},
		{path: "BUILD.bazel", want: "bazel"},
		{path: "defs.bzl", want: "bazel"},
		{path: "script.lua", want: "lua"},
		{path: "main.tf", want: "terraform"},
		{path: "terraform.tfvars", want: "terraform"},
		{path: "notes.txt", want: "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := DetectLanguage(tt.path); got != tt.want {
				t.Fatalf("DetectLanguage() = %q, want %q", got, tt.want)
			}
		})
	}
}
