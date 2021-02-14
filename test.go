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
	// unreadCounters, _ := getUnreadCounters(rClient)
	// printUnreadSubCounts(subList, unreadCounters)
	tagList, _ := getTagList(rClient)
	// printTagFolderList(tagList)
	printTagsXorFolders(tagList, "folders")

	// params := map[string]string{
	// 	"n": "5",
	// 	"s": "feed/http://www.osnews.com/files/recent.xml",
	// }

	// streamContents, _ := getStreamContents(rClient, params)
	// printStreamContents(streamContents)
}
