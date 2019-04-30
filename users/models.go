package users

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/raja/argon2pw"
)

type User struct {
	Id       uint      `json:"-"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	// UpdatedAt pq.NullTime `json:"updated_at" db:"updated_at"`
}

func NewUser(db *sqlx.DB, username string, email string, password string) User {
	// language=SQL
	qry := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *`
	result := db.QueryRowx(qry, username, email, HashPassword(password))

	var user User
	_ = result.StructScan(&user)
	return user
}

func HashPassword(plaintext string) string {
	hashed, err := argon2pw.GenerateSaltedHash(plaintext)
	if err != nil {
		log.Panicf("Hash generated returned error: %v", err)
	}

	return hashed
}
