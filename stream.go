package main

import (
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

// TODO: index stream content output per item
//   Command: inoreader stream summary <streamID> <itemID>
//   Use item index to print the following:
//   Feed name
//   Item title
//   Date Published
//   Canonical href
//   Author
//   Summary content
