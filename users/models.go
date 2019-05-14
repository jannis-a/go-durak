package users

import (
	"log"
	"time"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/utils"
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

type UserCreate struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func New(a *app.App, username string, email string, password string) User {
	// Hash the password
	hashedPassword, err := utils.Argon2Hash(password, a.Argon2Params)
	if err != nil {
		log.Panic(err)
	}

	// Insert data into database
	qry := `insert into users (username, email, password) values ($1, $2, $3) 
          returning id, username, email, joined_at`
	res := a.DB.QueryRow(qry, username, email, hashedPassword)

	// Scan result of insert into struct
	var user User
	err = res.Scan(&user.Id, &user.Username, &user.Email, &user.JoinedAt)
	if err != nil {
		log.Fatal(err)
	}

	return user
}
