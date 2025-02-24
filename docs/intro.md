# 博客系统的设计要点

博客是最适合编程入门的项目之一，因为它涵盖了许多基本概念，如数据库、用户认证、前端和后端框架等。

## 如何开始运行这个代码
1. 确保有go的加速代理，go需要运行1.22以上版本，建议到 [golang官网](https://golang.google.cn/dl/) 下载最新版本, 根据自己的需要选择对应的版本
2. 安装go之后，需要设置代理，可以使用：   
    `go env -w GOPROXY=https://goproxy.cn` 设置加速代理
3. 更新当前代码的所有依赖，这个动作类似于`mvn clean install`，可以使用：      
    `go mod download` 下载所有依赖

看一下现在的代码结构:
```bash
.
├── README.cn.md
├── README.md
├── admins.go
├── admins_test.go
├── cmd
│   ├── main.go
│   ├── miniblog.db
│   ├── static
│   └── templates
│       ├── 404.html
│       ├── admin.html
│       ├── edit.html
│       ├── index.html
│       └── post.html
├── docs
├── go.mod
├── go.sum
├── handlers.go
├── handlers_test.go
├── models.go
└── models_test.go
```
所有的代码都在一级目录下，`cmd`是程序的入口，对应`cmd/main.go`

所以要运行代码，需要进入`cmd`目录，然后运行`go run .`，这样就可以启动服务了

启动服务之后， 默认会运行在`http://localhost:8080`，可以在浏览器中输入这个地址查看博客

## 博客设计要点
我们要设计一个最简单的博客，包括以下几个功能：

1. 能发布文章，每个文章都有自己的标题和内容，还有最重要的路径，这个路径就是我们的`slug`, 用来标识这个文章，这样可以让文章的链接更加友好，不会出现大量:`?id=2`这样的链接
2. 能发表评论，普通用户可以发表评论，管理员可以删除评论
3. 有一个后台管理系统，可以查看所有的文章和评论，可以删除评论，可以编辑文章，网站上线后，这个后台系统只有管理员可以访问
4. 需要一个可以更换主题的模版，只需要修改html就可以实现改变外观，这样可以让网站更加个性化
5. 需要支持markdown编辑，文章的编写格式更加适合表达
6. 需要支持分页，当文章数量过多的时候，需要分页显示
7. 需要支持404页面，当用户访问不存在的页面的时候，需要显示一个友好的页面
8. 需要支持静态文件，比如图片，css，js等，这些文件需要单独存放，不需要放在数据库中
9. 需要支持表单提交，比如评论的提交，文章的提交，这些都需要通过表单提交

以上是我们的博客的设计要点，接下来我们会一步步实现这些功能

## 实现第一个main函数，实现第一个index路由

第一步就是我们要实现一个main函数，实现第一个路由，可以查看`cmd/main.go`的代码：

我们用gin的框架，启动了一个httpserver, 运行在`8080`端口，然后定义了一个路由，当用户访问`/`的时候，会显示一个`index.html`的页面

`cmd/main.go`的代码如下：

```go
	r := gin.Default()
	h := miniblog.NewHandlers(db, password)
	h.RegisterHandlers(r)
	log.Println()
	log.Println("🎉 Starting server at", addr)
	log.Println("🏆 by https://ruzhila.cn")
	log.Println()

	r.Run(addr)
```

这个就是入口函数，我们在这里注册了一个路由，然后启动了一个httpserver，这样就可以访问我们的博客了

`RegisterHandlers`是在`handlers.go`中，我们可以看到index的路由是这样注册的：

```go
func (h *Handlers) RegisterHandlers(engine *gin.Engine) {
	engine.FuncMap["markdown"] = markdownToHTML
	engine.FuncMap["strtime"] = strtimeFormat
	engine.LoadHTMLGlob(filepath.Join(h.templateDir, "*.html"))

	engine.GET("/", h.handleIndexPage)
    ...
}
```
这样就注册了一个路由，当用户访问`/`的时候，会调用`handleIndexPage`这个函数, 实现读取数据库，然后渲染页面的功能
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