# TODO go-inoreader

- [x] Add the remaining base API methods
- [x] Abstract some functions to avoid repetition
    - [x] Functions that send POST reqs
- [ ] Add .go file for higher-order functions that process base API methods and return useful results
    - [ ] Allow various output formats for streams/items: JSON, CSV
    - [ ] Add feature to format stream/item metadata and output into markdown file
    - [ ] Add feature to download OPML
    - [ ] Use [flag](https://golang.org/pkg/flag/) package to handle command line args
- [ ] Remove unneeded type definitions
- [x] Rename pkg/inoreader to pkg/api
- [ ] Add configuration settings
    - [x] OAuth2
    - [x] Save client data to $HOME/.local/share/go-inoreader.json
    - [ ] Integrate configuration with pkg/client and pkg/api
- [ ] Use ParseForm for handling url query parameters in pkg/inoreader/base.go