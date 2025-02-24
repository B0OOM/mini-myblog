# 用golang + gin + gorm编写的最简单博客

最简单的gin + gorm 博客,  总代码不超过1000行 🎉

[中文指导](./docs/README.md)

## 特点
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

## 如何运行

> 中文用户需要先配置一下goproxy.cn的加速代理
### 先增加go的包管理加速：goproyx.cn 配置

这一步是必须的，不然会导致下载依赖包失败

```bash
go env -w GOPROXY=https://goproxy.cn
```
## 运行代码 

```bash
cd mini-blog
go mod download
go get ./...

cd cmd
go run .
# 启动之后，可以访问 `http://localhost:8080` 查看博客
```

## 运行单元测试
```bash
cd mini-blog
go mod download
go get ./...

go test -v ./...
# 这样就可以查看测试的结果了

```

## 管理后台

启动之后，可以访问管理后台 `http://localhost:8080/admin`

默认的管理员账号是: `admin` 如果你启动的时候不指定密码(`-password`) 那么会随机生成一个密码

**密码在控制台中能看到**

# 🙋‍♀️学习更多100行实战项目
 想要学习更多100行实战项目，可关注公众号<strong>「入职啦」</strong>，或加入学习交流群，每日分享有趣实用的100行代码实战项目，更有老师在线答疑解惑，帮助你快速提升编程能力。

<div style="display: inline-block;text-align: center;">
   <div style="display: inline-block;">
     <h3>入职啦实战项目</h3>
     <img src="./cmd/static/projectQrcode.jpg" width="200" margin-right="100" alt="入职啦实战项目二维码" >
   </div>
   <div style="display: inline-block; margin-left: 30px;">
     <h3>入职啦微信公众号</h3>
     <img src="./cmd/static/weixinQrcode.png" width="200" alt="入职啦公众号二维码" />
   </div>
 </div>