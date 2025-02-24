package miniblog

import (
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

type CreateCommentForm struct {
	Email   string `form:"email" binding:"required,email"`
	Name    string `form:"name" binding:"required"`
	Content string `form:"content" binding:"required"`
}

type Handlers struct {
	db          *gorm.DB
	password    string
	templateDir string
}

func NewHandlers(db *gorm.DB, password string) *Handlers {
	return &Handlers{db: db, password: password, templateDir: "templates"}
}

func markdownToHTML(content string) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	result := string(markdown.Render(doc, renderer))
	return template.HTML(result)
}

func strtimeFormat(t reflect.Value, format string) string {
	f := t.MethodByName("Format")
	if !f.IsValid() {
		return ""
	}
	return f.Call([]reflect.Value{reflect.ValueOf(format)})[0].String()
}

func (h *Handlers) RegisterHandlers(engine *gin.Engine) {
	//engine.FuncMap["markdown"] = markdownToHTML
	//engine.FuncMap["strtime"] = strtimeFormat
	//engine.LoadHTMLGlob(filepath.Join(h.templateDir, "*.html"))

	engine.POST("/posts", h.handleIndexPage)
	engine.GET("/post/:slug", h.handlePost)
	engine.POST("/post/:slug/comments", h.handleCreateComment)
	//engine.StaticFS("/static", http.Dir("static"))
	engine.NoRoute(func(c *gin.Context) {
		//c.HTML(http.StatusNotFound, "404.html", nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	admin := engine.Group("/admin", gin.BasicAuth(gin.Accounts{"admin": h.password}))
	admin.GET("/api/posts", h.handleAdminIndexPage) //获取管理员页面文章列表
	//admin.GET("/create", h.handleAdminCreatePostPage) //创建空白文章模版
	admin.POST("/post", h.handleAdminCreatePost)             //创建文章
	admin.GET("/edit", h.handleAdminUpdatePostPage)          //返回更新文章内容
	admin.POST("/edit", h.handleAdminUpdatePost)             //更新文章
	admin.GET("/delete/post", h.handleAdminDeletePost)       //删除文章
	admin.GET("/delete/comment", h.handleAdminDeleteComment) //删除评论
}

func (h *Handlers) handleIndexPage(c *gin.Context) {
	//page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	//limit := 10

	var request struct {
		Page  int `json:"page" binding:"required,number,min=1"`
		Limit int `json:"limit" binding:"required,number,min=1,max=100"`
	}

	if err := c.ShouldBind(&request); err != nil {
		log.Println("handleIndexPage ShouldBind", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	page := int64(request.Page)
	limit := request.Limit

	posts, err := GetPosts(h.db, int(page), limit)
	if err != nil {
		log.Println("handleIndexPage GetPosts", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts":      posts,
		"pagination": GetPagination(h.db, int(page), limit),
	})
}

func (h *Handlers) handlePost(c *gin.Context) {
	post, err := GetPostBySlug(h.db, c.Param("slug"))
	if err != nil {
		log.Println("handlePost GetPostBySlug", err)
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (h *Handlers) handleCreateComment(c *gin.Context) {
	var form CreateCommentForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println("handleCreateComment ShouldBind", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	post, err := IsPostExist(h.db, c.Param("slug"))
	if err != nil {
		log.Println("handleCreateComment IsPostExist", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := CreateComment(h.db, post.ID, form.Email, form.Name, form.Content, c.ClientIP()); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println("Comment created for post", post.ID, "by", form.Name, form.Email, c.ClientIP(), form.Content)
	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}
