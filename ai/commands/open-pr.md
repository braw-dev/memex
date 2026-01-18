---
description: Open a GitHub Pull Request using the repository template, linked to the active issue.
---

# Open Pull Request

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Goal

Open a Pull Request for the current branch using the repository's PR template (`.github/pull_request_template.md`), filling it out concisely based on the current implementation and linking it to the active GitHub issue.

## Outline

1. **Identify Active Issue**:
   * Determine the issue number from `$ARGUMENTS`, the current branch name, or local specification files (e.g., `specs/*/spec.md`).
   * **ERROR**: If no issue number can be identified, ask the user. If they cannot provide one, stop.

2. **Gather Context**:
   * Fetch repository information (`owner`, `repo`) from the git remote.
   * Identify the current branch name. **ERROR** if on `main` or `master`.
   * Read the target issue details using the `user-github` MCP server (`issue_read`).
   * Summarize changes by comparing the current branch to the base branch (default: `main`).

3. **Prepare PR Content**:
   * Read `.github/pull_request_template.md`.
   * Populate the template using the gathered context.
   * Ensure the "Resolution" section correctly links to the issue (e.g., `Closes #123`).

4. **Execute PR Creation**:
   * Push the current branch to the remote (`git push -u origin HEAD`).
   * Create the PR using the `user-github` MCP server (`create_pull_request`).

5. **Completion**:
   * Report the PR URL to the user.

## Key Rules

* **Mandatory Issue**: A PR **MUST** be linked to an active GitHub issue. Do not proceed without one.
* **Template Strictness**: You **MUST** use the exact structure provided in `.github/pull_request_template.md`.
* **Conciseness**: Keep the PR summary and changes list concise but informative.
* **No Main Branch**: Never attempt to open a PR from `main` or `master`.
* **Verification**: Ensure all checklist items in the template are accurately marked based on the actual state of the implementation.

## Detailed Steps

### Step 1: Issue Identification

1. Check `$ARGUMENTS` for an issue reference (e.g., `#123`, `123`).
2. Run `git branch --show-current` and look for numeric patterns.
3. Search `specs/` for any `spec.md` files that are currently being worked on.
4. If multiple issues are found, ask for clarification.
5. If none are found, ask: "I couldn't identify an active issue for this PR. Which issue number (e.g., #123) should this PR link to?"

### Step 2: Context Retrieval

1. Run `git config --get remote.origin.url` to get the repo URL.
2. Parse owner and repo name.
3. Call `user-github` `issue_read` with the identified issue number.
4. Run `git diff --stat main...HEAD` to see a summary of changed files.
5. Check for `tasks.md` in the relevant `specs/` subdirectory to understand completed work.

### Step 3: Template Population

1. Read `.github/pull_request_template.md`.
2. **Title**: Format as `[type]: [Issue Title] (#[Issue Number])`.
3. **Summary**: Describe the "why" and "what" in 2 sentences max.
4. **Changes**: List the key technical changes.
5. **Resolution**: Use the keyword `Closes` followed by the issue number.

### Step 4: Submission

1. Run `git push -u origin HEAD`.
2. Invoke `user-github` `create_pull_request`:
   * `owner`: [parsed owner]
   * `repo`: [parsed repo]
   * `title`: [generated title]
   * `body`: [filled template]
   * `head`: [current branch]
   * `base`: "main"
3. Capture the output URL and display it prominently.
