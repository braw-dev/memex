# Tasks: Repo-Aware Scope Manager

**Input**: Design documents from `/specs/002-repo-scope-manager/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Add go-git/v5 dependency to go.mod
- [x] T002 Create ScopeContext entity in pkg/types/scope.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [x] T003 Create internal/proxy/scope.go with empty DetectScope function stub
- [x] T004 Create internal/proxy/middleware.go with ScopeMiddleware skeleton

---

## Phase 3: User Story 1 - Git Project Isolation (Priority: P1) 🎯 MVP

**Goal**: Identify the current project by its Git Remote Origin URL

**Independent Test**: Verify that running in a git repo produces a ScopeID matching the remote URL.

### Tests for User Story 1

- [x] T005 [P] [US1] Create unit tests for Git detection in tests/unit/scope_test.go

### Implementation for User Story 1

- [x] T006 [US1] Implement DetectScope with go-git logic in internal/proxy/scope.go
- [x] T007 [US1] Implement ScopeMiddleware to inject ScopeContext into request context in internal/proxy/middleware.go
- [x] T008 [US1] Register ScopeMiddleware in internal/proxy/handler.go (NewServer function)
- [x] T009 [US1] Log detected scope on startup/request for verification

**Checkpoint**: At this point, running the proxy in a git repo should show the remote URL in logs (if debug enabled)

---

## Phase 4: User Story 2 - Non-Git Project Isolation (Priority: P2)

**Goal**: Fallback to path hashing when no git remote is available

**Independent Test**: Verify that running in a non-git directory produces a SHA-256 hash of the path.

### Tests for User Story 2

- [x] T010 [P] [US2] Add unit tests for path fallback in tests/unit/scope_test.go

### Implementation for User Story 2

- [x] T011 [US2] Implement fallback logic in DetectScope in internal/proxy/scope.go
- [x] T012 [US2] Verify hash generation uses SHA-256 (crypto/sha256)

**Checkpoint**: At this point, running in any directory should produce a stable, isolated ScopeID

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T013 Verify compliance with Constitution Principle V (Middleware Order)
- [x] T014 Run quickstart verification scenarios
