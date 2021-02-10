package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

const timeFormLong = "Mon 2 Jan 2006 3:04 PM"

func printSubList(subList *SubscriptionList) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subscription", "URL"})

	for _, v := range subList.Subscriptions {
		table.Append([]string{v.Title, v.URL})
	}
	table.Render()
}

func printUnreadCounts(subList *SubscriptionList, unreadCounts *UnreadCounters) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subscription", "# Unread"})

	titlesIDs := make(map[string]string)
	for _, v := range subList.Subscriptions {
		titlesIDs[v.ID] = v.Title
	}

	for _, v := range unreadCounts.Unreadcounts {
		count, err := v.Count.Int64()
		if err != nil {
			log.Fatalln(err)
		}

		idPrefix := "user/-/state/com.google/"
		if count > 0 && v.ID != idPrefix+"reading-list" && v.ID != idPrefix+"starred" {
			table.Append([]string{titlesIDs[v.ID], strconv.FormatInt(count, 10)})
		}
	}
	table.Render()
}

func printTagFolderList(tagList *TagFolderList) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Label", "Type", "# Unread"})

	for _, v := range tagList.Tags {
		label := strings.Split(v.ID, "/")
		labelSuffix := label[len(label)-1]
		table.Append([]string{labelSuffix, v.Type, strconv.FormatInt(v.UnreadCount, 10)})
	}
	table.Render()
}

func printStreamContents(streamContents *StreamContents) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Feed", "Title", "Date"})
	table.SetColWidth(55)

	for _, v := range streamContents.Items {
		date := time.Unix(v.Published, 0).Format(timeFormLong)
		table.Append([]string{v.Origin.Title, v.Title, date})
	}
	table.Render()
}
