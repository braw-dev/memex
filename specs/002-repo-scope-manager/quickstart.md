# Quickstart: Verification Guide

## Prerequisites

- Build the binary: `go build -o memex ./cmd/proxy`

## Scenario 1: Git Repository (Standard)

1. Navigate to a directory that is a Git repository with a remote `origin`.

   ```bash
   cd /path/to/my-repo
   git remote -v # verify origin exists
   ```

2. Run Memex.

   ```bash
   ./memex
   ```

3. Check logs (if implemented) or verification output.
   - Expected: `Detected Scope: GitRemote <url>`

## Scenario 2: Non-Git Directory (Fallback)

1. Create a temporary directory.

   ```bash
   mkdir /tmp/memex-test
   cd /tmp/memex-test
   ```

2. Run Memex.

   ```bash
   /path/to/memex
   ```

3. Check logs.
   - Expected: `Detected Scope: PathHash <sha256>`

## Scenario 3: Isolation Test

1. Start Memex in Repo A. Make a request (mocked).
2. Start Memex in Repo B. Make the *same* request.
3. Verify that the second request does not hit the cache from the first request (if you can inspect the store).
   - *Note*: Without `!reset` implemented yet, manual store inspection might be needed, or trusting the logging output.
