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

	Unsubscribe struct {
		StreamID string `arg,required,name:"stream ID" help:"Unsubscribe from <stream ID>"`
	} `cmd:"" aliases:"un" help:"Unsubscribe from <stream ID>"`

	SetTitle struct {
		Title    string `arg,required,name:"title" short:"t" help:"New title of subscription"`
		StreamID string `arg,required,name:"stream ID" short:"s" help:"Stream ID of subscription to change"`
	} `cmd:"" help:"Change the title of a subscription"`

	AddToFolder struct {
		Folder   string `arg,required,name:"folder" short:"f" help:"Folder to add subscription to"`
		StreamID string `arg,required,name:"stream ID" short:"s" help:"Stream ID of subscription"`
	} `cmd:"" help:"Add subscription to folder"`

	RemFromFolder struct {
		Folder   string `arg,required,name:"folder" short:"f" help:"Folder to remove subscription from"`
		StreamID string `arg,required,name:"stream ID" short:"s" help:"Stream ID of subscription"`
	} `cmd:"" help:"Remove subscription from folder"`
}

// Run (*SubscriptionCmd) ---
func (s *SubscriptionCmd) Run(ctx *Context) error {
	// TODO: Refactor
	if s.List.All {
		printSubList(false)
	}

	if s.List.Unread {
		printSubList(true)
	}

	execUnsubscribe(s.Unsubscribe.StreamID)

	execSetSubTitle(s.SetTitle.Title, s.SetTitle.StreamID)

	execAddSubToFolder(s.AddToFolder.Folder, s.AddToFolder.StreamID)

	execRemSubFromFolder(s.RemFromFolder.Folder, s.RemFromFolder.StreamID)

	return nil
}

// TagsCmd ---
type TagsCmd struct {
	List struct {
		All    bool   `arg,required,name:"all" short:"a" xor:"list" help:"List all tags and/or folders"`
		Unread bool   `arg,required,name:"unread" short:"u" xor:"list" help:"List only unread tags and/or folders"`
		Type   string `arg,optional,name:"type" short:"t" placeholder:"tags|folders" help:"List either only tags or only folders"`
	} `cmd:"" aliases:"ls" help:"List tags and/or folders"`

	Rename struct {
		Src  string `arg,required,name:"src tag" short:"s" help:"Source tag name"`
		Dest string `arg,required,name:"dest tag" short:"d" help:"Destination tag name"`
	} `cmd:"" aliases:"mv" help:"Rename a tag"`

	Delete struct {
		Delete string `arg,required,name:"tag name" help:"Tag name to delete"`
	} `cmd:"" aliases:"rm" help:"Delete"`
}

// Run (*TagsCmd) ---
func (t *TagsCmd) Run(ctx *Context) error {
	// TODO: Refactor
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
	Unread bool   `arg,name:"Unread" xor:"read" help:"Mark <itemID> as unread"`
	Star   bool   `arg,name:"Star" short:"s" xor:"star" help:"Star or unstar <itemID>"`
	Unstar bool   `arg,name:"Unstar" xor:"star" help:"Unstar <itemID>"`
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
	Num           string `arg,optional,name:"num" short:"n" default:"5" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"n|o" short:"r" default:"n" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,name:"stream ID" short:"x" help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,name:"stream ID" short:"i" help:"Specify stream ID of target to include"`
	StreamID      string `arg,required,name:"stream ID" short:"s" help:"Specify stream ID of feed"`
	URLS          bool   `arg,optional,default:"false" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,default:"false" xor:"stream" help:"Display items with their IDs"`
}

// Run (*StreamCmd) ---
func (sc *StreamCmd) Run(ctx *Context) error {
	if sc.URLS {
		printStreamContentsWithURL(sc.Num, sc.Order, sc.ExcludeTarget, sc.IncludeTarget, sc.StreamID)
	} else if sc.IDS {
		printStreamContentsWithIDs(sc.Num, sc.Order, sc.ExcludeTarget, sc.IncludeTarget, sc.StreamID)
	} else {
		printStreamContentsWithDate(sc.Num, sc.Order, sc.ExcludeTarget, sc.IncludeTarget, sc.StreamID)
	}

	return nil
}

// ReadCmd ---
type ReadCmd struct {
	StreamID string `arg,required,name:"stream ID" short:"s" help:"Enter stream ID or feed URL to mark as read"`
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

	Subscription SubscriptionCmd `cmd:"" aliases:"sub" help:"Query subscriptions"`

	Tags TagsCmd `cmd:"" help:"Query tags and folders"`

	MarkItem MarkItemCmd `cmd:"" aliases:"mark" help:"Mark items as read/unread or starred/unstarred"`

	Stream StreamCmd `cmd:"" help:"Query stream contents"`

	Read ReadCmd `cmd:"" help:"Mark all items in stream as read"`
}

func main() {

	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
