package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/pbnjk/backend/personal-blog/article"
	"github.com/pbnjk/backend/personal-blog/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := article.LoadArticles(); err != nil {
		log.Fatal(err)
	}

	s := server.NewServer()
	s.Serve()
}
