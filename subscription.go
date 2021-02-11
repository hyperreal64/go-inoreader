package main

import (
	"log"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

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
