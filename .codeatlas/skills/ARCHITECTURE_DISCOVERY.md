# Architecture Discovery With CodeAtlas

Use this skill when you need to understand subsystem design, boundaries, data flow, or ownership.

## Objective

Build an architecture map through symbol and metadata retrieval.

The default path is:

```text
Locate entrypoints
  -> Query interfaces
  -> Query implementations
  -> Inspect imports
  -> Build subsystem map
```

Use CodeAtlas metadata before opening files broadly.

## Start With Entrypoints

Find symbols that initiate control flow:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(^|\.)(main|Run|Execute|Start|Serve|Listen|Process|Analyze)$'
```

Find command or application packages:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^(cmd|app|cli|internal)\.'
```

Open the likely entrypoints and read enough to identify orchestration symbols:

```bash
SYMBOL='cmd.Run'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' codeatlas.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' codeatlas.json)
START=$((LINE > 40 ? LINE - 40 : 1))
END=$((LINE + 160))
sed -n "${START},${END}p" "$FILE"
```

## Discover Interfaces

Architecture is often expressed in interfaces, abstract base types, protocols, providers, stores, clients, and handlers.

Query likely boundaries:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Interface|Provider|Store|Repository|Client|Service|Handler|Parser|Analyzer|Formatter|Backend|Driver)$'
```

Open interface-like symbols first. They define contracts and subsystem boundaries more clearly than leaf functions.

Then find implementations using related terms:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(FileStore|MemoryStore|HTTPClient|Parser|Analyzer|Formatter)'
```

## Use Imports After Locating Files

Once a relevant file is identified, inspect imports:

```bash
jq '.files["internal/analyzer.go"].imports' codeatlas.json
```

Imports help answer:

- Which packages are external boundaries?
- Which internal packages does this subsystem depend on?
- Does this package depend inward or outward?
- Is this file orchestration, infrastructure, or domain logic?

Do not use imports as the first step. Use `.idx` to find code, then `.files` to understand file-level relationships.

## Build A Subsystem Map

Create a compact map:

```text
Subsystem: analysis

Entrypoints:
- internal.AnalyzeRepo -> internal/analyzer.go:41

Contracts:
- internal.Parser -> internal/parser.go:12
- internal.Store -> internal/store.go:8

Implementations:
- internal.GoParser -> internal/go_parser.go:20
- internal.JSONStore -> internal/json_store.go:18

Data flow:
- AnalyzeRepo -> WalkFiles -> ParseFile -> IndexSymbols -> WriteIndex

External boundaries:
- filesystem via internal.WalkFiles
- JSON output via internal.WriteIndex
```

Keep the map at symbol granularity. Avoid describing every file in a directory.

## File Size And Centrality

Use LOC metadata to find central files:

```bash
jq -r '
.files
| to_entries[]
| "\(.value.loc)\t\(.key)"
' codeatlas.json | sort -nr | head -40
```

Use this to supplement symbol discovery. A large file may be central, but it may also be generated, legacy, or test-heavy. Correlate file size with symbols and imports before drawing conclusions.

## Package And Namespace Discovery

List symbol prefixes:

```bash
jq -r '.idx | keys[]' codeatlas.json | sed 's/\.[^.]*$//' | sort -u
```

Inspect symbols in a prefix:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^internal\.'
```

Compare packages by symbol count:

```bash
jq -r '.idx | keys[]' codeatlas.json \
| sed 's/\.[^.]*$//' \
| sort \
| uniq -c \
| sort -nr
```

Use this to identify dense subsystems, then retrieve entrypoints and interfaces from those subsystems.

## Architectural Questions

For "How does X work?", query X-related symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Analyze|Analysis|Analyzer)'
```

For "Where is X configured?", query configuration symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Config|Options|Settings|Flag|Env)'
```

For "Where does X enter the system?", query entrypoint symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Run|Execute|Serve|Handle|Process)'
```

For "What implements X?", query the interface name and common implementation suffixes:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Store|Repository|Provider|Client)'
```

## Avoid False Architecture

Do not infer architecture only from:

- folder names
- package names
- README claims
- import lists
- large files
- your memory of similar projects

Use actual symbol definitions and call paths. Open source code only at retrieved symbol locations.

## When To Use Text Search

Use `rg` when architecture depends on non-symbol artifacts:

- configuration files
- routes
- SQL migrations
- templates
- generated schemas
- environment variables
- build scripts

Example:

```bash
rg -n 'DATABASE_URL|/api/|CREATE TABLE|feature_flag'
```

After finding a relevant source file, inspect `.files[path].symbols` and continue through `.idx`.

## Freshness And Architecture

Architecture discovery is especially sensitive to stale indexes.

Regenerate CodeAtlas when:

- package ownership changed
- files were moved
- interfaces were redesigned
- many symbols were added or removed
- entrypoints changed

If `.idx` paths do not exist or line numbers are repeatedly wrong, stop relying on the snapshot until it is regenerated.

## Output Expectations

When reporting architecture, cite symbols and files:

```text
The analysis flow starts at internal.AnalyzeRepo in internal/analyzer.go:41.
It delegates file discovery to internal.WalkFiles and parsing to internal.ParseFile.
The index shape is written by internal.WriteIndex.
```

Avoid vague claims like "the analyzer package handles analysis" unless backed by retrieved symbols.
