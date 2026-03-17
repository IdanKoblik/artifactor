package product

import (
	"artifactor/internal/utils"
	"artifactor/pkg/types"
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newUploadRequest(t *testing.T, product, version, filename, content string) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if product != "" {
		require.NoError(t, writer.WriteField("product", product))
	}
	if version != "" {
		require.NoError(t, writer.WriteField("version", version))
	}
	if filename != "" {
		part, err := writer.CreateFormFile("file", filename)
		require.NoError(t, err)
		_, err = part.Write([]byte(content))
		require.NoError(t, err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func productWithToken(token string, perms types.TokenPermissions) *types.Product {
	return &types.Product{
		Name:     "myproduct",
		Tokens:   map[string]types.TokenPermissions{utils.Hash(token): perms},
		Versions: map[string]types.Version{},
	}
}

func TestHandleUpload_InvalidForm(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "", "", "", "")

	handler := &ProductHandler{Repo: &mockProductRepo{}}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleUpload_FetchProductError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(nil, errors.New("db error"))

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "db error")
	repo.AssertExpectations(t)
}

func TestHandleUpload_ProductNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(nil, nil)

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
	repo.AssertExpectations(t)
}

func TestHandleUpload_PermissionDenied_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")
	c.Set("admin", false)
	c.Set("token", "mytoken")

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(
		productWithToken("mytoken", types.TokenPermissions{Upload: true}), nil,
	)

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "permission denied")
	repo.AssertExpectations(t)
}

func TestHandleUpload_PermissionDenied_NoUploadPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")
	c.Set("admin", true)
	c.Set("token", "mytoken")

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(
		productWithToken("mytoken", types.TokenPermissions{Upload: false}), nil,
	)

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "permission denied")
	repo.AssertExpectations(t)
}

func TestHandleUpload_VersionAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")
	c.Set("admin", true)
	c.Set("token", "mytoken")

	product := productWithToken("mytoken", types.TokenPermissions{Upload: true})
	product.Versions["1.0.0"] = types.Version{Path: "/some/path", Checksum: "abc"}

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(product, nil)

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "already uploaded")
	repo.AssertExpectations(t)
}

func TestHandleUpload_AddVersionError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")
	c.Set("admin", true)
	c.Set("token", "mytoken")

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(
		productWithToken("mytoken", types.TokenPermissions{Upload: true}), nil,
	)
	repo.On("AddVersion", "myproduct", "1.0.0", "mytoken", true, mock.Anything).
		Return(errors.New("db error"))

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "db error")
	repo.AssertExpectations(t)
}

func TestHandleUpload_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = newUploadRequest(t, "myproduct", "1.0.0", "artifact.zip", "data")
	c.Set("admin", true)
	c.Set("token", "mytoken")

	repo := &mockProductRepo{}
	repo.On("FetchProduct", "myproduct").Return(
		productWithToken("mytoken", types.TokenPermissions{Upload: true}), nil,
	)
	repo.On("AddVersion", "myproduct", "1.0.0", "mytoken", true, mock.Anything).Return(nil)

	handler := &ProductHandler{Repo: repo}
	handler.HandleUpload(c)

	assert.Equal(t, http.StatusCreated, c.Writer.Status())
	repo.AssertExpectations(t)
}
