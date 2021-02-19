package main

type Context struct {
	Debug bool
}

type LoginCmd struct{}

type AddCmd struct {
	Add string `arg:"" required,placeholder:"streamID" name:"stream id" help:"Stream ID of subscription, i.e. user/-/label/Tech"`
}

func (a *AddCmd) Run(ctx *Context) error {
	execAddSub(a.Add)
	return nil
}

type SubscriptionCmd struct {
	List struct {
		All    bool `arg,required,name:"all" short:"a" xor:"list" help:"List all subscriptions"`
		Unread bool `arg,required,name:"unread" short:"u" xor:"list" help:"List only unread subscriptions"`
	} `cmd:"" aliases:"ls" help:"List subscriptions"`
	Edit struct {
		A bool `arg,required,placeholder:"edit|subscribe|unsubscribe" help:"Edit, subscribe, or unsubscribe"`
		S bool `arg,required,placeholder:"stream ID" help:"Stream ID of subscription feed"`
		T bool `arg,optional,placeholder:"title" help:"Change the subscription title. Omit to keep title."`
		F bool `arg,optional,placeholder:"folder" help:"Add subscription to folder. Use folder name like user/-/label/Tech"`
		R bool `arg,optional,placeholder:"folder" help:"Remove subscription from folder. Use folder name like user/-/label/Tech"`
	} `cmd:"" aliases:"ed" help:"Edit a subscription"`
}

// TODO Add Run function for each subcommand

type TagsCmd struct {
	List struct {
		All    string `arg,required,name:"all" short:"a" group:"list" xor:"list" help:"List all tags and/or folders"`
		Unread string `arg,required,name:"unread" short:"u" group:"list" xor:"list" help:"List only unread tags and/or folders"`
		Type   string `arg,optional,name:"type" short:"t" placeholder:"tags|folders" help:"List either only tags or only folders"`
	} `cmd:"" aliases:"ls" help:"List tags and/or folders"`

	Edit struct {
		A string `arg,optional,placeholder:"tag" help:"Tag to add"`
		R string `arg,optional,placeholder:"tag" help:"Tag to remove"`
		I string `arg,required,placeholder:"item ID" help:"Item ID of item to add/remove tag"`
	} `cmd:"" aliases:"ed" help:"Edit a tag on an item"`

	Rename struct {
		Src  string `arg,required,placeholder:"src tag" short:"s" help:"Source tag name"`
		Dest string `arg,required,placeholder:"dest tag" short:"d" help:"Destination tag name"`
	} `cmd:"" aliases:"mv" help:"Rename a tag"`

	Delete struct {
		Delete string `arg,required,placeholder:"tag name" help:"Tag name to delete"`
	} `cmd:"" aliases:"rm" help:"Delete"`
}

type StreamCmd struct {
	N  string `arg,optional,placeholder:"integer" default:"5" help:"Specify number of items returned"`
	R  string `arg,optional,placeholder:"n|o" help:"Specify order of items returned, i.e. newest or oldest"`
	Xt string `arg,optional,placeholder:"stream ID" help:"Specify stream ID of target to exclude"`
	It string `arg,optional,placeholder:"stream ID" help:"Specify stream ID of target to include"`
	S  string `arg,required,placeholder:"stream ID" help:"Specify stream ID of feed"`
	U  bool   `arg,optional,default:"false" help:"Display items with their urls"`
}

type ReadCmd struct {
	Read string `arg,required,placeholder:"stream ID" help:"stream ID to mark all as read"`
}

func main() {}
