package main

import (
	"context"
	"log"
	"net/http"
)

// Init --- Initiate Oauth flow
func Init() {
	ctx, cancel := context.WithCancel(context.Background())
	http.HandleFunc("/", handleInoreaderLogin)
	http.HandleFunc("/oauth/redirect", handleInoreaderCallback)
	http.HandleFunc("/go-inoreader", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r)
		cancel()
	})

	srv := &http.Server{Addr: ":8081"}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server error: %s", err.Error())
		}
	}()
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil && err != context.Canceled {
		log.Println(err)
	}
	log.Println("Done")
}

func main() {
	// Init()
	Test()
}
