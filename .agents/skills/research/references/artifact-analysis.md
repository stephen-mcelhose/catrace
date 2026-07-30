---
title: Artifact Analysis
description: How to extract high-value insights from existing documents and artifacts
---

# Artifact Analysis

When you need to deeply analyze existing documents, research, or plans to extract actionable insights.

## Core Principle

**Extract high-value insights, filter out noise.**

Your job is to identify main decisions, actionable recommendations, and critical technical details while filtering out tangential or outdated information.

## What to Extract

### Focus On

- **Decisions made**: "We decided to..."
- **Trade-offs analyzed**: "X vs Y because..."
- **Constraints identified**: "We must..." "We cannot..."
- **Lessons learned**: "We discovered that..."
- **Action items**: "Next steps..." "TODO..."
- **Technical specifications**: Specific values, configs, approaches

### Filter Out

- Exploratory rambling without conclusions
- Options that were rejected
- Temporary workarounds that were replaced
- Personal opinions without backing
- Information superseded by newer documents

## Retrieving Artifacts for Analysis

Use `RetrieveArtifact` to get the full content:

```
RetrieveArtifact
  artifact_id: "abc123"
  markConsumed: true
```

For **research** artifacts, the system automatically parses and returns structured data:

- `criticalFiles`: Array of `{ path, purpose, exports }` - key files identified
- `codeExamples`: Array of code snippets with descriptions
- `implementationHints`: Actionable implementation guidance
- `constraints`: Known limitations or requirements
- `openQuestions`: Unresolved questions from the research

This parsed structure can accelerate your analysis by pre-extracting key sections.

## Analysis Strategy

### Step 1: Read with Purpose

- Read the entire document first
- Identify the document's main goal
- Note the date and context
- Understand what question it was answering
- Check the parsed structure (if available) for quick extraction

### Step 2: Extract Strategically

Ask yourself:

- What decisions were made?
- What constraints were identified?
- What specific values/configs were chosen?
- What lessons were learned?
- What's still open/unresolved?

### Step 3: Filter Ruthlessly

Remove anything that is:

- Just exploring possibilities
- Personal musing without conclusion
- Clearly superseded
- Too vague to action
- Redundant with better sources

## Output Format

```markdown
## Analysis of: [Document Path]

### Document Context
- **Date**: [When written]
- **Purpose**: [Why this document exists]
- **Status**: [Still relevant / Implemented / Superseded?]

### Key Decisions
1. **[Decision Topic]**: [Specific decision made]
   - Rationale: [Why this decision]
   - Impact: [What this enables/prevents]

2. **[Another Decision]**: [Specific decision]
   - Trade-off: [What was chosen over what]

### Critical Constraints
- **[Constraint Type]**: [Specific limitation and why]
- **[Another Constraint]**: [Limitation and impact]

### Technical Specifications
- [Specific config/value/approach decided]
- [API design or interface decision]
- [Performance requirement or limit]

### Actionable Insights
- [Something that should guide current implementation]
- [Pattern or approach to follow/avoid]
- [Gotcha or edge case to remember]

### Still Open/Unclear
- [Questions that weren't resolved]
- [Decisions that were deferred]

### Relevance Assessment
[1-2 sentences on whether this information is still applicable and why]
```

## Quality Filters

### Include Only If

- It answers a specific question
- It documents a firm decision
- It reveals a non-obvious constraint
- It provides concrete technical details
- It warns about a real gotcha/issue

### Exclude If

- It's just exploring possibilities
- It's personal musing without conclusion
- It's been clearly superseded
- It's too vague to action
- It's redundant with better sources

## Example Transformation

### From Document
>
> "I've been thinking about rate limiting and there are so many options. We could use Redis, or maybe in-memory, or perhaps a distributed solution. Redis seems nice because it's battle-tested, but adds a dependency. In-memory is simple but doesn't work for multiple instances. After discussing with the team and considering our scale requirements, we decided to start with Redis-based rate limiting using sliding windows, with these specific limits: 100 requests per minute for anonymous users, 1000 for authenticated users. We'll revisit if we need more granular controls. Oh, and we should probably think about websockets too at some point."

### To Analysis

```markdown
### Key Decisions
1. **Rate Limiting Implementation**: Redis-based with sliding windows
   - Rationale: Battle-tested, works across multiple instances
   - Trade-off: Chose external dependency over in-memory simplicity

### Technical Specifications
- Anonymous users: 100 requests/minute
- Authenticated users: 1000 requests/minute
- Algorithm: Sliding window

### Still Open/Unclear
- Websocket rate limiting approach
- Granular per-endpoint controls
```

## Guidelines

- **Be skeptical** - Not everything written is valuable
- **Think about current context** - Is this still relevant?
- **Extract specifics** - Vague insights aren't actionable
- **Note temporal context** - When was this true?
- **Highlight decisions** - These are usually most valuable
- **Question everything** - Why should the reader care about this?

## Remember

You're a curator of insights, not a document summarizer. Return only high-value, actionable information that will actually help make progress.
