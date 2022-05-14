package subscription

import (
	"context"
	"testing"

	"github.com/hyperreal64/go-inoreader/config"
)

func TestQuickAddSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	params := map[string]string{
		"quickadd": "feed/https://fedoramagazine.org/feed/",
	}

	quickadd, err := QuickAddSubscription(rc, params)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", quickadd)
}

func TestEditSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	feedURL := "https://fedoramagazine.org/feed/"
	params := map[string]string{
		"ac": "unsubscribe",
		"s":  "feed/" + feedURL,
	}

	if err := EditSubscription(rc, params); err != nil {
		t.Fatal(err)
	}

	t.Logf("Successful unsubscribe from %s", feedURL)
}

func TestGetSubscriptionList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	sublist, err := GetSubscriptionList(rc)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", sublist)
}

func TestGetUnreadCounters(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	unreadcounters, err := GetUnreadCounters(rc)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", unreadcounters)
}
