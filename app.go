package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.HandleFunc("/", serveTemplate)

	http.ListenAndServe(":8080", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Grab template file
	layout := "layouts/main.html"
	page := filepath.Join("pages", filepath.Clean(r.URL.Path))

	// Return 404 if page doesn't exist
	info, err := os.Stat(page)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Redirect to index page if request is for a directory
	if info.IsDir() {
		http.Redirect(w, r, filepath.Join(filepath.Clean(r.URL.Path), "index.html"), 302)
		return
	}

	template, err := template.ParseFiles(layout, page)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
