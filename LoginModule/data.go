package main

import (
	"database/sql"

	"crypto/rand"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//User define an authenticated one
type User struct {
	UUID     string
	Username string
	Password string
	Fname    string
	Lname    string
	Email    string
	Errors   map[string]string
}

func saveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (uuid text, firstname text, lastname text, username text, email text, password text)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (uuid, firstname, lastname, username, email, password) values (?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.UUID, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

func userExists(u *User) bool {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	var ps, us string
	q, err := db.Query("select username, password from users where username = '" + u.Username + "'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&us, &ps)
	}
	pw := bcrypt.CompareHashAndPassword([]byte(ps), []byte(u.Password))
	if us == u.Username && pw == nil {
		return true
	}
	return false
}

func encryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

//UUID generate new unique identifier
func UUID() (id string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return
	}
	id = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
