package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func DownloadHTML(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func appendPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	url = strings.TrimRight(url, "/")
	return url
}

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2f µs", float64(d.Nanoseconds())/1_000)
	} else if d < time.Second {
		return fmt.Sprintf("%.2f ms", float64(d.Milliseconds()))
	} else {
		return fmt.Sprintf("%.2f s", d.Seconds())
	}
}

func PrintHistory(history []ResponseData) {
	fmt.Println("\n--- Proxy Request History ---")

	for i, res := range history {
		status := "MISS"
		if res.Cached {
			status = "HIT"
		}
		fmt.Printf("[%d] Status: %-4s | Latency: %-10s | URL: %s\n",
			i+1, status, res.Duration, res.URL)
	}
	fmt.Printf("Total Requests: %d\n", len(history))
}
