package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
    "time"
    "github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

var newsapi *news.Client

// Routes
func indexHandler(w http.ResponseWriter, r *http.Request) {
    tpl.Execute(w, nil)
}

func search(newsapi *news.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        u, err := url.Parse(r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        params := u.Query()
        searchQuey := params.Get("q")
        page := params.Get("page")
        if page == "" {
            page = "1"
        }

        fmt.Println("Search Query is: ", searchQuey)
        fmt.Println("Page is: ", page)
    }
}

func main() {
    // Loads env variables
    err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    apiKey := os.Getenv("NEWS_API_KEY")
    if apiKey == "" {
        log.Fatal("Env: API Key must be set")
    }

    myClient := &http.Client{Timeout: 10 * time.Second}
    newsapi := news.NewsClient(myClient, apiKey, 10)

    fs := http.FileServer(http.Dir("assets"))
    
    // Router
    mux := http.NewServeMux()
    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
    
    // Routes
    mux.HandleFunc("/", indexHandler)
    mux.HandleFunc("/search", search(newsapi))

    http.ListenAndServe(":"+port, mux)
}
