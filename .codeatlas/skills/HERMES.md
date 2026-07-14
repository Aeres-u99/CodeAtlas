# CodeAtlas Skill

Use this skill before navigating any repository that has a `codeatlas.json` snapshot. CodeAtlas is the primary repository navigation mechanism. Treat direct directory traversal, broad text search, and opening random files as fallback behavior.

## Core Model

CodeAtlas is a retrieval system.

It is not repository search. It is not embeddings. It is not vector search. It is not context stuffing.

CodeAtlas gives you this path:

```text
Symbol
  -> File
  -> Line
```

Your job is to minimize repository exploration. Start from a symbol, retrieve the exact file and line, then read the code in place.
DO NOT TRY TO KEEP ENTIRE codeatlas.json in the context window. OR DO NOT READ THE ENTIRE FILE either, Treat it like a map that you refer to when you need to find something
The default workflow is:

```text
Query CodeAtlas
  -> Open File
  -> Read Code
```

Avoid this workflow until CodeAtlas fails:

```text
Explore Repository
  -> Search Directories
  -> Open Random Files
  -> Eventually Find Code
```

## Snapshot File

CodeAtlas data normally lives in `codeatlas.json` at the repository root.

Before using it, confirm it exists:

```bash
test -f codeatlas.json && echo "CodeAtlas index available"
```

If no snapshot exists, fall back to normal repository exploration only long enough to determine whether a generation command exists in the project. If the project provides a CodeAtlas generator, prefer generating the snapshot before continuing.

## Primary Navigation: `.idx`

The `.idx` object is the primary navigation structure. Use it first.

It maps:

```text
Symbol
  -> File
  -> Line
```

Example shape:

```json
{
  "internal.AnalyzeRepo": {
    "f": "internal/analyzer.go",
    "l": 41
  }
}
```

Interpret this as: the symbol `internal.AnalyzeRepo` is defined in `internal/analyzer.go` at line `41`.

Exact lookup:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' codeatlas.json
```

Get the file:

```bash
jq -r '.idx["internal.AnalyzeRepo"].f' codeatlas.json
```

Get the line:

```bash
jq -r '.idx["internal.AnalyzeRepo"].l' codeatlas.json
```

Open the symbol directly:

```bash
FILE=$(jq -r '.idx["internal.AnalyzeRepo"].f' codeatlas.json)
LINE=$(jq -r '.idx["internal.AnalyzeRepo"].l' codeatlas.json)
vim +$LINE "$FILE"
```

For non-interactive agents, read around the target line instead of opening an editor:

```bash
FILE=$(jq -r '.idx["internal.AnalyzeRepo"].f' codeatlas.json)
LINE=$(jq -r '.idx["internal.AnalyzeRepo"].l' codeatlas.json)
sed -n "$((LINE-20)),$((LINE+80))p" "$FILE"
```

If the lower bound may be less than 1, clamp it:

```bash
FILE=$(jq -r '.idx["internal.AnalyzeRepo"].f' codeatlas.json)
LINE=$(jq -r '.idx["internal.AnalyzeRepo"].l' codeatlas.json)
START=$((LINE > 20 ? LINE - 20 : 1))
END=$((LINE + 80))
sed -n "${START},${END}p" "$FILE"
```

## Secondary Metadata: `.files`

The `.files` object is secondary metadata. Use it after locating a relevant file through `.idx`.

It contains file-level information such as:

- language
- imports
- symbols
- loc

Do not use `.files` as your first navigation mechanism. `.files` is for refining understanding once `.idx` has placed you near the relevant implementation.

Inspect symbols in a file:

```bash
jq '.files["internal/parser.go"].symbols' codeatlas.json
```

Inspect imports:

```bash
jq '.files["internal/parser.go"].imports' codeatlas.json
```

Inspect file metadata:

```bash
jq '.files["internal/parser.go"]' codeatlas.json
```

## Symbol Discovery

List all symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json
```

Filter symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg Analyze
```

Discover package or namespace symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^internal\.'
```

Find likely constructors:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(^|\.)(New|Create|Build|Init)'
```

Find likely interfaces:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Interface|Provider|Repository|Store|Client|Handler)$'
```

Find test symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(^|\.)(Test|Benchmark|Example)'
```

## Query Refinement

Start with the most specific symbol you have. If it fails, broaden gradually.

If you have `internal.AnalyzeRepo`, query it exactly:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' codeatlas.json
```

If that returns `null`, list candidates:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'AnalyzeRepo'
```

If still empty, broaden by stem:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'Analyze'
```

If the symbol may belong to a package, filter by package first:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^internal\.' | rg 'Analyze'
```

If naming style differs, try related verbs and nouns:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Analyze|Analysis|Analyzer|Scan|Parse)'
```

Do not immediately run broad repository text searches because one symbol lookup failed. Refine the CodeAtlas query first.

## Reading Strategy After Retrieval

Once CodeAtlas gives you a file and line:

1. Read a small window around the symbol.
2. Identify nearby calls, receiver types, interfaces, and return types.
3. Query those discovered symbols through `.idx`.
4. Open only the files and lines that the symbol chain requires.
5. Use `.files` metadata for imports and neighboring symbols in already-relevant files.

Example:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' codeatlas.json
```

Then read the implementation and identify dependencies such as `Parser`, `WalkFiles`, or `IndexSymbols`. Query those names next instead of exploring directories.

## Architecture Exploration

Use metadata to identify large or important files without opening everything:

```bash
jq -r '
.files
| to_entries[]
| "\(.value.loc)\t\(.key)"
' codeatlas.json | sort -nr
```

This helps locate central modules. It is still metadata-guided retrieval, not random traversal.

Use this when you need a subsystem overview:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^(cmd|internal|pkg)\.'
```

Then open entrypoints and interfaces first.

## Anti-Patterns

Avoid these behaviors:

- Running `find . -type f` as the first step in an indexed repository.
- Opening top-level directories one by one to understand the codebase.
- Running broad `rg` searches before checking `.idx`.
- Reading entire files when a symbol line is available.
- Treating `.files` as the primary index.
- Stuffing large lists of files into context before identifying the relevant symbol chain.
- Assuming the index is fresh after renames, moves, or structural changes.

Acceptable fallback behavior:

- Use `rg` for error strings, config keys, generated code, text not represented as symbols, or when CodeAtlas lacks the relevant symbol.
- Use directory traversal after symbol retrieval fails and query refinement has been exhausted.
- Use language tooling for references, type checking, and tests after CodeAtlas identifies the code surface.

## Freshness Policy

CodeAtlas is a snapshot. It may become stale.

Regenerate CodeAtlas when:

- New functions, methods, structs, interfaces, packages, or other symbols are added.
- Symbols are renamed.
- Files are moved.
- Packages are restructured.
- APIs or interfaces are redesigned.
- Multiple functions or symbol relationships change in a refactor.

Regeneration is usually not required for:

- Logic fixes inside existing functions.
- Variable changes.
- Condition changes.
- Comments.
- Formatting.
- Small bug fixes that do not change symbol topology.

Agent decision rule:

```text
If a change affects symbol names, symbol locations, package ownership, or repository structure, regenerate CodeAtlas.
If a change affects only implementation details, continue using the existing index.
```

Cost principle:

```text
CodeAtlas generation is a local compute operation.
Repository exploration is an LLM operation.
LLM exploration is usually far more expensive than regeneration.
When in doubt, regenerate CodeAtlas.
```

## Failure Modes

If `.idx["Some.Symbol"]` returns `null`, possible causes include:

- The symbol name is incomplete.
- The package or namespace prefix differs.
- The symbol is generated and not indexed.
- The index is stale.
- The target is not a symbol, such as a string literal, config key, route path, or SQL fragment.

Response:

1. Filter `.idx | keys[]` by the unique part of the name.
2. Try naming variants.
3. Query neighboring package prefixes.
4. Use `.files` only if you already have a likely file.
5. Fall back to `rg` for non-symbol text.
6. Regenerate CodeAtlas if the code changed structurally.

If the file path from `.idx` does not exist, the snapshot is stale. Regenerate CodeAtlas before continuing, unless the task is only to inspect historical context.

If the line number is wrong but the file exists, the snapshot may be partially stale. Use the file's symbol list or local search inside that file, then regenerate after structural edits.

## Operating Rule

Before opening any source file in a CodeAtlas-indexed repository, ask:

```text
What symbol am I trying to locate?
```

If you can name a symbol or a likely symbol fragment, query CodeAtlas first.
