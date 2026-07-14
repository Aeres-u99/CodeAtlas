# CodeAtlas Progress Summary

## Project Goal

**DISCLAIMER: PART OF THE ROADMAP IS AI GENERATED for the understanding of viability**

CodeAtlas is a lightweight semantic code mapper for AI-assisted development.
The goal is not to perform static analysis, call graph generation, embeddings, or code intelligence.
The goal is:

```text
Repository
    ↓
CodeAtlas
    ↓
Compact Semantic Map
    ↓
LLM decides what files to read
```

This reduces token usage and improves file selection accuracy.

---

# Current Architecture

## File Analysis

```go
AnalyzeFile(path string) (*FileAnalysis, error)
```

Produces:

```go
type FileAnalysis struct {
    FileInfo FileInfo
    Index    map[string]Location
}
```

Responsibilities:

* Detect language
* Extract imports via Tree-sitter
* Extract symbols via Universal CTags
* Build local symbol index

---

## Repository Analysis

```go
AnalyzeRepo(root string) (*AnalysisResult, error)
```

Produces:

```go
type AnalysisResult struct {
    Files map[string]FileInfo
    Index map[string]Location
}
```

Responsibilities:

* Walk repository
* Apply .codeatlasignore
* Analyze each file
* Merge results
* Build global symbol index

---

## Output Generation

```go
BuildOutput(result *AnalysisResult) Output
```

Produces final JSON.

---

# Symbol Extraction

Currently powered by Universal CTags.

Example:

```json
{
  "n": "AnalyzeRepo",
  "t": "fn",
  "l": 42
}
```

Supported symbol mappings:

| CTags    | CodeAtlas  |
| -------- | ------- |
| class    | cls     |
| function | fn      |
| member   | method  |
| variable | var     |
| struct   | struct  |
| package  | package |
| const    | const   |

---

# Import Extraction

Current State:

* Python supported
* Tree-sitter infrastructure exists

Next Phase:

Enable:

* Go
* JavaScript
* TypeScript
* Rust

Tree-sitter should remain focused on imports.

CTags should remain responsible for symbols.

---

# JSON Schema

Current compact schema:

```json
{
  "v": 1,
  "generated": "...",
  "files": {
    "internal/analyzer.go": {
      "lang": "go",
      "loc": 120,
      "imports": [],
      "symbols": []
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

Compression decisions:

```text
n = name
t = type
l = line

f = file
```

to reduce token consumption.

---

# .codeatlasignore

Implemented.

Uses:

```go
github.com/monochromegane/go-gitignore
```

Important lesson learned:

Correct constructor:

```go
gitignore.NewGitIgnore(".codeatlasignore", root)
```

NOT:

```go
gitignore.NewGitIgnore(root, ".codeatlasignore")
```

This bug caused ignore matching to fail.

---

## Current Ignore Rules

```text
.git/

sample/**
pkg/**
internal/grammar/**

output/**

*.pew
codeatlas

*.json
!codeatlas.json

*.md
```

---

## Results

Repository map reduced significantly.

Major reductions came from:

```text
.git/
internal/grammar/**
sample/**
pkg/**
```

Grammar repositories were the biggest source of noise.

---

# Current Repository Structure

Relevant source files:

```text
cmd/codeatlas/main.go

internal/
├── analyzer.go
├── consts.go
├── ctags.go
├── lang.go
├── output.go
├── parser.go
├── structs.go
├── symbols.go
```

Ignored:

```text
internal/grammar/**
sample/**
pkg/**
.git/**
```

---

# Design Decisions

## Accepted

### Use CTags For Symbols

Reason:

* Mature
* Multi-language
* Fast
* Reliable

### Use Tree-sitter For Imports

Reason:

* Better language-specific parsing
* Imports are grammar dependent

### Keep CodeAtlas Small

Avoid:

* Call graphs
* Static analysis
* LSP replacement
* Embeddings
* Semantic search

CodeAtlas should answer:

```text
"What does this repository look like?"
```

Nothing more.

---

# Future Roadmap

## Phase 2

Multi-language Tree-sitter support.

Enable:

* Go
* JavaScript
* TypeScript
* Rust

Create language registry:

```go
GetLanguage(lang string)
```

instead of hardcoded Python parser.

---

## Phase 3

Watcher Support.

Goal:

```text
File Change
    ↓
Regenerate Map
```

Initially:

Full rebuild acceptable.

Target:

3-4 minute rebuild acceptable.

---

## Phase 4

Incremental Updates.

Goal:

```text
Changed File
    ↓
AnalyzeFile()
    ↓
Update In-Memory Index
    ↓
Write Updated Map
```

Avoid full repository rebuilds.

---

# Benchmark Framework

Five benchmark categories defined.

## Test 1

Repository Reduction

Measures:

* LOC
* Files
* CodeAtlas size
* Estimated tokens

---

## Test 2

Symbol Discovery

Measures:

* Symbol coverage
* Correct file
* Correct line

---

## Test 3

Navigation Accuracy

Input:

* Tree
* CodeAtlas Map

No source code.

Measure:

* Correct file prediction

---

## Test 4

Retrieval Efficiency

Compare:

Tree Only

vs

Tree + CodeAtlas

Measure:

* Files opened
* Reduction percentage

---

## Test 5

AI Cost Benchmark

Compare:

Tree Only

vs

Tree + CodeAtlas

Measure:

* Prompt tokens
* Completion tokens
* Modification quality
* File selection accuracy

---

# Current Status

Completed:

✅ Repository Analysis

✅ File Analysis

✅ Universal CTags Integration

✅ Symbol Extraction

✅ Import Extraction (Python)

✅ Global Symbol Index

✅ Compact JSON Schema

✅ .codeatlasignore

✅ Repository Walking

✅ Output Generation

✅ Benchmark Framework

Next Immediate Task:

```text
Enable Go Tree-sitter
Enable JS Tree-sitter
Enable TS Tree-sitter
Enable Rust Tree-sitter
Build Language Registry
```

After that:

```text
Watcher
→ Incremental Updates
→ AI-DLC Integration
```

Use codemap.json, whenever you can, after making any change, ensure that you rerun the codeatlas to build the codemap yet again.
I would like to continue working on this further

### Milestone

Phase 2A
---------
✓ Language Registry
✓ Go Imports
✓ JS Imports
✓ TS Imports
✓ Rust Imports
✓ Duplicate Symbol Handling
✓ Exported/Public Symbol Detection

Phase 2B
---------
✓ Benchmark against 5 large repositories
✓ Measure token savings
✓ Measure file selection accuracy

Phase 3
---------
Watcher


---

 ---
  First, to address the pending statusline task: I was unable to complete it because file write permission was denied. The two
  things that need to be done manually are creating /home/kei/.claude/statusline-command.sh and updating
  /home/kei/.claude/settings.json as described in my previous message. For any future status line changes, please use this
  statusline-setup agent.
  
  ---
  Codemap Effort Estimate

  What the codemap gives me

  The codemap.json is a well-structured snapshot of the codeatlas codebase. It covers:

  - Every source file with language, line count, and import list
  - Every exported and unexported symbol with file and line number, fully qualified (e.g. internal.AnalyzeRepo ->
  internal/analyzer.go:41)
  - A flat idx index for instant symbol lookup

  The codebase itself is small and well-partitioned:

  ┌──────────────────────┬─────┬──────────────────────────────────────────────────┐
  │         File         │ LOC │                       Role                       │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/parser.go   │ 236 │ Tree-sitter import extraction per language       │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/analyzer.go │ 108 │ File/repo walking, gitignore, merging            │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/symbols.go  │ 109 │ Symbol building, index merging, public detection │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/structs.go  │ 51  │ All data types                                   │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/ctags.go    │ 45  │ ctags invocation                                 │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/output.go   │ 25  │ Output assembly                                  │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ internal/lang.go     │ 31  │ Language detection                               │
  ├──────────────────────┼─────┼──────────────────────────────────────────────────┤
  │ cmd/codeatlas/main.go   │ 55  │ CLI entry point                                  │
  └──────────────────────┴─────┴──────────────────────────────────────────────────┘

  Total meaningful Go code: ~660 LOC across 8 files.

  ---
  With codemap

  For most change tasks I would need zero or one follow-up file reads before editing:

  - Adding support for a new language in the parser: I immediately know the relevant functions are internal.GetLanguage
  (parser.go:21), internal.ExtractImports (parser.go:53), and the per-language helpers like internal.extractGoImports
  (parser.go:88). I go straight to the right lines.
  - Adding a new output field: internal.Output is at structs.go:23, internal.BuildOutput is at output.go:7 — two targeted
  reads, two edits.
  - Changing symbol visibility logic: internal.IsPublicSymbol is at symbols.go:77, internal.startsWithUpper at symbols.go:105
  — precise.
  - Understanding data flow end-to-end: the import graph in the codemap (cmd/codeatlas/main.go imports internal,
  internal/analyzer.go uses go-gitignore, internal/parser.go uses tree-sitter bindings) tells me the full dependency chain
  without reading any file.

  Effort overhead: near zero. I can locate any symbol in one lookup, understand the module boundary in seconds, and
  confidently edit without reading the whole codebase.

  ---
  Without codemap

  I would need to read every file before making any non-trivial change, because:

  - The internal package has 8 files and no obvious naming cues for which struct lives where (e.g. Symbol, CTag, FileInfo,
  AnalysisResult could be in any file).
  - The parser file at 236 LOC has 13 symbols — I would read it entirely to understand which per-language function to touch.
  - Cross-file relationships (e.g. BuildOutput consuming AnalyzeRepo output) are invisible until I read both files.

  Effort overhead: I would likely spend 3-5 extra tool calls reading context files before each substantive edit. For a small
  codebase like this (660 LOC) that is manageable but still roughly doubles the back-and-forth. For a larger codebase the cost
  would compound.

  ---
  Summary

  The codemap cuts navigation cost to near zero for this project. Given the codebase is already small, the absolute time saved
  per task is moderate (2-4 reads per change), but the benefit scales directly with codebase size. The symbol index (idx) is
  especially valuable — it is the equivalent of having ctags output pre-indexed and queryable without any tool invocation.
─────────────────────────────────────────────────────────────────────────────────────── Configure statusLine from shell PS1 ──

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  Settings  Status   Config   Usage   Stats

  Session

  Total cost:            $0.2797
  Total duration (API):  1m 36s
  Total duration (wall): 4m 56s
  Total code changes:    11 lines added, 1 line removed
  Usage by model:
     claude-sonnet-4-6:  711 input, 3.1k output, 302.6k cache read, 30.9k cache write ($0.2797)

  Current session
  ███                                                6% used
  Resets 12:20am (Asia/Kolkata)

  Current week (all models)
  ████████████████                                   32% used
  Resets Jun 21, 12:30am (Asia/Kolkata)

  What's contributing to your limits usage?
  Approximate, based on local sessions on this machine — does not include other devices or claude.ai

  Last 24h · these are independent characteristics of your usage, not a breakdown
─────────────────────────────────────────────────────────────────────────────────────── Configure statusLine from shell PS1 ──

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  Settings  Status   Config   Usage   Stats
  
  Session

  Total cost:            $0.2797
  Total duration (API):  1m 36s
  Total duration (wall): 5m 10s
  Total code changes:    11 lines added, 1 line removed
  Usage by model:
     claude-sonnet-4-6:  711 input, 3.1k output, 302.6k cache read, 30.9k cache write ($0.2797)

  Current session
  ███                                                6% used
  Resets 12:20am (Asia/Kolkata)

  Current week (all models)
  ████████████████                                   32% used
  Resets Jun 21, 12:30am (Asia/Kolkata)

  What's contributing to your limits usage?
  Approximate, based on local sessions on this machine — does not include other devices or claude.ai

  Last 24h · these are independent characteristics of your usage, not a breakdown
  I would need to read every file before making any non-trivial change, because:

  - The internal package has 8 files and no obvious naming cues for which struct lives where (e.g. Symbol, CTag, FileInfo,
  AnalysisResult could be in any file).
  - The parser file at 236 LOC has 13 symbols — I would read it entirely to understand which per-language function to touch.
  - Cross-file relationships (e.g. BuildOutput consuming AnalyzeRepo output) are invisible until I read both files.

  Effort overhead: I would likely spend 3-5 extra tool calls reading context files before each substantive edit. For a small
  codebase like this (660 LOC) that is manageable but still roughly doubles the back-and-forth. For a larger codebase the cost
  would compound.

  ---
  Summary

  The codemap cuts navigation cost to near zero for this project. Given the codebase is already small, the absolute time saved
  per task is moderate (2-4 reads per change), but the benefit scales directly with codebase size. The symbol index (idx) is
  especially valuable — it is the equivalent of having ctags output pre-indexed and queryable without any tool invocation.

