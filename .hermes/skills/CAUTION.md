## Critical Rule: Hermes Is A Map, Not Context

Never load the entire `hermes.json` into the context window.

Never read the entire file.

Never summarize the entire file.

Treat Hermes as a navigation map.

The purpose of Hermes is to locate information, not to become information.

Correct workflow:

```text
Question
    ↓
Query Hermes
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
Read Entire hermes.json
    ↓
Attempt To Understand Repository
```

---

## Preferred Usage

Use targeted queries:

```bash
jq -r '.idx["internal.AnalyzeRepo"]' hermes.json
```

```bash
jq -r '.idx | keys[]' hermes.json | rg Analyze
```

```bash
jq '.files["internal/parser.go"].symbols' hermes.json
```

Only retrieve the specific information required for the current task.

---

## Anti-Patterns

Never do:

```bash
cat hermes.json
```

Never do:

```bash
less hermes.json
```

Never do:

```bash
jq '.' hermes.json
```

Never do:

```bash
Read entire hermes.json
```

These operations defeat the purpose of Hermes and increase token consumption.

---

## Principle

Hermes should consume fewer tokens than repository exploration.

If an agent is reading large portions of hermes.json, it is using Hermes incorrectly.

Hermes is a map.

Follow the map.

Do not study the map.

```
```
