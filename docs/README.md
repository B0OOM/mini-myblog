# 入职啦博客中文指导

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

## 课程目录
### 项目介绍
[讲解一个博客系统的设计要点](./intro.md)
  - 如何开始运行这个代码
  - 博客设计要点
  - 实现第一个main函数，实现第一个index路由

### 文章管理功能开发
[设计文章系统，基于Markdown实现文章的编辑](./article.md)
  - 路由设计
  - 数据库的库表设计
  - 数据库的CRUD 
  - 模版渲染
  - markdown的使用
  
### 评论功能开发
[实现评论功能，熟悉表单提交](./comments.md)
  - 评论的设计
  - 为什么不用前后端分离提交表单
  - 如何减少不必要的数据库查询代码
  
### 后台管理功能开发
[实现管理员后台，方便做文章、评论的管理](./admin.md)
  - 后台路由, 后台应该包括什么功能
  - 后台的权限管理实现

### 项目部署和启动
[将自己的代码发布到线上，支持MySQL，更专业的服务发布上线](./devops.md)
  - 支持MySQL数据库
  - 如何部署到自己的服务器

## 如何扩展
[下一步，你可以尝试实现以下功能：](./nextstep.md)
  1. 实现复杂路由，实现sitemap.xml自动化生成
  2. 实现用户系统，支持多用户
  3. 实现前后端分离，提示操作体验
  4. 实现单元测试，提升开发效率
  5. 实现高级编辑器，支持图片上传