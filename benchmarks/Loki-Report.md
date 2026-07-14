# CodeAtlas Benchmark Report — Loki

## Repository Information

| Metric                  | Value                                    |
| ----------------------- | ---------------------------------------- |
| Repository              | Loki                                     |
| Commit SHA              | 66ac3f08235b581a8b7149c11129ea5c008ce8ef |
| Total Files             | 17,160                                   |
| Total LOC               | 510,562 (Go)                             |
| CodeAtlas Runtime (s)      | 11.6                                     |
| CodeAtlas JSON Size (KB)   | 32,797                                   |
| Full Map Tokens (ref)   | 8,396,000                                |

---

# Benchmark Task

## Task Definition

Add query duration to request logging.

Every query log entry should include total execution time in milliseconds.

---

## Benchmark Goal

Measure how effectively an LLM can identify the correct interaction between request handling, query execution, and logging subsystems.

---

## Ground Truth

Actual files modified:

* `pkg/logql/metrics.go` — RecordRangeAndInstantQueryMetrics adds duration_ms field
* `pkg/logql/engine.go` — execTime computed via time.Since(start), passed to metrics
* `pkg/querier/http.go` — HTTP handler wires query result to logging pipeline

---

# Baseline Run (Without CodeAtlas)

## Context Provided

Repository tree only.

No CodeAtlas map.

---

## Exploration Log

| Step | Action                                               | Files Opened | Tokens Consumed |
| ---- | ---------------------------------------------------- | ------------ | --------------- |
| 1    | Survey repository root structure                     | 0            | 4,200           |
| 2    | Explore pkg/ directory, identify logql and querier   | 0            | 5,800           |
| 3    | Search "query.*log\|log.*query" across pkg/          | 0            | 8,400           |
| 4    | Open pkg/querier/http.go (entry point inspection)    | 1            | 11,400          |
| 5    | Open pkg/logql/engine.go (query execution)           | 1            | 15,000          |
| 6    | Open pkg/logql/metrics.go (logging sink)             | 1            | 18,100          |
| 7    | Search "ExecTime\|execTime\|duration" to verify path | 0            | 19,800          |
| 8    | Open pkg/distributor/http.go (false lead)            | 1            | 23,100          |

---

## Results

| Metric                | Value  |
| --------------------- | ------ |
| Files Opened          | 4      |
| Search Operations     | 2      |
| Tool Calls            | 8      |
| Input Tokens          | 86,800 |
| Output Tokens         | 4,100  |
| Total Tokens          | 90,900 |
| Time To Correct File  | Step 5 |
| Correct File Selected | Yes    |

---

# CodeAtlas Run

## Context Provided

Repository tree

CodeAtlas JSON symbol index (queried for RecordRangeAndInstantQueryMetrics, ExecTime, query logging)

codeatlas.json

---

## Exploration Log

| Step | Action                                           | Files Opened | Tokens Consumed |
| ---- | ------------------------------------------------ | ------------ | --------------- |
| 1    | Query codeatlas.json for query duration symbols     | 0            | 5,380           |
| 2    | Open pkg/logql/metrics.go (direct navigation)    | 1            | 9,800           |
| 3    | Open pkg/logql/engine.go (confirm execTime path) | 1            | 13,200          |
| 4    | Verify duration_ms addition site                 | 0            | 14,100          |

---

## Results

| Metric                | Value  |
| --------------------- | ------ |
| Files Opened          | 2      |
| Search Operations     | 1      |
| Tool Calls            | 4      |
| Input Tokens          | 39,900 |
| Output Tokens         | 2,700  |
| Total Tokens          | 42,600 |
| Time To Correct File  | Step 2 |
| Correct File Selected | Yes    |

---

# CodeAtlas Generation Cost

| Metric              | Value     |
| ------------------- | --------- |
| Generation Time (s)      | 11.6      |
| Output Size (KB)         | 32,797    |
| Full Map Tokens (ref)    | 8,396,000 |
| CodeAtlas Query Tokens      | 880       |

---

# Exploration Collapse Ratio

Without CodeAtlas Steps: 8

---

With CodeAtlas Steps: 4

---

Formula:

(Without - With) / Without × 100

Result:

50.0 %

---

# Cost Analysis

| Metric          | Without CodeAtlas | With CodeAtlas |
| --------------- | -------------- | ----------- |
| Input Tokens    | 86,800         | 39,900      |
| Output Tokens   | 4,100          | 2,700       |
| Total Tokens    | 90,900         | 42,600      |
| Sonnet Cost ($) | $0.3219        | $0.1602     |

Savings:

$ 0.1617

Savings %:

50.2 %

---

# Final Comparison

| Metric               | Without CodeAtlas | With CodeAtlas | Improvement |
| -------------------- | -------------- | ----------- | ----------- |
| Files Opened         | 4              | 2           | -50.0%      |
| Search Operations    | 2              | 1           | -50.0%      |
| Tool Calls           | 8              | 4           | -50.0%      |
| Total Tokens         | 90,900         | 42,600      | -53.1%      |
| Time To Correct File | Step 5         | Step 2      | -60.0%      |
| Sonnet Cost          | $0.3219        | $0.1602     | -50.2%      |

---

# Conclusion

Navigation Reduction: 50.0% fewer steps; codeatlas.json symbol index surfaced RecordRangeAndInstantQueryMetrics and execTime directly, bypassing a false lead into pkg/distributor/http.go and two directory traversals.

Token Reduction: 53.1% fewer total tokens; skipping distributor/http.go and two search passes saves ~48,300 tokens. CodeAtlas grep query consumed only 880 tokens (vs 8.4M full map).

Cost Reduction: 50.2% cost reduction per query on Claude Sonnet ($0.1617 saved per run).

Time Reduction: 60.0% faster to correct file; step 2 vs step 5.

Observed Outcome: CodeAtlas located `RecordRangeAndInstantQueryMetrics` in `pkg/logql/metrics.go` and the `execTime` computation in `pkg/logql/engine.go` in 2 file reads vs 4 without CodeAtlas. The 880-token grep query (smallest of the three benchmarks) demonstrates precise retrieval for a tightly-named symbol — `execTime` and `RecordRangeAndInstantQueryMetrics` appear in few files, producing a minimal query payload from an 8.4M-token map.

---

# Retrieval Reduction Ratio

Query Amplification = (CodeAtlas Query Tokens / Full Map Tokens) × 100

| Metric                | Value     |
| --------------------- | --------- |
| Full Map Tokens       | 8,396,000 |
| CodeAtlas Query Tokens   | 880       |
| Query Amplification   | 0.010%    |
| Reduction             | 99.990%   |

CodeAtlas functions as a retrieval system, not a context-stuffing system. The grep payload delivered to the LLM is 9,541× smaller than the full map. An 8.4M-token repository was reduced to an 880-token retrieval payload for this task — the smallest query output across the three benchmarks, reflecting that `RecordRangeAndInstantQueryMetrics` and `execTime` are tightly scoped symbols.

---

# Cross-Repository Scaling Analysis

| Repository | Files  | LOC       | Runtime (s) | JSON (KB) | Query Tokens |
| ---------- | ------ | --------- | ----------- | --------- | ------------ |
| Terraform  | 5,411  | 667,641   | 13.8        | 7,773     | 5,340        |
| Loki       | 17,160 | 510,562   | 11.6        | 32,797    | 880          |
| Kubernetes | 30,536 | 5,068,787 | 66.2        | 57,493    | 2,059        |

Key Observation: Repository size increases 6× in files and 10× in Go LOC across benchmarks. Query token consumption does not scale proportionally — it is governed by symbol frequency across implementations, not repository size. CodeAtlas scales sub-linearly from a retrieval perspective.

---

# CodeAtlas Target Repository Size

CodeAtlas is designed primarily for medium and large repositories where repository exploration constitutes a significant portion of LLM token usage.

| Repository Size    | Expected Benefit    |
| ------------------ | ------------------- |
| < 100 files        | Not recommended     |
| 100–1,000 files    | Situational benefit |
| 1,000–10,000 files | Strong benefit      |
| 10,000+ files      | Primary target      |

Loki (17,160 files) sits in the strong benefit range. Navigation savings of ~50% and cost reduction of ~50.2% are consistent with expected behavior at this scale.

---

# Index Once, Query Many

CodeAtlas index generation is a one-time cost per repository. Retrieval savings are realized on every subsequent query.

| Metric                      | Value   |
| --------------------------- | ------- |
| Index Generation Time       | 11.6 s  |
| Token Savings per Query     | ~48,300 |
| Estimated Queries per Month | 50–200  |
| Break-even (queries)        | 1       |

At typical engineering team usage, the index pays for itself within the first query. The 11.6-second generation cost is incurred once per repository lifetime; it should be regenerated only when significant structural changes are made to the codebase.
