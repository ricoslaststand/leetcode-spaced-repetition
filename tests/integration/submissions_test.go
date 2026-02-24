package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveProblemSubmission_Valid(t *testing.T) {
	insertTestProblem(t, testProblem)
	clearSubmissions(t)

	body := `{"problemId":1,"confidenceLevel":"good","timeTaken":90}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/problems/submissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSaveProblemSubmission_InvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/problems/submissions", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSaveProblemSubmission_MissingConfidence(t *testing.T) {
	body := `{"problemId":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/problems/submissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetAllProblemSubmissions_Empty(t *testing.T) {
	clearSubmissions(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/submissions", nil)
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
		t.Fatalf("expected data array, got: %v", resp)
	}
	if len(data) != 0 {
		t.Fatalf("expected empty submissions, got %d", len(data))
	}
}

func TestGetAllProblemSubmissions_WithData(t *testing.T) {
	insertTestProblem(t, testProblem)
	clearSubmissions(t)

	// Save one submission via the API
	body := `{"problemId":1,"confidenceLevel":"easy"}`
	saveReq, _ := http.NewRequest(http.MethodPost, "/problems/submissions", strings.NewReader(body))
	saveReq.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(httptest.NewRecorder(), saveReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/submissions", nil)
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
		t.Fatalf("expected at least one submission, got: %v", resp)
	}
}

func TestGetProblemSubmissions_ByID(t *testing.T) {
	insertTestProblem(t, testProblem)
	clearSubmissions(t)

	body := `{"problemId":1,"confidenceLevel":"hard"}`
	saveReq, _ := http.NewRequest(http.MethodPost, "/problems/submissions", strings.NewReader(body))
	saveReq.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(httptest.NewRecorder(), saveReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/problems/1/submissions", nil)
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
		t.Fatalf("expected at least one submission for problem 1, got: %v", resp)
	}
}

// TestImportSubmissions_DatePreservation is the primary regression test for the
// date timezone bug: a submission with date "2025-06-15" must be stored as
// "2025-06-15", not shifted to "2025-06-14" by timezone conversion.
func TestImportSubmissions_DatePreservation(t *testing.T) {
	insertTestProblem(t, testProblem)
	clearSubmissions(t)

	const wantDate = "2025-06-15"
	xlsxBuf := makeXLSXFile(t, [][]string{
		{"1", wantDate, "1m30s", "3"},
	})

	w := httptest.NewRecorder()
	req := multipartFileRequest(t, "/problems/submissions/import", xlsxBuf.Bytes())
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("import returned %d: %s", w.Code, w.Body.String())
	}

	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	imported := result["imported"].(float64)
	if imported != 1 {
		t.Fatalf("expected 1 imported, got %v (errors: %v)", imported, result["errors"])
	}

	// Query the stored date directly from the database.
	var storedDate string
	err := testDB.QueryRow(`SELECT submission_date::text FROM problem_submissions LIMIT 1`).Scan(&storedDate)
	if err != nil {
		t.Fatalf("failed to query submission_date: %v", err)
	}

	// PostgreSQL returns DATE as YYYY-MM-DD.
	if !strings.HasPrefix(storedDate, wantDate) {
		t.Fatalf("date drift detected: want %s, got %s", wantDate, storedDate)
	}
}

func TestImportSubmissions_InvalidRow(t *testing.T) {
	insertTestProblem(t, testProblem)
	clearSubmissions(t)

	// Confidence "99" is not a valid level.
	xlsxBuf := makeXLSXFile(t, [][]string{
		{"1", "2025-06-15", "", "99"},
	})

	w := httptest.NewRecorder()
	req := multipartFileRequest(t, "/problems/submissions/import", xlsxBuf.Bytes())
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 with error list, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	errors, ok := result["errors"].([]any)
	if !ok || len(errors) == 0 {
		t.Fatalf("expected errors array with at least one entry, got: %v", result)
	}
}

func TestImportSubmissions_MissingFile(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/problems/submissions/import", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

// multipartFileRequest constructs a multipart/form-data POST request with the
// xlsx file content attached under the "file" field name.
func multipartFileRequest(t *testing.T, url string, xlsxData []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", "submissions.xlsx")
	if err != nil {
		t.Fatalf("multipartFileRequest: create form file: %v", err)
	}
	if _, err := part.Write(xlsxData); err != nil {
		t.Fatalf("multipartFileRequest: write data: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("multipartFileRequest: close writer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		t.Fatalf("multipartFileRequest: new request: %v", err)
	}
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary()))
	return req
}
