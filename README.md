# twx

A decentralised microblogging client based on the specs of [twtxt](https://dev.twtxt.net/), used to handle your twtxt file.

## Why another client
I love twtxt format for its simplicity and being a real decentralized communication protocol. 
But I see some issues that I want to focus. So i created my own twtxt client because I want to be able to make tests to evolve the format, on what I think lack to the protocol :

1. Decentralized ID
2. Tweet Encryption
3. Metadata (Profile) Encryption
4. Discoverability
5. Ability to create rooms
6. Ability to administrate rooms

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

### Todo ( more precise ) :

## Timeline : 
- replace on timeline @nick or @<nick url> and present it like : @nick@url
- highlight tags
- add mkdown integration + color
- use_pager: use a pager (less) to display your timeline (sorting not apply here)
- porcelain style output in an easy-to-parse format (without pretty things)
- use_abs_time: use absolute datetimes in your timeline
- timeline --web
- timeline --gemini
- timeline generates microformat v2 (https://microformats.org/wiki/microformats2) => html to see
- timeline with thread organized

## Tweet :
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
