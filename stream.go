package main

import (
	"context"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

const (
	starredTag = "user/-/state/com.google/starred"
	savedTag   = "user/-/state/com.google/saved-web-pages"
)

func getStreamContents(rc *resty.Client, params map[string]string) (*StreamContents, error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Get(streamContentsURL)

	if err != nil {
		return nil, err
	}

	streamContents := &StreamContents{}
	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), streamContents); err != nil {
		return nil, err
	}

	return streamContents, nil
}

func markAllAsRead(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/mark-all-as-read")
	if err != nil {
		return err
	}

	return nil
}

func printStreamContentsWithDate(n string, r string, xt string, it string, s string) error {

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	var (
		streamID  string
		tableData [][]string
	)

	if s == savedTag || s == starredTag {
		streamID = s
	} else {
		streamID = "feed/" + s
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
		"s":  streamID,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		return errors.Wrapf(err, "Could not get stream contents with parameters: %v", params)
	}

	for _, v := range streamContents.Items {
		tableData = append(tableData, []string{v.Origin.Title, v.Title, time.Unix(v.Published, 0).Format(timeFormLong)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Feed", "Title", "Date"})
	table.AppendBulk(tableData)
	table.Render()

	return nil
}

func printStreamContentsWithURL(n string, r string, xt string, it string, s string) error {

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	var (
		streamID    string
		tableHeader []string
		tableData   [][]string
		url         string
	)

	if s == savedTag || s == starredTag {
		streamID = s
	} else {
		streamID = "feed/" + s
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
		"s":  streamID,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		return errors.Wrapf(err, "Could not get stream contents with parameters %v", params)
	}

	if s == savedTag {
		tableHeader = []string{"Title", "URL"}

		for _, v := range streamContents.Items {
			for _, w := range v.Canonical {
				url = w.Href
			}
			tableData = append(tableData, []string{v.Title, url})
		}
	} else {
		tableHeader = []string{"Feed", "Title", "URL"}

		for _, v := range streamContents.Items {
			for _, w := range v.Canonical {
				url = w.Href
			}
			tableData = append(tableData, []string{v.Origin.Title, v.Title, url})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableHeader)
	table.AppendBulk(tableData)
	table.Render()

	return nil
}

func printStreamContentsWithIDs(n string, r string, xt string, it string, s string) error {

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	var (
		streamID    string
		tableHeader []string
		tableData   [][]string
	)

	if s == savedTag || s == starredTag {
		streamID = s
	} else {
		streamID = "feed/" + s
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
		"s":  streamID,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		return errors.Wrapf(err, "Could not get stream contents with parameters %v", params)
	}

	if s == savedTag {
		tableHeader = []string{"Title", "Item ID"}

		for _, v := range streamContents.Items {
			tableData = append(tableData, []string{v.Title, v.ID})
		}

	} else {
		tableHeader = []string{"Feed", "Title", "Item ID"}

		for _, v := range streamContents.Items {
			tableData = append(tableData, []string{v.Origin.Title, v.Title, v.ID})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableHeader)
	table.AppendBulk(tableData)
	table.Render()

	return nil
}

func execMarkStreamAsRead(streamURL string) error {

	streamID := "feed/" + streamURL

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	params := map[string]string{
		"s": streamID,
	}

	if err := markAllAsRead(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark stream as read: %s", streamID)
	}

	return nil
}
