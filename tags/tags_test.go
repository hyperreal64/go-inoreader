package tags

import (
	"context"
	"testing"

	"github.com/hyperreal64/go-inoreader/config"
)

func TestGetTagList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	taglist, err := GetTagList(rc)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", taglist)
}

func TestEditTag(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	params := map[string]string{
		"a": "user/-/state/com.google/starred",
		"i": "33050093431",
	}

	if err := EditTag(rc, params); err != nil {
		t.Fatal(err)
	}
}

func TestRenameTag(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()

	src := "linux"
	dest := "foss"
	params := map[string]string{
		"s": src,
		"d": dest,
	}

	if err := RenameTag(rc, params); err != nil {
		t.Fatal(err)
	}

	t.Logf("Renamed tag %s to %s", src, dest)
}
