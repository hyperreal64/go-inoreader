package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/tkanos/gonfig"
)

// Test ---
func Test() {
	cf := &Configuration{}
	if err := gonfig.GetConf(getCfgFilePath(), cf); err != nil {
		log.Fatalln("Error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.GetOAuthResponse(ctx)
	rClient := resty.NewWithClient(client)

	subList, _ := getSubscriptionList(rClient)
	fmt.Println(subList)
	// printSubList(subList)
	// unreadCounters, _ := getUnreadCounters(client)
	// printUnreadCounts(subList, unreadCounters)
	// tagList, _ := getTagList(client)
	// printTagFolderList(tagList)

	// params := &ContentsParams{
	// 	NumOfItems: "5",
	// 	StreamID:   "feed/http://www.osnews.com/files/recent.xml",
	// }

	// streamContents, _ := getStreamContents(client, params)
	// printStreamContents(streamContents)
}
