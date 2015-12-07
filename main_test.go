package tw

import (
	"os"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	ck := os.Getenv("TWITTER_CONSUMER_KEY")
	cs := os.Getenv("TWITTER_CONSUMER_SECRET")
	_, err := GetAccessToken(ck, cs)
	if err != nil {
		t.Errorf("Expected to succeed but failed: %s", err)
	}
}
