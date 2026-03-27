package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/sipeed/picoclaw/pkg/config"
)

func TestGetSystemVersionUsesPicoclawBinaryInfo(t *testing.T) {
	originalVersion := config.Version
	originalGitCommit := config.GitCommit
	originalBuildTime := config.BuildTime
	originalGoVersion := config.GoVersion
	originalFinder := findPicoclawBinaryForInfo
	originalRunner := runPicoclawVersionOutput
	t.Cleanup(func() {
		config.Version = originalVersion
		config.GitCommit = originalGitCommit
		config.BuildTime = originalBuildTime
		config.GoVersion = originalGoVersion
		findPicoclawBinaryForInfo = originalFinder
		runPicoclawVersionOutput = originalRunner
	})

	config.Version = "dev"
	config.GitCommit = ""
	config.BuildTime = ""
	config.GoVersion = ""

	findPicoclawBinaryForInfo = func() string { return "picoclaw" }
	runPicoclawVersionOutput = func(_ context.Context, _ string) (string, error) {
		return "🦞 picoclaw v1.2.3 (git: deadbeef)\n  Build: 2026-03-27T12:34:56Z\n  Go: go1.25.8\n", nil
	}

	h := NewHandler("")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/system/version", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got systemVersionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got.Version != "v1.2.3" {
		t.Fatalf("version = %q, want %q", got.Version, "v1.2.3")
	}
	if got.GitCommit != "deadbeef" {
		t.Fatalf("git_commit = %q, want %q", got.GitCommit, "deadbeef")
	}
	if got.BuildTime != "2026-03-27T12:34:56Z" {
		t.Fatalf("build_time = %q, want %q", got.BuildTime, "2026-03-27T12:34:56Z")
	}
	if got.GoVersion != "go1.25.8" {
		t.Fatalf("go_version = %q, want %q", got.GoVersion, "go1.25.8")
	}
}

func TestGetSystemVersionFallsBackToLauncherInfoWhenCommandFails(t *testing.T) {
	originalVersion := config.Version
	originalGitCommit := config.GitCommit
	originalBuildTime := config.BuildTime
	originalGoVersion := config.GoVersion
	originalFinder := findPicoclawBinaryForInfo
	originalRunner := runPicoclawVersionOutput
	t.Cleanup(func() {
		config.Version = originalVersion
		config.GitCommit = originalGitCommit
		config.BuildTime = originalBuildTime
		config.GoVersion = originalGoVersion
		findPicoclawBinaryForInfo = originalFinder
		runPicoclawVersionOutput = originalRunner
	})

	config.Version = "v9.9.9"
	config.GitCommit = "cafebabe"
	config.BuildTime = "2026-03-27T10:43:34+0000"
	config.GoVersion = "go1.25.8"

	findPicoclawBinaryForInfo = func() string { return "picoclaw" }
	runPicoclawVersionOutput = func(_ context.Context, _ string) (string, error) {
		return "", errors.New("binary unavailable")
	}

	h := NewHandler("")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/system/version", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got systemVersionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got.Version != config.Version {
		t.Fatalf("version = %q, want %q", got.Version, config.Version)
	}
	if got.GitCommit != config.GitCommit {
		t.Fatalf("git_commit = %q, want %q", got.GitCommit, config.GitCommit)
	}
	if got.BuildTime != config.BuildTime {
		t.Fatalf("build_time = %q, want %q", got.BuildTime, config.BuildTime)
	}
	if got.GoVersion != config.GoVersion {
		t.Fatalf("go_version = %q, want %q", got.GoVersion, config.GoVersion)
	}
}

func TestParsePicoclawVersionOutput(t *testing.T) {
	raw := "\u001b[1;31m████\u001b[0m\n🦞 picoclaw 18ec263 (git: 18ec2631)\n  Build: 2026-03-27T10:43:34+0000\n  Go: go1.25.8\n"
	got, ok := parsePicoclawVersionOutput(raw)
	if !ok {
		t.Fatal("parsePicoclawVersionOutput() should parse valid output")
	}
	if got.Version != "18ec263" {
		t.Fatalf("version = %q, want %q", got.Version, "18ec263")
	}
	if got.GitCommit != "18ec2631" {
		t.Fatalf("git_commit = %q, want %q", got.GitCommit, "18ec2631")
	}
	if got.BuildTime != "2026-03-27T10:43:34+0000" {
		t.Fatalf("build_time = %q, want %q", got.BuildTime, "2026-03-27T10:43:34+0000")
	}
	if got.GoVersion != "go1.25.8" {
		t.Fatalf("go_version = %q, want %q", got.GoVersion, "go1.25.8")
	}
}

func TestResolveSystemVersionInfoFallsBackRuntimeGoVersion(t *testing.T) {
	originalVersion := config.Version
	originalGitCommit := config.GitCommit
	originalBuildTime := config.BuildTime
	originalGoVersion := config.GoVersion
	originalFinder := findPicoclawBinaryForInfo
	originalRunner := runPicoclawVersionOutput
	t.Cleanup(func() {
		config.Version = originalVersion
		config.GitCommit = originalGitCommit
		config.BuildTime = originalBuildTime
		config.GoVersion = originalGoVersion
		findPicoclawBinaryForInfo = originalFinder
		runPicoclawVersionOutput = originalRunner
	})

	config.Version = "dev"
	config.GitCommit = ""
	config.BuildTime = ""
	config.GoVersion = ""

	findPicoclawBinaryForInfo = func() string { return "picoclaw" }
	runPicoclawVersionOutput = func(_ context.Context, _ string) (string, error) {
		return "picoclaw v1.0.0\n", nil
	}

	h := NewHandler("")
	got := h.resolveSystemVersionInfo()
	if got.GoVersion != runtime.Version() {
		t.Fatalf("go_version = %q, want runtime version %q", got.GoVersion, runtime.Version())
	}
}
