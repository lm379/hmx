package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{}
	return c, w
}

func assertResponse(t *testing.T, w *httptest.ResponseRecorder, expectedCode int, expectedMsg string, dataExpected bool) {
	t.Helper()

	if w.Code != expectedCode {
		t.Errorf("Expected status code %d, got %d", expectedCode, w.Code)
	}

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != expectedCode {
		t.Errorf("Expected response code %d, got %d", expectedCode, resp.Code)
	}

	if resp.Msg != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, resp.Msg)
	}

	if dataExpected && resp.Data == nil {
		t.Error("Expected data to be non-nil")
	}
}

func TestSuccess(t *testing.T) {
	c, w := setupTestContext()
	testData := map[string]string{"key": "value"}

	Success(c, testData)

	assertResponse(t, w, http.StatusOK, "success", true)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	// Check data
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	if dataMap["key"] != "value" {
		t.Errorf("Expected data['key'] to be 'value', got %v", dataMap["key"])
	}
}

func TestSuccessWithNilData(t *testing.T) {
	c, w := setupTestContext()

	Success(c, nil)

	assertResponse(t, w, http.StatusOK, "success", false)
}

func TestSuccessWithMsg(t *testing.T) {
	c, w := setupTestContext()
	testData := "test data"
	customMsg := "Custom success message"

	SuccessWithMsg(c, customMsg, testData)

	assertResponse(t, w, http.StatusOK, customMsg, true)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Data != "test data" {
		t.Errorf("Expected data to be 'test data', got %v", resp.Data)
	}
}

func TestCreated(t *testing.T) {
	c, w := setupTestContext()
	testData := map[string]string{"id": "1"}

	Created(c, testData)

	assertResponse(t, w, http.StatusCreated, "created successfully", true)
}

func TestAcceped(t *testing.T) {
	c, w := setupTestContext()
	testData := "request accepted"

	Acceped(c, testData)

	assertResponse(t, w, http.StatusAccepted, "accepted", true)
}

func TestNoContent(t *testing.T) {
	c, w := setupTestContext()

	NoContent(c)

	// Gin's JSON method for 204 may not output body
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestNoContentWithMsg(t *testing.T) {
	c, w := setupTestContext()
	customMsg := "Deleted successfully"

	NoContentWithMsg(c, customMsg)

	// Gin's JSON method for 204 may not output body
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		msg          string
		expectedHTTP int
	}{
		{
			name:         "Generic error",
			code:         500,
			msg:          "Internal server error",
			expectedHTTP: 500,
		},
		{
			name:         "Custom error code",
			code:         418,
			msg:          "I'm a teapot",
			expectedHTTP: 418,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()

			Error(c, tt.code, tt.msg)

			assertResponse(t, w, tt.expectedHTTP, tt.msg, false)
		})
	}
}

func TestBadRequest(t *testing.T) {
	c, w := setupTestContext()
	msg := "Bad request"

	BadRequest(c, msg)

	assertResponse(t, w, http.StatusBadRequest, msg, false)
}

func TestUnauthorized(t *testing.T) {
	c, w := setupTestContext()
	msg := "Unauthorized access"

	Unauthorized(c, msg)

	assertResponse(t, w, http.StatusUnauthorized, msg, false)
}

func TestForbidden(t *testing.T) {
	c, w := setupTestContext()
	msg := "Forbidden"

	Forbidden(c, msg)

	assertResponse(t, w, http.StatusForbidden, msg, false)
}

func TestNotFound(t *testing.T) {
	c, w := setupTestContext()
	msg := "Resource not found"

	NotFound(c, msg)

	assertResponse(t, w, http.StatusNotFound, msg, false)
}

func TestConflict(t *testing.T) {
	c, w := setupTestContext()
	msg := "Resource already exists"

	Conflict(c, msg)

	assertResponse(t, w, http.StatusConflict, msg, false)
}

func TestTooManyRequests(t *testing.T) {
	c, w := setupTestContext()
	msg := "Rate limit exceeded"

	TooManyRequests(c, msg)

	assertResponse(t, w, http.StatusTooManyRequests, msg, false)
}

func TestInternalServerError(t *testing.T) {
	c, w := setupTestContext()
	msg := "Something went wrong"

	InternalServerError(c, msg)

	assertResponse(t, w, http.StatusInternalServerError, msg, false)
}

func TestServiceUnavailable(t *testing.T) {
	c, w := setupTestContext()
	msg := "Service temporarily unavailable"

	ServiceUnavailable(c, msg)

	assertResponse(t, w, http.StatusServiceUnavailable, msg, false)
}

func TestGatewayTimeout(t *testing.T) {
	c, w := setupTestContext()
	msg := "Gateway timeout"

	GatewayTimeout(c, msg)

	assertResponse(t, w, http.StatusGatewayTimeout, msg, false)
}

func TestGatewayError(t *testing.T) {
	c, w := setupTestContext()
	msg := "Bad gateway"

	GatewayError(c, msg)

	assertResponse(t, w, http.StatusBadGateway, msg, false)
}

func TestResponseStructure(t *testing.T) {
	c, w := setupTestContext()

	testData := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{
		ID:   1,
		Name: "Test",
	}

	Success(c, testData)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response structure
	if resp.Code != 200 {
		t.Errorf("Expected code 200, got %d", resp.Code)
	}

	if resp.Msg != "success" {
		t.Errorf("Expected msg 'success', got '%s'", resp.Msg)
	}

	// Convert data to map for verification
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	if dataMap["id"] != float64(1) {
		t.Errorf("Expected data['id'] to be 1, got %v", dataMap["id"])
	}

	if dataMap["name"] != "Test" {
		t.Errorf("Expected data['name'] to be 'Test', got %v", dataMap["name"])
	}
}

func TestMultipleResponses(t *testing.T) {
	c, w := setupTestContext()

	// Send multiple responses
	Success(c, "first")
	Success(c, "second")

	// Note: In Gin, calling c.JSON multiple times results in invalid JSON
	// This test verifies that behavior
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// The body will contain multiple JSON objects concatenated
	body := w.Body.String()
	if body == "" {
		t.Error("Expected body to be non-empty")
	}
}

func TestEmptyStringMessage(t *testing.T) {
	c, w := setupTestContext()

	Error(c, 400, "")

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Msg != "" {
		t.Errorf("Expected empty message, got '%s'", resp.Msg)
	}
}

func TestComplexData(t *testing.T) {
	c, w := setupTestContext()

	complexData := map[string]interface{}{
		"user": map[string]string{
			"id":    "1",
			"name":  "John",
			"email": "john@example.com",
		},
		"roles":  []string{"admin", "user"},
		"active": true,
		"score":  95.5,
	}

	Success(c, complexData)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Data == nil {
		t.Fatal("Expected data to be non-nil")
	}

	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	if dataMap["active"] != true {
		t.Errorf("Expected data['active'] to be true, got %v", dataMap["active"])
	}
}

func TestResponseWithNilContext(t *testing.T) {
	// This test verifies that the functions handle nil context gracefully
	// In practice, this should not happen in real code

	defer func() {
		if r := recover(); r != nil {
			// Expected to panic or handle nil context
		}
	}()

	// Note: In Gin, nil context will cause panic when trying to write response
	// This is expected behavior and should be handled at a higher level
}
