### CodeAtlas
⏺ This is a great idea. Here's my recommendation before we go deeper:

  Core format - a two-layer map:
  - Layer 1 (always in context): flat symbol index - name, file, line, kind, signature only. ~1-3 tokens per symbol.
  - Layer 2 (on-demand): per-file scope blocks fetched via treesitter when a symbol is targeted.

  Main tradeoff: flat index vs graph. A flat index is cheap to load and sufficient for "where is X" queries. A graph (with call
  edges, import chains) enables "what does X depend on" but costs 3-5x more tokens to represent and maintain.

  Proposed JSON schema:
```
  {
    "v": 1,
    "generated": "<timestamp>",
    "files": {
      "path/to/file.py": {
        "lang": "python",
        "loc": 160,
        "imports": ["os", "boto3"],
        "symbols": [
          {"n": "MyClass", "t": "cls", "l": 10},
          {"n": "MyClass.method", "t": "fn", "l": 24, "sig": "(self, x: str) -> bool"},
          {"n": "helper", "t": "fn", "l": 80, "sig": "(items: list) -> dict"}
        ]
      }
    },
    "idx": {
      "MyClass": "path/to/file.py:10",
      "MyClass.method": "path/to/file.py:24",
      "helper": "path/to/file.py:80"
    }
  }
```

### How CodeAtlas would work:

  ctags run  →  symbol name, file, line, kind
  treesitter →  precise signature, scope end line, imports
  merge      →  write/patch codeatlas.json

  LLM workflow then becomes:
  1. Load codeatlas.json (~3k tokens for 3000 LOC)
  2. Look up symbol in idx → get file:line
  3. Read only lines l to scope_end via treesitter scope

  Key questions to nail down before building:
  - Should it be a single codeatlas.json at root, or per-directory shards (better for monorepos)?
  - Do you want call edges in v1 or defer to v2?
  - Which languages need to be supported first - Python/HCL/YAML for this repo, or broader?
