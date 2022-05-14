package stream

import (
	"context"
	"testing"

	"github.com/hyperreal64/go-inoreader/config"
)

func TestStreamContents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	params := map[string]string{
		"n": "5",
		"r": "n",
		"s": "feed/https://fedoramagazine.org/feed/",
	}

	sc, err := GetStreamContents(rc, params)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", sc)
}

func TestMarkAllAsRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	feedURL := "https://fedoramagazine.org/feed/"
	params := map[string]string{
		"s": "feed/" + feedURL,
	}

	if err := MarkAllAsRead(rc, params); err != nil {
		t.Fatal(err)
	}

	t.Logf("%s marked as read", feedURL)
}
