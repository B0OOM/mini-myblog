package miniblog

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) handleAdminIndexPage(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	limit := 10

	posts, err := GetPosts(h.db, int(page), limit)
	if err != nil {
		log.Println("handleAdminIndexPage GetPosts", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx := gin.H{
		"posts":      posts,
		"pagination": GetPagination(h.db, int(page), limit),
	}
	c.HTML(http.StatusOK, "admin.html", ctx)
}

// 提供一个空白的文章模版，用于创建新文章
func (h *Handlers) handleAdminCreatePostPage(c *gin.Context) {
	c.HTML(http.StatusOK, "edit.html", nil)
}

func (h *Handlers) handleAdminCreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBind(&post); err != nil {
		log.Println("handleAdminCreatePost ShouldBind", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println("post:", post)
	post.ID = 0
	err := CreatePost(h.db, &post)
	if err != nil {
		log.Println("handleAdminCreatePost CreatePost", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article creation successful"})
}

func (h *Handlers) handleAdminUpdatePostPage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	post, err := GetPostByID(h.db, uint(id))
	if err != nil {
		log.Println("handleAdminUpdatePostPage GetPostByID", err)
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"post": post,
	})
}

func (h *Handlers) handleAdminUpdatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBind(&post); err != nil {
		log.Println("handleAdminUpdatePost ShouldBind", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := UpdatePost(h.db, &post)
	if err != nil {
		log.Println("handleAdminUpdatePost UpdatePost", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

func (h *Handlers) handleAdminDeletePost(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("post_id"), 10, 64)
	DeletePostByID(h.db, uint(id))
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (h *Handlers) handleAdminDeleteComment(c *gin.Context) {
	postId, _ := strconv.ParseInt(c.Query("post_id"), 10, 64)
	commentID, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)

	DeleteCommentByID(h.db, uint(commentID))

	post, err := GetPostByID(h.db, uint(postId))
	if err != nil {
		log.Println("handleAdminDeleteComment GetPostByID", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"post": post,
	})
}
