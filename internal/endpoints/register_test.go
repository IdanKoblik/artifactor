package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httprequests "artifactor/pkg/http"
	"artifactor/pkg/tokens"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthHandlerRepo struct {
	mock.Mock
}

func (m *MockAuthHandlerRepo) FetchToken(rawToken string) (*tokens.Token, error) {
	args := m.Called(rawToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tokens.Token), args.Error(1)
}

func (m *MockAuthHandlerRepo) CreateToken(request *httprequests.CreateRequest) (string, error) {
	args := m.Called(request)
	return args.String(0), args.Error(1)
}

func TestHandleRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 401 when admin context is missing", func(t *testing.T) {
		mockRepo := new(MockAuthHandlerRepo)
		handler := &AuthHandler{Repo: mockRepo}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", nil)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Only admin allowed to register new tokens")
	})

	t.Run("returns 401 when admin is false", func(t *testing.T) {
		mockRepo := new(MockAuthHandlerRepo)
		handler := &AuthHandler{Repo: mockRepo}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", nil)
		c.Set("admin", false)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Only admin allowed to register new tokens")
	})

	t.Run("returns 400 when request body is missing", func(t *testing.T) {
		mockRepo := new(MockAuthHandlerRepo)
		handler := &AuthHandler{Repo: mockRepo}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("admin", true)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Missing request body")
	})

	t.Run("returns 400 when request body is invalid", func(t *testing.T) {
		mockRepo := new(MockAuthHandlerRepo)
		handler := &AuthHandler{Repo: mockRepo}

		body := []byte(`{invalid json}`)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("admin", true)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 500 when CreateToken fails", func(t *testing.T) {
		mockRepo := new(MockAuthHandlerRepo)
		mockRepo.On("CreateToken", mock.AnythingOfType("*http.CreateRequest")).
			Return("", assert.AnError)

		handler := &AuthHandler{Repo: mockRepo}

		requestBody := httprequests.CreateRequest{
			Admin:  true,
			Upload: true,
			Delete: false,
		}
		body, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("admin", true)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
		mockRepo.AssertCalled(t, "CreateToken", mock.AnythingOfType("*http.CreateRequest"))
	})

	t.Run("returns token when registration succeeds", func(t *testing.T) {
		expectedToken := "generated-token-123"
		mockRepo := new(MockAuthHandlerRepo)
		mockRepo.On("CreateToken", mock.AnythingOfType("*http.CreateRequest")).
			Return(expectedToken, nil)

		handler := &AuthHandler{Repo: mockRepo}

		requestBody := httprequests.CreateRequest{
			Admin:  true,
			Upload: true,
			Delete: false,
		}
		body, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("admin", true)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, expectedToken, w.Body.String())
		mockRepo.AssertCalled(t, "CreateToken", mock.AnythingOfType("*http.CreateRequest"))
	})

	t.Run("passes correct permissions to CreateToken", func(t *testing.T) {
		mockRepo := new(MockAuthHandlerRepo)
		mockRepo.On("CreateToken", mock.AnythingOfType("*http.CreateRequest")).
			Return("token", nil).Once()

		handler := &AuthHandler{Repo: mockRepo}

		requestBody := httprequests.CreateRequest{
			Admin:  false,
			Upload: true,
			Delete: true,
		}
		body, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("admin", true)

		handler.HandleRegister(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockRepo.AssertCalled(t, "CreateToken", mock.MatchedBy(func(req *httprequests.CreateRequest) bool {
			return req.Admin == false && req.Upload == true && req.Delete == true
		}))
	})
}
