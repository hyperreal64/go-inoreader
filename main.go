package main

import (
	"context"
	"log"
	"net/http"

	"github.com/tkanos/gonfig"
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

func main() {
	// Init()
	cf := &Configuration{}
	if err := gonfig.GetConf(GetCfgFilePath(), cf); err != nil {
		log.Fatalln("Error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.GetOAuthResponse(ctx)

	// subList, _ := getSubscriptionList(client)
	// unreadCounters, _ := getUnreadCounters(client)
	// printUnreadCounts(subList, unreadCounters)
	// tagList, _ := getTagList(client)
	// printTagFolderList(tagList)

	params := &ContentsParams{
		NumOfItems: "5",
		StreamID:   "feed/http://www.osnews.com/files/recent.xml",
	}

	streamContents, _ := getStreamContents(client, params)
	printStreamContents(streamContents)
}
