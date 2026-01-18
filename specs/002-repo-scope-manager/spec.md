# Feature Specification: Repo-Aware Scope Manager

**Feature Branch**: `002-repo-scope-manager`  
**Created**: 2026-01-18  
**Status**: Draft  
**Input**: User description: "Implement a mechanism to detect the 'Project Context' of the active session to isolate cache data (Multi-Tenancy on Localhost)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Git Project Isolation (Priority: P1)

As a developer working on multiple projects, I want my cache hits to be isolated to the current project so that I don't get irrelevant or wrong code suggestions from other projects.

**Why this priority**: This is the core value proposition of the feature. Without it, Memex is unsafe to use across multiple client projects.

**Independent Test**: Can be fully tested by initializing Memex in two different git repositories and verifying cache isolation.

**Acceptance Scenarios**:

1. **Given** Memex is running in Repo A, **When** a prompt is cached, **Then** the cache key includes Repo A's remote origin URL.
2. **Given** Memex is running in Repo B, **When** the same prompt is sent, **Then** it is a cache MISS (unless also cached in B).

---

### User Story 2 - Non-Git Project Isolation (Priority: P2)

As a developer working in a non-git directory, I want my cache to still be isolated by directory path so I can use Memex safely without git.

**Why this priority**: Ensures Memex works safely in scratchpads or early-stage projects before git initialization.

**Independent Test**: Can be tested by running Memex in a directory without a `.git` folder.

**Acceptance Scenarios**:

1. **Given** Memex is running in a directory without git, **When** it initializes, **Then** it uses the hashed absolute path as the Scope ID.
2. **Given** Memex is moved to a new directory, **When** it runs, **Then** it generates a different Scope ID.

---

### Edge Cases

- What happens when the git remote is changed while Memex is running? (Assume restart required or re-check on request?) -> Assumption: Scope is determined at startup or per-request if cheap.
- What happens if the directory is a git repo but has no remote? -> Should fall back to path hash or local git sha? -> Assumption: Use path hash if no remote.
- What happens if the user has no read permissions on `.git/config`? -> Fallback to path hash.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST attempt to detect the Git Remote Origin URL of the current working directory on initialization.
- **FR-002**: If a Git Remote Origin URL is found, System MUST use it as the Scope ID.
- **FR-003**: If no Git Remote Origin is found (or not a git repo), System MUST use a unique cryptographic hash of the absolute path of the current working directory as the Scope ID.
- **FR-004**: System MUST include the Scope ID in every cache key generation process (as a salt or prefix).
- **FR-005**: System MUST log the detected Scope ID on startup for debugging (masked if necessary, though public repos are public).

### Key Entities *(include if feature involves data)*

- **Scope Context**: Represents the unique identity of the current project.
  - `ScopeID`: String (Remote URL or Path Hash)
  - `Source`: Enum (GitRemote, PathHash)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of generated cache keys include the correct Scope ID for the active environment.
- **SC-002**: Identical prompts sent from two different project contexts result in 0% cache collision (complete isolation).
- **SC-003**: Scope detection adds < 5ms latency to the startup process (or first request).
