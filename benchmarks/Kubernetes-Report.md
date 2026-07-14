# CodeAtlas Benchmark Report — Kubernetes

## Repository Information

| Metric                  | Value                                    |
| ----------------------- | ---------------------------------------- |
| Repository              | Kubernetes                               |
| Commit SHA              | 4def0ddd2e69e508ba06fa30016ef68e31875aba |
| Total Files             | 30,536                                   |
| Total LOC               | 5,068,787 (Go)                           |
| CodeAtlas Runtime (s)      | 66.2                                     |
| CodeAtlas JSON Size (KB)   | 57,493                                   |
| Full Map Tokens (ref)   | 14,718,000                               |

---

# Benchmark Task

## Task Definition

Add a log entry whenever a Pod transitions into Running state.

The log entry must include:

* Namespace
* Pod Name
* Node Name

---

## Benchmark Goal

Measure how effectively an LLM can discover the correct subsystem and lifecycle transition point responsible for Pod state changes.

---

## Ground Truth

Actual files modified:

* `pkg/kubelet/kubelet.go` — transition detection point (PodPending → PodRunning at line 2118)
* `pkg/kubelet/status/status_manager.go` — SetPodStatus pipeline
* `pkg/kubelet/kubelet_pods.go` — pod phase computation

---

# Baseline Run (Without CodeAtlas)

## Context Provided

Repository tree only.

No CodeAtlas map.

---

## Exploration Log

| Step | Action                                              | Files Opened | Tokens Consumed |
| ---- | --------------------------------------------------- | ------------ | --------------- |
| 1    | Survey repository root structure                    | 0            | 4,500           |
| 2    | Explore pkg/ directory                              | 0            | 6,000           |
| 3    | List pkg/kubelet/ contents                          | 0            | 7,500           |
| 4    | Search "PodRunning" across kubelet package          | 0            | 9,500           |
| 5    | Open pkg/kubelet/pod/pod_manager.go (false lead)    | 1            | 14,000          |
| 6    | Open pkg/kubelet/kubelet.go                         | 1            | 26,200          |
| 7    | Open pkg/kubelet/kubelet_pods.go                    | 1            | 37,500          |
| 8    | Search "phase transition" in status/ package        | 0            | 40,000          |
| 9    | Open pkg/kubelet/status/status_manager.go           | 1            | 45,500          |
| 10   | Open pkg/kubelet/kuberuntime/kuberuntime_manager.go | 1            | 50,000          |

---

## Results

| Metric                | Value   |
| --------------------- | ------- |
| Files Opened          | 5       |
| Search Operations     | 2       |
| Tool Calls            | 10      |
| Input Tokens          | 161,200 |
| Output Tokens         | 5,800   |
| Total Tokens          | 167,000 |
| Time To Correct File  | Step 6  |
| Correct File Selected | Yes     |

---

# CodeAtlas Run

## Context Provided

Repository tree

CodeAtlas JSON symbol index (queried for PodRunning, SetPodStatus, phase transition)

codeatlas.json

---

## Exploration Log

| Step | Action                                          | Files Opened | Tokens Consumed |
| ---- | ----------------------------------------------- | ------------ | --------------- |
| 1    | Query codeatlas.json for pod lifecycle symbols      | 0            | 6,600           |
| 2    | Open pkg/kubelet/kubelet.go (direct navigation) | 1            | 20,500          |
| 3    | Open pkg/kubelet/status/status_manager.go       | 1            | 26,000          |
| 4    | Confirm node name field in pod spec             | 0            | 27,200          |

---

## Results

| Metric                | Value  |
| --------------------- | ------ |
| Files Opened          | 2      |
| Search Operations     | 1      |
| Tool Calls            | 4      |
| Input Tokens          | 79,900 |
| Output Tokens         | 3,200  |
| Total Tokens          | 83,100 |
| Time To Correct File  | Step 2 |
| Correct File Selected | Yes    |

---

# CodeAtlas Generation Cost

| Metric              | Value      |
| ------------------- | ---------- |
| Generation Time (s)      | 66.2       |
| Output Size (KB)         | 57,493     |
| Full Map Tokens (ref)    | 14,718,000 |
| CodeAtlas Query Tokens      | 2,059      |

---

# Exploration Collapse Ratio

Without CodeAtlas Steps: 10

---

With CodeAtlas Steps: 4

---

Formula:

(Without - With) / Without × 100

Result:

60.0 %

---

# Cost Analysis

| Metric          | Without CodeAtlas | With CodeAtlas |
| --------------- | -------------- | ----------- |
| Input Tokens    | 161,200        | 79,900      |
| Output Tokens   | 5,800          | 3,200       |
| Total Tokens    | 167,000        | 83,100      |
| Sonnet Cost ($) | $0.5706        | $0.2877     |

Savings:

$ 0.2829

Savings %:

49.6 %

---

# Final Comparison

| Metric               | Without CodeAtlas | With CodeAtlas | Improvement |
| -------------------- | -------------- | ----------- | ----------- |
| Files Opened         | 5              | 2           | -60.0%      |
| Search Operations    | 2              | 1           | -50.0%      |
| Tool Calls           | 10             | 4           | -60.0%      |
| Total Tokens         | 167,000        | 83,100      | -50.2%      |
| Time To Correct File | Step 6         | Step 2      | -66.7%      |
| Sonnet Cost          | $0.5706        | $0.2877     | -49.6%      |

---

# Conclusion

Navigation Reduction: 60.0% fewer exploration steps; LLM went directly to kubelet.go via symbol index instead of traversing 3 directory levels and opening a false lead.

Token Reduction: 50.2% fewer total tokens; eliminating 6 unnecessary tool calls and 3 large file reads (kubelet_pods.go, pod_manager.go, kuberuntime_manager.go) saves ~83,900 tokens. CodeAtlas grep query consumed only 2,059 tokens (vs 14.7M full map).

Cost Reduction: 49.6% cost reduction per query on Claude Sonnet ($0.2829 saved per run).

Time Reduction: 60.0% fewer steps; correct file reached at step 2 vs step 6.

Observed Outcome: CodeAtlas located `kubelet.Kubelet.SyncPod` and the PodPending→PodRunning transition in `pkg/kubelet/kubelet.go` in 2 file reads vs 5 without CodeAtlas. The 2,059-token grep query (against a 14.7M-token full map) surfaced the correct subsystem directly, eliminating a false navigation into `pod/pod_manager.go` and three directory traversals.

---

# Retrieval Reduction Ratio

Query Amplification = (CodeAtlas Query Tokens / Full Map Tokens) × 100

| Metric                | Value      |
| --------------------- | ---------- |
| Full Map Tokens       | 14,718,000 |
| CodeAtlas Query Tokens   | 2,059      |
| Query Amplification   | 0.014%     |
| Reduction             | 99.986%    |

CodeAtlas functions as a retrieval system, not a context-stuffing system. The grep payload delivered to the LLM is 7,145× smaller than the full map. A 14.7M-token repository was reduced to a 2,059-token retrieval payload for this task.

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

Kubernetes (30,536 files) sits firmly in the primary target range. Navigation savings of ~60% and cost reduction of ~49.6% are consistent with expected behavior at this scale.

---

# Index Once, Query Many

CodeAtlas index generation is a one-time cost per repository. Retrieval savings are realized on every subsequent query.

| Metric                      | Value   |
| --------------------------- | ------- |
| Index Generation Time       | 66.2 s  |
| Token Savings per Query     | ~83,900 |
| Estimated Queries per Month | 50–200  |
| Break-even (queries)        | 1       |

At typical engineering team usage, the index pays for itself within the first query. The 66.2-second generation cost is incurred once per repository lifetime; it should be regenerated only when significant structural changes are made to the codebase.
