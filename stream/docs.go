// Package stream provides streaming kernel estimation and policy analysis
// for JSONL agent event traces.
//
// The core pipeline:
//
//  1. [IngestJSONL] reads a JSONL file or stdin, decoding each line into an [Event].
//  2. [Estimator] accumulates EMA-weighted transition counts as events arrive via [Estimator.ObserveNext].
//  3. [Estimator.Kernel] materialises a *catrace.Kernel from current counts.
//  4. [TraceByName] optionally folds the kernel onto named boundary states.
//  5. [AnalyzeAndSuggest] ranks single-row perturbations that minimise mean first-passage time to a target state.
//
// [Monitor] wraps two Estimators at different half-lives to track behavioural drift via Jensen-Shannon divergence.
package stream
