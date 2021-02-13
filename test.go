package main

import (
	"context"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/tkanos/gonfig"
)

// Test ---
func Test() {
	cf := &config{}
	if err := gonfig.GetConf(getCfgFilePath(), cf); err != nil {
		log.Fatalln("Error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.getOAuthResponse(ctx)
	rClient := resty.NewWithClient(client)

	// subList, _ := getSubscriptionList(rClient)
	// printSubList(subList)
	// unreadCounters, _ := getUnreadCounters(client)
	// printUnreadCounts(subList, unreadCounters)
	// tagList, _ := getTagList(client)
	// printTagFolderList(tagList)

	params := map[string]string{
		"n": "5",
		"s": "feed/http://www.osnews.com/files/recent.xml",
	}

	streamContents, _ := getStreamContents(rClient, params)
	printStreamContents(streamContents)
}
