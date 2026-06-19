<div align="center">

<img src="assets/logo.png" alt="Hermes Logo" width="220"/>

# Hermes <img src="assets/wing.png" width="32"/>

### Symbolic Code Index. Instant Navigation.

**A deterministic code indexing and retrieval system for AI-assisted development.**

Hermes builds a rich, queryable symbol index of your entire codebase so LLMs can find the right code—fast. Instead of exploring millions of tokens, ask Hermes and jump directly to what matters.


<p align="center">
  <img src="https://img.shields.io/badge/status-stable-brightgreen" />
  <img src="https://img.shields.io/badge/license-MIT-blue" />
  <img src="https://img.shields.io/badge/built%20with-Go-00ADD8" />
  <img src="https://img.shields.io/badge/parser-Tree--sitter-green" />
</p>

</div>

---

# ✨ Key Benefits

| 🎯 Direct Navigation                                                    | ⚡ High Efficiency                                                | 💰 Lower Cost                                    | 🔍 Works at Scale                           |
| ----------------------------------------------------------------------- | ---------------------------------------------------------------- | ------------------------------------------------ | ------------------------------------------- |
| Jump directly to implementations without exploring the repository tree. | Reduce repository exploration through targeted symbol retrieval. | Fewer searches, fewer tool calls, lower latency. | Designed for medium and large repositories. |

---

# 📖 What is Hermes?

Hermes is a repository indexing and retrieval system that produces a structured symbol index (`hermes.json`) of your codebase.

You can then query this index using simple tools such as:

* grep
* jq
* ripgrep
* Claude Code
* Aider
* Custom agents

Instead of feeding your entire repository to an LLM, Hermes retrieves only the relevant functions, types, files, and locations.

```text
Repository
    │
    ▼
 Hermes Index
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

Hermes parses the repository using Tree-sitter and extracts:

* Functions
* Methods
* Structs
* Interfaces
* Constants
* Variables
* Imports
* File Locations

```bash
hermes -input . > hermes.json
```

Output:

```text
hermes.json
```

---

### 2️⃣ Query

Search for symbols using standard UNIX tools.

```bash
grep -C3 "ConfigureProvider" hermes.json
```

or

```bash
jq '.symbols[] | select(.name=="ConfigureProvider")' hermes.json
```

---

### 3️⃣ Navigate

Open the exact files and locations returned by Hermes.

No blind repository exploration required.

---

# 🚀 Why Hermes Matters

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

Hermes transforms that into:

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

---

# 📊 Benchmarks

Hermes has been benchmarked against three real-world repositories of different scales.

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

Hermes is a retrieval system, not a context-stuffing system.

| Repository | Full Map Tokens | Query Tokens |
| ---------- | --------------- | ------------ |
| Kubernetes | 14,718,000      | 2,059        |
| Loki       | 8,396,000       | 880          |
| Terraform  | 1,990,000       | 5,340        |

Hermes consistently reduces retrieval payloads by over **99.7%**.

The entire repository is never sent to the model.

Only the relevant symbols are retrieved.

---

# 🔎 First Lookup Accuracy

| Repository | First Lookup Correct  |
| ---------- | --------------------- |
| Kubernetes | ✅ Yes                 |
| Loki       | ✅ Yes                 |
| Terraform  | ❌ Required Refinement |

Hermes performs best when queried using specific symbols.

Broad interfaces and common method names may require one or more refinement steps.

---

# 📦 When To Use Hermes

Hermes is optimized for medium and large repositories.

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
  hermes.json
    │
    ├── Query #1
    ├── Query #2
    ├── Query #3
    ├── Query #4
    └── Query #N
```

The generation cost is paid once. Or, if post Query #1 the syntax or structure has changed significantly, regenerate the hermes.json to ensure that the map is not stale.

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
Hermes Index
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

Generate an index:

```bash
hermes -input . > hermes.json
```

Search for a symbol:

```bash
grep -C3 "MyFunction" hermes.json
```

Retrieve with jq:

```bash
jq '.symbols[] | select(.name=="MyFunction")' hermes.json
```

Open the returned file and line.

Done.

---

# 🛣 Roadmap

* [ ] Incremental indexing
* [ ] Multi-language support
* [ ] Cross-reference graph
  
  

---

# 🤝 Contributing

Contributions, bug reports, and feature requests are welcome.

Please open an issue or submit a pull request.

---

# 📄 License

TBU
