#tw
A new and simple Twitter API client.

## Usage

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rendon/tw"
)

func rotateKeys(tc *tw.Client) {
	ck := os.Getenv("EXTRA_TWITTER_CONSUMER_KEY")
	cs := os.Getenv("EXTRA_TWITTER_CONSUMER_SECRET")
	if err := tc.SetKeys(ck, cs); err != nil {
		log.Fatalf("Error setting keys: %s", err)
	}
	log.Printf("New keys set, go!\n")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Create new client
	tc := tw.NewClient()

	// Set keys
	ck := os.Getenv("TWITTER_CONSUMER_KEY")
	cs := os.Getenv("TWITTER_CONSUMER_SECRET")
	if err := tc.SetKeys(ck, cs); err != nil {
		log.Fatalf("Failed to setup client")
	}

	// GET users/show
	user, err := tc.GetUsersShow("twitterdev")
	if err != nil {
		log.Fatalf("Failed to obtain user: %s", err)
	}

	fmt.Printf("User ID: %d\n", user.ID)
	fmt.Printf("User name: %s\n", user.ScreenName)

	// GET followers/ids
	fmt.Printf("First follower IDs:\n")
	followers := tc.GetFollowersIdsByID(2244994945, 5)
	var ids []uint64
	for i := 0; i < 10; i++ {
		err = followers.Next(&ids)
		if err != nil {
			if err == tw.ErrTooManyRequests {
				rotateKeys(tc)
				i--
			} else if err == tw.ErrNoMorePages {
				break
			} else {
				log.Fatal(err)
			}
		}
		for _, id := range ids {
			fmt.Printf("%d\n", id)
		}
	}
	fmt.Println()

	// GET friends/ids
	fmt.Printf("First friends IDs:\n")
	friends := tc.GetFriendsIdsByID(191541009, 50)
	for i := 0; i < 10; i++ {
		err = friends.Next(&ids)
		if err != nil {
			if err == tw.ErrTooManyRequests {
				rotateKeys(tc)
				i--
			} else if err == tw.ErrNoMorePages {
				break
			} else {
				log.Fatal(err)
			}
		}
		for _, id := range ids {
			fmt.Printf("%d\n", id)
		}
	}
	fmt.Println()
}
```

## Detecting "Too Many Requests"
The client was designed to use multiple sets of keys, if you reach the rate limit (`err` will be equal to `ErrTooManyRequests`), just call `SetKeys()` with a new set of keys and keep going.

```go
// GET users/show
user, err := tc.UsersShow("twitterdev")
if err == tw.ErrTooManyRequests {
    // Set a new set of keys and try again
}
```
