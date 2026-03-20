package product

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleListProducts_Unauthorized_NoAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/product/list", nil)

	handler := &ProductHandler{Repo: &mockProductRepo{}}
	handler.HandleListProducts(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleListProducts_Unauthorized_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/product/list", nil)
	c.Set("admin", false)

	handler := &ProductHandler{Repo: &mockProductRepo{}}
	handler.HandleListProducts(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleListProducts_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/product/list", nil)
	c.Set("admin", true)

	repo := &mockProductRepo{}
	repo.On("ListProducts").Return(nil, errors.New("db error"))

	handler := &ProductHandler{Repo: repo}
	handler.HandleListProducts(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "db error")
	repo.AssertExpectations(t)
}

func TestHandleListProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/product/list", nil)
	c.Set("admin", true)

	names := []string{"productA", "productB"}
	repo := &mockProductRepo{}
	repo.On("ListProducts").Return(names, nil)

	handler := &ProductHandler{Repo: repo}
	handler.HandleListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "productA")
	assert.Contains(t, w.Body.String(), "productB")
	repo.AssertExpectations(t)
}
