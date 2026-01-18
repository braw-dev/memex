---
description: "Task list for HTTP Reverse Proxy implementation"
---

# Tasks: HTTP Reverse Proxy

**Input**: Design documents from `/specs/001-http-proxy/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are included to verify each user story independently as defined in the "Independent Test" sections of the spec.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project directories (`cmd/proxy`, `internal/proxy`, `pkg/types`, `tests/integration`, `tests/unit`)
- [x] T002 [P] Verify `go.mod` dependencies (ensure `knadh/koanf/v2` is available)
- [x] T003 Create `main.go` entrypoint skeleton in `cmd/proxy/main.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Create `ProxyConfig` struct and `ConfigLoader` interface in `internal/proxy/config.go`
- [x] T005 Implement configuration loading using `koanf` in `internal/proxy/config.go`
- [x] T006 [P] Define `ProxyRequest`, `ProxyResponse`, and `SchemaType` types in `pkg/types/proxy.go`
- [x] T007 Create basic `ProxyHandler` structure and `NewProxyHandler` factory in `internal/proxy/handler.go`
- [x] T008 Integrate config loading into `cmd/proxy/main.go` to initialize server

**Checkpoint**: Foundation ready - basic server structure exists and can load config.

---

## Phase 3: User Story 1 - Transparent Passthrough (Priority: P1) 🎯 MVP

**Goal**: Forward traffic unchanged if no AI schema detected (foundational capability).

**Independent Test**: Configure HTTP client to use proxy; verify non-AI requests reach destination unchanged.

### Implementation for User Story 1

- [x] T009 [US1] Implement `ServeHTTP` using `httputil.ReverseProxy` in `internal/proxy/handler.go`
- [x] T010 [US1] Implement `HandleCONNECT` for HTTPS tunneling in `internal/proxy/handler.go`
- [x] T011 [US1] Configure `ReverseProxy` Director to handle standard headers and upstream URL parsing in `internal/proxy/handler.go`
- [x] T012 [US1] Create integration test for HTTP passthrough in `tests/integration/proxy_test.go`
- [x] T013 [US1] Create integration test for HTTPS CONNECT tunneling in `tests/integration/proxy_test.go`

**Checkpoint**: At this point, the proxy acts as a standard transparent proxy.

---

## Phase 4: User Story 2 & 3 - Schema Detection (Priority: P2)

**Goal**: Detect Anthropic (`/v1/messages`) and OpenAI (`/v1/chat/completions`) requests.

**Independent Test**: Send requests to specific paths; verify proxy identifies schema type (via logs/debug).

### Implementation for User Stories 2 & 3

- [x] T014 [US2] [US3] Implement `SchemaDetector` logic (URL path matching) in `internal/proxy/schema.go`
- [x] T015 [US2] [US3] Update `ServeHTTP` to use `SchemaDetector` and store `SchemaType` in request context in `internal/proxy/handler.go`
- [x] T016 [P] [US2] [US3] Add debug logging to print detected schema in `internal/proxy/handler.go` (guarded by config.Debug)
- [x] T017 [P] [US2] [US3] Create unit tests for `SchemaDetector` in `tests/unit/schema_test.go`

**Checkpoint**: Proxy now correctly identifies and logs AI traffic types.

---

## Phase 5: User Story 4 - Standard Proxy Configuration (Priority: P3)

**Goal**: Support standard proxy environment variables and headers.

**Independent Test**: Verify `HTTP_PROXY` env vars work with standard clients.

### Implementation for User Story 4

- [x] T018 [US4] Verify and ensure `Proxy-Authorization` header handling in `internal/proxy/handler.go`
- [x] T019 [US4] Add documentation on using `HTTP_PROXY`/`HTTPS_PROXY` env vars to `README.md` or `quickstart.md`
- [x] T020 [P] [US4] Add integration test case verifying proxy works when client uses standard proxy headers in `tests/integration/proxy_test.go`

**Checkpoint**: Proxy is fully compliant with standard client configuration methods.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T021 [P] Implement graceful shutdown handling in `cmd/proxy/main.go`
- [x] T022 Ensure error handling (502/504) returns JSON responses where appropriate in `internal/proxy/errors.go`
- [x] T023 Run performance benchmark to verify <5ms overhead (simple latency check script)
- [x] T024 Validate `quickstart.md` instructions against implementation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Phase 1. Blocks all stories.
- **US1 (Phase 3)**: Depends on Phase 2. MVP.
- **US2/3 (Phase 4)**: Depends on Phase 3 (needs working proxy to add detection).
- **US4 (Phase 5)**: Independent of US2/3, depends on Phase 3.
- **Polish (Phase 6)**: Run after feature completion.

### Parallel Opportunities

- T002 (Verify deps) can run with T001.
- T006 (Define types) can run with T004/T005 (Config).
- T016 (Debug logging) and T017 (Unit tests) can run in parallel with detection logic.
- T020 (Header tests) can run in parallel with docs.

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 & 2.
2. Complete Phase 3 (Passthrough).
3. **STOP and VALIDATE**: Verify generic HTTP/HTTPS browsing works through proxy.

### Incremental Delivery

1. MVP (Passthrough) -> Deployable as generic proxy.
2. Add Schema Detection -> Enabler for future caching features.
3. Polish -> Production readiness.
