# Command structure

* `inoreader` - Main command. If no args supplied, display usage info
    * `login` - Initiate OAuth, save info to config
    * `help` - Display help info
    * `sub` - Query subscriptions; requires only one of below args
        * `list` - List subscriptions
        * `unread` - List only unread subs with counts
    * `tags` - Query tags and folders; requires only one of below args
        * `list` - List all tags and folders; optional below args
            * `-t|--type <tags|folders>` - Specify and list only tags or folders (string)
        * `unread` - List only folders containing unread items with counts
    * `stream <streamID>` - Query stream contents; requires at least one of below args
        * `-n <int>` - Specify numbers of items returned
        * `-r <newest|oldest>` - Specify order of items displayed
        * `-xt <streamID>` - Specify streamID of target to exclude
        * `-it <streamID>` - Specify streamID of target to include
        * `-s <itemID>` - Display summary of item
    * `output <streamID|itemID>` - Fetch and save content to output format; requires at least one of below args
        * `-m` - Export to Markdown
        * `-c` - Export to CSV
        * `-j` - Export to JSON
