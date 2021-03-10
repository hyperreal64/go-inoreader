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

func quickAddSubscription(rc *resty.Client, params map[string]string) (*QuickAdd, error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Post(addSubURL)
	if err != nil {
		return nil, err
	}

	quickAdd := &QuickAdd{}
	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &quickAdd); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON object: %v", quickAdd)
	}

	return quickAdd, nil
}

func editSubscription(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(editSubURL)
	if err != nil {
		return err
	}

	return nil
}

func getSubscriptionList(rc *resty.Client) (*SubscriptionList, error) {

	resp, err := rc.R().Get(subListURL)
	if err != nil {
		return nil, err
	}

	subList := &SubscriptionList{}
	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &subList); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON object: %v", subList)
	}

	return subList, nil
}

func getUnreadCounters(rc *resty.Client) (*UnreadCounters, error) {

	resp, err := rc.R().Get(unreadCountersURL)
	if err != nil {
		return nil, err
	}

	unreadCounters := &UnreadCounters{}
	if err = resty.Unmarshalc(rc, "application/json", resp.Body(), &unreadCounters); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON object %v", unreadCounters)
	}

	return unreadCounters, nil
}

func printSubList(onlyUnread bool) error {

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	subList, err := getSubscriptionList(rClient)
	if err != nil {
		return errors.Wrap(err, "Unable to get subscription list")
	}

	unreadCounts, err := getUnreadCounters(rClient)
	if err != nil {
		return errors.Wrap(err, "Unable to get unread counters")
	}

	var (
		tableData   [][]string
		tableHeader []string
		tableFooter []string
		hasFooter   bool
	)

	if onlyUnread {
		idTitleMap := make(map[string]string)
		idURLMap := make(map[string]string)
		for _, v := range subList.Subscriptions {
			idTitleMap[v.ID] = v.Title
			idURLMap[v.ID] = v.URL
		}

		var totalUnread int64
		for _, v := range unreadCounts.Unreadcounts {

			count, err := v.Count.Int64()
			if err != nil {
				return errors.Wrapf(err, "Unable to convert %T to Int64", v.Count)
			}

			if strings.Contains(v.ID, "state/com.google/reading-list") {
				totalUnread = count
			}

			var (
				idPrefix    string = "state/com.google/"
				labelPrefix string = "label/"
			)

			if count > 0 {
				if !strings.Contains(v.ID, idPrefix) && !strings.Contains(v.ID, labelPrefix) {
					tableData = append(tableData, []string{idTitleMap[v.ID], idURLMap[v.ID], strconv.FormatInt(count, 10)})
				}
			}
		}

		tableHeader = []string{"Subscription", "URL", "# Unread"}
		hasFooter = true
		tableFooter = []string{"", "Total unread:", strconv.FormatInt(totalUnread, 10)}

	} else {
		for _, v := range subList.Subscriptions {
			tableData = append(tableData, []string{v.Title, v.URL})
		}

		tableHeader = []string{"Subscription", "URL"}
		hasFooter = false
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableHeader)
	if hasFooter {
		table.SetFooter(tableFooter)
		table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	}
	table.AppendBulk(tableData)
	table.Render()

	return nil
}

func execAddSub(url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"quickadd": streamID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

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

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Unable to unsubscribe from subscription %s", url)
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

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Unable to set title %s on subscription %s", title, url)
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

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Unable to add subscription %s to folder %s", url, folder)
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

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := editSubscription(rClient, params); err != nil {
		return errors.Wrapf(err, "Unable to remove subscription %s from folder %s", url, folder)
	}

	return nil
}
