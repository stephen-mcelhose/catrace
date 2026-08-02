# How to Use and Create Skills in csgdaa-code (Updated)

**Objective:** Learn how to extend your AI assistant's capabilities by navigating the existing skill library and collaborating with the assistant to build new skills.

## Prerequisites
*   **Active csgdaa-code session:** You are currently in a chat interface.
*   **Workspace:** You should have a project open where you want to add functionality.

---

## 1. Discovering Available Skills
To see what your assistant can already do, check the library.

1.  Type `/skills` in the chat input and press **Enter**.
2.  **Observe the Output:** A table will appear showing the Name, Source, and Description of available skills.

**How to Verify:** If you see a list of tools like `answers`, `commit`, and `review`, the command was successful.

---

## 2. Using a Skill
You can trigger a skill directly by name or let the AI auto-trigger it based on your request.

1.  **Direct Trigger:** Type `/[skill-name]` (e.g., `/commit`). The agent will confirm: "Activated skill: [name]".
2.  **Natural Language Trigger:** Simply ask for a task (e.g., "Review my code"). If the skill `description` matches, the AI will activate it automatically.

---

## 3. Creating a Custom Skill (Collaborative Process)
In csgdaa-code, skill creation is a collaborative process where the AI performs the file operations based on your requirements.

### Step 1: Invoke the Skill Creator
Type `/skill-creator` and press **Enter**. This puts the AI into a specialized mode with the tools and knowledge needed to build skills.

### Step 2: Describe Your Goal
Tell the AI what you want the skill to do. 
*Example: "I want a how-to-writer skill that will reference rubric in order to help guide a user to create how-to documentation"*

### Step 3: AI-Led Implementation
The AI will then execute the necessary steps for you:
1.  **Create Directories:** It will set up the folder structure (e.g., `.csgdaa-code/skills/famous-quotes/scripts`).
2.  **Create Scripts:** It will write any necessary logic (e.g., a `curl` command in a `.sh` script).
3.  **Configure Metadata:** It will create the `SKILL.md` file with the correct name, description, and `allowed-tools`.

### Step 4: Verification
Run `/skills` again. Your new skill should now appear in the list. To test it, call it by name (e.g., `/how-to-writer`).

---

## 4. Advanced Example: Web-Enabled Skill
Skills can do more than just text processing. They can interact with the web or local system via existing commands.

**Example `SKILL.md` for a web call:**
```markdown
---
name: famous-quotes
description: Fetch a random famous quote from an online API. Use when the user asks for inspiration, a quote, or says "tell me something profound".
version: "0.1.1"
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
---

# Famous Quotes

This skill fetches a random inspirational quote from the ZenQuotes API using direct system tools.

## Workflow

1.  **Fetch Quote**: Use `curl` to call the ZenQuotes API.
2.  **Process Output**: Use `jq` to format the JSON response into a readable quote.
3.  **Present to User**: Display the formatted quote and author.

## Implementation

The following command pipeline is used to retrieve the quote:

```bash
curl -s https://zenquotes.io/api/random | jq -r '.[0] | "\"\(.q)\" — \(.a)"'
```
```

---

## Troubleshooting Common Issues

| Issue | Solution |
| :--- | :--- |
| **Skill doesn't appear in `/skills`** | Ensure the `SKILL.md` file exists and has valid YAML frontmatter. |
| **Permission Denied** | If your skill uses a script, the AI must run `chmod +x` on it during creation. |
| **Agent doesn't use the skill** | The `description` in `SKILL.md` is critical. It must clearly describe *when* the AI should use the skill. |
