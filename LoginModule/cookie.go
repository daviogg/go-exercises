package main

import "net/http"

func setSession(u *User, w http.ResponseWriter) {
	value := map[string]string{
		"name": u.Username,
		"pass": u.Password,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getUserName(r *http.Request) (username string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
		}
	}
	return username
}

func clearSession(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func getMsg(w http.ResponseWriter, r *http.Request, name string) (msg string) {
	if cookie, err := r.Cookie(name); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(name, cookie.Value, &cookieValue); err == nil {
			msg = cookieValue[name]
			clearSession(w, name)
		}
	}

	return msg
}

func setMsg(w http.ResponseWriter, name string, msg string) {
	value := map[string]string{
		name: msg,
	}
	if encoded, err := cookieHandler.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}
