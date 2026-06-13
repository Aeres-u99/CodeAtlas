# Hermes

Hermes is a lightweight code mapping utility designed for AI-assisted development workflows.

Instead of feeding an entire repository to an LLM, Hermes generates a compact semantic map of the codebase containing files, imports, symbols, and symbol locations. The resulting map can be used alongside the repository tree and targeted source files to dramatically reduce context requirements.

Hermes follows a simple philosophy:

* Keep it small.
* Keep it fast.
* Keep it language-agnostic.
* Generate structure, not explanations.
* Let the LLM decide what code to read next.

## Why Hermes?

Large repositories quickly exceed practical context limits.

Most AI coding assistants spend tokens reading files that are irrelevant to the requested change.

Hermes generates a compact index that answers:

* What files exist?
* What language is each file?
* What symbols are defined?
* Where are those symbols located?
* What imports does each file use?

This allows an AI workflow such as:

Repository Tree
→ Hermes Map
→ User Request
→ Targeted File Reads
→ Code Modification

Instead of:

Repository
→ Dump Everything Into Context
→ Hope For The Best

## Features

### Repository Analysis

Analyze an entire repository and generate a unified semantic map.

### File Analysis

Analyze individual files.

### Symbol Extraction

Uses Universal CTags to extract:

* Functions
* Methods
* Classes
* Structs
* Variables
* Constants
* Packages

### Import Extraction

Uses Tree-sitter to extract language-specific imports.

### Cross-File Index

Builds a global symbol index:

```json
{
  "AnalyzerRepo": {
    "f": "internal/analyzer.go",
    "l": 42
  }
}
```

### Language Detection

Currently supports:

* Go
* Python
* JavaScript
* TypeScript
* Rust
* C
* C++
* Java
* Lua
* Bazel

### Ignore Support

Hermes supports a `.hermesignore` file for excluding:

* Generated code
* Vendor dependencies
* Build artifacts
* Grammar repositories
* Test repositories
* Large irrelevant directories

Example:

```text
.git/
vendor/
node_modules/
internal/grammar/**

*.json
!hermes.json
```

### Compact JSON Schema

Hermes intentionally uses short field names to reduce token consumption.

Example:

```json
{
  "n": "AnalyzeRepo",
  "t": "fn",
  "l": 42
}
```

Where:

| Field | Meaning     |
| ----- | ----------- |
| n     | Symbol Name |
| t     | Symbol Type |
| l     | Line Number |

## Installation

### Requirements

* Go 1.24+
* Universal CTags

Install CTags:

Ubuntu:

```bash
sudo apt install universal-ctags
```

Arch:

```bash
sudo pacman -S ctags
```

macOS:

```bash
brew install universal-ctags
```

Build:

```bash
make
```

## Usage

Analyze a file:

```bash
./hermes -input internal/analyzer.go
```

Analyze a repository:

```bash
./hermes -input .
```

Output:

```json
{
  "files": {
    "internal/analyzer.go": {
      "lang": "go",
      "symbols": [...]
    }
  },
  "idx": {
    "AnalyzeRepo": {
      "f": "internal/analyzer.go",
      "l": 42
    }
  }
}
```

## Example AI Workflow

Generate a repository map:

```bash
./hermes -input . > hermes.json
```

Provide the following to the LLM:

1. Repository tree
2. hermes.json
3. User request

The LLM can then determine which files are relevant before reading source code.

This reduces unnecessary context consumption and improves modification accuracy.

## Current Status

Implemented:

* Repository analysis
* File analysis
* Universal CTags integration
* Tree-sitter parsing
* Symbol indexing
* Import extraction
* `.hermesignore`
* Compact JSON output

Planned:

* Additional language parsers
* Incremental updates
* Filesystem watcher
* Persistent in-memory index
* LSP integration

## Philosophy

Hermes is intentionally narrow in scope.

It does not:

* Build call graphs
* Perform static analysis
* Generate embeddings
* Index documentation
* Replace language servers

Hermes exists to answer one question efficiently:

"What does this repository look like?"
