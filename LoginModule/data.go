package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//User define an user for login
type User struct {
	UUID     string
	Fname    string
	Lname    string
	Username string
	Email    string
	Password string
}

func saveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (firstname text, lastname text, username text, email text, password text)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (firstname, lastname, username, email, password) values (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

func userExists(u *User) bool {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	var ps, us string
	q, err := db.Query("select username, password from users where username = '" + u.Username + "' and password = '" + u.Password + "'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&us, &ps)
	}
	if us == u.Username && ps == u.Password {
		return true
	}
	return false
}
