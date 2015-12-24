package tw

import (
	"log"
	"os"
	"strings"
	"testing"
)

var (
	ck string
	cs string
	tc *Client
)

func init() {
	ck = os.Getenv("TWITTER_CONSUMER_KEY")
	cs = os.Getenv("TWITTER_CONSUMER_SECRET")
}

func setup() {
	tc = NewClient()
	if err := tc.SetKeys(ck, cs); err != nil {
		log.Fatalf("Failed to setup client")
	}
}

func TestGetAccessToken(t *testing.T) {
	_, err := GetBearerAccessToken(ck, cs)
	if err != nil {
		t.Errorf("Expected to succeed but failed: %s", err)
	}
}

func TestGetUsersShow(t *testing.T) {
	setup()

	user, err := tc.GetUsersShow("twitterdev")
	if err != nil {
		t.Fatalf("Failed to obtain user: %s", err)
	}

	if user.ID != 2244994945 {
		t.Errorf("Expected ID to be 2244994945, got %v", user.ID)
	}
}

func TestGetUsersShowByID(t *testing.T) {
	setup()

	user, err := tc.GetUsersShowByID(int64(2244994945))
	if err != nil {
		t.Fatalf("Failed to obtain user: %s", err)
	}

	screenName := user.ScreenName
	if strings.ToLower(screenName) != "twitterdev" {
		t.Errorf("Expected user to be %q, got %q", "twitterdev", screenName)
	}
}

func TestTweetsByID(t *testing.T) {
	setup()

	tweets, err := tc.GetTweetsByID(2244994945, 10)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(tweets) != 10 {
		log.Fatalf("Expected to get 10 tweets, got %d", len(tweets))
	}
}

// TODO: find a reproducible way to test rate limits
func TestTooMuchRequests(t *testing.T) {
	ckTMR := os.Getenv("TWITTER_CONSUMER_KEY_TMR")
	csTMR := os.Getenv("TWITTER_CONSUMER_SECRET_TMR")
	tc := NewClient()
	if err := tc.SetKeys(ckTMR, csTMR); err != nil {
		t.Fatalf("Failed to setup client")
	}
	log.Printf("Too Much Requests...")
	for i := 0; i < 200; i++ {
		log.Printf("Request #%d", i)
		_, err := tc.GetUsersShow("twitterdev")
		if err != nil {
			if err != ErrTooManyRequests {
				t.Fatalf("Expected %s, got %s", ErrTooManyRequests, err)
			} else {
				break
			}
		}
	}
}
