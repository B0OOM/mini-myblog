package main

import (
	"flag"
	"log"

	"math/rand"

	"github.com/gin-gonic/gin"
	miniblog "github.com/ruzhila/mini-blog"
)

func RandText(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	var text = make([]byte, length)
	for i := range text {
		text[i] = chars[rand.Intn(len(chars))]
	}
	return string(text)
}

func main() {
	var addr string
	var dsn string
	var password string
	var dbdriver string = "sqlite"

	flag.StringVar(&addr, "addr", ":8080", "HTTP Serve address")
	flag.StringVar(&dsn, "db", "miniblog.db", "Database file name or connection string")
	flag.StringVar(&dbdriver, "dbdriver", "sqlite", "Database driver, sqlite or mysql")
	flag.StringVar(&password, "password", "", "Admin password")
	flag.Parse()

	//if password == "" {
	//	password = RandText(8)
	//}

	//修改为固定密码
	password = "2053097205"

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Addr        =", addr)
	log.Println("DSN         =", dsn)
	log.Println("DBdriver    =", dbdriver)
	log.Println()

	db, err := miniblog.ConnectDatabase(dsn, dbdriver)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	h := miniblog.NewHandlers(db, password)
	h.RegisterHandlers(r)

	log.Println()
	log.Println("🎉 Starting server at", addr)
	log.Println("🔧 Admin account: admin, password:", "\033[31m", password, "\033[0m", "visit /admin to manage posts")
	log.Println("🏆 by https://ruzhila.cn")
	log.Println()

	r.Run(addr)
}
