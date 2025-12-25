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
	history    = make([]ResponseData, 0)
	mu         sync.Mutex
)

func homehandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	clear(cacheStore)

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func URLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var res ResponseData
	url := r.FormValue("url")
	url = appendPrefix(url)

	start := time.Now()
	_, ok := cacheStore[url]
	res.URL = url
	if ok {
		res.Cached = true
	} else {
		res.HTML, err = DownloadHTML(url)
		if err != nil {
			log.Printf("failed to download page to db: %v", err)
		} else {
			mu.Lock()
			cacheStore[url] = res.HTML
			mu.Unlock()
		}
	}

	res.Duration = prettyDuration(time.Since(start))
	history = append(history, res)
	tmpl.Execute(w, history)
}

func main() {
	http.HandleFunc("/", homehandler)
	http.HandleFunc("/fetch", URLHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
