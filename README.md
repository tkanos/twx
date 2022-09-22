# twx

A decentralised microblogging client based on the specs of [twtxt](https://dev.twtxt.net/), used to handle your twtxt file.



# Todo :
## Timeline : 
- replace on timeline @nick or @<nick url> and present it like : @nick@url
- highlight tags
- add mkdown integration + color
- use_pager: use a pager (less) to display your timeline (sorting not apply here)
- porcelain style output in an easy-to-parse format (without pretty things)
- limit_timeline: limit amount of tweets shown in your timeline
- use_abs_time: use absolute datetimes in your timeline
- timeline --web
- timeline --gemini
- timeline generates microformat v2 (https://microformats.org/wiki/microformats2) => html to see
- timeline with thread organized

## Tweet :
- Tweet @nick will replace before to append
- character_warning: warn when composed tweet has more characters
- date + Encryption 

## Cache :
- use_cache : cache remote twtxt files locally (with parquet)
- timeline_update_interval: time in seconds cache is considered up-to-date

## New Commands :
- Thread command
- Tag Command (like Thread actually)
- see users's Tweets Command
- User Profile (metadata) Command

## Hook :
- yarn connection post (with keeping connection)
- github integration (nick/twtxt.txt)
- drive integration
- one drive integration 

## CI/CD :
- cmd line test : https://github.com/google/go-cmdtest or other script
- https://github.com/orlangure/gocovsh
- Unit Tests
- Test on Apple
- https://goreleaser.com/
- circleCI


## README.md
- do README file
- do a demo like : https://github.com/orlangure/gocovsh

## Web-Server :
- local web server
- local gemini server