package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProblems_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := resp["data"].([]any)
	if !ok {
		t.Fatalf("expected data array in response, got: %v", resp)
	}
	if len(data) != 0 {
		t.Fatalf("expected empty data, got %d items", len(data))
	}
}

func TestGetProblems_WithData(t *testing.T) {
	insertTestProblem(t, testProblem)
	insertTestTopic(t, testProblem.ID, "array")

	// GetProblems always filters by topics, so a topics param is required to get results.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems?topics=array", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := resp["data"].([]any)
	if !ok || len(data) == 0 {
		t.Fatalf("expected at least one problem in data, got: %v", resp)
	}
}

func TestGetProblemByID_Found(t *testing.T) {
	insertTestProblem(t, testProblem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/1", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["title"] != "Two Sum" {
		t.Fatalf("expected title 'Two Sum', got %v", resp["title"])
	}
}

func TestGetProblemByID_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/9999", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetProblemByID_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/abc", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetAllProblemTopics_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/topics", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, ok := resp["topics"]; !ok {
		t.Fatalf("expected 'topics' key in response, got: %v", resp)
	}
}

func TestGetProblemByID_WithTopics(t *testing.T) {
	insertTestProblem(t, testProblem)
	insertTestTopic(t, testProblem.ID, "array")
	insertTestTopic(t, testProblem.ID, "hash-table")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/1", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	topics, ok := resp["topics"].([]any)
	if !ok {
		t.Fatalf("expected topics array, got: %v", resp["topics"])
	}
	if len(topics) != 2 {
		t.Fatalf("expected 2 topics, got %d: %v", len(topics), topics)
	}
}
