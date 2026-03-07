package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"artifactor/internal/repository"
	httprequests "artifactor/pkg/http"
	"artifactor/pkg/tokens"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FetchToken(rawToken string) (*tokens.Token, error) {
	args := m.Called(rawToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tokens.Token), args.Error(1)
}

func (m *MockAuthRepository) CreateToken(request *httprequests.CreateRequest) (string, error) {
	args := m.Called(request)
	return args.String(0), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 401 when authorization header is missing", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		handler := AuthMiddleware(mockRepo)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header missing")
		mockRepo.AssertNotCalled(t, "FetchToken")
	})

	t.Run("returns 401 when token fetch returns error", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		mockRepo.On("FetchToken", "invalid_token").Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request.Header.Set(API_HEADER, "invalid_token")

		handler := AuthMiddleware(mockRepo)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
		mockRepo.AssertCalled(t, "FetchToken", "invalid_token")
	})

	t.Run("returns 401 when token is not found", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		mockRepo.On("FetchToken", "nonexistent_token").Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request.Header.Set(API_HEADER, "nonexistent_token")

		handler := AuthMiddleware(mockRepo)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid api token")
		mockRepo.AssertCalled(t, "FetchToken", "nonexistent_token")
	})

	t.Run("continues to next handler when token is valid", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		mockRepo.On("FetchToken", "valid_token").Return(&tokens.Token{
			Data: "token_data",
			Permissions: tokens.TokenPermissions{
				Admin:  true,
				Upload: true,
				Delete: false,
			},
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request.Header.Set(API_HEADER, "valid_token")

		handler := AuthMiddleware(mockRepo)
		handler(c)

		admin, exists := c.Get("admin")
		assert.True(t, exists)
		assert.True(t, admin.(bool))
		mockRepo.AssertCalled(t, "FetchToken", "valid_token")
	})

	t.Run("continues to next handler with admin false when token has no admin permission", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		mockRepo.On("FetchToken", "user_token").Return(&tokens.Token{
			Data: "token_data",
			Permissions: tokens.TokenPermissions{
				Admin:  false,
				Upload: true,
				Delete: false,
			},
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request.Header.Set(API_HEADER, "user_token")

		handler := AuthMiddleware(mockRepo)
		handler(c)

		admin, exists := c.Get("admin")
		assert.True(t, exists)
		assert.False(t, admin.(bool))
	})
}

var _ repository.AuthRepoInterface = (*MockAuthRepository)(nil)
