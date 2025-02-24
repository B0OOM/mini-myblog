# 实现评论功能，熟悉表单提交

实现了路由、模版渲染之后，我们就需要实现表单提交，完成用户评论的功能。

## 评论的设计

评论的设计是比较简单的，每个文章都有自己的评论，所以我们根据前面的库表设计，只需要让`Comment`这个对象有外键`PostID`，就可以关联到对应的文章。

```go
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

获取文章的时候，会把相应的评论一起读取出来，这样可以减少很多读取的代码，方便开发。
```go
func GetPosts(db *gorm.DB, page, limit int) ([]Post, error) {
	if page < 1 {
		page = 1
	}
	if limit < 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	var posts []Post

	tx := db.Offset(offset).Limit(limit).Preload("Comments").Order("position DESC").Order("updated_at DESC")
	if err := tx.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
```
其中的`Preload("Comments")`就是把文章的评论一起读取出来。这个是gorm的一个特性，可以一次性读取出来，减少很多读取的代码。

## 为什么不用前后端分离提交表单

我们这次项目没有采用前后端分离的方案实现表单提交，给大家看一下最朴实的表单提交是什么样子的， 可以看一下 `post.html`

```html
<form action="" method="post">
	<input type="hidden" name="slug" value="{{ .post.Slug }}">
	<p class="text-xl font-semibold">留言</p>
	<div class="flex items-center gap-x-6 mt-4 py-2">
		<p class="hang">
			<label for="name">* 显示名称</label>
			<input type="text" name="name" id="name" required class="round">
		</p>
		<p class="hang">
			<label for="email">* 联系邮箱地址(不会显示)</label>
			<input type="email" name="email" id="email" required class="round">
		</p>
	</div>

	<p class="hang">
		<label for="content">* 留言内容</label>
		<textarea name="content" id="content" required class="round h-20"></textarea>
	</p>
	<p class="flex justify-end w-full">
		<button type="submit" class="bg-gray-900 text-gray-100 px-3 py-1.5 rounded-md">发布评论</button>
	</p>
</form>
```
这个是HTML自带的表单，通过`POST`的方式提交到服务器，然后服务器处理这个表单，把数据存储到数据库中。

```go

func (h *Handlers) RegisterHandlers(engine *gin.Engine) {
	...
	engine.POST("/post/:slug", h.handleCreateComment)
	...
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
	c.Redirect(http.StatusFound, "/post/"+c.Param("slug"))
}
```

通过当前页面实现评论的提交，这样可以减少很多前后端分离的代码，让大家更容易理解。

1. 先通过`c.ShouldBind(&form)`把表单的数据绑定到`CreateCommentForm`这个结构体上
2. 然后通过`IsPostExist`判断这个文章是否存在
3. 最后通过`CreateComment`创建评论，存储到数据库中

## 如何减少不必要的数据库查询代码

我们在`GetPosts`的时候，通过`Preload("Comments")`一次性把文章的评论一起读取出来，这样可以减少很多读取的代码，方便开发。

通过ShoudBind把表单的数据绑定到结构体上，表单的校验是由`gin`自带实现，这样可以减少很多校验的代码。

最后，当我们提交了表单之后，页面会重新刷新，重新执行模版渲染的流程，这样就不需要前端去重新读取comment渲染的代码，可以有效减少接口和代码的复杂度。

当然不是前后端分离也会带来更多的数据传输和数据库访问，这个是需要权衡的，不同的项目有不同的需求，需要根据实陵情况来选择。

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