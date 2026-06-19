module github.com/Aeres-u99/hermes/v2

go 1.26.3

require (
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00
	github.com/tree-sitter/go-tree-sitter v0.25.0
	github.com/tree-sitter/tree-sitter-go v0.23.4
	github.com/tree-sitter/tree-sitter-javascript v0.23.1
	github.com/tree-sitter/tree-sitter-python v0.23.6
	github.com/tree-sitter/tree-sitter-rust v0.23.2
	github.com/tree-sitter/tree-sitter-typescript v0.0.0
)

require github.com/mattn/go-pointer v0.0.1 // indirect

replace github.com/tree-sitter/tree-sitter-typescript => ./internal/grammar/tree-sitter-typescript
