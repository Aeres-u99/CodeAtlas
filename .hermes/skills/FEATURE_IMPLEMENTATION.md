# Feature Implementation With Hermes

Use this skill when adding functionality to an existing repository.

## Objective

Implement the feature by first locating the existing pattern.

The default path is:

```text
Locate existing implementation
  -> Query related symbols
  -> Discover extension points
  -> Inspect neighboring implementations
  -> Modify code
```

Do not start writing new code until you have found the closest existing behavior through Hermes.

## Start From Existing Behavior

Translate the requested feature into likely existing symbols.

Examples:

- "add JSON output" -> `Output`, `Formatter`, `JSON`, `Render`, `Encode`
- "add a CLI flag" -> `Command`, `Flags`, `Options`, `Config`, `Run`
- "support a new parser" -> `Parser`, `Parse`, `Language`, `Detect`, `Analyze`
- "add a repository backend" -> `Store`, `Repository`, `Client`, `Provider`

Query Hermes:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Output|Formatter|Render|Encode|JSON)'
```

Open the strongest match:

```bash
SYMBOL='internal.NewFormatter'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' hermes.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' hermes.json)
START=$((LINE > 30 ? LINE - 30 : 1))
END=$((LINE + 140))
sed -n "${START},${END}p" "$FILE"
```

## Find Extension Points

Look for symbols that represent extension boundaries:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Interface|Provider|Registry|Factory|Builder|Options|Config|Handler|Formatter|Parser|Store|Client)'
```

Find constructors:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(^|\.)(New|Create|Build|Init)'
```

Find registration functions:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Register|Registry|Add|Mount|Install)'
```

Open the relevant interface or factory before editing implementations.

## Inspect Neighboring Implementations

After locating a relevant file, inspect its symbols:

```bash
jq '.files["internal/formatter.go"].symbols' hermes.json
```

Inspect imports:

```bash
jq '.files["internal/formatter.go"].imports' hermes.json
```

Find symbols with the same naming family:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Formatter|Format|Render)'
```

Read one or two neighboring implementations to copy established patterns:

- error handling
- logging
- validation
- dependency injection
- test style
- naming conventions
- configuration flow

Do not infer project style from memory. Retrieve it.

## Implementation Planning

Before editing, produce a compact symbol-based plan:

```text
Existing pattern:
- symbol -> file:line

Extension point:
- symbol -> file:line

Files likely to edit:
- file because symbol X owns behavior Y

Tests:
- existing test symbol -> file:line
- new/updated test surface

Hermes freshness:
- regeneration required? yes/no
```

This plan should come from Hermes queries and code reads, not directory guesses.

## Adding New Symbols

If the feature requires a new function, method, type, interface, or package, implement it normally, then regenerate Hermes after the edit if the workflow depends on continued symbol retrieval.

Reason: the current snapshot cannot contain symbols added after generation.

Examples that require regeneration:

- adding `internal.RenderJSON`
- adding `type JSONFormatter struct`
- adding `func NewJSONFormatter(...)`
- creating a new `internal/output` package
- moving a handler to another file

Examples that usually do not require regeneration:

- changing logic inside `internal.Render`
- adding a branch to an existing function
- changing a validation condition
- updating comments or formatting

## Tests Through Hermes

Find existing tests by symbol:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Test.*Formatter|Formatter.*Test|Test.*Output|Output.*Test)'
```

Open the closest test:

```bash
SYMBOL='internal.TestFormatterJSON'
FILE=$(jq -r --arg s "$SYMBOL" '.idx[$s].f' hermes.json)
LINE=$(jq -r --arg s "$SYMBOL" '.idx[$s].l' hermes.json)
START=$((LINE > 30 ? LINE - 30 : 1))
END=$((LINE + 140))
sed -n "${START},${END}p" "$FILE"
```

If there are no test symbols, inspect nearby file metadata and use targeted file search only for tests:

```bash
rg --files | rg '(_test\.go|test_|\.test\.|spec\.)'
```

Then return to symbol-driven investigation.

## Command And CLI Features

Find command setup:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Command|Execute|Run|Flags|Options|Config)'
```

Look for existing flag registration:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Flag|Flags|Option|Bind)'
```

Open the command symbol, then trace into the execution symbol. Add new CLI behavior where existing flags and options are already wired, unless the code clearly separates parsing from execution.

## API And Handler Features

Find handlers and routes:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Handler|Serve|Route|Controller|Endpoint)'
```

If route strings are not symbolized, use targeted search:

```bash
rg -n 'GET |POST |PUT |DELETE |/api/'
```

After finding the route file, inspect symbols:

```bash
jq '.files["internal/server/routes.go"].symbols' hermes.json
```

Then query handler symbols through `.idx`.

## Data Model Features

Find model and serialization types:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Model|DTO|Request|Response|Config|Options|Record|Entity)'
```

Open the type definition and then inspect related constructors, validators, encoders, and tests.

Avoid adding fields without tracing where the type is constructed, serialized, validated, and tested.

## Failure Modes

If you cannot find a similar implementation, broaden by concept:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Render|Write|Encode|Print|Format)'
```

If no symbols match because the feature is cross-cutting, find entrypoints and extension boundaries first:

```bash
jq -r '.idx | keys[]' hermes.json | rg '(Run|Execute|Serve|Process|Analyze)'
```

If Hermes points to files that no longer exist, regenerate before implementing.

If the feature adds many symbols, regenerate midway before continuing symbol-based navigation.

## Completion Rule

Before finishing:

1. Re-run or inspect the relevant tests.
2. Regenerate Hermes if symbol topology changed.
3. Use the refreshed index for any final navigation.
4. Report edited files and verification performed.

The feature is not complete if the code works but the Hermes snapshot is stale after structural changes.
