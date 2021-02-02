package main

import (
	"context"
	"log"
	"net/http"

	"github.com/hyperreal64/go-inoreader/config"
)

// Init --- Initiate Oauth flow
func Init() {
	ctx, cancel := context.WithCancel(context.Background())
	http.HandleFunc("/", config.HandleInoreaderLogin)
	http.HandleFunc("/oauth2/redirect", config.HandleInoreaderCallback)
	http.HandleFunc("/go-inoreader", func(w http.ResponseWriter, r *http.Request) {
		config.ServeTemplate(w, r)
		cancel()
	})

	srv := &http.Server{Addr: ":53682"}
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
	Init()
}
