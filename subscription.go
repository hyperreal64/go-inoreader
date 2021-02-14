package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// inoreader sub list
// print subscriptions with url
func printSubList(subList *SubscriptionList) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subscription", "URL"})
	table.SetColMinWidth(0, 40)

	for _, v := range subList.Subscriptions {
		table.Append([]string{v.Title, v.URL})
	}
	table.Render()
}

// inoreader sub unread
// print only subscriptions with unread items,
// display unread count
func printUnreadSubCounts(subList *SubscriptionList, unreadCounts *UnreadCounters) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subscription", "# Unread"})
	table.SetColMinWidth(0, 40)

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
			table.Append([]string{titleString, strconv.FormatInt(count, 10)})
		}
	}
	table.Render()
}
