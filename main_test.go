package tw

import (
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
	tc := NewClient(ck, cs)
	if err := tc.Setup(); err != nil {
		t.Errorf("Failed to setup client")
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
