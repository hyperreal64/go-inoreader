package main

import (
	"context"
	"log"

	"github.com/tkanos/gonfig"
)

// Test ---
func Test() {
	cf := &Configuration{}
	if err := gonfig.GetConf(GetCfgFilePath(), cf); err != nil {
		log.Fatalln("Error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.GetOAuthResponse(ctx)

	Test()
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
