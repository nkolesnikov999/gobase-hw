package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
)

func withTempWorkDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "cli-tests-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		os.RemoveAll(dir)
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(orig)
		_ = os.RemoveAll(dir)
	})
	return dir
}

func writeTempJSON(t *testing.T, dir, name string, data any) string {
	t.Helper()
	path := filepath.Join(dir, name)
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}
	return path
}

func readBinsFile(t *testing.T) BinsList {
	t.Helper()
	b, err := os.ReadFile(binsFileName)
	if err != nil {
		t.Fatalf("failed to read bins file: %v", err)
	}
	var bl BinsList
	if err := json.Unmarshal(b, &bl); err != nil {
		t.Fatalf("failed to unmarshal bins list: %v", err)
	}
	return bl
}

func TestCLI_CreateBin(t *testing.T) {
	work := withTempWorkDir(t)
	_ = work

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const binID = "cli-create-1"
	httpmock.RegisterResponder("POST", "https://api.jsonbin.io/v3/b", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(201, fmt.Sprintf(`{"record":{"a":1},"metadata":{"id":"%s","createdAt":"2024-01-01T00:00:00Z","private":false}}`, binID)), nil
	})

	jsonFile := writeTempJSON(t, ".", "payload.json", map[string]any{"a": 1})
	createBin(jsonFile, "test-bin")

	bl := readBinsFile(t)
	if len(bl.Bins) != 1 || bl.Bins[0].ID != binID || bl.Bins[0].Name != "test-bin" {
		t.Fatalf("unexpected bins.json content: %+v", bl)
	}
	// cleanup bin list
	if err := removeBinFromList(binID); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}

func TestCLI_GetBin(t *testing.T) {
	_ = withTempWorkDir(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const binID = "cli-get-1"
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, fmt.Sprintf(`{"record":{"name":"alpha"},"metadata":{"id":"%s"}}`, binID)), nil
	})

	getBin(binID)
}

func TestCLI_UpdateBin(t *testing.T) {
	_ = withTempWorkDir(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const binID = "cli-update-1"
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, fmt.Sprintf(`{"record":{"name":"beta"},"metadata":{"id":"%s","createdAt":"2024-01-01T00:00:00Z","private":false}}`, binID)), nil
	})

	jsonFile := writeTempJSON(t, ".", "payload.json", map[string]any{"name": "beta"})
	updateBin(binID, jsonFile)
}

func TestCLI_DeleteBin(t *testing.T) {
	_ = withTempWorkDir(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const binID = "cli-delete-1"
	// Prepare local bins.json containing the bin
	bl := BinsList{Bins: []BinRecord{{ID: binID, Name: "to-delete"}}}
	if b, err := json.MarshalIndent(bl, "", "  "); err != nil {
		t.Fatalf("marshal bins list: %v", err)
	} else if err := os.WriteFile(binsFileName, b, 0644); err != nil {
		t.Fatalf("write bins file: %v", err)
	}

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, fmt.Sprintf(`{"metadata":{"id":"%s","versionsDeleted":1},"message":"Bin deleted successfully"}`, binID)), nil
	})

	deleteBin(binID)
	bl2 := readBinsFile(t)
	if len(bl2.Bins) != 0 {
		t.Fatalf("expected bins list to be empty after delete, got: %+v", bl2)
	}
}

func TestCLI_CRUD_Grouped(t *testing.T) {
	_ = withTempWorkDir(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const binID = "cli-group-1"

	// Create
	httpmock.RegisterResponder("POST", "https://api.jsonbin.io/v3/b", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(201, fmt.Sprintf(`{"record":{"a":1},"metadata":{"id":"%s","createdAt":"2024-01-01T00:00:00Z","private":false}}`, binID)), nil
	})
	jsonFile := writeTempJSON(t, ".", "create.json", map[string]any{"a": 1})
	createBin(jsonFile, "group-bin")

	// Get
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, fmt.Sprintf(`{"record":{"a":1},"metadata":{"id":"%s"}}`, binID)), nil
	})
	getBin(binID)

	// Update
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, fmt.Sprintf(`{"record":{"a":2},"metadata":{"id":"%s","createdAt":"2024-01-01T00:00:00Z","private":false}}`, binID)), nil
	})
	updFile := writeTempJSON(t, ".", "update.json", map[string]any{"a": 2})
	updateBin(binID, updFile)

	// Delete
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, fmt.Sprintf(`{"metadata":{"id":"%s","versionsDeleted":1},"message":"Bin deleted successfully"}`, binID)), nil
	})
	deleteBin(binID)

	// Ensure bins.json cleaned
	if _, err := os.Stat(binsFileName); err == nil {
		b := readBinsFile(t)
		if len(b.Bins) != 0 {
			t.Fatalf("expected no bins after cleanup, got: %+v", b)
		}
	}
}
