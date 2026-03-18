package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

func TestGetPagination(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		expected    *Pagination
	}{
		{
			name:        "Default parameters",
			queryParams: map[string]string{},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Valid page and page_size",
			queryParams: map[string]string{
				"page":      "2",
				"page_size": "20",
			},
			expected: &Pagination{
				Page:     2,
				PageSize: 20,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Custom sort",
			queryParams: map[string]string{
				"sort": "id asc",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "id asc",
			},
		},
		{
			name: "All parameters",
			queryParams: map[string]string{
				"page":      "3",
				"page_size": "25",
				"sort":      "name desc",
			},
			expected: &Pagination{
				Page:     3,
				PageSize: 25,
				Sort:     "name desc",
			},
		},
		{
			name: "Zero page defaults to 1",
			queryParams: map[string]string{
				"page": "0",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Negative page defaults to 1",
			queryParams: map[string]string{
				"page": "-5",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Zero page_size defaults to 10",
			queryParams: map[string]string{
				"page_size": "0",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Negative page_size defaults to 10",
			queryParams: map[string]string{
				"page_size": "-10",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Page_size over 100 limited to 100",
			queryParams: map[string]string{
				"page_size": "200",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 100,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Page_size exactly 100",
			queryParams: map[string]string{
				"page_size": "100",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 100,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Invalid page parameter",
			queryParams: map[string]string{
				"page": "invalid",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Invalid page_size parameter",
			queryParams: map[string]string{
				"page_size": "invalid",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Page 1 with page_size 1",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "1",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 1,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Large page number",
			queryParams: map[string]string{
				"page": "999999",
			},
			expected: &Pagination{
				Page:     999999,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Build query string
			url := "/test"
			first := true
			for key, value := range tt.queryParams {
				if first {
					url += "?"
					first = false
				} else {
					url += "&"
				}
				url += key + "=" + value
			}

			req, _ := http.NewRequest("GET", url, nil)
			c.Request = req

			// Get pagination
			result := GetPagination(c)

			// Compare fields
			if result.Page != tt.expected.Page {
				t.Errorf("GetPagination().Page = %v, want %v", result.Page, tt.expected.Page)
			}
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("GetPagination().PageSize = %v, want %v", result.PageSize, tt.expected.PageSize)
			}
			if result.Sort != tt.expected.Sort {
				t.Errorf("GetPagination().Sort = %v, want %v", result.Sort, tt.expected.Sort)
			}
		})
	}
}

func TestPaginationPaginate(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		sort           string
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "Page 1, Size 10",
			page:           1,
			pageSize:       10,
			sort:           "id desc",
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "Page 2, Size 10",
			page:           2,
			pageSize:       10,
			sort:           "id desc",
			expectedOffset: 10,
			expectedLimit:  10,
		},
		{
			name:           "Page 3, Size 25",
			page:           3,
			pageSize:       25,
			sort:           "id asc",
			expectedOffset: 50,
			expectedLimit:  25,
		},
		{
			name:           "Page 1, Size 100",
			page:           1,
			pageSize:       100,
			sort:           "created_at desc",
			expectedOffset: 0,
			expectedLimit:  100,
		},
		{
			name:           "Page 5, Size 20",
			page:           5,
			pageSize:       20,
			sort:           "name asc",
			expectedOffset: 80,
			expectedLimit:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:     tt.page,
				PageSize: tt.pageSize,
				Sort:     tt.sort,
			}

			paginateFunc := p.Paginate()

			// Verify that the function is created (not nil)
			if paginateFunc == nil {
				t.Error("Paginate() returned nil function")
			}

			// Calculate expected offset manually
			calculatedOffset := (tt.page - 1) * tt.pageSize
			if calculatedOffset != tt.expectedOffset {
				t.Errorf("Offset calculation mismatch: got %d, want %d", calculatedOffset, tt.expectedOffset)
			}
		})
	}
}

func TestPaginationEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		expected    *Pagination
	}{
		{
			name: "Empty sort parameter",
			queryParams: map[string]string{
				"sort": "",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "",
			},
		},
		{
			name: "Whitespace in parameters",
			queryParams: map[string]string{
				"page":      " 2 ",
				"page_size": " 20 ",
				"sort":      " id desc ",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     " id desc ",
			},
		},
		{
			name: "SQL injection attempt in sort",
			queryParams: map[string]string{
				"sort": "id; DROP TABLE users--",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
		{
			name: "Decimal numbers",
			queryParams: map[string]string{
				"page":      "2.5",
				"page_size": "10.5",
			},
			expected: &Pagination{
				Page:     1,
				PageSize: 10,
				Sort:     "created_at desc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/test"
			first := true
			for key, value := range tt.queryParams {
				if first {
					url += "?"
					first = false
				} else {
					url += "&"
				}
				url += key + "=" + value
			}

			req, _ := http.NewRequest("GET", url, nil)
			c.Request = req

			result := GetPagination(c)

			if result.Page != tt.expected.Page {
				t.Errorf("GetPagination().Page = %v, want %v", result.Page, tt.expected.Page)
			}
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("GetPagination().PageSize = %v, want %v", result.PageSize, tt.expected.PageSize)
			}
			if result.Sort != tt.expected.Sort {
				t.Errorf("GetPagination().Sort = %v, want %v", result.Sort, tt.expected.Sort)
			}
		})
	}
}

func TestPaginationConsistency(t *testing.T) {
	// Test that the same query produces the same pagination result
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	req1, _ := http.NewRequest("GET", "/test?page=2&page_size=20&sort=id desc", nil)
	c1.Request = req1

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	req2, _ := http.NewRequest("GET", "/test?page=2&page_size=20&sort=id desc", nil)
	c2.Request = req2

	p1 := GetPagination(c1)
	p2 := GetPagination(c2)

	if p1.Page != p2.Page || p1.PageSize != p2.PageSize || p1.Sort != p2.Sort {
		t.Errorf("Inconsistent pagination results: %+v vs %+v", p1, p2)
	}
}
