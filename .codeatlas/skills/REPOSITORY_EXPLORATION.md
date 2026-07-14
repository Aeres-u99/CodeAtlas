# Repository Exploration With CodeAtlas

Use this skill when entering an unfamiliar repository. Your goal is to build a useful mental model while opening the fewest files possible.

## Objective

Understand the repository through symbol retrieval, not directory wandering.

The default path is:

```text
Identify likely entrypoints
  -> Query CodeAtlas
  -> Open implementations
  -> Expand through discovered symbols
  -> Build mental model
```

Delay repository traversal as long as possible.

## First Checks

Confirm the CodeAtlas snapshot exists:

```bash
test -f codeatlas.json && echo "CodeAtlas index available"
```

Inspect the top-level symbol space:

```bash
jq -r '.idx | keys[]' codeatlas.json | head -100
```

This is not a request to read 100 files. It is a quick look at naming structure: packages, namespaces, commands, handlers, services, and tests.

## Entrypoint Discovery

Start with symbols that commonly represent application entrypoints:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(^|\.)(main|Main|Run|Execute|Start|Serve|Listen|Handler|Command)$'
```

Find command packages:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^(cmd|cli|commands)\.'
```

Find server or worker entrypoints:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Server|Worker|Daemon|Scheduler|Consumer|Processor|Handler)'
```

Open only the strongest candidates. For each candidate:

```bash
SYMBOL='internal.Server.Run'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' codeatlas.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' codeatlas.json)
START=$((LINE > 30 ? LINE - 30 : 1))
END=$((LINE + 120))
sed -n "${START},${END}p" "$FILE"
```

## Build The First Map

After opening an entrypoint, extract the next symbols from code:

- constructors
- interfaces
- concrete service types
- handlers
- repositories
- clients
- parser or analyzer functions
- package-level orchestration functions

For each discovered symbol, return to CodeAtlas:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'NewAnalyzer|Analyzer|Repository'
```

Do not follow imports blindly. Imports tell you dependencies, but `.idx` tells you where code lives.

## Use `.files` Only After `.idx`

Once an entrypoint file is known, inspect metadata:

```bash
jq '.files["cmd/codeatlas/main.go"]' codeatlas.json
```

Inspect symbols in that file:

```bash
jq '.files["cmd/codeatlas/main.go"].symbols' codeatlas.json
```

Inspect imports:

```bash
jq '.files["cmd/codeatlas/main.go"].imports' codeatlas.json
```

Use imports to identify subsystem boundaries, then query subsystem symbols through `.idx`.

## Package-Level Survey

If the repository uses package-style symbol names, survey packages through symbol prefixes:

```bash
jq -r '.idx | keys[]' codeatlas.json | sed 's/\.[^.]*$//' | sort -u
```

Then inspect symbols in a package:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^internal\.analyzer\.'
```

If symbols are flatter, filter by terms:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Analyze|Parse|Index|Config|Command|Server)'
```

## File Importance Without Reading Files

Use `loc` to find large files that may represent central modules:

```bash
jq -r '
.files
| to_entries[]
| "\(.value.loc)\t\(.key)"
' codeatlas.json | sort -nr | head -30
```

This is useful after entrypoint retrieval, not before it. Large files are not automatically important; they are candidates to correlate with symbols already discovered.

## Mental Model Template

Build a concise model as you retrieve:

```text
Entrypoint:
- symbol -> file:line

Primary flow:
- symbol -> symbol -> symbol

Core data types:
- type/interface -> file:line

External boundaries:
- HTTP/CLI/DB/filesystem/API clients

Tests:
- relevant test symbols -> file:line
```

Keep the map symbol-based. Avoid listing every directory.

## Refinement Strategy

If entrypoint discovery is noisy:

1. Filter by likely top-level packages such as `cmd`, `internal`, `pkg`, `app`, `server`, or `cli`.
2. Prefer `Run`, `Execute`, `Serve`, `Start`, and constructors.
3. Prefer symbols referenced by multiple entrypoints.
4. Inspect file metadata for imports only after a candidate file is selected.
5. Use tests to confirm behavior after the main flow is known.

Example:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^cmd\.' | rg '(Run|Execute|main|Command)'
```

## When To Fall Back

Fallback to repository exploration only when:

- `codeatlas.json` is missing.
- `.idx` does not contain enough useful symbols after refinement.
- You need non-symbol assets such as schemas, templates, migrations, config files, or documentation.
- File paths in `.idx` are stale.
- The task asks about repository layout itself rather than code behavior.

Even then, keep exploration targeted:

```bash
rg --files | rg '(schema|migration|config|template|README|Makefile)'
```

## Exploration Anti-Patterns

Avoid:

- Reading `README` first when the task is implementation-oriented and CodeAtlas exists.
- Listing every directory to infer architecture.
- Opening all files in a package before querying symbols.
- Using import graphs as a substitute for symbol lookup.
- Treating a large file list as understanding.

Good exploration produces a short symbol chain and a small set of opened code windows.

## Freshness During Exploration

If the repository has changed since `codeatlas.json` was generated, exploration can become misleading.

Regenerate CodeAtlas when you observe:

- indexed file paths that no longer exist
- symbols whose definitions moved
- package names that differ from current code
- new entrypoints not present in `.idx`

When unsure, prefer regeneration. Local index generation is cheaper than LLM-driven wandering through stale structure.
