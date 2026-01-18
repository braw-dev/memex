# Feature Specification: Initialize Embedded Storage Engines

**Feature Branch**: `004-init-embedded-storage`  
**Created**: 2026-01-18  
**Status**: Draft  
**Input**: Store: Initialize Embedded BadgerDB and DuckDB Storage (Issue #4)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Application Initialization with Fresh Storage (Priority: P1)

As a user running the Memex application, I want the storage system to initialize automatically without configuration so that I can start using the application immediately.

**Why this priority**: Essential for application startup. Without storage, the app cannot function.

**Independent Test**: Can be tested by deleting `.memex` folder and running the initialization code.

**Acceptance Scenarios**:

1. **Given** no `.memex` directory exists in the project root, **When** the application storage layer initializes, **Then** a `.memex` directory is created.
2. **Given** the `.memex` directory is created, **When** initialization completes, **Then** valid BadgerDB data files exist within it.
3. **Given** the `.memex` directory is created, **When** initialization completes, **Then** valid DuckDB data files exist within it.
4. **Given** the user has no external databases installed, **When** the application starts, **Then** no errors related to missing database connections occur.

---

### User Story 2 - Application Restart with Existing Storage (Priority: P2)

As a user restarting the Memex application, I want existing data to be preserved and accessible so that my cache and logs are persistent.

**Why this priority**: Persistence is the core requirement of this feature.

**Independent Test**: Can be tested by writing data, restarting the app, and reading the data back.

**Acceptance Scenarios**:

1. **Given** an existing `.memex` directory with valid data, **When** the application storage layer initializes, **Then** it opens the existing databases without error.
2. **Given** an existing `.memex` directory, **When** the application starts, **Then** it does not overwrite or lose existing data.

---

### User Story 3 - Graceful Failure Handling (Priority: P3)

As a user with a misconfigured environment (e.g., permissions issues), I want clear error messages if storage cannot be initialized so that I can fix the issue.

**Why this priority**: Improves user experience and debuggability.

**Independent Test**: Can be tested by removing write permissions from the project root and running initialization.

**Acceptance Scenarios**:

1. **Given** the project root is read-only, **When** the application storage layer initializes, **Then** it returns a descriptive error about permission denied.
2. **Given** the database files are locked by another process, **When** the application initializes, **Then** it returns a descriptive error about the lock.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST initialize a hidden `.memex` directory in the project root (or configured location) if it does not exist.
- **FR-002**: System MUST initialize BadgerDB v4 for key-value storage within the `.memex` directory.
- **FR-003**: System MUST initialize DuckDB for OLAP/logging within the `.memex` directory.
- **FR-004**: System MUST NOT require any external database processes (like PostgreSQL or Redis) to be running.
- **FR-005**: System MUST configure BadgerDB options to be optimized for large binary payloads (code blocks).
- **FR-006**: System MUST configure DuckDB settings appropriate for structured event logging.
- **FR-007**: System MUST provide a unified mechanism to close both storage engines gracefully on application shutdown.

### Key Entities *(include if feature involves data)*

- **Store**: The abstraction layer managing lifecycle of both embedded databases.
- **KVStore**: The interface/wrapper for BadgerDB operations.
- **LogStore**: The interface/wrapper for DuckDB operations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Storage initialization completes in under 500ms on standard hardware (excluding first-time creation overhead if significant).
- **SC-002**: Application successfully starts and writes a test record to both BadgerDB and DuckDB in a clean environment without user intervention.
- **SC-003**: Database files are contained entirely within the `.memex` directory.
- **SC-004**: Zero external dependencies required (verified by running in a clean container/environment).
