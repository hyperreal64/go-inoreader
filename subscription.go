package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func printSubList(onlyUnread bool) error {

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	subList, err := getSubscriptionList(rClient)
	if err != nil {
		return errors.Wrap(err, "Could not get subscription list")
	}

	unreadCounts, err := getUnreadCounters(rClient)
	if err != nil {
		return errors.Wrap(err, "Could not get unread counters")
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
				return errors.Wrapf(err, "Could not convert %T to Int64", v.Count)
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

	return nil
}

func execAddSub(url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"quickadd": streamID,
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	quickAdd, err := quickAddSubscription(rClient, params)
	if err != nil {
		return err
	}

	if quickAdd.NumResults == 0 {
		return errors.New("Please check if the URL is correct")
	}

	return nil
}

func execUnsubscribe(url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "unsubscribe",
		"s":  streamID,
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not unsubscribe from subscription %s", url)
	}

	return nil
}

func execSetSubTitle(title string, url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"t":  title,
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not set title %s on subscription %s", title, url)
	}

	return nil
}

func execAddSubToFolder(folder string, url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"a":  folder,
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not add subscription %s to folder %s", url, folder)
	}

	return nil
}

func execRemSubFromFolder(folder string, url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"r":  folder,
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not remove subscription %s from folder %s", url, folder)
	}

	return nil
}
