package article

import (
	"encoding/json"
	"fmt"
	"os"
)

type Blog struct {
	Articles []Article
}

type Article struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	Content       string `json:"content"`
}

var blog Blog

func LoadArticles() error {
	f, err := os.ReadFile("articles.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(f, &blog); err != nil {
		return err
	}

	return nil
}

func GetBlog() Blog {
	return blog
}

func GetArticleByID(id int) (Article, error) {
	if len(blog.Articles) < id {
		return Article{}, fmt.Errorf("no article matches ID '%d'", id)
	}

	return blog.Articles[id-1], nil
}
