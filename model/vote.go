package model

// Vote representation
type Vote struct {
	User    uint `json:"user"    db:"user_ref"`
	Article uint `json:"article" db:"article_ref"`
	Value   int  `json:"value"   db:"value"`
}
