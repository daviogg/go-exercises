package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()

func indexPage(w http.ResponseWriter, r *http.Request) {
	u := &User{}
	tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html")
	err := tmpl.ExecuteTemplate(w, "base", u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	pass := r.FormValue("password")
	u := &User{Username: name, Password: pass}
	redirect := "/"
	if name != "" && pass != "" {
		if userExists(u) {
			setSession(u, w)
			redirect = "/example"
		}
	}
	http.Redirect(w, r, redirect, 302)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func examplePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("base.html", "index.html", "internal.html")
	username := getUserName(r)
	if username != "" {
		err := tmpl.ExecuteTemplate(w, "base", &User{Username: username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("signup.html", "index.html", "base.html")
		u := &User{}
		tmpl.ExecuteTemplate(w, "base", u)
	case "POST":
		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("email")
		un := r.FormValue("userName")
		pass := r.FormValue("password")

		u := &User{Fname: f, Lname: l, Email: em, Username: un, Password: pass}
		saveData(u)
		http.Redirect(w, r, "/", 302)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/example", examplePage)
	router.HandleFunc("/signup", signup).Methods("POST", "GET")
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
