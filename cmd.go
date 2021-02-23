package main

import (
	"github.com/alecthomas/kong"
)

// Context ---
type Context struct {
	Debug bool
}

// LoginCmd ---
type LoginCmd struct{}

// Run (*LoginCmd) ---
func (l *LoginCmd) Run(ctx *Context) error {
	Init()
	return nil
}

// AddCmd ---
type AddCmd struct {
	Add string `arg:"" required,name:"stream id" help:"Stream ID of subscription, i.e. user/-/label/Tech"`
}

// Run (*AddCmd) ---
func (a *AddCmd) Run(ctx *Context) error {
	execAddSub(a.Add)
	return nil
}

// SubscriptionCmd ---
type SubscriptionCmd struct {
	List struct {
		All    bool `arg,required,name:"all" short:"a" xor:"list" help:"List all subscriptions"`
		Unread bool `arg,required,name:"unread" short:"u" xor:"list" help:"List only unread subscriptions"`
	} `cmd:"" aliases:"ls" help:"List subscriptions"`
	Edit struct {
		A string `arg,required,name:"action" short:"a" help:"Edit, subscribe, or unsubscribe"`
		S string `arg,required,name:"stream ID" short:"s" help:"Stream ID of subscription feed"`
		T string `arg,optional,name:"title" short:"t" help:"Change the subscription title. Omit to keep title."`
		F string `arg,optional,name:"folder" short:"f" help:"Add subscription to folder. Use folder name like user/-/label/Tech"`
		R string `arg,optional,name:"folder" short:"r" help:"Remove subscription from folder. Use folder name like user/-/label/Tech"`
	} `cmd:"" aliases:"ed" help:"Edit a subscription"`
}

// Run (*SubscriptionCmd) ---
func (s *SubscriptionCmd) Run(ctx *Context) error {
	if s.List.All {
		printSubList(false)
	}

	if s.List.Unread {
		printSubList(true)
	}

	execEditSub(s.Edit.A,
		s.Edit.S,
		s.Edit.T,
		s.Edit.F,
		s.Edit.R)

	return nil
}

// TagsCmd ---
type TagsCmd struct {
	List struct {
		All    bool   `arg,required,name:"all" short:"a" group:"list" xor:"list" help:"List all tags and/or folders"`
		Unread bool   `arg,required,name:"unread" short:"u" group:"list" xor:"list" help:"List only unread tags and/or folders"`
		Type   string `arg,optional,name:"type" short:"t" placeholder:"tags|folders" help:"List either only tags or only folders"`
	} `cmd:"" aliases:"ls" help:"List tags and/or folders"`

	Rename struct {
		Src  string `arg,required,placeholder:"src tag" short:"s" help:"Source tag name"`
		Dest string `arg,required,placeholder:"dest tag" short:"d" help:"Destination tag name"`
	} `cmd:"" aliases:"mv" help:"Rename a tag"`

	Delete struct {
		Delete string `arg,required,placeholder:"tag name" help:"Tag name to delete"`
	} `cmd:"" aliases:"rm" help:"Delete"`
}

// Run (*TagsCmd) ---
func (t *TagsCmd) Run(ctx *Context) error {
	if t.List.All {
		printTagsFolders(false, t.List.Type)
	}

	if t.List.Unread {
		printTagsFolders(true, t.List.Type)
	}

	execRenameTag(t.Rename.Src, t.Rename.Dest)

	execDelTag(t.Delete.Delete)

	return nil
}

// MarkItemCmd ---
type MarkItemCmd struct {
	Read   bool   `arg,name:"Read" short:"r" xor:"read" help:"Mark <itemID> as read"`
	Unread bool   `arg,name:"Unread" short:"ur" xor:"read" help:"Mark <itemID> as unread"`
	Star   bool   `arg,name:"Star" short:"s" xor:"star" help:"Star or unstar <itemID>"`
	Unstar bool   `arg,name:"Unstar" short:"us" xor:"star" help:"Unstar <itemID>"`
	ItemID string `arg,required,name:"Item" help:"Item ID"`
}

// Run (*MarkItemCmd) ---
func (m *MarkItemCmd) Run(ctx *Context) error {

	if m.Read {
		execEditTagRead(m.ItemID, true)
	}

	if m.Unread {
		execEditTagRead(m.ItemID, false)
	}

	if m.Star {
		execEditTagStar(m.ItemID, true)
	}

	if m.Unstar {
		execEditTagStar(m.ItemID, false)
	}

	return nil
}

// StreamCmd ---
type StreamCmd struct {
	N  string `arg,optional,placeholder:"integer" default:"5" help:"Specify number of items returned"`
	R  string `arg,optional,placeholder:"n|o" help:"Specify order of items returned, i.e. newest or oldest"`
	Xt string `arg,optional,placeholder:"stream ID" help:"Specify stream ID of target to exclude"`
	It string `arg,optional,placeholder:"stream ID" help:"Specify stream ID of target to include"`
	S  string `arg,required,placeholder:"stream ID" help:"Specify stream ID of feed"`
	U  bool   `arg,optional,default:"false" xor:"stream" help:"Display items with their URLs"`
	ID bool   `arg,optional,default:"false" xor:"stream" help:"Display items with their IDs"`
}

// Run (*StreamCmd) ---
func (sc *StreamCmd) Run(ctx *Context) error {
	if sc.U {
		printStreamContentsWithURL(sc.N, sc.R, sc.Xt, sc.It, sc.S)
	} else if sc.ID {
		printStreamContentsWithIDs(sc.N, sc.R, sc.Xt, sc.It, sc.S)
	} else {
		printStreamContentsWithDate(sc.N, sc.R, sc.Xt, sc.It, sc.S)
	}

	return nil
}

// ReadCmd ---
type ReadCmd struct {
	StreamID string `arg,required,placeholder:"stream ID" help:"Enter stream ID or feed URL to mark as read"`
}

// Run (*ReadCmd) ---
func (r *ReadCmd) Run(ctx *Context) error {
	execMarkStreamAsRead(r.StreamID)

	return nil
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Login LoginCmd `cmd:"" help:"Login and initiate Oauth flow"`

	Add AddCmd `cmd:"" help:"Add a subscription"`

	Subscription SubscriptionCmd `cmd:"" help:"Query subscriptions"`

	Tags TagsCmd `cmd:"" help:"Query tags and folders"`

	MarkItem MarkItemCmd `cmd:"" help:"Mark items as read/unread or starred/unstarred"`

	Stream StreamCmd `cmd:"" help:"Query stream contents"`

	Read ReadCmd `cmd:"" help:"Mark all items in stream as read"`
}

func main() {

	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
