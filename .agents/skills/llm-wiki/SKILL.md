---
name: llm-wiki
description: >
  Maintain a compounding LLM-wiki knowledge base following the Karpathy pattern — a directory of
  interlinked Markdown files where the LLM acts as the programmer and the wiki is the codebase.
  Use this skill whenever the user wants to ingest a new source into their wiki, query the wiki
  for a synthesized answer, or run a lint pass to clean up contradictions, orphans, and stale
  claims. Trigger on: "/llm-wiki", "add this to the wiki", "ingest this", "update the wiki with",
  "query the wiki about", "lint the wiki", "update my wiki", "add to my knowledge base", "put this
  in the wiki", "what does my wiki say about", or when the user shares a URL/file alongside any
  mention of a wiki or knowledge base. Also trigger on "/llm-wiki init" when the user wants
  to bootstrap a new wiki from scratch.
version: "1.0.0"
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
  - Glob
  - Grep
  - WebFetch
  - WebSearch
  - AskUserQuestion
---

# llm-wiki — LLM Wiki Maintainer

A compounding personal knowledge base. Three layers:

1. **Raw sources** (immutable) — articles, PDFs, transcripts, notes. The LLM reads; never writes.
2. **Wiki** — a directory of interlinked `.md` files, synthesized and maintained by the LLM.
3. **Schema** — an `AGENTS.md` in the wiki root defining its conventions and domain.

> "Obsidian is the IDE; the LLM is the programmer; the wiki is the codebase." — Karpathy

Reference: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f

---

## Step 0 — Locate the wiki root

Determine `WIKI_ROOT` before doing anything else:

1. **Explicit path** — user passed a path argument (e.g., `/llm-wiki ~/wiki ingest …`). Use it.
   - If the path exists and contains `index.md` + `log.md`: confirmed, proceed.
   - If the path does not exist (or is empty): offer to initialize it there, then run **Init flow**.
2. **Session context** — a wiki root was established earlier in this conversation. Use it.
3. **Auto-discover** — check for an `index.md` + `log.md` pair in:
   - Current working directory
   - `~/wiki`
   - `~/notes/wiki`
   - `~/obsidian`
   - `~/Documents/wiki`
4. **Not found** — ask the user once:

```
Where does your wiki live? (path, or "new" to initialize one)
```

If the user says "new" or provides a path with no existing wiki, run the **Init flow** below, then proceed with the requested operation.

---

## Init flow — bootstrapping a new wiki

When `WIKI_ROOT` doesn't exist yet:

1. Confirm the directory path with the user.
2. Create the directory.
3. Write `AGENTS.md` using the template in **Appendix A**.
4. Write `index.md`:

```markdown
# Wiki Index

| Page | Summary | Tags |
| ---- | ------- | ---- |
```

5. Write `log.md`:

```markdown
# Wiki Log

<!-- Append-only. Never edit existing entries. -->
```

6. Append the init entry to `log.md`:

```
## [YYYY-MM-DD] init | wiki initialized at <WIKI_ROOT>
```

---

## Operation: ingest

**Trigger phrases:** "ingest this", "add this to the wiki", "add to my wiki", "update the wiki with", "put this in my knowledge base", user shares URL/file + mentions wiki.

### Steps

1. **Read the source.**
   - If it's a URL: `WebFetch` it.
   - If it's a file path: `Read` it. Never write to it.

2. **Discuss takeaways.** Briefly surface (in chat) the 3–5 key ideas from the source. This is the human's chance to redirect focus before writes happen.

3. **Determine the slug.** Derive a kebab-case slug from the source title (e.g., `transformer-architecture`, `karpathy-llm-wiki`). This becomes `WIKI_ROOT/<slug>.md`.

4. **Write the source summary page.**

   ```markdown
   ---
   title: <Title>
   tags: [<relevant>, <keywords>]
   sources: [<slug-of-raw-source>]
   updated: <YYYY-MM-DD>
   ---

   # <Title>

   <2–4 paragraph synthesis of the source. Cross-reference related wiki pages with [[Wikilinks]].
   Focus on insight, not transcription.>

   ## Key Points

   - …

   ## Sources

   - `<raw-source-path-or-url>`
   ```

5. **Propagate to related pages.** Read `index.md` to identify pages that relate to this source's topics. For each related page:
   - Read the page.
   - Determine what new information this source adds, contradicts, or confirms.
   - Edit the page to integrate the new information (new section, updated claim, new cross-reference). A single ingest commonly touches 5–15 pages.

6. **Update `index.md`.** Add a row for the new page:

   ```markdown
   | [[<slug>]] | <one-line summary> | <tags> |
   ```

7. **Append to `log.md`:**

   ```
   ## [YYYY-MM-DD] ingest | <title>
   ```

---

## Operation: query

**Trigger phrases:** "query the wiki", "what does my wiki say about", "ask the wiki", "look up in the wiki", "search the wiki for".

### Steps

1. **Read `index.md`** to get a map of the wiki. Identify the 3–6 pages most relevant to the question.

2. **Read those pages.** If a page links to others via `[[Wikilink]]`, follow the links if they seem relevant.

3. **Synthesize an answer.** Write a cited response that traces every claim back to a specific wiki page (and transitively to a raw source). Format:

   > <Answer paragraph.> ([[Page Name]])

4. **Write-back rule.** If the synthesis produced a genuinely new insight, comparison table, or map that isn't captured anywhere in the wiki — write it as a new page (slug: `synthesis-<topic>-<date>`). Don't write back routine answers.

5. **Append to `log.md`:**

   ```
   ## [YYYY-MM-DD] query | <question (truncated to 80 chars)>
   ```

---

## Operation: lint

**Trigger phrases:** "lint the wiki", "clean up the wiki", "audit the wiki", "check the wiki for stale content".

### Steps

1. **Inventory the wiki.** `Glob("**/*.md", WIKI_ROOT)` to get all pages. Read `index.md`.

2. **Check each page for:**
   - **Orphans** — pages with no inbound `[[wikilink]]` from any other page. List them.
   - **Contradictions** — claims that conflict with claims in other pages (require reading both).
   - **Stale claims** — statements marked with a date or tied to a source where a newer source supersedes them.
   - **Missing cross-references** — a page mentions a concept that has its own page but doesn't link to it.
   - **Index gaps** — pages not listed in `index.md`.

3. **Fix what you can:**
   - Add missing `[[wikilinks]]` inline.
   - Add orphans to `index.md` if missing.
   - Mark contradictions with a `> ⚠️ Contradiction with [[Other Page]] — needs resolution.` callout.
   - Update stale claims if a newer source in the wiki provides the correct information.

4. **Report what you cannot fix** — contradictions that require human judgment, gaps addressable only by new sources.

5. **Append to `log.md`:**

   ```
   ## [YYYY-MM-DD] lint | <N pages checked, M issues found, K fixed>
   ```

---

## Wiki conventions

| Convention            | Rule                                                                                       |
| --------------------- | ------------------------------------------------------------------------------------------ |
| **Page slugs**        | `kebab-case.md` derived from the concept title                                             |
| **Frontmatter**       | Every page must have `title`, `tags`, `sources`, `updated`                                 |
| **Cross-references**  | Use `[[Page Slug]]` wikilinks. Never use relative paths.                                   |
| **Sources section**   | Every page ends with `## Sources` listing all raw docs it was derived from                 |
| **Raw sources**       | Always read-only. LLM never writes to source files.                                        |
| **index.md**          | Updated on every write operation. One row per page.                                        |
| **log.md**            | Append-only. Never rewrite or delete entries.                                              |
| **Synthesis over transcription** | Pages should integrate and cross-reference, not just summarize one source.    |
| **Dates**             | `date -u +%Y-%m-%d` at execution time                                                      |

---

## Rules

- **Never modify raw sources.** If the user points at a source file, read it; never edit it.
- **Propagate aggressively.** Treat ingest like a code refactor — ripple changes to all affected pages.
- **Index is always current.** Every write operation must update `index.md` before finishing.
- **Log is append-only.** Always append new entries at the bottom. Never rewrite or delete existing entries.
- **Cite in every page.** Every claim traceable to a source must say which source.
- **Synthesize, don't transcribe.** A wiki page that just copies a source verbatim has no value.
- **One wiki per session.** If the user mentions a second wiki location, confirm before switching `WIKI_ROOT`.

---

## Appendix A — `AGENTS.md` template

Write this to `WIKI_ROOT/AGENTS.md` when initializing a new wiki. The user should edit it to fill in their domain.

```markdown
# Wiki Schema

This wiki is maintained by an LLM using the llm-wiki skill
(https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f).

## Domain

<What is this wiki about? What topics does it cover?>

## Conventions

- **Page slugs**: kebab-case (e.g., `transformer-architecture.md`)
- **Frontmatter**: `title`, `tags`, `sources` (raw source slugs), `updated` (YYYY-MM-DD)
- **Cross-references**: `[[Page Slug]]` wikilinks
- **Sources section**: every page ends with `## Sources` listing its raw inputs

## Operations

Run these via the `llm-wiki` skill:

- `ingest <source>` — read a new source, write a summary page, propagate to related pages
- `query <question>` — synthesize an answer from wiki pages, optionally write back
- `lint` — audit for orphans, contradictions, stale claims, missing links

## Raw Sources

Raw source files live in `raw/`. They are immutable — the LLM reads them but never writes to them.

## index.md

Structured catalog of all wiki pages. Updated on every write operation.

## log.md

Append-only chronological log. Format: `## [YYYY-MM-DD] operation | detail`
```
