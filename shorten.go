package main

import (
    "fmt"
    "log"
    "net/http"
    "math/rand"
    "time"
)

var (
    port = "8080"
    urlMap = make(map[string]string)
)

func main() {
    http.HandleFunc("/shorten", shortenURL)
    http.HandleFunc("/r/", redirectURL)

    fmt.Printf("Server starting on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    originalURL := r.URL.Query().Get("url")
    if originalURL == "" {
        http.Error(w, "URL is required", http.StatusBadRequest)
        return
    }

    shortened := generateShortLink()
    urlMap[shortened] = originalURL

    fmt.Fprintf(w, "Shortened URL: http://localhost:%s/r/%s\n", port, shortened)
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
    shortened := r.URL.Path[len("/r/"):]
    originalURL, ok := urlMap[shortened]
    if !ok {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    http.Redirect(w, r, originalURL, http.StatusFound)
}

func generateShortLink() string {
    rand.Seed(time.Now().UnixNano())
    letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, 6)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
