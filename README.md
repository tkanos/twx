# twx

A decentralised microblogging client based on the specs of [twtxt](https://dev.twtxt.net/), used to handle your twtxt file.

## Why another client
I love twtxt format for it's simplicity and being real decentralized communication protocol. 
But I see some issues that I want to focus. So i craeted my own twtxt client because I was not happy with the existing one, and to be able to make tests to evolve the format on what I think lack to the protocol :

1. Having a nice Client, I can evolve and do pocs
2. Decentralized ID
3. Tweet Encryption
4. Metadata (Profile) Encryption
5. Discoverability

## Already Done :
twx can :
- post a tweet
- reply to a tweet
- create a timeline
- follow people
- unfollow people
- preHook and postHook


## Todo :
twx can't *yet* :
- Have a pretty Timeline (in many format)
- Have caching
- Have plugin hooks (to github / gdrive /one drive / ...)
- Lot of useful commands (thread, tags, profile, ....)
- create a local webserver
- Encryption
- Discoverability
- Decentralized ID

### Todo in more precise way :

## Timeline : 
- replace on timeline @nick or @<nick url> and present it like : @nick@url
- highlight tags
- add mkdown integration + color
- use_pager: use a pager (less) to display your timeline (sorting not apply here)
- porcelain style output in an easy-to-parse format (without pretty things)
- use_abs_time: use absolute datetimes in your timeline
- timeline_show_ascii_images Show images on ascii
- timeline --web
- timeline --gemini
- timeline generates microformat v2 (https://microformats.org/wiki/microformats2) => html to see
- timeline with thread organized

## Tweet :
- Tweet @nick will replace before to append
- character_warning: warn when composed tweet has more characters
- Encryption 

## Cache :
- use_cache : cache remote twtxt files locally (with parquet)
- timeline_update_interval: time in seconds cache is considered up-to-date

## New Commands :
- Thread command
- Tag Command (like Thread actually)
- see users's Tweets Command
- User Profile (metadata) Command
- Mentions and Replied Thread Commands
- Discover Command (download the followings of users that you follow ....)
- Should we separate feed from follow ?
- version console interactive (IRC like)

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
