package model

// Article representation
type Article struct {
	ID     uint   `json:"article_id" db:"article_id"`
	Author uint   `json:"author"     db:"author_ref"`
	Title  string `json:"title"      db:"title"`
	Body   string `json:"body"       db:"body"`
	Score  int    `json:"score"      db:"score"`
}
