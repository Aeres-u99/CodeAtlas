## Critical Rule: CodeAtlas Is A Map, Not Context

Never load the entire `codeatlas.json` into the context window.

Never read the entire file.

Never summarize the entire file.

Treat CodeAtlas as a navigation map.

The purpose of CodeAtlas is to locate information, not to become information.

Correct workflow:

```text
Question
    ↓
Query CodeAtlas
    ↓
Locate Symbol
    ↓
Open Target File
    ↓
Read Implementation
```

Incorrect workflow:

```text
Question
    ↓
Read Entire codeatlas.json
    ↓
Attempt To Understand Repository
```

---

## Preferred Usage

Use targeted queries:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' codeatlas.json
```

```bash
jq -r '.idx | keys[]' codeatlas.json | rg Analyze
```

```bash
jq '.files["internal/parser.go"].symbols' codeatlas.json
```

Only retrieve the specific information required for the current task.

---

## Anti-Patterns

Never do:

```bash
cat codeatlas.json
```

Never do:

```bash
less codeatlas.json
```

Never do:

```bash
jq '.' codeatlas.json
```

Never do:

```bash
Read entire codeatlas.json
```

These operations defeat the purpose of CodeAtlas and increase token consumption.

---

## Principle

CodeAtlas should consume fewer tokens than repository exploration.

If an agent is reading large portions of codeatlas.json, it is using CodeAtlas incorrectly.

CodeAtlas is a map.

Follow the map.

Do not study the map.

```
```
