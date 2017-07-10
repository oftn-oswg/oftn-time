package main

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var config Config
var store *sessions.CookieStore

const kSessionName = "sess"

func loadConfiguration() error {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Load configuration
	err := loadConfiguration()
	if err != nil {
		log.Println("Error while loading configuration")
		log.Fatal(err)
	}
	initLogin()

	// Load session store
	store = sessions.NewCookieStore([]byte(config.SessionSecret))

	// Load server
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.HandleFunc("/", serveMain)
	http.HandleFunc("/login", serveLogin)
	http.HandleFunc("/login/github", serveLoginGithub)
	http.HandleFunc("/logout", serveLogout)
	http.HandleFunc("/callback", serveLoginCallback)

	port := config.Port
	if port == 0 {
		port = 8080
	}

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func serveMain(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, _ := store.Get(r, kSessionName)
	_, ok := session.Values["id"]
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}

	template, err := template.ParseFiles("templates/main.html", "templates/dashboard.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = template.Execute(w, session.Values)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
