# 设计文章系统，基于Markdown实现文章的编辑

## 课程目标
这个章节，我们讲一下一个博客最基本的路由、库表、模版渲染的基本功能。

也就是定义了一篇文章是如何存储、读取并且被浏览器访问的。

## 路由设计
路由设计是所有网站产品设计的第一步，需要规划好URL访问路径，方便管理和维护。

一个博客提供三种基本的访问需求：
- 首页
- 文章页
- 静态资源

我们的路由设计如下：
- 首页：`/`
- 文章页：`/post/:slug`
- 静态资源：`/static/*`

这样可以确保我们的网站能够正常访问，当我们要访问一个文章的时候，可以通过`/post/:slug`的方式访问。
比如`/post/welcome` 就会访问到`welcome`这篇文章。

参考`handlers.go`的代码：
```go
func (h *Handlers) RegisterHandlers(engine *gin.Engine) {
	
	engine.GET("/", h.handleIndexPage)
	engine.GET("/post/:slug", h.handlePost)
	engine.POST("/post/:slug", h.handleCreateComment)
	engine.StaticFS("/static", http.Dir("static"))
	engine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})
    ....
}
```

然后我们还留了一个404页面，当用户访问不存在的页面的时候，会显示一个404页面。

## 数据库的库表设计

我们的博客需要存储文章和评论，所以我们需要两个表，一个是文章表，一个是评论表。

我们定义了一个简单的文章表，包括以下字段：代码可以参考`models.go`
```go
type Post struct {
	ID        uint      `form:"id" gorm:"primarykey"`
	CreatedAt time.Time `form:"-"`
	UpdatedAt time.Time `form:"-"`
	Slug      string    `form:"slug" gorm:"unique;size:200"`
	Title     string    `form:"title" gorm:"size:200"`
	Desc      string    `form:"desc" gorm:"size:2000"`
	Keywords  string    `form:"keywords" gorm:"size:200"`
	Content   string    `form:"content"`
	Position  int       `form:"position"`
	Author    string    `form:"author" gorm:"size:200"`
	Comments  []Comment `form:"-"`
}

type Comment struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	PostID    uint
	Email     string `gorm:"size:200"`
	Author    string `gorm:"size:200"`
	Content   string
	IP        string `gorm:"size:64"`
}
```

其中`Post`表存储文章的信息，`Comment`表存储评论的信息。
- **Post** 这个又个`slug`是全局唯一的，用来标识这篇文章，也就是不允许重复，我们可以用`gorm:"unique"`来定义这个字段是唯一的。
  - 每个文章有对应的评论，所以我们在`Post`表中定义了一个`Comments`字段，当读取Post的时候，我们可以一次性把所有的相关的评论读取出来。
  - 这样可以减少很多读取的代码，方便开发
- **Comment** 这个表存储评论的信息，包括评论的作者、内容、IP等信息。
  - Comment有个`PostID`字段，用来关联到对应的文章，这样我们可以通过`PostID`来读取这篇文章的所有评论。

每个gorm的表，都有个`CreatedAt`和`UpdatedAt`字段，用来记录创建时间和更新时间，这两个字段是不需要用户输入的，是gorm自动维护的。

另外，很多java的程序员喜欢存储时间戳，一般不建议存储时间戳，因为会导致用sql操作的时候，需要再换算时间戳，也不方便阅读，都是建议存储DateTime类型，这样可以方便的操作时间。

在orm的时候，也会自动换成`time.Time`类型，方便操作。

## 数据库的CRUD
使用了orm之后，操作数据库就不需要写繁琐的sql语句了，只需要调用orm的方法就可以了。

具体可以查看`models.go`中的代码，我们定义了一些基本的CRUD操作，比如：
- `GetPostBySlug` 通过`slug`读取文章
- `GetPosts` 按照分页取文章列表

另外我们在CRUD之前，其实做了一个初始化库表的动作，在`ConnectDatabase`中：

```go
func ConnectDatabase(dbfile, dbdriver string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if dbdriver == "" || dbdriver == "sqlite" {
		dialector = sqlite.Open(dbfile)
	} else if dbdriver == "mysql" {
		dialector = mysql.Open(dbfile)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	isInit := !db.Migrator().HasTable(&Post{})
	err = db.AutoMigrate(&Post{}, &Comment{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connected", dbfile)

	if isInit {
        // 初始化默认的文章
		// 
    }
```

gorm会自动根据结构体生成对应的表，如果表不存在，会自动创建，如果表存在，会自动更新表结构。

我们在创建表之前可以判断是不是第一次生成，这样就可以把默认的文章生成到库里面，这个功能是非常常见，并且有用的技巧，特别是生成初始化数据的时候。

## 模版渲染
文章读取之后，我们就需要对内容进行渲染

我们采用了go内置的`html/template`包，这个包是一个非常强大的模版渲染引擎，可以方便的渲染html页面。

我们看一下`handlers.go`中的代码：
```go
func (h *Handlers) handleIndexPage(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	limit := 10

	posts, err := GetPosts(h.db, int(page), limit)
	if err != nil {
		log.Println("handleIndexPage GetPosts", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts":      posts,
		"pagination": GetPagination(h.db, int(page), limit),
	})
}
```

从db中获取了文章之后，我们调用`c.HTML`方法，这个方法会把`index.html`模版渲染出来，同时，会给渲染的模版传递一个`gin.H`的map，这个map中包含了需要渲染的数据：
- `posts` 文章列表
- `pagination` 分页信息

着两个字段对应会在`index.html`中使用，我们看一下`index.html`的代码：
```html
 {{ range .posts }}
    <li class="border-b py-4 px-4 mt-2 group rounded bg-white shadow">
        <div>
            <a href="/post/{{ .Slug }}">{{.Title}}</a>
		</div>
	</li>
{{ end }}
```
这个模版语法是go的模版语法，`{{ .posts }}`表示遍历`posts`这个数组，然后访问每个文章的`Slug`和`Title`字段。这样就完成了模版的渲染工作

## markdown的使用

我们的文章是采用markdown格式存储的，这样可以方便的编辑文章，在模版中如何渲染markdown呢？

1. 我们需要依赖一个markdown的库，比如`github.com/gomarkdown/markdown` 用于将markdown转换成html
2. 渲染的动作我们可以写在go中，也可以通过扩展模版函数的方式，直接在模版中调用：
   
```go

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

func (h *Handlers) RegisterHandlers(engine *gin.Engine) {
	engine.FuncMap["markdown"] = markdownToHTML
	...
}
```
我们注册了一个`markdown`的函数，然后在模版中可以直接调用这个函数(查看`post.html`的代码)：
```html
 <div class="markdown-content py-2">{{ markdown .post.Content }}</div>
```

这样就可以把markdown转换成html，然后渲染到页面上了。



课程内容学习遇到问题，可以联系我们老师进行沟通（微信号：jinti2000），我们会及时更新课程内容。

<div style="display: inline-block;text-align: center;">
   <div style="display: inline-block;">
	 <h3>入职啦实战项目</h3>
	 <img src="../cmd/static/projectQrcode.jpg" width="200" margin-right="100" alt="入职啦实战项目二维码" >
   </div>
   <div style="display: inline-block; margin-left: 30px;">
	 <h3>入职啦微信公众号</h3>
	 <img src="../cmd/static/weixinQrcode.png" width="200" alt="入职啦公众号二维码" />
   </div>
 </div>

[课程介绍](./README.md)