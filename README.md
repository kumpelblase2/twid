# twid

Twitch recently started requiring user authentication on most API endpoints, even for endpoints that previously were accessible with just a client ID.
Unfortunately, that makes it no longer feasable to write small scripts or programs for it as you'd now need to take care of authenticating the user first.
Reason being that you'd either need to setup a server to accept the redirection request or provide a website the then transmits the created token to your local program.
Neither of which are particularly easy to do in scripts or at least add a lot of unnecessary cruft.

Thus, instead you can just put this tool in front to get a token for you and continue on from there. I use it inside my follower selection script inside my [dotfiles](https://github.com/kumpelblase2/dotfiles/blob/master/bin/twitch-select).

# Usage
Place the `twid` executable somewhere (preferably inside your path) and run it like this:
```bash
token=$(twid <client_id> <client_secret>)
```
You can then use `token` as the access token for making api requests. Sadly, `twid` needs the client secret too since that is required by Twitch to actually get an access token.

# Building

To build the tool, make sure you have the go toolchain installed and on the path, then just run the following:
```shell
go build .
```
Which should leave you with the `twid` executable. Make sure it's executabe ( `chmod +x twid` )!
