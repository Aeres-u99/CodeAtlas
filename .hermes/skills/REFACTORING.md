# Refactoring With Hermes

Use this skill when performing structural changes, cleanup, extraction, renaming, moving code, or changing contracts.

## Objective

Make safe changes by reasoning from symbols and impact scope before editing.

The default path is:

```text
Locate symbol
  -> Discover related symbols
  -> Identify impact scope
  -> Inspect neighboring implementations
  -> Apply changes
```

Refactoring should be symbol-driven. Do not start by opening directories and manually scanning files.

## Locate The Primary Symbol

Start with the symbol being refactored:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' hermes.json
```

Open the implementation:

```bash
SYMBOL='internal.AnalyzeRepo'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' hermes.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' hermes.json)
START=$((LINE > 40 ? LINE - 40 : 1))
END=$((LINE + 160))
sed -n "${START},${END}p" "$FILE"
```

If the exact symbol is missing:

```bash
jq -r '.idx | keys[]' hermes.json | rg 'AnalyzeRepo|Analyze|Analyzer'
```

Do not edit until you have the correct symbol definition.

## Discover Related Symbols

Find the naming family:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Analyze|Analyzer|Analysis)'
```

Inspect symbols in the same file:

```bash
jq '.files["internal/analyzer.go"].symbols' hermes.json
```

Inspect imports:

```bash
jq '.files["internal/analyzer.go"].imports' hermes.json
```

Look for tests:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Test.*Analyze|Analyze.*Test|Benchmark.*Analyze)'
```

This gives an initial impact surface before broad searching.

## Identify Impact Scope

Hermes maps definitions, not all references. Use Hermes to locate definitions, then use language tooling or targeted text search for references.

Definition scope:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(AnalyzeRepo|Analyzer)'
```

Reference scope:

```bash
rg -n '\bAnalyzeRepo\b'
```

For method names that may collide, include receiver or type context when possible:

```bash
rg -n '\.AnalyzeRepo\(|func .*AnalyzeRepo'
```

Use the language server, compiler, or test runner for final validation. Hermes is the navigation layer, not a type checker.

## Neighboring Implementations

Before extracting or renaming, inspect similar code:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Parse|Parser|Index|Indexer|Walk|Walker)'
```

Open only the relevant neighboring symbols. You are looking for project conventions:

- constructor shape
- interface placement
- error handling
- naming style
- test layout
- package ownership
- dependency direction

Refactors should preserve local style unless the task explicitly changes it.

## Rename Workflow

For a symbol rename:

1. Locate the current symbol through `.idx`.
2. Find related symbols and tests.
3. Use targeted reference search or language tooling to update references.
4. Run tests or type checks.
5. Regenerate Hermes because symbol names changed.
6. Use the regenerated index for any follow-up navigation.

Example:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' hermes.json
rg -n '\bAnalyzeRepo\b'
```

After renaming to `AnalyzeRepository`, the old lookup must not be trusted. Regenerate.

## Move Workflow

For moving files or packages:

1. Locate all symbols in the source file:

```bash
jq '.files["internal/analyzer.go"].symbols' hermes.json
```

2. Inspect imports of the source file:

```bash
jq '.files["internal/analyzer.go"].imports' hermes.json
```

3. Find references to moved symbols:

```bash
rg -n '\bAnalyzeRepo\b|\bAnalyzer\b'
```

4. Move code using project conventions.
5. Update imports and package ownership.
6. Run tests or type checks.
7. Regenerate Hermes because file paths and symbol locations changed.

File moves always make the old snapshot stale.

## Extraction Workflow

For extracting helper functions or types:

1. Locate the source symbol.
2. Read only the relevant implementation window.
3. Identify the coherent block to extract.
4. Check neighboring helpers in the same file:

```bash
jq '.files["internal/analyzer.go"].symbols' hermes.json
```

5. Choose a name consistent with existing symbols:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(extract|build|parse|collect|write)'
```

6. Add the new helper.
7. Regenerate Hermes if the new helper is a symbol that future navigation should see.

## Interface Refactor Workflow

For interface or contract changes:

1. Query the interface symbol:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Store|Repository|Parser|Analyzer|Client)$'
```

2. Open the interface.
3. Find implementations by naming family.
4. Find constructors and wiring code.
5. Find tests for the contract and implementations.
6. Apply changes across the full implementation set.
7. Run type checks and tests.
8. Regenerate Hermes because contract topology changed.

Do not change an interface after inspecting only one implementation.

## Safety Checklist

Before editing:

- Primary symbol located through `.idx`.
- Related symbols listed.
- Tests or validation commands identified.
- Reference search strategy chosen.
- Regeneration need decided.

After editing:

- Tests or type checks run where feasible.
- Hermes regenerated if symbol topology changed.
- Final navigation uses the regenerated snapshot if further code reading is needed.

## Freshness Rules

Regenerate Hermes after refactors that affect:

- symbol names
- symbol locations
- file paths
- package names
- interface definitions
- constructor or factory topology
- new extracted functions or types
- deleted symbols

Regeneration is usually unnecessary for:

- local condition cleanup
- variable rename inside a function
- comment edits
- formatting
- implementation-only simplification that preserves symbols

When in doubt, regenerate. The cost of local generation is usually lower than debugging stale navigation.

## Failure Modes

If `.idx` finds the symbol but reference search finds no uses, the symbol may be an entrypoint, dynamically referenced, used externally, or dead. Do not delete it without checking tests, exports, public API, and build configuration.

If many unrelated references match, narrow with package prefixes, receiver names, or call syntax.

If file metadata is missing for a path from `.idx`, the index is inconsistent. Regenerate before refactoring.

If tests fail after a rename, query the old and new names in Hermes and with `rg` to find missed definitions or references.

## Refactoring Principle

Hermes helps you keep the refactor small. Use it to move from one exact symbol to the next exact symbol. Large context dumps are a sign that the refactor scope has not been defined precisely enough.
