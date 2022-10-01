# twx config

twx uses a simple toml configuration file. By default placed :
- (linux)
- (linux)
- (linux)
- (mac)

 It’s recommended to use `twx quickstart` to create it. 

Here’s an example conf file, showing every currently supported option:
```
[twtxt]
nick = "Nick"
twtfile = "~/twtxt.txt"
twturl = http://example.org/twtxt.txt
check_following = True
use_pager = False
use_cache = True
porcelain = False
disclose_identity = False
character_warning = 140
limit_timeline = 20
timeline_update_interval = 10
timeout = 5.0
sorting = "descending"
pre_tweet_hook = "scp buckket@example.org:~/public_html/twtxt.txt {twtfile}"
post_tweet_hook = "scp {twtfile} buckket@example.org:~/public_html/twtxt.txt"
# post_tweet_hook = "tail -1 {twtfile} | cut -f2 | sed -e 's/^/twt=/'| curl -s -d @- -d 'name=foo' -d 'password=bar' http://htwtxt.plomlompom.com/feeds"
# post_tweet_hook = "aws s3 cp {twtfile} s3://mybucket.org/twtxt.txt --acl public-read --storage-class REDUCED_REDUNDANCY --cache-control 'max-age=60,public'"

[following]
bob = "https://example.org/bob.txt"
alice = "https://example.org/alice.txt"

[posthook]
yarn_url = "https://twtxt.net/"
```

## [twtxt]

| Option: | Type: | Default: | Help: |
|---|---|---|---|
|nick|TEXT|   |your nick, will be displayed in your timeline|
|twtfile|PATH|   |path to your local twtxt file|
|twturl|TEXT|   |URL to your public twtxt file|
|check_following|BOOL|True|try to resolve URLs when listing followings|
|use_pager|BOOL|False|use a pager (less) to display your timeline |
|use_cache|BOOL|True|cache remote twtxt files locally|
|porcelain|BOOL|False|style output in an easy-to-parse format |
|disclose_identity|BOOL|False|include nick and twturl in twtxt’s user-agent and metadata |
|character_warning|INT|None|warn when composed tweet has more characters|
|limit_timeline|INT|20|limit amount of tweets shown in your timeline|
|timeline_update_interval|INT|10|time in seconds cache is considered up-to-date|
|timeout|FLOAT|5.0|maximal time a http request is allowed to take|
|sorting|TEXT|descending|sort timeline either descending or ascending  |
|use_abs_time|BOOL|False|use absolute datetimes in your timeline|
|pre_tweet_hook|TEXT|   |command to be executed before tweeting|
|post_tweet_hook|TEXT|   |command to be executed after tweeting|

pre_tweet_hook and post_tweet_hook are very useful if you want to push your twtxt file to a remote (web) server. Check the example above tho see how it’s used with scp.

## [followings]

This section holds all your followings as nick, URL pairs. You can edit this section manually or use the follow/unfollow commands of twtxt for greater comfort.

## [posthook]

depending of the post hook plugin you want to use, you need to add some config.
The key should be compose of the pluginNAme_Key example : YARN_URL, or GITHUB_SSH, ......
