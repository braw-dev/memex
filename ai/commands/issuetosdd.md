---
description: Implement a GitHub issue using the full Spec Driven Development workflow (Specify -> Plan -> Tasks -> Implement).
---

# GitHub Issue to Implementation with Spec Driven Development

## User Input

    $ARGUMENTS

You **MUST** consider the user input before proceeding (if not empty).

## Goal

Take a GitHub issue number and optional context, validate the issue content, and then orchestrate the complete Spec Driven Development (SDD) lifecycle to implement the requested feature.

## Outline

### Parse Input

* Extract the **Issue Number** from `$ARGUMENTS` (look for formats like `#123` or just `123`).
* If there is no issue number, check to see if the user is asking for "next unassigned" or similar.
* Treat the rest of `$ARGUMENTS` as **Additional Context**.
* If no issue number is found or user has not asked for next unassigned issue or there are no unassigned issues, **ERROR**: "Please provide a GitHub issue number (e.g., #123 or just 123) to proceed or a valid issue to work on."

### Fetch Issue Details

* Get the current repository `owner` and `repo` by running:

    `git config --get remote.origin.url`

* Parse the output (e.g., `https://github.com/owner/repo.git` -> owner: `owner`, repo: `repo`).
* Call the `issue_read` tool from the `user-github` MCP server:
  * `server_name`: `user-github` (or verify via `list_tools`)
  * `tool_name`: `issue_read`
  * `arguments`:
  * `method`: "get"
  * `owner`: [parsed owner]
  * `repo`: [parsed repo]
  * `issue_number`: [parsed number (as integer)]
    * **Error Handling**: If the tool fails or issue is not found, stop and report the error to the user.

### Validate & Clarify

* Analyze the **Issue Title**, **Issue Body**, and **Additional Context**.
* **Assessment**: Does this combined context provide a clear goal and sufficient requirements to build a specification?
  * *Clear Goal*: Is it obvious what needs to be built or fixed?
  * *Scope*: Are the boundaries reasonably defined?
* **Decision**:
  * **IF CLEAR**: Proceed to Step 4.
  * **IF UNCLEAR**:
  1. Identify specific missing information (e.g., "What is the expected output format?", "Which user role is this for?").
  2. **STOP and Ask**: Present the findings to the user and ask for clarification.
  3. **Wait**: Do not proceed until the user provides the missing details.
  4. **Resume**: Once answered, treat the answer as more **Additional Context** and re-validate.

### Orchestrate SDD Workflow

* Assign the GitHub issue to me
* Execute the following commands in sequence.
* **Crucial**: Pass the full context (Issue Title + Body + User Context) to the first command (`/speckit.specify`).

#### Phase 1: Specification (`/speckit.specify`)

**Action**: Run `/speckit.specify` with the argument: `"[Issue Title] [Issue Body] [User Context]"`
**Observation**: This command will create a new feature branch and generate `spec.md`.
**Verification**:

* Ensure the command completed successfully.
* Check that you are on the new branch.
* Check that `checklists/requirements.md` indicates a PASS (or at least no critical failures).
     **Commit**:
  * **SAFETY CHECK**: Verify current branch is **NOT** `main` or `master`.
  * Run `git add .`
  * Run `git commit -m "feat: initial specification from issue #[Issue Number]"`

#### Phase 2: Planning (`/speckit.plan`)

**Action**: Run `/speckit.plan`. (No arguments needed; it reads the local `spec.md`).
**Observation**: This generates `plan.md`, `data-model.md`, `research.md`.
**Verification**: Ensure `plan.md` exists and contains a valid plan.
**Commit**:

* **SAFETY CHECK**: Verify current branch is **NOT** `main` or `master`.
* Run `git add .`
* Run `git commit -m "docs: add technical plan and research"`

#### Phase 3: Task Generation (`/speckit.tasks`)

**Action**: Run `/speckit.tasks`.
**Observation**: This generates `tasks.md` based on the plan and spec.
**Verification**: Ensure `tasks.md` exists and contains checklist items.
**Commit**:

* **SAFETY CHECK**: Verify current branch is **NOT** `main` or `master`.
* Run `git add .`
* Run `git commit -m "docs: generate implementation tasks"`

#### Phase 4: Analysis (`/speckit.analyze`)

**Action**: Run `/speckit.analyze`.
**Observation**: This checks for consistency across artifacts.
**Gate**:

* If the report shows **CRITICAL** issues or **Constitution Violations**: **STOP**. Display the report and ask the user to resolve them before proceeding.
* If issues are LOW/MEDIUM: Proceed automatically.

#### Phase 5: Implementation (`/speckit.implement`)

**Action**: Run `/speckit.implement`.
**Observation**: This executes the tasks in `tasks.md`.
**Completion**: Wait for all tasks to finish.
**DO NOT COMMIT**: Wait for user approval

1. **Final Report**:
   * Summarize the work done.
   * Link to the branch name.
   * List any remaining manual steps (e.g., "Run manual verification test").
