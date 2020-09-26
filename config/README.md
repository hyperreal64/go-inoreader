# config

See [douglasmakey/oauth2-example](https://github.com/douglasmakey/oauth2-example)

* `config.go`: handles configuration file, such as determining its location in the user's home directory, writing to the file, and retrieving the data for use with Inoreader API calls. 
* The configuration file will be written to the following locations:
  * macOS, Linux, *BSD: `$HOME/.local/share`
  * Windows: `$env:LOCALAPPDATA\go-inoreader.json`
* The configuration file will look like the example below:

```json
{
    "userId": "<user id>",
    "userName": "<username>",
    "Oauth2Response": {
        "access_token": "<access token>",
        "token_type": "Bearer",
        "refresh_token": "<refresh token>",
        "expiry": "2020-09-18T16:51:44.180363409-05:00"
    }
}
```

* `init.go`: initiates the OAuth2 flow for access and refresh tokens.

* `oauthutil.go`: the OAuth handlers and their helper functions.