package dto

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestAdminListArticleReqToDTOPublishedRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/admin/article?published_from=2026-10-01&published_to=2026-10-02", nil)

	result, err := AdminListArticleReqToDTO(c)
	if err != nil {
		t.Fatalf("AdminListArticleReqToDTO err: %v", err)
	}
	from := time.Date(2026, 10, 1, 0, 0, 0, 0, time.Local).Unix()
	to := time.Date(2026, 10, 3, 0, 0, 0, 0, time.Local).Unix()
	if result.PublishedFrom == nil || *result.PublishedFrom != from {
		t.Fatalf("published_from = %v, want %d", result.PublishedFrom, from)
	}
	if result.PublishedTo == nil || *result.PublishedTo != to {
		t.Fatalf("published_to = %v, want %d", result.PublishedTo, to)
	}
}

func TestAdminListArticleReqToDTORejectsInvalidPublishedDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/admin/article?published_from=2026/10/01", nil)

	if _, err := AdminListArticleReqToDTO(c); err == nil {
		t.Fatal("expected invalid published date error")
	}
}
