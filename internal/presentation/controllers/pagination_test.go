package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func createPaginationContext(queryString string) *gin.Context {
	req := httptest.NewRequest(http.MethodGet, "/test?"+queryString, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}

func TestParsePagination_Defaults(t *testing.T) {
	c := createPaginationContext("")
	params := parsePagination(c)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.Limit)
}

func TestParsePagination_CustomValues(t *testing.T) {
	c := createPaginationContext("page=3&limit=50")
	params := parsePagination(c)
	assert.Equal(t, 3, params.Page)
	assert.Equal(t, 50, params.Limit)
}

func TestParsePagination_LimitCappedAt100(t *testing.T) {
	c := createPaginationContext("limit=200")
	params := parsePagination(c)
	// limit=200 fails the l <= 100 check, so it stays at default 20
	assert.Equal(t, 20, params.Limit)
}

func TestParsePagination_InvalidPage(t *testing.T) {
	c := createPaginationContext("page=-1")
	params := parsePagination(c)
	// page=-1 fails the p > 0 check, so it stays at default 1
	assert.Equal(t, 1, params.Page)
}

func TestParsePagination_NonNumeric(t *testing.T) {
	c := createPaginationContext("page=abc&limit=xyz")
	params := parsePagination(c)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.Limit)
}
