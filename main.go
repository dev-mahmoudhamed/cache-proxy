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
	tmpl       = template.Must(template.ParseFiles("index.html"))
	cacheStore = make(map[string]string)
	//history    = make([]ResponseData, 0)
	mu sync.Mutex
)

func homehandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// prevent browser caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// clear server cache safely
	mu.Lock()
	clear(cacheStore)
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
	url := appendPrefix(r.FormValue("url"))
	var res ResponseData

	start := time.Now()
	res.URL = url

	mu.Lock()
	html, found := cacheStore[url]
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
		cacheStore[url] = res.HTML
		mu.Unlock()
	}

	res.Duration = prettyDuration(time.Since(start))
	//history = append(history, res)
	// tmpl.Execute(w, history)

	err := tmpl.Execute(w, []ResponseData{res})
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
