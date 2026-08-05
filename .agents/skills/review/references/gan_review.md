# GAN (adversarial) review — catrace

Use this reference **only in GAN mode** (user asked, or merge-critical code).

Goal: **break** the change.

## Stance

- Assume the implementation is wrong until proven otherwise.
- Prefer concrete failures over style nits.
- Walkthroughs are not GAN (use peer-code-review only when asked).

## When GAN mode applies

- User says GAN / adversarial review, or
- Diff includes package-root `*.go`, `experiments/*/main.go`, or merge-ready runnable PR

Otherwise use normal review output in the skill.

## Procedure

1. Diff scope vs base.
2. Threaten claimed behavior (tests, walkthrough, PR).
3. Attack checklist (domain-specific):
   - Kernel rows not stochastic / negative after boost
   - Joint index order mismatches names
   - Kronecker used where coupling was claimed
   - `Stationary` / `EntropyRate` on non-ergodic chain
   - `Trace` when hidden set contains recurrent states
   - MFPT self-loop convention vs docs
   - Silent ignored errors in examples
4. Prove with test or exact repro when possible.
5. Report attacks table + required fixes + residual risk.

## Shell

Prefer `jq` / `gh --jq` when reading issue/PR JSON during review.
