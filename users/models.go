package users

import (
	"database/sql"
	"time"
)

type UserPub struct {
	Id       uint      `json:"id"`
	Username string    `json:"username"`
	JoinedAt time.Time `json:"joined_at"`
}

type User struct {
	UserPub
	Email    string `json:"email"`
	Password string `json:"-"`
}

func New(db *sql.DB, username string, email string, password string) User {
	var user User
	qry := `insert into users (username, email, password) values ($1, $2, $3) returning id, username, email, joined_at`
	res := db.QueryRow(qry, username, email, password)
	err := res.Scan(&user.Id, &user.Username, &user.Email, &user.JoinedAt)
	if err != nil {
		println(err.Error())
	}
	return user
}
