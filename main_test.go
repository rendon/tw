package tw

import (
	"log"
	"os"
	"testing"
)

var (
	ck string
	cs string
)

func init() {
	ck = os.Getenv("TWITTER_CONSUMER_KEY")
	cs = os.Getenv("TWITTER_CONSUMER_SECRET")
}

func TestGetAccessToken(t *testing.T) {
	_, err := GetBearerAccessToken(ck, cs)
	if err != nil {
		t.Errorf("Expected to succeed but failed: %s", err)
	}
}

func TestUsersShow(t *testing.T) {
	tc := NewClient()
	if err := tc.SetKeys(ck, cs); err != nil {
		t.Fatalf("Failed to setup client")
	}

	data, err := tc.UsersShow("twitterdev")
	if err != nil {
		t.Fatalf("Failed to obtain user: %s", err)
	}

	id := data["id_str"].(string)
	if id != "2244994945" {
		t.Errorf("Expected ID to be 2244994945, got %v", id)
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
		_, err := tc.UsersShow("twitterdev")
		if err != nil {
			if err != ErrTooManyRequests {
				t.Fatalf("Expected %s, got %s", ErrTooManyRequests, err)
			} else {
				break
			}
		}
	}
}
