package model

// User representation
type User struct {
	ID        uint      `json:"user_id"   db:"user_id"`
	Username  string    `json:"username"  db:"username"`
	Firstname string    `json:"firstname" db:"firstname"`
	Lastname  string    `json:"lastname"  db:"lastname"`
	Articles  []Article `json:"articles"  db:"articles"`
}
