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
	msg := getMsg(w, r, "message")
	var u = &User{}
	u.Errors = make(map[string]string)
	if msg != "" {
		u.Errors["message"] = msg
		render(w, "signin", u)
	} else {
		u := &User{}
		render(w, "signin", u)
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
		} else {
			setMsg(w, "message", "please signup or enter a valid username and password!")
		}
	} else {
		setMsg(w, "message", "Username or Password field are empty!")
	}
	http.Redirect(w, r, redirect, 302)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w, "session")
	http.Redirect(w, r, "/", 302)
}

func examplePage(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	if username != "" {
		render(w, "internal", &User{Username: username})
	} else {
		setMsg(w, "message", "Please login first!")
		http.Redirect(w, r, "/", 302)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &User{}
		render(w, "signup", u)
	case "POST":
		u := &User{Fname: r.FormValue("fName"),
			Lname:    r.FormValue("lName"),
			Email:    r.FormValue("email"),
			Username: r.FormValue("userName"),
			Password: r.FormValue("password"),
		}
		saveData(u)
		http.Redirect(w, r, "/", 302)
	}
}

func render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseGlob("*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(w, name, data)
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
