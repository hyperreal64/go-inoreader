package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func getTagList(rc *resty.Client) (*TagFolderList, error) {

	resp, err := rc.R().Get(tagListURL)
	if err != nil {
		return nil, err
	}

	tagList := &TagFolderList{}
	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), tagList); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", tagList)
	}

	return tagList, nil
}

func renameTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/rename-tag")
	if err != nil {
		return err
	}

	return nil
}

// func deleteTag(rc *resty.Client, tagName string) error {

// 	_, err := rc.R().
// 		SetQueryParams(map[string]string{
// 			"s": tagName,
// 		}).
// 		Post(baseURL + "/disable-tag")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func editTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/edit-tag")
	if err != nil {
		return err
	}

	return nil
}

func printTagsFolders(onlyUnread bool, option string) error {

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	tagList, err := getTagList(rClient)
	if err != nil {
		return errors.Wrap(err, "Could not get tags list")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{option, "# unread"})

	switch option {
	case "tags":
		for _, v := range tagList.Tags {
			if onlyUnread {
				if v.Type == "tag" && v.UnreadCount > 0 {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			} else {
				if v.Type == "tag" {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			}
		}

	case "folders":
		for _, v := range tagList.Tags {
			if onlyUnread {
				if v.Type == "folder" && v.UnreadCount > 0 {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			} else {
				if v.Type == "folder" {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			}
		}
	}
	table.Render()

	return nil
}

// TODO: Notify user when item cannot be marked unread due to `timestampUsec` being older than `firstitemsec` of its feed
func execEditTagRead(itemID string, markRead bool) error {

	var params = map[string]string{}
	if markRead {
		params = map[string]string{
			"a": "user/-/state/com.google/read",
			"i": itemID,
		}
	} else {
		params = map[string]string{
			"r": "user/-/state/com.google/read",
			"i": itemID,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := editTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark item %s as read", itemID)
	}

	return nil
}

func execEditTagStar(itemID string, starred bool) error {
	var params = map[string]string{}
	if starred {
		params = map[string]string{
			"a": "user/-/state/com.google/starred",
			"i": itemID,
		}
	} else {
		params = map[string]string{
			"r": "user/-/state/com.google/starred",
			"i": itemID,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := editTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark item %s as %s", itemID, starred)
	}

	return nil
}

// func execEditTagSaved(url string, saved bool) error {

// 	var params = map[string]string{}
// 	if saved {
// 		params = map[string]string{
// 			"a": "user/-/state/com.google/saved-web-pages",
// 			"i": url,
// 		}
// 	} else {
// 		params = map[string]string{
// 			"r": "user/-/state/com.google/saved-web-pages",
// 			"i": url,
// 		}
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	rClient := oauth2RestyClient(ctx)
// 	defer cancel()

// 	if err := editTag(rClient, params); err != nil {
// 		return errors.Wrapf(err, "Could not mark item %s as %s", url, saved)
// 	}

// 	return nil
// }

func execRenameTag(src string, dest string) error {

	params := map[string]string{
		"s":    src,
		"dest": dest,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := renameTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not rename tag %s to %s", src, dest)
	}

	return nil
}

func execDelTag(tagName string) error {

	params := map[string]string{"s": tagName}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := renameTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not delete tag %s", tagName)
	}

	return nil
}
