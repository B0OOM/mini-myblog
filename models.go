package miniblog

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const defaultPostTitle = "这是我的第一篇文章"
const defaultPostSlug = "welcome"
const defaultPostAuthor = "kui"
const defaultPostKeywords = "ruzhila, blog, golang, gin"
const defautlPostDesc = "这是我的第一篇文章，你可以编辑或删除它。采用ruzhila.cn的开源代码。"
const defaultPostContent = `# 👏 👏 欢迎

这是默认的第一篇文章，你可以编辑或删除它。

这是一个简单的博客系统，采用Golang和Gin框架开发。你可以在[GitHub](https://github.com/ruzhila/min-blog) 上找到源代码。

🏆 博客基于Golang, gin, gorm 实现，整体代码不到1000行
- 🎉 没有使用前端框架，前端页面使用纯HTML和CSS，不需要前后端分离
- 👍 82%的代码覆盖率，简单好理解的单元测试入门教程
- 📚 数据库支持sqlite和MySQL，单机就可以部署，也可以直接部署到生产环境上
- 🎁 不需要写SQL, 用orm完成所有的工作
- ✍️ 支持Markdown编辑
- 🔥 支持评论
- 📖 支持分页
- 🔧 自带一个管理后台
- 🚀 代码简单，适合新手学习，gin路由、模版、表单、认证一次性学会

用来学习Golang和Gin框架非常合适，也可以用来搭建自己的博客。


欢迎访问 入职啦 https://ruzhila.cn ， 我们提供简历优化和编程课程培训，帮你找到心仪好工作。
`

type Post struct {
	ID        uint      `form:"id" gorm:"primarykey"` // 主键
	CreatedAt time.Time `form:"-"`
	UpdatedAt time.Time `form:"-"`
	Slug      string    `form:"slug" gorm:"unique;size:200"` // 文章的URL
	Title     string    `form:"title" gorm:"size:200"`       // 文章标题
	Desc      string    `form:"desc" gorm:"size:2000"`       // 文章描述
	Keywords  string    `form:"keywords" gorm:"size:200"`    // 文章关键字
	Content   string    `form:"content"`                     // 文章内容
	Position  int       `form:"position"`                    // 文章排序
	Author    string    `form:"author" gorm:"size:200"`      // 文章作者
	Comments  []Comment `form:"-"`                           // 评论
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

type Pagination struct {
	Total int //总记录
	Page  int //当前页面
	Limit int
	Prev  int //上一页页码
	Next  int //下一页页码
	Last  int //最后一页页码
}

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
		defaultPost := Post{
			Title:    defaultPostTitle,
			Slug:     defaultPostSlug,
			Author:   defaultPostAuthor,
			Desc:     defautlPostDesc,
			Keywords: defaultPostKeywords,
			Content:  defaultPostContent,
		}
		err = CreatePost(db, &defaultPost)
		if err != nil {
			return nil, err
		}
		log.Println("Default post created")
	}
	return db, nil
}

// 创建文章
func CreatePost(db *gorm.DB, post *Post) error {
	return db.Create(post).Error
}

// 更新文章
func UpdatePost(db *gorm.DB, post *Post) error {
	return db.Save(post).Error
}

// 根据slug获取文章
func GetPostBySlug(db *gorm.DB, slug string) (Post, error) {
	var post Post
	if err := db.Where("slug", slug).Preload("Comments").First(&post).Error; err != nil {
		return Post{}, err
	}
	return post, nil
}

// 判断文章是否存在
func IsPostExist(db *gorm.DB, slug string) (Post, error) {
	var post Post
	if err := db.Where("slug", slug).First(&post).Error; err != nil {
		return Post{}, err
	}
	return post, nil
}

// 根据id获取文章
func GetPostByID(db *gorm.DB, id uint) (Post, error) {
	var post Post
	if err := db.Where("id", id).Preload("Comments").First(&post).Error; err != nil {
		return Post{}, err
	}
	return post, nil
}

// 创建评论
func CreateComment(db *gorm.DB, postID uint, email, author, content, ip string) error {
	comment := Comment{
		PostID:  postID,
		Email:   email,
		Author:  author,
		Content: content,
		IP:      ip,
	}
	return db.Create(&comment).Error
}

// 删除文章
func DeletePostByID(db *gorm.DB, postID uint) error {
	return db.Where("id", postID).Delete(&Post{}).Error
}

// 删除评论
func DeleteCommentByID(db *gorm.DB, commentID uint) error {
	return db.Where("id", commentID).Delete(&Comment{}).Error
}

// 获取文章列表
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

// 获取文章总数
func GetPostsCount(db *gorm.DB) int {
	var total int64
	db.Model(&Post{}).Count(&total)
	return int(total)
}

func GetPagination(db *gorm.DB, page, limit int) (p Pagination) {
	p.Total = GetPostsCount(db)
	if page < 1 {
		page = 1
	}
	p.Page = page
	p.Prev = page - 1
	p.Limit = limit

	if p.Prev < 1 {
		p.Prev = 1
	}
	p.Next = page + 1
	if p.Total <= int(page)*limit {
		p.Next = page
	}
	p.Last = (p.Total + limit - 1) / limit
	if p.Last < 1 {
		p.Last = 1
	}
	return
}
