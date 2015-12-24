package tw

import (
	"strings"
	"time"
)

type twitterTime struct {
	value time.Time
}

func (t twitterTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.value.Format(time.RubyDate) + `"`), nil
}

func (t twitterTime) UnmarshalJSON(data []byte) error {
	ts := strings.Trim(string(data), `"`)
	var err error
	t.value, err = time.Parse(time.RubyDate, ts)
	return err
}

type User struct {
	ID              uint64      `json:"id"                bson:"id"`
	Name            string      `json:"name"              bson:"name"`
	ScreenName      string      `json:"screen_name"       bson:"screen_name"`
	Description     string      `json:"description"       bson:"description"`
	ProfileImageURL string      `json:"profile_image_url" bson:"profile_image_url"`
	Location        string      `json:"location"          bson:"location"`
	Lang            string      `json:"lang"              bson:"lang"`
	TimeZone        string      `json:"time_zone"         bson:"time_zone"`
	URL             string      `json:"url"               bson:"url"`
	Protected       bool        `json:"protected"         bson:"protected"`
	Verified        bool        `json:"verified"          bson:"verified"`
	FriendsCount    int         `json:"friends_count"     bson:"friends_count"`
	ListedCount     int         `json:"listed_count"      bson:"listed_count"`
	FavouritesCount int         `json:"favourites_count"  bson:"favourites_count"`
	FollowersCount  int         `json:"followers_count"   bson:"followers_count"`
	StatusesCount   int         `json:"statuses_count"    bson:"statuses_count"`
	CreatedAt       twitterTime `json:"created_at"        bson:"created_at"`
}
