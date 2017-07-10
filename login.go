package main

import (
	"encoding/hex"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

const (
	kGithubAuthorizeURL = "https://github.com/login/oauth/authorize"
	kGithubTokenURL     = "https://github.com/login/oauth/access_token"
)

var auth *oauth2.Config

func initLogin() {
	auth = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  kGithubAuthorizeURL,
			TokenURL: kGithubTokenURL,
		},
	}
}

func serveLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, kSessionName)
	delete(session.Values, "id")
	delete(session.Values, "token")
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

func serveLogin(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles("templates/main.html", "templates/login.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = template.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

}

func serveLoginGithub(w http.ResponseWriter, r *http.Request) {
	// Generate random state variable for this sign-in request
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	// Store this state variable in the session
	session, _ := store.Get(r, kSessionName)
	session.Values["state"] = state
	session.Save(r, w)

	// Send state variable along with request
	url := auth.AuthCodeURL(state)

	http.Redirect(w, r, url, 302)
}

func serveLoginCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, kSessionName)

	if r.URL.Query().Get("state") != session.Values["state"] {
		// Also a possible CSRF
		http.Error(w, "Cookies not enabled", 400)
		return
	}

	token, err := auth.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Could not get GitHub token", 500)
		return
	}

	if !token.Valid() {
		http.Error(w, "Retrieved invalid GitHub token", 500)
		return
	}

	client := github.NewClient(auth.Client(oauth2.NoContext, token))
	user, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		http.Error(w, "Could not get authenticated user information", 500)
		return
	}

	session.Values["id"] = user.ID
	session.Values["token"] = token.AccessToken
	delete(session.Values, "state")
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
