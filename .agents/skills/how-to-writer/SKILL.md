---
name: how-to-writer
description: Specialized tool for evaluating and improving "how-to" documentation and step-by-step guides. Unlike general writing tools, this focuses on technical clarity, prerequisite identification, and verification steps for beginners. Use when a user asks to "review my how-to", "score this guide", or "help me write a tutorial".
version: "0.1.0"
allowed-tools: Read
---

# How-To Writer Helper

This skill helps you write high-quality "how-to" documentation for beginners. It uses a structured evaluation rubric to score existing drafts or provide suggestions for improvement.

## Workflow

1.  **Analyze Content**: Read the how-to guide provided by the user.
2.  **Consult Rubric**: Reference `references/rubric.md` for the scoring criteria and `references/example-guide.md` for a "Gold Standard" example of a high-scoring guide.
3.  **Perform Task**:
    *   **If Scoring**: Assign a score (1-4) for each of the 6 categories in the rubric. Calculate the total score out of 24.
    *   **If Suggesting**: Provide specific, actionable suggestions for improvement based on categories where the writing could be better (e.g., clarity, structure, verification).
4.  **Format Output**: Present the evaluation in a clear, formatted table or list.

## Rubric Categories

1.  **Clarity and Simplicity of Language**: Simple, direct, jargon-free.
2.  **Step-by-Step Structure**: Small, manageable, numbered steps.
3.  **Goal and Prerequisite Clarity**: Clear outcome and required setup listed upfront.
4.  **Verification and Feedback**: Expected outcomes after steps and at the end.
5.  **Visuals and Concrete Examples**: High-quality snippets and examples.
6.  **Troubleshooting and Common Issues**: Proactive error handling.

## Example Scoring Output

| Category | Score | Observations |
| :--- | :--- | :--- |
| Clarity | 3 | Mostly clear, but uses "reconciliation loop" without definition. |
| Structure | 4 | Excellent use of numbered steps. |
| ... | ... | ... |
| **Total** | **18/24** | **Good effort, focus on jargon and troubleshooting.** |

## Example Suggestion Output

> **Suggestion for Step 3**: Instead of "Configure the ingress", use "Step 3: Create the Ingress resource by running `kubectl apply -f ingress.yaml`. You should see the message `ingress.networking.k8s.io/my-app created`." (Improves **Verification and Feedback**)
