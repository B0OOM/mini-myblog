package miniblog

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownRender(t *testing.T) {
	c := markdownToHTML(`# Hello, world!`)
	assert.Equal(t, c, template.HTML("<h1 id=\"hello-world\">Hello, world!</h1>\n"))
}
func TestTimeFmt(t *testing.T) {
	now := time.Now()
	{
		tm := strtimeFormat(reflect.ValueOf(now), "2006-01-02")
		assert.Equal(t, tm, fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day()))
	}
	{
		tm := strtimeFormat(reflect.ValueOf(now), "15:04")
		assert.Equal(t, tm, fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()))
	}
}

func TestPostPage(t *testing.T) {
	db, _ := ConnectDatabase(":memory:", "")
	h := NewHandlers(db, "hello")
	h.templateDir = "cmd/templates" // For unit test
	t.Run("Access index page", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)
		h.handleIndexPage(c)
		assert.Equal(t, w.Code, 200)
		body := w.Body.String()
		assert.Contains(t, body, defaultPostTitle)
	})

	t.Run("Access not exist slug", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)
		c.AddParam("slug", "hello")
		h.handlePost(c)
		assert.Equal(t, w.Code, http.StatusNotFound)
		assert.Contains(t, w.Body.String(), "抱歉找不到这篇文章啦")
	})
	t.Run("Access exist slug", func(t *testing.T) {
		CreatePost(db, &Post{
			Slug:  "hello",
			Title: "Hello",
			Desc:  "Hello, world",
		})
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		c.AddParam("slug", "hello")
		h.RegisterHandlers(r)
		h.handlePost(c)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}

func TestCreateComment(t *testing.T) {
	db, _ := ConnectDatabase(":memory:", "")
	h := NewHandlers(db, "hello")
	h.templateDir = "cmd/templates" // For unit test

	{
		w := httptest.NewRecorder()
		body := url.Values{}
		data := body.Encode()
		req, _ := http.NewRequest("POST", "/post/welcome", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		_, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusBadRequest)
	}
	{
		body := url.Values{}
		body.Add("name", "bob")
		body.Add("email", "bob@ruzhila.cn")
		body.Add("content", "Hello, comment")
		data := body.Encode()
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)

			req, _ := http.NewRequest("POST", "/post/not-exists", bytes.NewBufferString(data))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusNotFound)
		}
		{
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			h.RegisterHandlers(r)
			//
			req, _ := http.NewRequest("POST", "/post/welcome", bytes.NewBufferString(data))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusFound)

			p, _ := GetPostBySlug(db, "welcome")
			assert.Equal(t, len(p.Comments), 1)
		}
	}
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/post/welcome", nil)
		_, r := gin.CreateTestContext(w)
		h.RegisterHandlers(r)
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Contains(t, w.Body.String(), "Hello, comment")
	}
}
