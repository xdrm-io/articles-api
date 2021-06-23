package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/xdrm-io/articles-api/model"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// DB using sqlite3 driver
type DB struct {
	db *sqlx.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS user(
	user_id INTEGER,
	username VARCHAR(100) NOT NULL,
	firstname VARCHAR(200) NOT NULL,
	lastname VARCHAR(200) NOT NULL,

	PRIMARY KEY(user_id)
);
CREATE TABLE IF NOT EXISTS article(
	article_id INTEGER,
	title VARCHAR(255) NOT NULL,
	body text NOT NULL,
	author_ref INTEGER,

	PRIMARY KEY(article_id),
	FOREIGN KEY(author_ref) REFERENCES user(user_id)
);
CREATE TABLE IF NOT EXISTS vote(
	user_ref INTEGER,
	article_ref VARCHAR(100) NOT NULL,
	value INTEGER NOT NULL,

	PRIMARY KEY(user_ref, article_ref),
	FOREIGN KEY(user_ref) REFERENCES user(user_id),
	FOREIGN KEY(article_ref) REFERENCES article(article_id),
	CHECK( value = 1 OR value = -1 )
);
`

// Open the database
func (s *DB) Open() error {
	db, err := sqlx.Connect("sqlite3", "./local.db")
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInitFailed, err)
	}
	s.db = db

	// create tables
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInitFailed, err)
	}
	_, err = tx.Exec(schema)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInitFailed, err)
	}
	tx.Commit()
	return nil
}

// Close the databse
func (s *DB) Close() error {
	return s.db.Close()
}

// CreateUser ...
func (s *DB) CreateUser(username, firstname, lastname string) (*model.User, error) {
	user := model.User{
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
	}

	// insert user
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, ErrTransaction
	}
	res, err := tx.NamedExec(`INSERT INTO
		user(username, firstname, lastname)
		VALUES(:username,:firstname,:lastname);`, &user)
	if err != nil {
		return nil, ErrCreate
	}

	// find inserted id
	insertedID, err := res.LastInsertId()
	if err != nil {
		return nil, ErrNotFound
	}
	tx.Commit()
	user.ID = uint(insertedID)

	// return user info
	return &user, nil
}

// UpdateUser ...
func (s *DB) UpdateUser(id uint, username, firstname, lastname *string) (*model.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, ErrTransaction
	}

	// update fields
	if username != nil {
		_, err := tx.Exec(`UPDATE user SET username = ? WHERE user_id = ?;`, *username, id)
		if err != nil {
			return nil, ErrUpdate
		}

	}
	if firstname != nil {
		_, err := tx.Exec(`UPDATE user SET firstname = ? WHERE user_id = ?;`, *firstname, id)
		if err != nil {
			return nil, ErrUpdate
		}
	}
	if lastname != nil {
		_, err := tx.Exec(`UPDATE user SET lastname = ? WHERE user_id = ?;`, *lastname, id)
		if err != nil {
			return nil, ErrUpdate
		}
	}

	// select updated user
	user := model.User{}
	err = tx.Get(&user, `SELECT user_id, username, firstname, lastname FROM user WHERE user_id = $1;`, id)
	if err != nil {
		return nil, ErrNotFound
	}
	tx.Commit()

	return &user, nil
}

// DeleteUser ...
func (s *DB) DeleteUser(id uint) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return ErrTransaction
	}

	res, err := tx.Exec("DELETE FROM user WHERE user_id = ?;")
	if err != nil {
		return ErrDelete
	}

	mustBeOne, err := res.RowsAffected()
	if err != nil || mustBeOne < 1 {
		return ErrDelete
	}
	tx.Commit()

	return nil
}

// GetUserByID ...
func (s *DB) GetUserByID(id uint) (*model.User, error) {
	user := model.User{}

	err := s.db.Get(&user, `SELECT user_id, username, firstname, lastname FROM user WHERE user_id = $1;`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, ErrUnexpected
	}

	return &user, nil
}

// GetAllUsers ...
func (s *DB) GetAllUsers() ([]model.User, error) {
	users := []model.User{}
	err := s.db.Select(&users, "SELECT user_id, username, firstname, lastname FROM user;")
	if err != nil {
		return nil, ErrNotFound
	}
	return users, nil
}

// CreateArticle ...
func (s *DB) CreateArticle(title, body string, author uint) (*model.Article, error) {
	article := model.Article{
		Title:  title,
		Body:   body,
		Author: author,
	}

	// insert article
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, ErrTransaction
	}

	res, err := tx.NamedExec(`INSERT INTO
		article(title, body, author_ref)
		VALUES(:title,:body,:author_ref);`, &article)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, ErrCreate
	}

	// find inserted id
	insertedID, err := res.LastInsertId()
	if err != nil {
		return nil, ErrNotFound
	}
	tx.Commit()
	article.ID = uint(insertedID)

	return &article, nil
}

// UpdateArticle ...
func (s *DB) UpdateArticle(id uint, title *string, body *string) (*model.Article, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, ErrTransaction
	}

	if title != nil {
		_, err := tx.Exec(`UPDATE article SET title = ? WHERE article_id = ?;`, *title, id)
		if err != nil {
			return nil, ErrUpdate
		}
	}
	if body != nil {
		_, err := tx.Exec(`UPDATE article SET body = ? WHERE article_id = ?;`, *body, id)
		if err != nil {
			return nil, ErrUpdate
		}
	}

	// select updated article
	article, err := s.GetArticleByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	tx.Commit()

	return article, nil
}

// DeleteArticle ...
func (s *DB) DeleteArticle(id uint) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return ErrTransaction
	}

	res, err := tx.Exec("DELETE FROM article WHERE article_id = ?;")
	if err != nil {
		return ErrDelete
	}

	mustBeOne, err := res.RowsAffected()
	if err != nil || mustBeOne < 1 {
		return ErrDelete
	}
	tx.Commit()

	return nil
}

// GetArticleByID ...
func (s *DB) GetArticleByID(id uint) (*model.Article, error) {
	article := model.Article{}

	err := s.db.Get(&article, `SELECT article_id, title, body, author_ref, COALESCE(SUM(vote.value),0) as score
		FROM article
		LEFT OUTER JOIN vote ON article.article_id = vote.article_ref
		WHERE article_id = $1
		GROUP BY article_id, title, body, author_ref
		LIMIT 1;`, id)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, ErrNotFound
	}

	return &article, nil
}

// GetArticlesByAuthor ...
func (s *DB) GetArticlesByAuthor(id uint) ([]model.Article, error) {
	articles := []model.Article{}

	err := s.db.Select(&articles, `SELECT article_id, title, body, author_ref, COALESCE(SUM(vote.value),0) as score
		FROM article
		LEFT OUTER JOIN vote ON article.article_id = vote.article_ref
		WHERE author_ref = $1
		GROUP BY article_id, title, body, author_ref;`, id)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, ErrNotFound
	}

	return articles, nil
}

// GetAllArticles ...
func (s *DB) GetAllArticles() ([]model.Article, error) {
	articles := []model.Article{}

	err := s.db.Select(&articles, `SELECT article_id, title, body, author_ref, COALESCE(SUM(vote.value),0) as score
		FROM article
		LEFT OUTER JOIN vote ON article.article_id = vote.article_ref
		GROUP BY article_id, title, body, author_ref;`)
	if err != nil {
		return nil, ErrNotFound
	}

	return articles, nil
}

// UpVote ...
func (s *DB) UpVote(user uint, article uint) (*model.Vote, error) {
	// update if alreay exists
	vote := model.Vote{
		User:    user,
		Article: article,
		Value:   1,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, ErrTransaction
	}

	// alreay exists -> update
	row := tx.QueryRow("SELECT user_ref, article_ref, value FROM vote WHERE article_ref=? AND user_ref=?;", vote.Article, vote.User)
	if row != nil && row.Scan() != sql.ErrNoRows {
		_, err := tx.Exec(`UPDATE vote SET value = ? WHERE article_ref = ? AND user_ref = ?;`, vote.Value, vote.Article, vote.User)
		if err != nil {
			return nil, ErrUpdate
		}

		tx.Commit()
		return &vote, nil
	}

	// create vote
	_, err = tx.Exec("INSERT INTO vote(user_ref, article_ref, value) VALUES(?, ?, ?);", vote.User, vote.Article, vote.Value)
	if err != nil {
		return nil, ErrCreate
	}
	tx.Commit()

	return &vote, nil
}

// DownVote ...
func (s *DB) DownVote(user uint, article uint) (*model.Vote, error) {
	// update if alreay exists
	vote := model.Vote{
		User:    user,
		Article: article,
		Value:   -1,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, ErrTransaction
	}

	// alreay exists -> update
	row := tx.QueryRow("SELECT user_ref, article_ref, value FROM vote WHERE article_ref=? AND user_ref=?;", vote.Article, vote.User)
	if row != nil && row.Scan() != sql.ErrNoRows {
		_, err := tx.Exec(`UPDATE vote SET value = ? WHERE article_ref = ? AND user_ref = ?;`, vote.Value, vote.Article, vote.User)
		if err != nil {
			return nil, ErrUpdate
		}

		tx.Commit()
		return &vote, nil
	}

	// create vote
	_, err = tx.Exec("INSERT INTO vote(user_ref, article_ref, value) VALUES(?, ?, ?);", vote.User, vote.Article, vote.Value)
	if err != nil {
		return nil, ErrCreate
	}

	tx.Commit()
	return &vote, nil
}
