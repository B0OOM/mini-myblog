# 实现管理员后台，方便做文章、评论的管理

所有的网站产品都应该提供管理后台，而不是通过修改数据库的方式来管理数据。毕竟产品上线之后，应该由产品经理、运营人员来管理数据，而不是开发人员。

这样可以更方便的管理数据，更安全的管理数据。

## 后台路由, 后台应该包括什么功能

一个博客后台应该包括以下功能：
1. 文章管理 可以发布和编辑文章
2. 评论管理 可以删除评论
3. 用户管理 可以查看用户信息
4. 资源管理 可以查看静态资源，比如图片，css，js等

这个版本我们就实现了`1`和`2`这两个功能，后续版本会实现`3`和`4`这两个功能。

先定义了admin的路由:
```go
func (h *Handlers) RegisterHandlers(engine *gin.Engine) {
	admin := engine.Group("/admin", gin.BasicAuth(gin.Accounts{"admin": h.password}))
	admin.GET("/", h.handleAdminIndexPage)
	admin.GET("/create", h.handleAdminCreatePostPage)
	admin.POST("/create", h.handleAdminCreatePost)
	admin.GET("/edit", h.handleAdminUpdatePostPage)
	admin.POST("/edit", h.handleAdminUpdatePost)
	admin.GET("/delete/post", h.handleAdminDeletePost)
	admin.GET("/delete/comment", h.handleAdminDeleteComment)
}

```

这里，实现了最基本的文章、评论的管理工作，每个都有对应的代码实现，大家可以自己看看实现

不同的是，发现多了`GET /create`和`POST /create`这两个路由，需要一个`GET`的路由来显示页面，然后一个`POST`的路由来处理表单提交。

## 后台的权限管理实现

多了一个`gin.BasicAuth`的方法，这个方法是用来做基本的认证的，只有用户名和密码正确的时候，才能访问这个路由。

因为admin是一个`Group`，所以这个Group下面的所有路由都需要认证。 这样可以确保只有权限的人才能访问管理后台

h.password是在`main.go`中传入的，可以查看main.go关于密码的生成, 虽然是一个小的demo,但是仍然要求不允许用简单的密码，这是最基本的安全常识:

```go
    var password string
	flag.StringVar(&password, "password", "", "Admin password")
	flag.Parse()

	if password == "" {
		password = RandText(8)
	}
```

也就是如果启动的时候，指定`-password`那么就可以用你指定的密码，如果没有指定，那么就会生成一个随机的密码。

这样可以确保任何时候都不会用简单的密码。


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