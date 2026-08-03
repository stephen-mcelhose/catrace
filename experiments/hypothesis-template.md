# Hypothesis: [short title]

> Fill in all sections marked [TODO] **before running catrace**. Fill in Results and Verdict after.

## Claim

[TODO: One or two sentences stating the architectural assertion in plain language.
Must be falsifiable — the numbers should be able to go the other way.]

**Example:** "Mechanism X is the primary recovery driver; mechanism Y contributes but is not the bottleneck."

## Context

- **Pattern:** [Which agentic pattern this addresses — e.g., Self-Healing / Adaptive (14)]
- **Related example:** [Which catrace example, if any, this extends]
- **Motivation:** [Why this question matters for architectural decision-making]

## Variant definitions

| Variant | Label | Description |
|---------|-------|-------------|
| A | [name] | [What is different about this variant — which mechanism is stronger/weaker] |
| B | [name] | [What is different about this variant] |

*Keep topology identical across variants. Only kernel parameter values should differ.*

## Variable kernel entries

List the specific matrix entries that differ between variants. Everything else is held equal.

| Kernel | Entry (from, to) | Variant A value | Variant B value | Why this encodes the claim |
|--------|-----------------|-----------------|-----------------|---------------------------|
| [TODO] | [TODO] | [TODO] | [TODO] | [TODO] |

## Predictions

State the expected direction for each metric *before running*. Circle back and check after.

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| Stationary mass | π([target state]) | A > B | [TODO] |
| Recovery speed | MFPT([bad state] → [good state]) | A < B | [TODO] |
| Predictability | H(J) bits/step | A < B | [TODO] |
| Observable health | π_trace([coarse state]) | A > B | [TODO] |

*Remove rows for metrics that don't apply.*

## Verdict rule

[TODO: How many metrics must agree to call the claim supported?]

**Recommended:** majority (≥2 of 3) for three-metric comparisons; unanimous for two-metric comparisons where metrics are closely related.

---

## Results

> Fill in after running.

| Metric | Variant A | Variant B | Predicted direction | Correct? |
|--------|-----------|-----------|--------------------|----|
| π([target state]) | | | A > B | |
| MFPT([bad] → [good]) | | | A < B | |
| H(J) | | | A < B | |
| π_trace([coarse]) | | | A > B | |

## Verdict

**Claim:** [supported / not supported / trade-off]
**Metrics in agreement:** [N/N]

## Interpretation

[What the numbers mean in architectural terms. If the verdict is "trade-off", describe what each variant wins on and under what operating conditions each would be preferred.]
