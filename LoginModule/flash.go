package main

import (
	"encoding/base64"
	"net/http"
	"time"
)

func encode(s []byte) string {
	return base64.URLEncoding.EncodeToString(s)
}

func decode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

func setMsg(w http.ResponseWriter, name string, msg []byte) {
	c := &http.Cookie{Name: name, Value: encode(msg)}
	http.SetCookie(w, c)
}

func getMsg(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return nil, err
	}
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(0, 1)}
	http.SetCookie(w, dc)
	return value, nil
}
