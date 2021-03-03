package main

import (
	"context"
	"encoding/json"
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

	getRestyErr string = "Could not get resty client"

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
