---
description: Create a GitHub issue following the Spec Driven Development feature request template.
---

# Create SDD-Compliant GitHub Issue

## User Input

    $ARGUMENTS

You **MUST** consider the user input before proceeding (if not empty).

## Goal

Create a new GitHub issue that adheres to the `.github/ISSUE_TEMPLATE/feature_request_spec_driven.md` template. This ensures the issue is ready for automated implementation via Spec Driven Development (SDD).

## Outline

### 1. Identify Repository (Speed First)

* Determine the repository `owner` and `repo` name.
* **Method A (Fastest)**: Run `gh repo view --json owner,name` and parse the JSON.
* **Method B (Fallback)**: Run `git remote get-url origin` and parse the `owner/repo` from the URL.

### 2. Parse User Input

* Extract the **Title** and **Core Objective** from `$ARGUMENTS`.
* If `$ARGUMENTS` is too brief, search the codebase or recent context to understand the intent.
* **Clarification**: If the intent is still ambiguous, ask the user for more details before drafting.

### 3. Draft Issue Content

* Read the template file: `.github/ISSUE_TEMPLATE/feature_request_spec_driven.md`.
* Populate the template fields:
  * **🎯 Objective**: Clear, high-level goal.
  * **📖 Context & Reasoning**: Why is this being built? How does it align with the Memex Constitution?
  * **🛠 Technical Specification**: Initial thoughts on target files, expected logic, and data model changes.
  * **✅ Requirements**: Actionable acceptance criteria (checkboxes).
  * **🧪 Verification Plan**: How to test (automated & manual).
  * **⚠️ Constraints & Edge Cases**: Mention Go-only, single-binary, and other constitutional invariants.

### 4. Create the Issue

* **Attempt 1: GitHub MCP Server**
  * Use the `issue_write` tool from the `user-github` MCP server.
  * `method`: "create"
  * `owner`: [identified owner]
  * `repo`: [identified repo]
  * `title`: [drafted title]
  * `body`: [drafted body from template]
* **Attempt 2: GitHub CLI (Fallback)**
  * If the MCP tool fails (e.g., 403 or 401), run:
      `gh issue create --title "[title]" --body "[body]"`

### 5. Final Display

* **CRITICAL**: Prominently display the URL of the created issue.
* Suggest the next step: `To implement this, run: /implement-issue-with-sdd #[issue_number]`

## Best Practices for Company-of-One

* **Structured Metadata**: Ensure requirements are in a checkbox list so agents can track progress.
* **Constitution Alignment**: Always reference how the feature respects the core principles (Single Binary, Go-Only, etc.).
* **Clear Boundaries**: Define what is *out of scope* to prevent agent drift.
