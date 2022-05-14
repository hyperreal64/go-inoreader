package userinfo

import (
	"context"
	"testing"

	"github.com/hyperreal64/go-inoreader/config"
)

func TestUserInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	userinfo, err := GetUserInfo(rc)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", userinfo)
}
