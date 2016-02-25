package tw

import (
	"fmt"
	"log"
	"os"
)

func ExampleNewClient() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Create new client
	ck := os.Getenv("TWITTER_CONSUMER_KEY")
	cs := os.Getenv("TWITTER_CONSUMER_SECRET")
	tc := NewClient(ck, cs)

	// Set keys
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
		if err == ErrEndOfList {
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
		if err == ErrEndOfList {
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
