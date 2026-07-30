---
title: Web Research
description: How to conduct effective web research for current information
---

# Web Research

When you need information beyond what's in the codebase - documentation, best practices, current standards.

## Core Principle

**Find authoritative, current information with proper attribution.**

Use web search for modern information that wouldn't be in your training data or the codebase.

## When to Use Web Research

- Looking for best practices and current standards
- Researching library/framework APIs and documentation
- Finding security considerations
- Understanding industry conventions
- Comparing approaches with community consensus

## Search Strategy

### Step 1: Analyze the Query

Before searching, identify:

- Key terms and concepts
- Relevant source types (docs, tutorials, Q&A)
- Multiple search angles
- Related terms and synonyms

### Step 2: Execute Strategic Searches

Start broad, then refine:

**Broad initial search:**

```
WebSearch
  query: "TypeScript authentication best practices 2025"
```

**Specific technical query:**

```
WebSearch
  query: "zod discriminated union validation"
  num_results: 5
```

**WebSearch Parameters:**

- `query`: Search query (required, min 2 characters)
- `num_results`: Maximum results to return (optional, default: 5, max: 10)

**Returns:** Array of `{ title, url }` objects from search results.

### Step 3: Fetch and Validate

Fetch full content from promising results:

```
WebFetch
  url: "https://docs.example.com/auth"
  prompt: "Extract authentication patterns and recommendations"
```

**WebFetch Parameters:**

- `url`: Full URL to fetch (required, must be http:// or https://)
- `prompt`: What information to extract (required)

**How it works:**

1. First tries Google AI's `urlContext` tool for direct access
2. Falls back to HTTP fetch + HTML-to-Markdown conversion if needed
3. AI summarizes content based on your prompt
4. Returns summary (not raw content) to preserve context

**Returns:** `{ success, url, summary, method }` where method is `"urlContext"` or `"fallback"`

Prioritize:

- Official documentation
- Reputable sources (major tech blogs, established libraries)
- Recent content (check dates)

### Step 4: Synthesize Findings

Organize findings by:

- Authority level (official docs > tutorials > forums)
- Relevance to specific query
- Publication date (note when information might be outdated)

## Output Format

```markdown
## Web Research: [Topic]

### Summary
[2-3 paragraphs synthesizing key findings]

### Official Documentation

#### [Source Title](url)
**Authority**: Official documentation
**Date**: [Publication/update date]

**Key Information**:
- Point 1
- Point 2

**Relevant Quote**:
> "Direct quote with context"

### Community Best Practices

#### [Source Title](url)
**Authority**: [Reputable blog / Major tutorial site]
**Date**: [Publication date]

**Key Insights**:
- Insight 1
- Insight 2

### Conflicting Information
[Note any disagreements between sources and why]

### Recommendations
Based on research:
1. [Recommendation with rationale]
2. [Recommendation with rationale]

### Information Gaps
- [What couldn't be found]
- [What needs more research]

### Sources
- [Title](url) - Brief description
- [Title](url) - Brief description
```

## Quality Standards

### Sources Must Be

- **Accurate**: Direct citations, verifiable claims
- **Relevant**: Specific to the query
- **Current**: Dates noted, recency considered
- **Authoritative**: Official sources prioritized

### Always Include

- Publication dates
- Direct quotes with attribution
- Source URLs
- Confidence level in findings

## Search Tips

- Start with 2-3 well-crafted searches before fetching
- Initially retrieve 3-5 most promising pages
- Refine search terms based on initial results
- Use search operators effectively (`site:`, quotes, `-exclude`)
- Explore varied formats (tutorials, docs, Q&A, forums)

## What NOT to Do

- Don't rely on a single source
- Don't ignore publication dates
- Don't present opinions as facts
- Don't skip attribution
- Don't assume first result is best

## Remember

You're gathering external knowledge to complement codebase understanding. Always attribute sources and note when information might be outdated.
