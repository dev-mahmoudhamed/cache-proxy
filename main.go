package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type ResponseData struct {
	Cached   bool
	HTML     string
	Duration string
	URL      string
}

var (
	tmpl = template.Must(template.ParseFiles("index.html"))
	mu   sync.Mutex
)

func homehandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	session, sessionID := getSession(r)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// prevent browser caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")

	// clear server cache safely
	mu.Lock()
	clear(session.Cache)
	session.History = make([]ResponseData, 0)
	mu.Unlock()

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, sessionID := getSession(r)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	url := appendPrefix(r.FormValue("url"))
	var res ResponseData

	start := time.Now()
	res.URL = url

	mu.Lock()
	html, found := session.Cache[url]
	mu.Unlock()

	if found {
		res.Cached = true
		res.HTML = html
	} else {
		var err error
		res.HTML, err = DownloadHTML(url)
		if err != nil {
			log.Printf("failed to download page to db: %v", err)
		}
		mu.Lock()
		session.Cache[url] = res.HTML
		mu.Unlock()
	}

	res.Duration = prettyDuration(time.Since(start))
	mu.Lock()
	session.History = append(session.History, res)
	historyToRender := make([]ResponseData, len(session.History))
	copy(historyToRender, session.History)
	mu.Unlock()

	err := tmpl.Execute(w, historyToRender)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homehandler)
	http.HandleFunc("/fetch", fetchHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
