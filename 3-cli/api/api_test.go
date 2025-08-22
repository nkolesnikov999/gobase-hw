package api

import (
	"cli/config"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func newServiceWithKey(key string) *APIService {
	return NewAPIService(&config.Config{Key: key})
}

func TestCreateBin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const url = "https://api.jsonbin.io/v3/b"

	tests := []struct {
		name       string
		key        string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "success created",
			key:        "test-key",
			statusCode: http.StatusCreated,
			body:       `{"record":{"foo":"bar"},"metadata":{"id":"bin123","createdAt":"2024-01-01T00:00:00Z","private":false}}`,
			wantErr:    false,
		},
		{
			name:       "error missing key",
			key:        "",
			statusCode: http.StatusCreated,
			body:       `{"metadata":{"id":"irrelevant"}}`,
			wantErr:    true,
		},
		{
			name:       "error non-2xx",
			key:        "test-key",
			statusCode: http.StatusInternalServerError,
			body:       `{"message":"server error"}`,
			wantErr:    true,
		},
		{
			name:       "error invalid json",
			key:        "test-key",
			statusCode: http.StatusCreated,
			body:       `{invalid`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			// For missing key case, we do not need to register a responder because request isn't sent.
			if tc.key != "" {
				var gotContentType, gotKey string
				httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
					gotContentType = req.Header.Get("Content-Type")
					gotKey = req.Header.Get("X-Master-Key")
					return httpmock.NewStringResponse(tc.statusCode, tc.body), nil
				})

				svc := newServiceWithKey(tc.key)
				resp, err := svc.CreateBin(`{"foo":"bar"}`)
				if (err != nil) != tc.wantErr {
					t.Fatalf("unexpected error state: err=%v wantErr=%v", err, tc.wantErr)
				}
				if err != nil {
					return
				}
				if gotContentType != "application/json" {
					t.Errorf("missing or wrong Content-Type header: %q", gotContentType)
				}
				if gotKey != tc.key {
					t.Errorf("X-Master-Key header mismatch: got=%q want=%q", gotKey, tc.key)
				}
				if resp == nil || resp.Metadata.ID == "" {
					t.Errorf("expected valid response with ID, got: %+v", resp)
				}
			} else {
				svc := newServiceWithKey(tc.key)
				resp, err := svc.CreateBin(`{"foo":"bar"}`)
				if err == nil {
					t.Fatalf("expected error for missing key, got resp=%+v", resp)
				}
			}
		})
	}
}

func TestGetBin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	binID := "bin-get-1"
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID)

	tests := []struct {
		name       string
		key        string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "success get",
			key:        "test-key",
			statusCode: http.StatusOK,
			body:       `{"record":{"name":"alpha"},"metadata":{"id":"` + binID + `"}}`,
			wantErr:    false,
		},
		{
			name:       "error missing key",
			key:        "",
			statusCode: http.StatusOK,
			body:       `{"metadata":{"id":"` + binID + `"}}`,
			wantErr:    true,
		},
		{
			name:       "error non-200",
			key:        "test-key",
			statusCode: http.StatusNotFound,
			body:       `{"message":"not found"}`,
			wantErr:    true,
		},
		{
			name:       "error invalid json",
			key:        "test-key",
			statusCode: http.StatusOK,
			body:       `{invalid`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			if tc.key != "" {
				var gotKey string
				httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
					gotKey = req.Header.Get("X-Master-Key")
					return httpmock.NewStringResponse(tc.statusCode, tc.body), nil
				})
				svc := newServiceWithKey(tc.key)
				resp, err := svc.GetBin(binID)
				if (err != nil) != tc.wantErr {
					t.Fatalf("unexpected error state: err=%v wantErr=%v", err, tc.wantErr)
				}
				if err != nil {
					return
				}
				if gotKey != tc.key {
					t.Errorf("X-Master-Key header mismatch: got=%q want=%q", gotKey, tc.key)
				}
				if resp == nil {
					t.Fatalf("expected non-nil response")
				}
				// basic shape check
				meta, ok := resp["metadata"].(map[string]interface{})
				if !ok || meta["id"] == "" {
					t.Errorf("metadata.id missing in response: %+v", resp)
				}
			} else {
				svc := newServiceWithKey(tc.key)
				_, err := svc.GetBin(binID)
				if err == nil {
					t.Fatalf("expected error for missing key")
				}
			}
		})
	}
}

func TestUpdateBin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	binID := "bin-upd-1"
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID)

	tests := []struct {
		name       string
		key        string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "success update",
			key:        "test-key",
			statusCode: http.StatusOK,
			body:       `{"record":{"name":"updated"},"metadata":{"id":"` + binID + `","createdAt":"2024-01-01T00:00:00Z","private":false}}`,
			wantErr:    false,
		},
		{
			name:       "error missing key",
			key:        "",
			statusCode: http.StatusOK,
			body:       `{"metadata":{"id":"` + binID + `"}}`,
			wantErr:    true,
		},
		{
			name:       "error non-200",
			key:        "test-key",
			statusCode: http.StatusBadRequest,
			body:       `{"message":"bad request"}`,
			wantErr:    true,
		},
		{
			name:       "error invalid json",
			key:        "test-key",
			statusCode: http.StatusOK,
			body:       `{invalid`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			if tc.key != "" {
				var gotContentType, gotKey string
				httpmock.RegisterResponder("PUT", url, func(req *http.Request) (*http.Response, error) {
					gotContentType = req.Header.Get("Content-Type")
					gotKey = req.Header.Get("X-Master-Key")
					return httpmock.NewStringResponse(tc.statusCode, tc.body), nil
				})
				svc := newServiceWithKey(tc.key)
				resp, err := svc.UpdateBin(binID, `{"name":"updated"}`)
				if (err != nil) != tc.wantErr {
					t.Fatalf("unexpected error state: err=%v wantErr=%v", err, tc.wantErr)
				}
				if err != nil {
					return
				}
				if gotContentType != "application/json" {
					t.Errorf("missing or wrong Content-Type header: %q", gotContentType)
				}
				if gotKey != tc.key {
					t.Errorf("X-Master-Key header mismatch: got=%q want=%q", gotKey, tc.key)
				}
				if resp == nil || resp.Metadata.ID == "" {
					t.Errorf("expected valid response with ID, got: %+v", resp)
				}
			} else {
				svc := newServiceWithKey(tc.key)
				_, err := svc.UpdateBin(binID, `{"name":"updated"}`)
				if err == nil {
					t.Fatalf("expected error for missing key")
				}
			}
		})
	}
}

func TestDeleteBin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	binID := "bin-del-1"
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID)

	tests := []struct {
		name       string
		key        string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "success delete",
			key:        "test-key",
			statusCode: http.StatusOK,
			body:       `{"metadata":{"id":"` + binID + `","versionsDeleted":1},"message":"deleted"}`,
			wantErr:    false,
		},
		{
			name:       "error missing key",
			key:        "",
			statusCode: http.StatusOK,
			body:       `{"metadata":{"id":"` + binID + `"}}`,
			wantErr:    true,
		},
		{
			name:       "error non-200",
			key:        "test-key",
			statusCode: http.StatusNotFound,
			body:       `{"message":"not found"}`,
			wantErr:    true,
		},
		{
			name:       "error invalid json",
			key:        "test-key",
			statusCode: http.StatusOK,
			body:       `{invalid`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			if tc.key != "" {
				var gotKey string
				httpmock.RegisterResponder("DELETE", url, func(req *http.Request) (*http.Response, error) {
					gotKey = req.Header.Get("X-Master-Key")
					return httpmock.NewStringResponse(tc.statusCode, tc.body), nil
				})
				svc := newServiceWithKey(tc.key)
				resp, err := svc.DeleteBin(binID)
				if (err != nil) != tc.wantErr {
					t.Fatalf("unexpected error state: err=%v wantErr=%v", err, tc.wantErr)
				}
				if err != nil {
					return
				}
				if gotKey != tc.key {
					t.Errorf("X-Master-Key header mismatch: got=%q want=%q", gotKey, tc.key)
				}
				if resp == nil || resp.Metadata.ID == "" {
					t.Errorf("expected valid delete response with ID, got: %+v", resp)
				}
			} else {
				svc := newServiceWithKey(tc.key)
				_, err := svc.DeleteBin(binID)
				if err == nil {
					t.Fatalf("expected error for missing key")
				}
			}
		})
	}
}

func TestCRUD_Flow_Grouped(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	key := "test-key"
	svc := newServiceWithKey(key)

	binID := "group-123"

	// Create
	httpmock.RegisterResponder("POST", "https://api.jsonbin.io/v3/b", func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("X-Master-Key") != key || req.Header.Get("Content-Type") != "application/json" {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"message":"bad headers"}`), nil
		}
		body := map[string]any{
			"record":   map[string]any{"name": "alpha"},
			"metadata": map[string]any{"id": binID, "createdAt": "2024-01-01T00:00:00Z", "private": false},
		}
		b, _ := json.Marshal(body)
		return httpmock.NewBytesResponse(http.StatusCreated, b), nil
	})

	// Get
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("X-Master-Key") != key {
			return httpmock.NewStringResponse(http.StatusUnauthorized, `{"message":"unauthorized"}`), nil
		}
		return httpmock.NewStringResponse(http.StatusOK, `{"record":{"name":"alpha"},"metadata":{"id":"`+binID+`"}}`), nil
	})

	// Update
	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("X-Master-Key") != key || req.Header.Get("Content-Type") != "application/json" {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"message":"bad headers"}`), nil
		}
		return httpmock.NewStringResponse(http.StatusOK, `{"record":{"name":"beta"},"metadata":{"id":"`+binID+`","createdAt":"2024-01-01T00:00:00Z","private":false}}`), nil
	})

	// Delete
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("X-Master-Key") != key {
			return httpmock.NewStringResponse(http.StatusUnauthorized, `{"message":"unauthorized"}`), nil
		}
		return httpmock.NewStringResponse(http.StatusOK, `{"metadata":{"id":"`+binID+`","versionsDeleted":1},"message":"deleted"}`), nil
	})

	// Ensure cleanup even if assertions fail later
	t.Cleanup(func() {
		_, _ = svc.DeleteBin(binID)
	})

	created, err := svc.CreateBin(`{"name":"alpha"}`)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Metadata.ID != binID {
		t.Fatalf("unexpected created id: %s", created.Metadata.ID)
	}

	got, err := svc.GetBin(binID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	meta := got["metadata"].(map[string]interface{})
	if meta["id"].(string) != binID {
		t.Fatalf("unexpected get id: %v", meta["id"])
	}

	updated, err := svc.UpdateBin(binID, `{"name":"beta"}`)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Metadata.ID != binID {
		t.Fatalf("unexpected updated id: %s", updated.Metadata.ID)
	}

	deleted, err := svc.DeleteBin(binID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if deleted.Metadata.ID != binID {
		t.Fatalf("unexpected deleted id: %s", deleted.Metadata.ID)
	}
}
