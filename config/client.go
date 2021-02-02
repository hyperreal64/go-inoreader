package config

import (
	"context"
	"fmt"
	"log"
	"net/http"

	api "github.com/hyperreal64/go-inoreader/api"

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

	unreadCounters := &api.UnreadCounters{}
	if err := api.GetUnreadCounters(client, unreadCounters); err != nil {
		log.Fatalln(err)
	}

	for _, v := range unreadCounters.Unreadcounts {
		fmt.Println(v)
	}
}
