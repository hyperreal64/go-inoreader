package client

import (
	"fmt"
	"log"

	"github.com/hyperreal64/go-inoreader/config"
	"github.com/hyperreal64/go-inoreader/pkg/api"
	"github.com/pkg/errors"
)

var (
	cf = &config.Configuration{}
	ctx, cancel = context.WithCancel(context.Background())
	client = cf.GetOAuthResponse(ctx)
)

func getConfig() error {

	if err := cf.GetConfigContent(); err != nil {
		return errors.Wrap(err, "Could not get configuration content")
	}

	return nil
}

// ListUnreadCounters ---
func ListUnreadCounters(client *http.Client) {

	unreadCounters := &api.UnreadCounters{}
	if err := api.GetUnreadCounters(client, unreadCounters)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range unreadCounters.Unreadcounts {
		fmt.Println(v)
	}
}
