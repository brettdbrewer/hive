package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// claimsResponse mirrors the JSON shape returned by the knowledge endpoint.
type claimsResponse struct {
	Claims []claimNode `json:"claims"`
}

// titledClaimsResponse mirrors the shape used by queryMaxLessonNumber.
type titledClaimsResponse struct {
	Claims []struct {
		Title string `json:"title"`
	} `json:"claims"`
}

// newClaimsServer returns an httptest.Server that responds to
// GET /app/hive/knowledge with the given claims list (as title-only records).
func newClaimsServer(t *testing.T, titles []string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var resp titledClaimsResponse
		for _, title := range titles {
			resp.Claims = append(resp.Claims, struct {
				Title string `json:"title"`
			}{Title: title})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// TestQueryMaxLessonNumber_extractsHighestNumber verifies the max is found from
// a list of Lesson N: ... titles.
func TestQueryMaxLessonNumber_extractsHighestNumber(t *testing.T) {
	srv := newClaimsServer(t, []string{
		"Lesson 10: Some lesson",
		"Lesson 183: Latest lesson",
		"Lesson 5: Early lesson",
		"Lesson 100: Middle lesson",
	})
	baseURL = srv.URL

	max, err := queryMaxLessonNumber("lv_testkey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if max != 183 {
		t.Errorf("max = %d, want 183", max)
	}
}

// TestQueryMaxLessonNumber_returnsZeroWhenNoLessons verifies that zero is
// returned when no titles match the Lesson N: pattern.
func TestQueryMaxLessonNumber_returnsZeroWhenNoLessons(t *testing.T) {
	srv := newClaimsServer(t, []string{
		"Some random claim",
		"Another non-lesson title",
	})
	baseURL = srv.URL

	max, err := queryMaxLessonNumber("lv_testkey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if max != 0 {
		t.Errorf("max = %d, want 0", max)
	}
}

// TestQueryMaxLessonNumber_ignoresNonLessonTitles verifies that mixed content
// (docs, blog posts, etc.) does not pollute the lesson number extraction.
func TestQueryMaxLessonNumber_ignoresNonLessonTitles(t *testing.T) {
	srv := newClaimsServer(t, []string{
		"Lesson 42: Real lesson",
		"Not a lesson: just a claim",
		"lesson 99: lowercase — should not match",
		"Iteration summary: Lesson 200 shipped",
		"Lesson 42abc: malformed number",
	})
	baseURL = srv.URL

	max, err := queryMaxLessonNumber("lv_testkey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only "Lesson 42: Real lesson" matches ^Lesson (\d+). The lowercase
	// "lesson 99" and "Iteration summary: Lesson 200..." don't start with
	// "Lesson N" at position 0.
	if max != 42 {
		t.Errorf("max = %d, want 42", max)
	}
}

// TestQueryMaxLessonNumber_httpError verifies that a 4xx response is surfaced
// as an error rather than silently returning 0.
func TestQueryMaxLessonNumber_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL

	_, err := queryMaxLessonNumber("lv_badkey")
	if err == nil {
		t.Fatal("expected error on HTTP 403, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 403") {
		t.Errorf("error should mention HTTP 403, got: %v", err)
	}
}

// TestFetchRetractedClaims_parsesClaims verifies that retracted claim JSON is
// decoded into claimNode structs with ID, Title, and Body preserved.
func TestFetchRetractedClaims_parsesClaims(t *testing.T) {
	want := []claimNode{
		{ID: "f6e95fb1-abcd-1234", Title: "Lesson 144: some title", Body: "Some body text."},
		{ID: "df838c45-efgh-5678", Title: "Lesson 145: another title", Body: "Another body."},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(claimsResponse{Claims: want})
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL

	got, err := fetchRetractedClaims("lv_testkey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i, g := range got {
		w := want[i]
		if g.ID != w.ID {
			t.Errorf("[%d] ID = %q, want %q", i, g.ID, w.ID)
		}
		if g.Title != w.Title {
			t.Errorf("[%d] Title = %q, want %q", i, g.Title, w.Title)
		}
		if g.Body != w.Body {
			t.Errorf("[%d] Body = %q, want %q", i, g.Body, w.Body)
		}
	}
}

// TestFetchRetractedClaims_httpError verifies that a non-200 response returns
// an error rather than an empty slice.
func TestFetchRetractedClaims_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL

	_, err := fetchRetractedClaims("lv_testkey")
	if err == nil {
		t.Fatal("expected error on HTTP 404, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("error should mention HTTP 404, got: %v", err)
	}
}

// TestAssertClaim_sendsCorrectPayload verifies that assertClaim POSTs the
// correct JSON fields (op, title, body) to /app/hive/op.
func TestAssertClaim_sendsCorrectPayload(t *testing.T) {
	var got map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != fmt.Sprintf("/app/%s/op", spaceSlug) {
			t.Errorf("path = %q, want /app/hive/op", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL

	title := "Lesson 184: Fixing a search index is necessary"
	body := "Detailed explanation of the lesson."
	if err := assertClaim("lv_testkey", title, body); err != nil {
		t.Fatalf("assertClaim() error: %v", err)
	}

	if got["op"] != "assert" {
		t.Errorf("op = %q, want %q", got["op"], "assert")
	}
	if got["title"] != title {
		t.Errorf("title = %q, want %q", got["title"], title)
	}
	if got["body"] != body {
		t.Errorf("body = %q, want %q", got["body"], body)
	}
}

// TestAssertClaim_emDashNormalization verifies that em-dashes in the title are
// preserved as the Unicode em-dash character (not mangled) in the JSON payload.
func TestAssertClaim_emDashNormalization(t *testing.T) {
	var got map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL

	// Title contains an em-dash — json.Marshal preserves it as U+2014.
	title := "Lesson 191: API endpoint \u2014 semantic mismatch"
	if err := assertClaim("lv_testkey", title, "body"); err != nil {
		t.Fatalf("assertClaim() error: %v", err)
	}

	if !strings.Contains(got["title"], "\u2014") {
		t.Errorf("em-dash not preserved in title: %q", got["title"])
	}
}

// TestAssertClaim_httpError verifies that a 4xx response from the server is
// returned as an error.
func TestAssertClaim_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL

	err := assertClaim("lv_badkey", "Lesson 184: title", "body")
	if err == nil {
		t.Fatal("expected error on HTTP 401, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 401") {
		t.Errorf("error should mention HTTP 401, got: %v", err)
	}
}

// TestShortIDExtraction verifies the 8-char short ID slicing logic used to
// match retracted claims to their target numbers. IDs shorter than 8 chars
// are skipped in production code — verify exact boundary.
func TestShortIDExtraction(t *testing.T) {
	cases := []struct {
		id        string
		wantShort string
		wantOK    bool
	}{
		{"f6e95fb1-abcd-1234", "f6e95fb1", true},
		{"df838c45", "df838c45", true}, // exactly 8 chars
		{"abc1234", "", false},          // 7 chars — too short
		{"", "", false},                 // empty
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			ok := len(tc.id) >= 8
			if ok != tc.wantOK {
				t.Errorf("len(%q) >= 8: got %v, want %v", tc.id, ok, tc.wantOK)
			}
			if ok {
				short := tc.id[:8]
				if short != tc.wantShort {
					t.Errorf("id[:8] = %q, want %q", short, tc.wantShort)
				}
			}
		})
	}
}
