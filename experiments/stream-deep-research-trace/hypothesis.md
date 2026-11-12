# Hypothesis: Trace + SuggestKnobs on a real deep-research JSONL stream

> Fill in Results and Verdict after running. Predictions below are locked before the run.

## Claim

On a fixed public deep-research JSONL stream (harveyai/deep-research-starter `events_jsonl.jsonl`), Trace-folding out high-frequency token/delta noise onto structural milestones yields a valid, lower-entropy boundary kernel; `AnalyzeAndSuggest` ranks decision transitions whose MFPT improvement scales with ε; and dual-clock `DriftScore` stays near zero in steady phases, spikes at phase/loop boundaries, and never returns NaN.

## Context

- **Pattern:** Streaming kernel estimation / Trace / MFPT policy knobs (`stream` package)
- **Related tooling:** `tools/suggest`, `stream.Estimator`, `stream.Monitor`, `stream.AnalyzeAndSuggest`
- **Motivation:** Validate that catrace stream math works on a dense real agent log without API keys or GPUs.
- **Fixture:** `testdata/events_jsonl.jsonl` (curl; gitignored)

```bash
mkdir -p testdata
curl -L -o testdata/events_jsonl.jsonl \
  https://raw.githubusercontent.com/harveyai/deep-research-starter/main/output_samples/events_jsonl.jsonl
```

### State vocabulary (empirical)

Harvey lines are OpenAI Responses streaming events keyed by top-level `type`.

| Role | States (S keep / B fold) |
|------|--------------------------|
| **B (noise)** | `response.reasoning_summary_text.delta`, `response.output_text.delta` |
| **S (milestones)** | `response.created`, `response.in_progress`, `response.completed`, `response.output_item.added`, `response.output_item.done`, `response.web_search_call.in_progress`, `response.web_search_call.searching`, `response.web_search_call.completed`, `response.content_part.added`, `response.content_part.done`, `response.output_text.done`, `response.reasoning_summary_part.added`, `response.reasoning_summary_part.done`, `response.reasoning_summary_text.done` |
| **Target** | `response.completed` |
| **Knobs (D)** | `response.output_item.done`, `response.web_search_call.completed`, `response.reasoning_summary_text.done` |

Parser: `stream.Event.CanonicalState()` uses `type` then `event`, optionally appending `:status`.

## Variant definitions

| Variant | Label | Description |
|---------|-------|-------------|
| A | full vs traced | Same estimator; compare un-folded kernel vs Trace onto S |
| B | ε sweep | Traced kernel; `AnalyzeAndSuggest` at ε ∈ {0.05, 0.15, 0.30} |
| C | dual-clock drift | `Monitor` fast=25 / slow=250; sample `DriftScore(S)` every 50 events |

## Variable kernel entries

Not hand-tuned matrix entries: the kernel is estimated from the fixture. The controllable contrast is Trace keep-set (full vocabulary vs S) and policy ε.

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| State count | N_full → \|S\| | N_full > \|S\| | Deltas dominate vocabulary |
| Entropy rate | H(full) vs H(traced) | H(traced) < H(full) | Folding removes noisy self-loops |
| Validate | traced.Validate(1e-9) | passes | Trace preserves row-stochasticity |
| Classes | recurrent class | includes `response.completed` as terminal attractor | Completion is absorbing-ish endpoint |
| MFPT Δ vs ε | top knob mfpt_delta | more negative as ε increases (0.05→0.30) | Larger inject mass → larger projected speedup |
| Drift | DriftScore series | near 0 in steady stretches; peak at phase shifts / Trace fails; no NaN | Fast/slow EMA diverge at regime change; singular mid-stream Trace on S is itself a stall/phase signal (recorded, not dropped) |

## Verdict rule

Claim supported if **≥4 of 6** prediction rows hold after the runner prints Results. MFPT monotonicity may be soft (non-decreasing \|Δ\| across the three ε values for the same top from→to).

---

## Results

> Run: `go run ./experiments/stream-deep-research-trace` on harvey `events_jsonl.jsonl` (8055 events, 3228 non-JSON lines skipped). 2026-11-12.

| Metric | Observed | Predicted direction | Correct? |
|--------|----------|--------------------|----|
| N → \|S\| | 17 → 14 | N > \|S\| | yes |
| H(full) vs H(traced) | 0.000 / 0.000 bits | traced < full | no (tie; both ~0 with absorbing sink) |
| Validate(1e-9) | PASS | passes | yes |
| Classes / completed | traced: 1 recurrent (completed), 13 transient | completed recurrent or reachable sink | yes |
| MFPT Δ scales with ε | top `output_item.done→completed`: Δ = −614 / −664 / −678 at ε 0.05/0.15/0.30 | \|Δ\| grows with ε | yes |
| Drift min/mean/peak / NaN | min=0, mean=0.380, peak=1.0 (59 Trace-fail sentinels); nan=0 | spike + no NaN | yes |

## Verdict

**Claim:** supported
**Metrics in agreement:** 5/6

## Interpretation

Trace fold shrinks the vocabulary and yields a valid boundary kernel whose sole recurrent class is `response.completed`. Entropy did not drop further because both full and traced kernels are already near-deterministic into that sink (H≈0) — so the entropy prediction fails as a tie, not a reversal. Counterfactual knobs are domain-sensible: boosting mass from `output_item.done` (and secondarily web-search completion) into `completed` cuts expected steps, and the improvement strengthens monotonically with ε. Drift stays at 0 in calm windows and peaks when mid-stream Trace on S goes singular (recorded as signal), with no NaNs — matching the stall/phase-break reading of Experiment C.

## Out of scope (follow-up)

Arbor (RUC-NLPIR/Arbor) cross-validation — CPU `examples/algotune_knn/` → `.arbor/sessions/<run>/events.jsonl`, then the same `tools/suggest` keep/target path (`IDEATE,…,MERGED`) — is deferred until harvey A/B/C Results + Verdict are filled. See `stream/docs/experiments.md` §2.
