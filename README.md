# A mini blog by golang + gin + gorm

A simple & clean blog by golang. Loc is less 1000 🎉

[中文用户](./README.cn.md)
  
# Features
- 🐶 GIN + GORM basic example
- 🤭 Builtin golang template
- 🎉 Admin console
- 📚 Database with sqlite
- 👌 Markdown content
- 👍 82% code cover

## How to run

```
cd mini-blog
go mod download
go get ./...

cd cmd
go run .
```

After start, visit `http://localhost:8080` to see the blog

## Admin console

Default admin account: `admin`, the password is random genrate if without `-password`

The password is output on startup message:

Visit `http://localhost:8080/admin`
