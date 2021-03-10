package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

const (
	baseURL           = "https://www.inoreader.com/reader/api/0"
	userInfoURL       = baseURL + "/user-info"
	addSubURL         = baseURL + "/subscription/quickadd"
	editSubURL        = baseURL + "/subscription/edit"
	unreadCountersURL = baseURL + "/unread-count"
	subListURL        = baseURL + "/subscription/list"
	tagListURL        = baseURL + "/tag/list?types=1&counts=1"
	streamContentsURL = baseURL + "/stream/contents"
	itemIDsURL        = baseURL + "/stream/items/ids"
	streamPrefsURL    = baseURL + "/preference/stream/list"
	streamPrefsSetURL = baseURL + "/preference/stream/set"

	timeFormLong string = "Mon 2 Jan 2006 3:04 PM"
)

func printUserInfo() error {

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	resp, err := rClient.R().Get(userInfoURL)
	if err != nil {
		return errors.Wrap(err, "Could not get user info")
	}

	userInfo := &UserInfo{}
	if err := json.Unmarshal(resp.Body(), userInfo); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", userInfo)
	}

	data := [][]string{
		{"User ID", userInfo.UserID},
		{"Username", userInfo.UserName},
		{"Profile ID", userInfo.UserProfileID},
		{"Email", userInfo.UserEmail},
		{"Blogger User", strconv.FormatBool(userInfo.IsBloggerUser)},
		{"Sign-up Date", time.Unix(userInfo.SignupTimeSec, 0).Format(timeFormLong)},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	return nil
}

func printCmdExamples() {

	data := [][]string{
		{"List all subscriptions", "go-inoreader list subscriptions -a"},
		{"List only unread subscriptions", "go-inoreader list subscriptions -u"},
		{"List all tags and folders", "go-inoreader list tags -a"},
		{"List only tags and folders with unread items", "go-inoreader list tags -u"},
		{"List only folders", "go-inoreader list tags --type=folders -a"},
		{"List only folders with unread items", "go-inoreader list tags --type=folders -u"},
		{"List the 5 newest items in stream with their URLs", "go-inoreader list stream <feed URL> -n 5 -r n -u"},
		{"List the 10 oldest items in stream with their Inoreader item IDs", "go-inoreader list stream <feed URL> -n 10 -r o -i"},
		{"List the 3 most recently starred items with their timestamps", "go-inoreader list starred -n 3 -r n -d"},
		{"List all web pages saved to Inoreader with their URLs", "go-inoreader list web-pages -u"},
		{"Mark an item as read", "go-inoreader mark-item --read --item-id=\"000000067ca9449f\""},
		{"Star an item", "go-inoreader mark-item --star --item-id=\"000000067ca9449f\""},
		{"Mark all items in feed as read", "go-inoreader mark-stream-read <feed URL>"},
		{"Add a subscription to Inoreader", "go-inoreader subscription add <feed URL>"},
		{"Unsubscribe from subscription", "go-inoreader subscription unsubscribe <feed URL>"},
		{"Change title of subscription to 'Linux News'", "go-inoreader subscription set-title --url=<feed URL> --title=\"Linux News\""},
		{"Add 'Fedora Magazine' subscription to the 'FOSS' folder", "go-inoreader subscription add-to-folder --folder=\"FOSS\" --url=\"https://fedoramagazine.org/feed/\""},
		{"Rename a user-defined tag from 'linux' to 'foss'", "go-inoreader tags --src=\"linux\" --dest=\"foss\""},
		{"Print user information", "go-inoreader user-info"},
		{"Login and authorize Inoreader account access", "go-inoreader login"},
	}

	for _, v := range data {
		fmt.Printf("%s:\n", v[0])
		fmt.Printf("%s\n\n", v[1])
	}
}
