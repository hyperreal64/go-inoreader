package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/tkanos/gonfig"
)

func config2Client() (*resty.Client, error) {
	cf := &config{}
	if err := gonfig.GetConf(getCfgFilePath(), cf); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.getOAuthResponse(ctx)

	return resty.NewWithClient(client), nil
}

func printSubList(subList *SubscriptionList, onlyUnread bool, unreadCounts *UnreadCounters) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColMinWidth(0, 40)

	if onlyUnread == true {
		titlesIDs := make(map[string]string)
		for _, v := range subList.Subscriptions {
			titlesIDs[v.ID] = v.Title
		}

		for _, v := range unreadCounts.Unreadcounts {

			count, err := v.Count.Int64()
			if err != nil {
				log.Fatalln(err)
			}

			var (
				titleString string
				idPrefix    string = "state/com.google/"
				labelPrefix string = "label/"
			)

			if count > 0 {
				if strings.Contains(v.ID, idPrefix) || strings.Contains(v.ID, labelPrefix) {
					label := strings.Split(v.ID, "/")
					titleString = label[len(label)-1]
				} else {
					titleString = titlesIDs[v.ID]
				}
				table.SetHeader([]string{"Subscription", "# Unread"})
				table.Append([]string{titleString, strconv.FormatInt(count, 10)})
			}
		}
	} else {
		for _, v := range subList.Subscriptions {
			table.SetHeader([]string{"Subscription", "URL"})
			table.Append([]string{v.Title, v.URL})
		}
	}
	table.Render()
}

func execAddSub(streamID string) (string, error) {

	params := map[string]string{
		"quickadd": streamID,
	}

	rClient, err := config2Client()
	if err != nil {
		return "", errors.Wrap(err, "Could not get resty client")
	}

	quickAdd, err := quickAddSubscription(rClient, params)
	if err != nil {
		return "", err
	}

	if quickAdd.NumResults < 1 {
		return "", errors.Wrapf(err, "Could not add subscription %s\n", quickAdd.StreamName)
	}

	return fmt.Sprintf("Successfully added subscription: %s\n", quickAdd.StreamName), nil
}

func execEditSub(action string, streamID string, title string, folderAdd string, folderRem string) (string, error) {

	var params = make(map[string]string)
	params["ac"] = action
	params["s"] = streamID

	if title != "" {
		params["t"] = title
	}

	if folderAdd != "" {
		params["a"] = folderAdd
	}

	if folderRem != "" {
		params["r"] = folderRem
	}

	rClient, err := config2Client()
	if err != nil {
		return "", errors.Wrap(err, "Could not get resty client")
	}

	if err := editSubscription(rClient, params); err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully edited subscription: %s\n", title), nil
}

func printTagsFolders(tagList *TagFolderList, onlyUnread bool, option string) {

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
}

// TODO: getItemID()
func execEditTag(addTag string, remTag string, itemID string) (string, error) {

	params := map[string]string{
		"a": addTag,
		"r": remTag,
		"i": itemID,
	}

	rClient, err := config2Client()
	if err != nil {
		return "", errors.Wrap(err, "Could not get resty client")
	}

	if err := editTag(rClient, params); err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully edited tags for item: %s\n", itemID), nil
}

func execRenameTag(src string, dest string) (string, error) {

	params := map[string]string{
		"s":    src,
		"dest": dest,
	}

	rClient, err := config2Client()
	if err != nil {
		return "", errors.Wrap(err, "Could not get resty client")
	}

	if err := renameTag(rClient, params); err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully renamed tag from %s to %s\n", src, dest), nil
}

func execDelTag(tagName string) (string, error) {

	params := map[string]string{"s": tagName}

	rClient, err := config2Client()
	if err != nil {
		return "", errors.Wrap(err, "Could not get resty client")
	}

	if err := renameTag(rClient, params); err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully deleted tag: %s\n", tagName), nil
}

const timeFormLong = "Mon 2 Jan 2006 3:04 PM"

// TODO: printStreamContents with parameters:
//   Command: inoreader stream 	(--help)
//   Number of items 			(-n <int>)
//   Order 						(-r <n|o>)
//   Exclude target 			(-xt <streamID>)
//   Include target 			(-it <streamID>)
func printStreamContents(streamContents *StreamContents, withURL bool) {

	table := tablewriter.NewWriter(os.Stdout)

	if withURL == true {
		table.SetHeader([]string{"Feed", "Title", "URL"})

		var url string
		for _, v := range streamContents.Items {
			for _, w := range v.Canonical {
				url = w.Href
			}
			table.Append([]string{v.Origin.Title, v.Title, url})
		}
	} else {
		table.SetColMinWidth(1, 50)
		table.SetHeader([]string{"Feed", "Title", "Date"})

		for _, v := range streamContents.Items {
			table.Append([]string{v.Origin.Title, v.Title, time.Unix(v.Published, 0).Format(timeFormLong)})
		}
	}
	table.Render()
}
