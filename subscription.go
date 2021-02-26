package main

import (
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

func execAddSub(url string) {

	streamID := "feed/" + url
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
}

func execUnsubscribe(url string) {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "unsubscribe",
		"s":  streamID,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := editSubscription(rClient, params); err != nil {
		log.Fatalln(err)
	}
}

func execSetSubTitle(title string, url string) {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"t":  title,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := editSubscription(rClient, params); err != nil {
		log.Fatalln(err)
	}
}

func execAddSubToFolder(folder string, url string) {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"a":  folder,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := editSubscription(rClient, params); err != nil {
		log.Fatalln(err)
	}
}

func execRemSubFromFolder(folder string, url string) {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"r":  folder,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := editSubscription(rClient, params); err != nil {
		log.Fatalln(err)
	}
}
