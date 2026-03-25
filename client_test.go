package ryft

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetermineBaseURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		secretKey string
		wantURL   string
		wantError bool
	}{
		{
			name:      "sandbox key",
			secretKey: "sk_sandbox_123",
			wantURL:   sandboxBaseURL,
		},
		{
			name:      "live key",
			secretKey: "sk_live_123",
			wantURL:   liveBaseURL,
		},
		{
			name:      "invalid key",
			secretKey: "pk_test_123",
			wantError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			gotURL, err := determineBaseURL(testCase.secretKey)
			if testCase.wantError {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("determineBaseURL returned error: %v", err)
			}

			if gotURL != testCase.wantURL {
				t.Fatalf("determineBaseURL = %q, want %q", gotURL, testCase.wantURL)
			}
		})
	}
}

func TestEventsListSetsAccountHeader(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/events")
		}
		if got := r.Header.Get("Account"); got != "ac_test_123" {
			t.Fatalf("Account header = %q, want %q", got, "ac_test_123")
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Fatalf("limit query = %q, want %q", got, "50")
		}
		if got := r.URL.Query().Get("ascending"); got != "false" {
			t.Fatalf("ascending query = %q, want %q", got, "false")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[]}`)
	}))
	defer server.Close()

	client, err := NewClient(Config{
		SecretKey: "sk_sandbox_123",
		BaseURL:   server.URL,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if _, err := client.Events.List(context.Background(), EventListParams{ListParams: ListParams{Ascending: false, Limit: 50}}, WithAccount("ac_test_123")); err != nil {
		t.Fatalf("Events.List returned error: %v", err)
	}
}

func TestFilesCreateUsesPDFContentType(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "evidence.pdf")
	if err := os.WriteFile(filePath, []byte("%PDF-1.4\nhello\n%%EOF\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data; boundary=") {
			t.Fatalf("request Content-Type = %q, want multipart/form-data", got)
		}

		if err := r.ParseMultipartForm(1024 * 1024); err != nil {
			t.Fatalf("ParseMultipartForm returned error: %v", err)
		}

		if got := r.FormValue("category"); got != "Evidence" {
			t.Fatalf("category = %q, want %q", got, "Evidence")
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile returned error: %v", err)
		}
		defer file.Close()

		if header.Header.Get("Content-Type") != "application/pdf" {
			t.Fatalf("multipart file Content-Type = %q, want %q", header.Header.Get("Content-Type"), "application/pdf")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"file_123","category":"Evidence","createdTimestamp":1}`)
	}))
	defer server.Close()

	client, err := NewClient(Config{
		SecretKey: "sk_sandbox_123",
		BaseURL:   server.URL,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	file, err := client.Files.Create(context.Background(), CreateFileRequest{
		FilePath: filePath,
		Category: "Evidence",
	})
	if err != nil {
		t.Fatalf("Files.Create returned error: %v", err)
	}

	if file.ID != "file_123" {
		t.Fatalf("file.ID = %q, want %q", file.ID, "file_123")
	}
}
