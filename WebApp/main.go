package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Page define a page
type Page struct {
	Title string
	Body  []byte
}

var (
	tmplView = template.Must(template.New("test").ParseFiles("base.html", "test.html", "index.html"))
	tmplEdit = template.Must(template.New("edit").ParseFiles("base.html", "edit.html", "index.html"))
	db, _    = sql.Open("sqlite3", "cache/web.db")
	createDB = "create table if not exists pages (title text, body blob, timestamp text)"
)

func (p *Page) saveCache() error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	f := "cache/" + p.Title + ".txt"
	db.Exec(createDB)
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into pages (title, body, timestamp) values (?, ?, ?)")
	_, err := stmt.Exec(p.Title, p.Body, timestamp)
	tx.Commit()
	ioutil.WriteFile(f, p.Body, 0600)
	return err
}

func load(title string) (*Page, error) {
	f := "cache/" + title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadSource(title string) (*Page, error) {
	var name string
	var body []byte
	q, err := db.Query("select title, body from pages where title = '" + title + "' order by timestamp Desc limit 1")
	if err != nil {
		return nil, err
	}
	for q.Next() {
		q.Scan(&name, &body)
	}
	return &Page{Title: name, Body: body}, nil
}

func view(w http.ResponseWriter, r *http.Request) {

	title := r.URL.Path[len("/test/"):]
	p, err := loadSource(title)
	if err != nil {
		p, _ = load(title)
	}
	if p.Title == "" {
		p, _ = load(title)
	}

	tmplView.ExecuteTemplate(w, "base", p)
	//t, _ := template.ParseFiles("test.html")
	//t.Execute(w, p)
}

func edit(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadSource(title)
	if err != nil {
		p, _ = load(title)
	}
	if p.Title == "" {
		p, _ = load(title)
	}

	tmplEdit.ExecuteTemplate(w, "base", p)
	//t, _ := template.ParseFiles("edit.html")
	//t.Execute(w, p)
}

func save(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.saveCache()
	http.Redirect(w, r, "/test/"+title, http.StatusFound)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/test/", view)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/save/", save)
	http.ListenAndServe(":8000", nil)
}
