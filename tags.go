package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func printTagsFolders(onlyUnread bool, option string) error {

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	tagList, err := getTagList(rClient)
	if err != nil {
		return errors.Wrap(err, "Could not get tags list")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{option, "# unread"})

	switch option {
	case "tags":
		for _, v := range tagList.Tags {
			if onlyUnread == true {
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
			if onlyUnread == true {
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

	var params = make(map[string]string)
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

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := editTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark item %s as read", itemID)
	}

	return nil
}

func execEditTagStar(itemID string, star bool) error {
	var params = make(map[string]string)
	if star {
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

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := editTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark item %s as starred", itemID)
	}

	return nil
}

func execRenameTag(src string, dest string) error {

	params := map[string]string{
		"s":    src,
		"dest": dest,
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := renameTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not rename tag %s to %s", src, dest)
	}

	return nil
}

func execDelTag(tagName string) error {

	params := map[string]string{"s": tagName}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := renameTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not delete tag %s", tagName)
	}

	return nil
}
