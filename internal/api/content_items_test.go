package api

import "testing"

func TestContentItemListParamsValuesIncludesSourceType(t *testing.T) {
	values := ContentItemListParams{
		Status:     "processed",
		Search:     "hello",
		SourceType: "inbound_email",
		Page:       2,
		PerPage:    25,
	}.Values()

	if got := values.Get("status"); got != "processed" {
		t.Fatalf("status = %q, want processed", got)
	}
	if got := values.Get("search"); got != "hello" {
		t.Fatalf("search = %q, want hello", got)
	}
	if got := values.Get("source_type"); got != "inbound_email" {
		t.Fatalf("source_type = %q, want inbound_email", got)
	}
	if got := values.Get("page"); got != "2" {
		t.Fatalf("page = %q, want 2", got)
	}
	if got := values.Get("items"); got != "25" {
		t.Fatalf("items = %q, want 25", got)
	}
}
