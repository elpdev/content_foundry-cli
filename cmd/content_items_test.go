package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestContentItemsListInboundEmailAddsEmailColumns(t *testing.T) {
	var gotSourceType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"test-token","expires_in":3600}`))
		case "/api/v1/content_items":
			gotSourceType = r.URL.Query().Get("source_type")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"content_items": [{
					"id": 42,
					"title": "Inbox story",
					"source_id": 7,
					"status": "processed",
					"fetched_at": "2026-04-09T04:00:00Z",
					"metadata": {
						"from": "reporter@example.com",
						"received_at": "2026-04-09T03:55:00Z",
						"attachment_count": 2
					},
					"assets": []
				}],
				"pagination": {"page":1,"items":20,"pages":1,"count":1}
			}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	output := executeRootCommand(t, server.URL, "--format", "text", "content-items", "list", "--source-type", "inbound_email")
	if gotSourceType != "inbound_email" {
		t.Fatalf("source_type query = %q, want inbound_email", gotSourceType)
	}

	for _, want := range []string{
		"From: reporter@example.com",
		"Received At: 2026-04-09T03:55:00Z",
		"Attachments: 2",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\n%s", want, output)
		}
	}
}

func TestContentItemsListWithoutInboundEmailKeepsDefaultColumns(t *testing.T) {
	var gotSourceType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"test-token","expires_in":3600}`))
		case "/api/v1/content_items":
			gotSourceType = r.URL.Query().Get("source_type")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"content_items": [{
					"id": 42,
					"title": "Inbox story",
					"source_id": 7,
					"status": "processed",
					"fetched_at": "2026-04-09T04:00:00Z",
					"metadata": {
						"from": "reporter@example.com",
						"received_at": "2026-04-09T03:55:00Z",
						"attachment_count": 2
					},
					"assets": []
				}],
				"pagination": {"page":1,"items":20,"pages":1,"count":1}
			}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	output := executeRootCommand(t, server.URL, "--format", "text", "content-items", "list")
	if gotSourceType != "" {
		t.Fatalf("source_type query = %q, want empty", gotSourceType)
	}

	for _, unwanted := range []string{"From:", "Received At:", "Attachments:"} {
		if strings.Contains(output, unwanted) {
			t.Fatalf("output unexpectedly contains %q\n%s", unwanted, output)
		}
	}
}

func TestContentItemsShowIncludesEmailMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"test-token","expires_in":3600}`))
		case "/api/v1/content_items/42":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"content_item": {
					"id": 42,
					"title": "Inbox story",
					"source_id": 7,
					"source_url": "imap://inbox/42",
					"status": "processed",
					"fetched_at": "2026-04-09T04:00:00Z",
					"created_at": "2026-04-09T04:01:00Z",
					"metadata": {
						"from": "reporter@example.com",
						"received_at": "2026-04-09T03:55:00Z",
						"attachment_count": 2
					},
					"assets": []
				},
				"drafts": []
			}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	output := executeRootCommand(t, server.URL, "--format", "text", "content-items", "show", "42")

	for _, want := range []string{
		"From: reporter@example.com",
		"Received At: 2026-04-09T03:55:00Z",
		"Attachment Count: 2",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\n%s", want, output)
		}
	}
}

func executeRootCommand(t *testing.T, baseURL string, args ...string) string {
	t.Helper()

	tempDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tempDir)

	cfg = &config.Config{
		BaseURL:        baseURL,
		ClientID:       "client-id",
		SecretKey:      "secret-key",
		DefaultBrandID: 1,
	}
	apiClient = api.NewClient(cfg)
	formatter = nil
	outFormat = "table"
	verbose = false
	brandFlag = ""

	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		cfg = nil
		apiClient = nil
		formatter = nil
		outFormat = "table"
		verbose = false
		brandFlag = ""
	}()

	rootCmd.SetArgs(args)
	resetCommandFlags(rootCmd)

	stdout := captureStdout(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("root command failed: %v", err)
		}
	})

	if _, err := os.Stat(filepath.Join(tempDir, "content_foundry", "token.toml")); err != nil {
		t.Fatalf("expected token cache to be written: %v", err)
	}

	return stdout
}

func resetCommandFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})
	for _, child := range cmd.Commands() {
		resetCommandFlags(child)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	outputCh := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	fn()
	_ = w.Close()
	_ = r.Close()

	return <-outputCh
}
