package article

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Blog struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	Content       string `json:"content"`
}

var blog Blog

// Loads the articles from the JSON file
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

// Saves the articles to the JSON file
func SaveArticles() error {
	j, err := json.Marshal(&blog)
	if err != nil {
		return err
	}

	return os.WriteFile("articles.json", j, 0644)
}

// Adds an article to the blog
func AddArticle(title, content string) {
	latestArticleID := 0
	if n := len(blog.Articles); n > 0 {
		latestArticleID = blog.Articles[n-1].ID
	}

	currentTime := getCurrentTime()
	a := Article{
		ID:            latestArticleID + 1,
		Title:         title,
		Content:       content,
		DatePublished: currentTime,
		DateModified:  currentTime,
	}

	blog.Articles = append(blog.Articles, a)
	SaveArticles()
}

// Retrieves an article by its ID
func GetArticleByID(id int) (Article, error) {
	idx, err := findArticleByID(id)
	if err != nil {
		return Article{}, err
	}

	return blog.Articles[idx], nil
}

// Edits an article by its ID
func EditArticleByID(id int, title, content string) error {
	idx, err := findArticleByID(id)
	if err != nil {
		return err
	}

	a := &blog.Articles[idx]
	if title != "" {
		a.Title = title
	}

	if content != "" {
		a.Content = content
	}

	a.DateModified = getCurrentTime()
	fmt.Println(title, content)

	SaveArticles()
	return nil
}

// Deletes an article by its ID
func DeleteArticleByID(id int) error {
	idx, err := findArticleByID(id)
	if err != nil {
		return err
	}

	// Removes at index
	blog.Articles = append(blog.Articles[:idx], blog.Articles[idx+1:]...)
	SaveArticles()

	return nil
}

// Returns the blog
func GetBlog() Blog {
	return blog
}

// Performs binary search for article
//
// We don't just index into the array because the adding & removal of articles
// might eventually make so that the IDs no longer correspond directly to the
// indices
func findArticleByID(id int) (int, error) {
	articles := blog.Articles

	l, r, idx := 0, len(articles)-1, len(articles)/2
	for l <= r {
		idx = l + ((r - l) / 2)
		a := articles[idx]
		if a.ID < id {
			l = idx + 1
		} else if a.ID > id {
			r = idx - 1
		} else {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("no such article with ID #'%d'!", id)
}

// Gets the current time, formatted as ISO 8601 time
func getCurrentTime() string {
	now := time.Now()
	return now.Format("2006-01-02T15:04:05Z")
}
