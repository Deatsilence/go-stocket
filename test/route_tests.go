package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Deatsilence/go-stocket/pkg/middleware"
	"github.com/Deatsilence/go-stocket/routes"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Authenticate())
	return r
}

func TestUserRoutes(t *testing.T) {
	r := setupRouter()
	routes.UserRoutes(r)

	t.Run("GetUsers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUser", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestProductRoutes(t *testing.T) {
	r := setupRouter()
	routes.ProductRoutes(r)

	t.Run("AddAProduct", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/products/add", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("DeleteAProduct", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/products/delete/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetProducts", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetProduct", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("SearchByBarcodePrefix", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products/search", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateAProduct", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/products/update/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateSomePropertiesOfProduct", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/api/products/updatepartially/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestPasswordRoutes(t *testing.T) {
	r := setupRouter()
	routes.PasswordRoutes(r)

	t.Run("RequestPasswordReset", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/passwordreset/request", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ResetPassword", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/passwordreset/confirm", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ChangePassword", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/passwordreset/changepassword", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthRoutes(t *testing.T) {
	r := setupRouter()
	routes.AuthRoutes(r)

	t.Run("VerifyEmail", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/users/verifyemail", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("SignUp", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/users/signup", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Login", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/users/login", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/users/logout", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
