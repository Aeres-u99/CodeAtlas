# Bug Investigation With CodeAtlas

Use this skill when debugging failures, errors, regressions, crashes, test failures, or incorrect behavior.

## Objective

Find the root cause quickly by turning failure evidence into symbol lookups.

The default path is:

```text
Extract symbols from failure evidence
  -> Query CodeAtlas
  -> Open implementation
  -> Trace execution path
  -> Identify failure point
```

Do not start by browsing directories. Bugs usually expose names: functions, methods, types, interfaces, files, packages, stack frames, test names, or error-producing code.

## Evidence Priority

Extract lookup terms in this order:

1. Function and method names from stack traces.
2. Test names from failing test output.
3. Type, struct, interface, or class names.
4. Handler, command, or job names.
5. Package or namespace names.
6. Error strings and log messages.
7. File paths and line numbers.

Symbols are preferred because CodeAtlas indexes symbols directly.

## Stack Trace Workflow

Given a stack trace frame such as:

```text
internal.AnalyzeRepo(...)
    internal/analyzer.go:41
```

Query exactly:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' codeatlas.json
```

Open around the symbol:

```bash
SYMBOL='internal.AnalyzeRepo'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' codeatlas.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' codeatlas.json)
START=$((LINE > 40 ? LINE - 40 : 1))
END=$((LINE + 140))
sed -n "${START},${END}p" "$FILE"
```

If the exact symbol is missing:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'AnalyzeRepo'
```

Then broaden:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Analyze|Analyzer)'
```

## Failing Test Workflow

If a test fails, look up the test symbol first:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'TestAnalyzeRepo'
```

Open the test:

```bash
SYMBOL='internal.TestAnalyzeRepo'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' codeatlas.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' codeatlas.json)
START=$((LINE > 20 ? LINE - 20 : 1))
END=$((LINE + 120))
sed -n "${START},${END}p" "$FILE"
```

From the test body, extract the production symbols under test and query them through `.idx`.

Do not read every test file in the package. The failing test usually names the behavior and the production entrypoint.

## Error String Workflow

Error strings may not be indexed as symbols. Still check symbol context first if the error includes a function, type, or package name.

If only an error string is available, use targeted text search:

```bash
rg -n 'failed to analyze repository'
```

Once the file is found, use CodeAtlas metadata to inspect nearby symbols:

```bash
jq '.files["internal/analyzer.go"].symbols' codeatlas.json
```

Then switch back to `.idx` for any discovered symbols.

## Trace Execution Symbol By Symbol

After opening the failing symbol:

1. Identify the immediate call that can produce the failure.
2. Query that callee through CodeAtlas.
3. Read a narrow window around the callee.
4. Repeat until the failure condition is visible.

Example chain:

```text
internal.AnalyzeRepo
  -> internal.ParseFiles
  -> internal.ParseFile
  -> internal.extractSymbols
```

Each arrow should come from code you read, not guesswork.

Query each discovered symbol:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'ParseFiles|ParseFile|extractSymbols'
```

## Interface And Implementation Bugs

If the failure is caused by a contract mismatch, locate the interface first:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'Analyzer|Parser|Store|Repository|Client'
```

Open the interface definition, then locate implementations:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Analyzer|Parser|Store)'
```

Use `.files` imports to identify which implementation the failing entrypoint wires in:

```bash
jq '.files["cmd/codeatlas/main.go"].imports' codeatlas.json
```

Then follow constructor symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'New.*Analyzer|New.*Store|New.*Parser'
```

## Configuration, Routes, And Non-Symbol Failures

CodeAtlas is symbol-first. Some bug causes are not symbols:

- config keys
- route paths
- SQL text
- template names
- environment variables
- feature flags
- serialization tags

For these, use targeted text search:

```bash
rg -n 'CODEATLAS_CONFIG|/api/analyze|repository_id'
```

After finding a file, inspect symbols in that file:

```bash
jq '.files["internal/config.go"].symbols' codeatlas.json
```

Then continue with symbol lookup.

## Determine Impact Scope

Once you identify a likely faulty symbol, find related symbols:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg 'Analyze|Analyzer|Analysis'
```

Inspect neighboring symbols in the file:

```bash
jq '.files["internal/analyzer.go"].symbols' codeatlas.json
```

Check relevant tests:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '(Test.*Analyze|Analyze.*Test)'
```

Use language tooling or test commands after symbol retrieval has scoped the problem.

## Root Cause Notes

Keep notes symbol-first:

```text
Observed failure:
- failing test/log/trace

Likely failing symbol:
- symbol -> file:line

Execution path:
- caller -> callee -> condition

Cause:
- precise state or condition

Fix surface:
- symbols/files to edit

Tests:
- symbols or test command
```

This prevents debugging from drifting into broad repository reading.

## Common Failure Modes

If the stack trace line differs from the CodeAtlas line, the index may be stale. Query by symbol name and inspect the current file. Regenerate CodeAtlas after structural changes.

If the symbol is absent, the failure may involve generated code, dynamically invoked code, or a stale snapshot. Try narrowed `rg`, then regenerate if symbols have changed.

If many symbols match, add package prefixes, receiver names, or test names:

```bash
jq -r '.idx | keys[]' codeatlas.json | rg '^internal\.' | rg 'Analyze'
```

If no symbols match an error string, use targeted text search for the string and immediately switch back to CodeAtlas once a containing file or symbol is known.

## Freshness Rule For Bug Fixes

Regeneration is usually not required for a small bug fix inside an existing function.

Regenerate CodeAtlas if the fix:

- adds a new function, method, struct, interface, or package
- renames a symbol
- moves files
- changes package ownership
- performs a larger refactor

When unsure, regenerate. A fresh local index is cheaper than debugging against stale symbol locations.
