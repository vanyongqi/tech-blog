package controller

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONAcceptsLargeAdminPayload(t *testing.T) {
	var target struct {
		Content string `json:"contentMarkdown"`
	}

	largeMarkdown := strings.Repeat("0123456789abcdef", 1500)
	body := `{"contentMarkdown":"` + largeMarkdown + `"}`
	request := httptest.NewRequest("POST", "/api/admin/posts/example", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	if err := decodeJSON(recorder, request, &target); err != nil {
		t.Fatalf("expected large admin payload to decode, got error: %v", err)
	}
	if target.Content != largeMarkdown {
		t.Fatalf("expected decoded content to match input")
	}
}
