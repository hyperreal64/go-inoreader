package main

import (
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

const timeFormLong = "Mon 2 Jan 2006 3:04 PM"

func printStreamContentsWithDate(n string, r string, xt string, it string, s string) error {

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
		"s":  s,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		return errors.Wrapf(err, "Could not get stream contents with parameters: %v", params)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Feed", "Title", "Date"})

	for _, v := range streamContents.Items {
		table.Append([]string{v.Origin.Title, v.Title, time.Unix(v.Published, 0).Format(timeFormLong)})
	}
	table.Render()

	return nil
}

func printStreamContentsWithURL(n string, r string, xt string, it string, s string) error {

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
		"s":  s,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		return errors.Wrapf(err, "Could not get stream contents with parameters %v", params)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Feed", "Title", "URL"})

	var url string
	for _, v := range streamContents.Items {
		for _, w := range v.Canonical {
			url = w.Href
		}
		table.Append([]string{v.Origin.Title, v.Title, url})
	}
	table.Render()

	return nil
}

func printStreamContentsWithIDs(n string, r string, xt string, it string, s string) error {

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
		"s":  s,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		return errors.Wrapf(err, "Could not get stream contents with parameters %v", params)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Feed", "Title", "Item ID"})
	table.SetColMinWidth(1, 50)

	for _, v := range streamContents.Items {
		table.Append([]string{v.Origin.Title, v.Title, v.ID})
	}
	table.Render()

	return nil
}

func execMarkStreamAsRead(streamID string) error {

	if strings.HasPrefix(streamID, "http") {
		streamID = "feed/" + streamID
	}

	rClient, err := config2Client()
	if err != nil {
		return errors.Wrap(err, getRestyErr)
	}

	params := map[string]string{
		"s": streamID,
	}

	if err := markAllAsRead(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark stream as read: %s", streamID)
	}

	return nil
}
