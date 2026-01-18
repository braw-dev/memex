package unit

import (
	"os"
	"testing"

	"github.com/braw-dev/memex/internal/proxy"
	"github.com/braw-dev/memex/pkg/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func TestDetectScope_GitRemote(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "memex-scope-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Init git repo
	r, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatal(err)
	}

	// Create remote
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/example/repo.git"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test DetectScope
	scope, err := proxy.DetectScope(tmpDir)
	if err != nil {
		t.Fatalf("DetectScope failed: %v", err)
	}
	if scope == nil {
		t.Fatal("Expected scope to be returned, got nil")
	}

	if scope.Type != types.ScopeTypeGitRemote {
		t.Errorf("Expected ScopeTypeGitRemote, got %v", scope.Type)
	}
	if scope.ID != "https://github.com/example/repo.git" {
		t.Errorf("Expected ID 'https://github.com/example/repo.git', got '%s'", scope.ID)
	}
}

func TestDetectScope_NoGit(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "memex-scope-test-nogit")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test DetectScope
	scope, err := proxy.DetectScope(tmpDir)
	if err != nil {
		t.Fatalf("DetectScope failed: %v", err)
	}
	if scope == nil {
		t.Fatal("Expected scope to be returned, got nil")
	}

	if scope.Type != types.ScopeTypePathHash {
		t.Errorf("Expected ScopeTypePathHash, got %v", scope.Type)
	}
	if scope.ID == "" {
		t.Error("Expected ID to be non-empty")
	}
}
