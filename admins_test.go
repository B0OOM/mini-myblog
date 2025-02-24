package miniblog

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminAccess(t *testing.T) {
	db, _ := ConnectDatabase(":memory:", "")
	h := NewHandlers(db, "hello")
	h.password = "world"
	h.templateDir = "cmd/templates" // For unit test
	t.Run("Access unauthorized admin page", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)

		req, _ := http.NewRequest("GET", "/admin/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusUnauthorized)
	})
	t.Run("Access admin page", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)

		req, _ := http.NewRequest("GET", "/admin/", nil)
		req.SetBasicAuth("admin", h.password)
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusOK)
	})
	t.Run("Create Post via amdin", func(t *testing.T) {
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)

			req, _ := http.NewRequest("GET", "/admin/create", nil)
			req.SetBasicAuth("admin", h.password)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusOK)
		}
		// dup create
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)
			body := url.Values{}
			body.Add("slug", "welcome")
			data := body.Encode()
			req, _ := http.NewRequest("POST", "/admin/create", bytes.NewBufferString(data))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("admin", h.password)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusInternalServerError)
		}
		var newId uint
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)
			body := url.Values{}
			body.Add("slug", "test")
			body.Add("content", "Test content")
			body.Add("title", "Test title")
			body.Add("desc", "Test desc")
			body.Add("keywords", "Go, Gin, Gorm")
			body.Add("position", "1")
			data := body.Encode()
			req, _ := http.NewRequest("POST", "/admin/create", bytes.NewBufferString(data))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("admin", h.password)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusFound)

			post, err := GetPostBySlug(db, "test")
			assert.Nil(t, err)
			assert.Equal(t, post.Title, "Test title")
			newId = post.ID
		}
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)
			getReq, _ := http.NewRequest("GET", fmt.Sprintf("/admin/edit?id=%d", newId), nil)
			getReq.SetBasicAuth("admin", h.password)
			r.ServeHTTP(w, getReq)
			assert.Equal(t, w.Code, http.StatusOK)
		}
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)
			body := url.Values{}
			body.Add("id", fmt.Sprintf("%d", newId))
			body.Add("slug", "test")
			body.Add("content", "Test content-2")
			data := body.Encode()
			req, _ := http.NewRequest("POST", "/admin/edit", bytes.NewBufferString(data))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("admin", h.password)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusFound)

			post, err := GetPostBySlug(db, "test")
			assert.Nil(t, err)
			assert.Equal(t, post.Content, "Test content-2")
		}
	})

	t.Run("Delete comment", func(t *testing.T) {
		p, _ := IsPostExist(db, "welcome")
		CreateComment(db, p.ID, "bob@ruzhil.cn", "bob", "Hello, comment", "test")

		p, _ = GetPostByID(db, p.ID)
		assert.Equal(t, len(p.Comments), 1)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/delete/comment?post_id=%d&comment_id=%d", p.ID, p.Comments[0].ID), nil)
		req.SetBasicAuth("admin", h.password)
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)

		p, _ = GetPostByID(db, p.ID)
		assert.Equal(t, len(p.Comments), 0)
	})
	t.Run("Delete post", func(t *testing.T) {
		p, _ := IsPostExist(db, "welcome")
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/delete/post?post_id=%d", p.ID), nil)
		req.SetBasicAuth("admin", h.password)
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusFound)

		_, err := IsPostExist(db, "welcome")
		assert.NotNil(t, err)
	})
}
