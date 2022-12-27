package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	rssreader "github.com/Alexander96/rssreader"
	jwt "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("best-kept-secret")

var users = map[string]string{
	"admin": "admin",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type RSSInput struct {
	Urls []string `json:"urls"`
}
type RSSOutput struct {
	Items []rssreader.RssItem
}

func Login(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok {
		fmt.Println("Error parsing basic auth")
		w.WriteHeader(401)
		return
	}
	expectedPassword, ok := users[user]

	if !ok || expectedPassword != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("token", tokenString)
	w.Header().Set("content-type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "success"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func HandleRSS(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var rssInput RSSInput
	errDec := json.NewDecoder(r.Body).Decode(&rssInput)
	if errDec != nil {
		http.Error(w, errDec.Error(), http.StatusBadRequest)
		return
	}

	items, errRss := rssreader.Parse(rssInput.Urls)
	if errRss != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

	output := RSSOutput{Items: items}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(output)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("token", tokenString)
	w.Header().Set("content-type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "success"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/rss", HandleRSS)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/refresh", Refresh)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}
