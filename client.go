package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

var (
	cf          = &Configuration{}
	ctx, cancel = context.WithCancel(context.Background())
	client      = cf.GetOAuthResponse(ctx)
)

func getConfig() error {

	if err := cf.GetConfigContent(); err != nil {
		return errors.Wrap(err, "Could not get configuration content")
	}

	return nil
}

// ListUnreadCounters ---
func ListUnreadCounters(client *http.Client) {

	unreadCounters := &UnreadCounters{}
	if err := GetUnreadCounters(client, unreadCounters); err != nil {
		log.Fatalln(err)
	}

	for _, v := range unreadCounters.Unreadcounts {
		fmt.Println(v)
	}
}
