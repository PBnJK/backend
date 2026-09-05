package main

import (
	"log"

	"github.com/pbnjk/backend/personal-blog/article"
	"github.com/pbnjk/backend/personal-blog/server"
)

func main() {
	if err := article.LoadArticles(); err != nil {
		log.Fatal(err)
	}

	server.Run()
}
