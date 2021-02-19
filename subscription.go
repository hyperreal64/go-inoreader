package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func printSubList(onlyUnread bool) {

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	subList, err := getSubscriptionList(rClient)
	if err != nil {
		log.Fatalln(err)
	}

	unreadCounts, err := getUnreadCounters(rClient)
	if err != nil {
		log.Fatalln(err)
	}

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

func execAddSub(streamID string) string {

	params := map[string]string{
		"quickadd": streamID,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	quickAdd, err := quickAddSubscription(rClient, params)
	if err != nil {
		log.Fatalln(err)
	}

	if quickAdd.NumResults < 1 {
		log.Fatalln(err)
	}

	return fmt.Sprintf("Successfully added subscription: %s\n", quickAdd.StreamName)
}

func execEditSub(action string, streamID string, title string, folderAdd string, folderRem string) string {

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
		log.Fatalln(err)
	}

	if err := editSubscription(rClient, params); err != nil {
		log.Fatalln(err)
	}

	return fmt.Sprintf("Successfully edited subscription: %s\n", title)
}
