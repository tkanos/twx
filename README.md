# twx

A decentralised microblogging client based on the specs of [twtxt](https://dev.twtxt.net/), used to handle your twtxt file.


Todo :
Before to push
- Merge follow and unfollow difference = f(); save file on Post() + check size before to rename or last tweet ?
- PreHook and PostHook
- Tweet @nick will replace before to append
- replace on timeline  @nick or @<nick url> and present it like : @nick@url
- Add readme.md for configuration
- do README file

After Push :
 - add mkdown integration + color

- Use toml configuration options
    //timeout: maximal time a http request is allowed to  adn refactor http client

    //use_pager: use a pager (less) to display your timeline
	//porcelain style output in an easy-to-parse format
	//character_limit: shorten incoming tweets with more characters
	//character_warning: warn when composed tweet has more characters
	//limit_timeline: limit amount of tweets shown in your timeline
	//use_abs_time: use absolute datetimes in your timeline

    //use_cache : cache remote twtxt files locally (with parquet)
    //timeline_update_interval: time in seconds cache is considered up-to-date

- read a Thread command
- read a Tag Command (like Thread actually)
- see tweets an user if pager more style command
- see the profile of one user (metadata) command

- yarn connection post (with keeping connection)

- timeline generates microformat v2 (https://microformats.org/wiki/microformats2) => html to see
- timeline with thread organized
- Encryption
- unit test coverage
- cmd line test : https://github.com/google/go-cmdtest or other script
- surge ? ngrok ?
- https://github.com/orlangure/gocovsh
- do a demo like : https://github.com/orlangure/gocovsh
- https://goreleaser.com/
- ci/cd
- timeline --web
- timeline gemini
- local web server
- local gemini server
- github integration (nick.twtxt.txt)
- drive integration
- one drive integration 