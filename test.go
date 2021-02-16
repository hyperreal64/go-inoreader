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
		log.Fatalln(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.getOAuthResponse(ctx)
	rClient := resty.NewWithClient(client)

	// TODO: Testing (properly)
}
