package main

import (
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

const timeFormLong = "Mon 2 Jan 2006 3:04 PM"

// TODO: printStreamContents with parameters:
//   Command: inoreader stream 	(--help)
//   Number of items 			(-n <int>)
//   Order 						(-r <n|o>)
//   Exclude target 			(-xt <streamID>)
//   Include target 			(-it <streamID>)
func printStreamContents(n string, r string, xt string, it string, withURL bool) {

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	params := map[string]string{
		"n":  n,
		"r":  r,
		"xt": xt,
		"it": it,
	}

	streamContents, err := getStreamContents(rClient, params)
	if err != nil {
		log.Fatalln(err)
	}

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
