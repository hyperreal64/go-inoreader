package config

import (
	"context"
	"log"
	"net/http"
)

// Init --- Initiate Oauth flow
func Init() {
	ctx, cancel := context.WithCancel(context.Background())
	http.HandleFunc("/", HandleInoreaderLogin)
	http.HandleFunc("/oauth2/redirect", HandleInoreaderCallback)
	http.HandleFunc("/go-inoreader", func(w http.ResponseWriter, r *http.Request) {
		ServeTemplate(w, r)
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
