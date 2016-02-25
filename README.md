#tw
[![GoDoc](https://godoc.org/github.com/rendon/tw?status.svg)](https://godoc.org/github.com/rendon/tw)

[![Codeship](https://codeship.com/projects/a8d73b50-be36-0133-01f5-327e6e59f642/status?branch=master)](https://godoc.org/github.com/rendon/tw)

A new and simple Twitter API client.

**NOTE:** This client does not implement the entire Twitter API, even more, it only implements functions for a few endpoints, those that I need in my current project. You're welcome to contribute or ask for features.

## Usage
The basic usage is:

Create a new client:
```go
tc := tw.NewClient(ck, cs)
```
Get access token before any query:
```go
if err := tc.GetBearerAccessToken(); err != nil {
    log.Fatalf("Failed to obtain access token: %s", err)
}
```

Go for it!
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rendon/tw"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Create new client
	ck := os.Getenv("TWITTER_CONSUMER_KEY")
	cs := os.Getenv("TWITTER_CONSUMER_SECRET")
	tc := tw.NewClient(ck, cs)

	if err := tc.GetBearerAccessToken(); err != nil {
		log.Fatalf("Failed to obtain access token: %s", err)
	}

	// GET users/show
	user, err := tc.GetUsersShow("twitterdev")
	if err != nil {
		log.Fatalf("Failed to obtain user: %s", err)
	}

	fmt.Printf("User ID: %d\n", user.ID)
	fmt.Printf("User name: %s\n", user.ScreenName)

	// GET statuses/user_timeline
	tweets, err := tc.GetTweets("twitterdev", 5)
	if err != nil {
		log.Fatalf("Failed to obtain tweets: %s", err)
	}

	fmt.Printf("=================== Tweets =================================\n")
	for i := range tweets {
		fmt.Printf("%s\n", tweets[i].Text)
		fmt.Printf("--------------------------------------------------------\n")
	}
	fmt.Println()

	// GET followers/ids
	fmt.Printf("First follower IDs:\n")
	followers := tc.GetFollowersIdsByID(2244994945, 5)
	var ids []int64
	for i := 0; i < 5; i++ {
		err = followers.Next(&ids)
		if err == tw.ErrEndOfList {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, id := range ids {
			fmt.Printf("%d\n", id)
		}
	}
	fmt.Println()

	// GET friends/ids
	fmt.Printf("First friends IDs:\n")
	friends := tc.GetFriendsIdsByID(191541009, 50)
	for i := 0; i < 5; i++ {
		err = friends.Next(&ids)
		if err == tw.ErrEndOfList {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, id := range ids {
			fmt.Printf("%d\n", id)
		}
	}
	fmt.Println()
}
```

## Detecting "Too Many Requests"
You can detect a rate limit error like so:

```go
if err != nil && err.Error() == tw.ErrMsgTooManyRequests {
    // Do something about it
}
```

Or use type assertion. The actual error in this case is of type `tw.RateLimitError`, which contains a `ResetTime` field with the time at which the access will be renewed. You can do something like this in the meanwhile:

```go
time.Sleep(err.ResetTime.Sub(time.Now()))
```

Originally this client was designed to use multiple sets of keys, rotating keys when detecting `Too Many Requests` errors. Now the recommended way is to create a list of clients, one per set of credentials and use some sort of balancing algorithm to decide which client to use.
