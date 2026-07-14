<div align="center">

<img src="assets/logo.png" alt="CodeAtlas Logo" width="220"/>

# CodeAtlas 🗺️ 

### Symbolic Code Index. Instant Navigation.

**A deterministic code indexing and retrieval system for AI-assisted development.**


In Greek mythology, CodeAtlas was the messenger of the gods, known for speed, travel, and delivering information between distant places, Most importantly he was the guide of travelers, the messenger of the gods, and the navigator between worlds.
CodeAtlas follows the same philosophy for codebases.
It does not attempt to replace understanding.
Instead, it guides developers and AI agents through large repositories, helping them navigate directly to the symbols, files, and implementations they need.
The repository already contains the answer.

CodeAtlas builds a rich, queryable symbol index of your entire codebase so LLMs can find the right code—fast. Instead of exploring millions of tokens, ask CodeAtlas and jump directly to what matters.


<p align="center">
  <img src="https://img.shields.io/badge/status-stable-brightgreen" />
  <img src="https://img.shields.io/badge/license-MIT-blue" />
  <img src="https://img.shields.io/badge/built%20with-Go-00ADD8" />
  <img src="https://img.shields.io/badge/parser-Tree--sitter-green" />
</p>

</div>

---
## 📚 Table of Contents

* [🛠️ Installation](#-installation)
* [✨ Key Benefits](#-key-benefits)
* [📖 What is CodeAtlas?](#-what-is-codeatlas)
* [⚙️ How It Works](#️-how-it-works)
  * [1️⃣ Index Generation](#1️⃣-index-generation)
  * [2️⃣ Query](#2️⃣-query)
  * [3️⃣ Navigate](#3️⃣-navigate)
  * [📦 .codeatlasignore](#-codeatlasignore)
* [🚀 Why CodeAtlas Matters](#-why-codeatlas-matters)
  * [Running Tests](#running-tests)
  * [Cross-Platform Builds](#cross-platform-builds)
* [📊 Benchmarks](#-benchmarks)
  * [Cross-Repository Summary](#cross-repository-summary)
  * [Detailed Reports](#detailed-reports)
  * [Retrieval Reduction](#retrieval-reduction)
* [🔎 First Lookup Accuracy](#-first-lookup-accuracy)
* [📦 When To Use CodeAtlas](#-when-to-use-codeatlas)
* [🔄 Index Once. Query Many.](#-index-once-query-many)
* [🏗 Architecture](#-architecture)
* [🧭 Design Principles](#-design-principles)
  * [Retrieval Over Stuffing](#retrieval-over-stuffing)
  * [Deterministic](#deterministic)
  * [Transparent](#transparent)
  * [Tool Friendly](#tool-friendly)
  * [Language Aware](#language-aware)
* [⚡ Quick Start](#-quick-start)
* [🛣 Roadmap](#-roadmap)
* [🤝 Contributing](#-contributing)
* [📄 License](#-license)

---

# Installation

## Prerequisites

CodeAtlas requires:

- Go 1.24+
- Universal Ctags https://github.com/universal-ctags/ctags?tab=GPL-2.0-1-ov-file (GPL 2.0)
- jq
- ripgrep / grep (optional)

### Install Universal Ctags

#### Ubuntu / Debian

```bash
sudo apt install universal-ctags
```

#### Arch Linux

```bash
sudo pacman -S universal-ctags
```

#### macOS

```bash
brew install universal-ctags
```

Verify installation:

```bash
ctags --version
```

Expected output should contain:

```text
Universal Ctags
```
#### Linux / macOS

Install the latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/Aeres-u99/codeatlas/master/install.sh | bash
```

If `~/.local/bin` is not on your `PATH`, the installer will tell you how to add it.

Verify the installation:

```bash
codeatlas --help
```

```bash
codeatlas --help
```

#### Manual Installation

Download the appropriate binary from the GitHub Releases page.

Rename it to `codeatlas`, make it executable, and place it somewhere on your `PATH`:

```bash
chmod +x codeatlas-linux-amd64
mv codeatlas-linux-amd64 codeatlas
sudo mv codeatlas /usr/local/bin/
```

Verify:

```bash
codeatlas --help
```

## Upgrading

Run the installer again or Download manually and install:

```bash
curl -fsSL https://raw.githubusercontent.com/Aeres-u99/codeatlas/master/install.sh | bash
```

---

## Build CodeAtlas

```bash
git clone https://github.com/Aeres-u99/codeatlas.git
cd codeatlas
make
```

---

## Running Tests
> This is pending, but we do plan to add tests.
> Sadly this wasn't made using TDD
```bash
make test
```

---

## Cross-Platform Builds

CodeAtlas uses CGo (via Tree-sitter), so cross-compilation requires a C cross-compiler for the target architecture — `GOARCH` alone is not enough.

### Linux arm64 (from Linux amd64)

Install the cross-compiler:

#### Ubuntu / Debian

```bash
sudo apt install gcc-aarch64-linux-gnu
```

#### Arch Linux

```bash
sudo pacman -S aarch64-linux-gnu-gcc
```

Build:

```bash
make build-linux-arm64
```

Output: `dist/codeatlas-linux-arm64`

Transfer the binary to your ARM machine and install:

```bash
scp dist/codeatlas-linux-arm64 user@arm-host:~
ssh user@arm-host 'sudo install ~/codeatlas-linux-arm64 /usr/local/bin/codeatlas'
```

---

### Linux amd64

```bash
make build-linux-amd64
```

Output: `dist/codeatlas-linux-amd64`

---

### macOS (run natively on the target Mac)

Darwin cross-compilation from Linux requires the macOS SDK and is not supported here. Run these on the Mac itself.

#### Apple Silicon (arm64)

```bash
make build-darwin-arm64
```

Output: `dist/codeatlas-darwin-arm64`

#### Intel (amd64)

```bash
make build-darwin-amd64
```

Output: `dist/codeatlas-darwin-amd64`

---

### Build all Linux targets at once

```bash
make build-all
```

Produces both `dist/codeatlas-linux-amd64` and `dist/codeatlas-linux-arm64`.

---

## Install CodeAtlas

```bash
sudo install codeatlas /usr/local/bin/codeatlas
```

Verify installation:

```bash
codeatlas --help
```

---

# ✨ Key Benefits

| 🎯 Direct Navigation                                                    | ⚡ High Efficiency                                                | 💰 Lower Cost                                    | 🔍 Works at Scale                           |
| ----------------------------------------------------------------------- | ---------------------------------------------------------------- | ------------------------------------------------ | ------------------------------------------- |
| Jump directly to implementations without exploring the repository tree. | Reduce repository exploration through targeted symbol retrieval. | Fewer searches, fewer tool calls, lower latency. | Designed for medium and large repositories. |

---

# 📖 What is CodeAtlas?

CodeAtlas is a repository indexing and retrieval system that produces a structured symbol index (`codeatlas.json`) of your codebase.

You can then query this index using simple tools such as:

* grep
* jq
* ripgrep
* Claude Code
* Aider
* Custom agents

Instead of feeding your entire repository to an LLM, CodeAtlas retrieves only the relevant functions, types, files, and locations.

```text
Repository
    │
    ▼
 CodeAtlas Index
    │
    ▼
 grep / jq / rg
    │
    ▼
 Relevant Symbols
    │
    ▼
      LLM
```

---

# ⚙️ How It Works

### 1️⃣ Index Generation

CodeAtlas parses the repository using Tree-sitter and extracts:

* Functions
* Methods
* Structs
* Interfaces
* Constants
* Variables
* Imports
* File Locations

```bash
codeatlas -input .
```

Output:

```text
codeatlas.json
```

---

### 2️⃣ Query

Search for symbols using standard UNIX tools.

```bash
grep -C3 "ConfigureProvider" codeatlas.json
```

or

```bash
jq '.idx["ConfigureProvider"]' codeatlas.json
```

---

### 3️⃣ Navigate

Open the exact files and locations returned by CodeAtlas.

No blind repository exploration required.

---

#### 📦 .codeatlasignore

CodeAtlas supports a `.codeatlasignore` file to exclude files and directories that provide little or no value for code navigation.

Ignoring generated code, vendored dependencies, documentation, test fixtures, repository metadata, and build artifacts can significantly reduce index size, indexing time, and query noise.

A well-tuned `.codeatlasignore` helps CodeAtlas focus on implementation code rather than auxiliary repository content.

### Example

```gitignore
# Repository metadata
.git/**
.github/**

# Dependencies
vendor/
node_modules/

# Documentation
docs/
*.md

# Generated files
**/zz_generated.*
**/*_generated.go

# Test fixtures
**/testdata/
**/*_test.go

# Build artifacts
dist/
build/
tmp/

# IDE files
.idea/
.vscode/
```

### Why Use `.codeatlasignore`?

- ⚡ Faster indexing
- 📉 Smaller index size
- 🎯 Better retrieval precision
- 💰 Lower token consumption
- 🔍 Reduced search noise

CodeAtlas is designed around the principle that not all repository content is equally valuable for navigation. A carefully curated `.codeatlasignore` allows the index to focus on the code that matters most.

> **Tip:** For large repositories such as Kubernetes, Loki, or Terraform, excluding `.git/**`, vendored dependencies, generated files, and documentation can dramatically reduce index size while preserving the vast majority of implementation-relevant symbols.

---

# 🚀 Why CodeAtlas Matters

Modern LLM workflows spend the majority of their time and tokens exploring repositories.

Typical workflow:

```text
Search
 ↓
Open File
 ↓
Wrong File
 ↓
Search Again
 ↓
Open More Files
 ↓
Find Implementation
```

CodeAtlas transforms that into:

```text
Lookup Symbol
 ↓
Open Correct File
 ↓
Done
```

Benefits:

* Fewer file reads
* Lower token usage
* Fewer tool calls
* Faster navigation
* Lower latency
* Better agent performance

## Generate Your First Index

```bash
codeatlas -input .
```

This generates:

```text
codeatlas.json
```

which can then be queried using:

```bash
jq -r '.idx | keys[]' codeatlas.json
```

or

```bash
jq -r '.idx["main.main"]' codeatlas.json
```

---
# 📊 Benchmarks

CodeAtlas has been benchmarked against three real-world repositories of different scales.

| Repository | Files  | LOC  | Runtime (s) | Full Map Tokens | Query Tokens | Retrieval Reduction |
| ---------- | ------ | ---- | ----------- | --------------- | ------------ | ------------------- |
| Kubernetes | 30,536 | 5.0M | 66.2        | 14.7M           | 2,059        | 99.986%             |
| Loki       | 17,160 | 510k | 11.6        | 8.4M            | 880          | 99.990%             |
| Terraform  | 5,411  | 667k | 13.8        | 2.0M            | 5,340        | 99.732%             |

---

## Cross-Repository Summary

| Repository | Exploration Collapse | Cost Reduction |
| ---------- | -------------------- | -------------- |
| Kubernetes | 60.0%                | 49.6%          |
| Loki       | 50.0%                | 50.2%          |
| Terraform  | 62.5%                | 7.2%           |

---

## Detailed Reports

* [Kubernetes Benchmark](benchmarks/Kubernetes-Report.md)
* [Loki Benchmark](benchmarks/Loki-Report.md)
* [Terraform Benchmark](benchmarks/Terraform-Report.md)

---

## Retrieval Reduction

CodeAtlas is a retrieval system, not a context-stuffing system.

| Repository | Full Map Tokens | Query Tokens |
| ---------- | --------------- | ------------ |
| Kubernetes | 14,718,000      | 2,059        |
| Loki       | 8,396,000       | 880          |
| Terraform  | 1,990,000       | 5,340        |

CodeAtlas consistently reduces retrieval payloads by over **99.7%**.

The entire repository is never sent to the model.

Only the relevant symbols are retrieved.

---

# 🔎 First Lookup Accuracy

| Repository | First Lookup Correct  |
| ---------- | --------------------- |
| Kubernetes | ✅ Yes                 |
| Loki       | ✅ Yes                 |
| Terraform  | ❌ Required Refinement |

CodeAtlas performs best when queried using specific symbols.

Broad interfaces and common method names may require one or more refinement steps.

---

# 📦 When To Use CodeAtlas

CodeAtlas is optimized for medium and large repositories.

| Repository Size      | Expected Benefit | Recommendation            |
| -------------------- | ---------------- | ------------------------- |
| < 100 files          | Low              | Generally not recommended |
| 100 – 1,000 files    | Medium           | Situational benefit       |
| 1,000 – 10,000 files | High             | Strong benefit            |
| 10,000+ files        | Very High        | Primary target            |

Smaller repositories may not see enough benefit to offset indexing overhead.

---

# 🔄 Index Once. Query Many.

Index generation is a one-time operation.

The resulting index can be reused across:

* Feature development
* Bug investigation
* Architecture discovery
* Code reviews
* AI-assisted programming
* Documentation generation

```text
Repository
    │
    ▼
 Generate Index
    │
    ▼
  codeatlas.json
    │
    ├── Query #1
    ├── Query #2
    ├── Query #3
    ├── Query #4
    └── Query #N
```

The generation cost is paid once. Or, if post Query #1 the syntax or structure has changed significantly, regenerate the codeatlas.json to ensure that the map is not stale.

The retrieval benefit compounds over time.

---

# 🏗 Architecture

```text
Repository
    │
    ▼
Tree-sitter Parser
    │
    ▼
CodeAtlas Index
    │
    ├── Functions
    ├── Methods
    ├── Structs
    ├── Interfaces
    ├── Variables
    ├── Constants
    └── Imports
    │
    ▼
grep / jq / rg
    │
    ▼
LLM / Agent
```

---

# 🧭 Design Principles

### Retrieval Over Stuffing

Never send the entire repository to the model.

Retrieve only the relevant implementation details.

### Deterministic

Same repository → same index.

No embeddings.

No probabilistic ranking.

### Transparent

Indexes are human-readable and versionable.

You can inspect, diff, and audit every result.

### Tool Friendly

Works with standard UNIX tooling:

* grep
* jq
* rg
* awk
* sed

### Language Aware

Built on Tree-sitter.

Designed to understand source structure rather than plain text.

---

# ⚡ Quick Start

CodeAtlas Help
```
❯ codeatlas --help
codeatlas - The Code Map you will Ever need!
For custom CTAGS path use Environment Variable CODEATLAS_CTAGS for custom path
  -input string
    	Code to Parse (default "code.py")
  -output string
    	Code to Parse (default "codeatlas.json")
Have Fun

```

Generate an index:

```bash
codeatlas -input .
```

Search for a symbol:

```bash
grep -C3 "MyFunction" codeatlas.json
```

Retrieve with jq:

```bash
jq '.idx["MyFunction"]' codeatlas.json
```

Open the returned file and line.

Done.

---

# 🛣 Roadmap

* [ ] [TODO](TODO.md)
* [ ] Incremental indexing
* [ ] Multi-language support
* [ ] Cross-reference graph

---

# 🤝 Contributing

Contributions, bug reports, and feature requests are welcome.

Please open an issue or submit a pull request.

---

# 📄 License

MIT
