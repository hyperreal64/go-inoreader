# TODO 

As of 2021-02-11 05:15:23

- [x] Add base API methods
- [ ] `cmd/go-inoreader/main.go`
    - [ ] Use [flag](https://golang.org/pkg/flag/) package to handle command-line args
    - [ ] Implement [command structure](cli-structure.md)
- [ ] Remove unneeded type definitions
- [ ] Add configuration settings
    - [x] OAuth2 flow
    - [ ] Make HTML template conditional on success/failure of OAuth2 flow
    - [x] Save configuration info to `$XDG_DATA_HOME/go-inoreader.json` on Unix/Linux and `%APPDATA%\go-inoreader.json` on Windows
- [x] Use [go-querystring](https://github.com/google/go-querystring) for more type-safe query parameter handling
    - [x] Add type definitions for base API methods that use query parameters
- [ ] Support various output formats for streams/items: JSON, CSV
- [ ] Support formatting stream/item metadata and output for markdown
- [ ] Support downloading OPML for subscriptions list
