package tw

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type FollowersPage struct {
	IDs            []uint64 `json:"ids"`
	NextCursor     int64    `json:"next_cursor"`
	PreviousCursor int64    `json:"previous_cursor"`
}

type RubyDate struct {
	value time.Time
}

func (t *RubyDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.value.Format(time.RubyDate) + `"`), nil
}

func (t *RubyDate) UnmarshalJSON(data []byte) error {
	ts := strings.Trim(string(data), `"`)
	var err error
	t.value, err = time.Parse(time.RubyDate, ts)
	return err
}

func (t *RubyDate) GetBSON() (interface{}, error) {
	return t.value, nil
}

func (t *RubyDate) SetBSON(raw bson.Raw) error {
	return raw.Unmarshal(&t.value)
}

type User struct {
	ID              uint64   `json:"id"                bson:"_id"`
	Name            string   `json:"name"              bson:"name"`
	ScreenName      string   `json:"screen_name"       bson:"screen_name"`
	Description     string   `json:"description"       bson:"description"`
	ProfileImageURL string   `json:"profile_image_url" bson:"profile_image_url"`
	Location        string   `json:"location"          bson:"location"`
	Lang            string   `json:"lang"              bson:"lang"`
	TimeZone        string   `json:"time_zone"         bson:"time_zone"`
	URL             string   `json:"url"               bson:"url"`
	Protected       bool     `json:"protected"         bson:"protected"`
	Verified        bool     `json:"verified"          bson:"verified"`
	FriendsCount    int      `json:"friends_count"     bson:"friends_count"`
	ListedCount     int      `json:"listed_count"      bson:"listed_count"`
	FavouritesCount int      `json:"favourites_count"  bson:"favourites_count"`
	FollowersCount  int      `json:"followers_count"   bson:"followers_count"`
	StatusesCount   int      `json:"statuses_count"    bson:"statuses_count"`
	CreatedAt       RubyDate `json:"created_at"        bson:"created_at"`
}

type Tweet struct {
	ID           uint64   `json:"id"                 bson:"_id"`
	Text         string   `json:"text"               bson:"text"`
	Retweeted    bool     `json:"retweeted"          bson:"retweeted"`
	RetweetCount uint     `json:"retweet_count"      bson:"retweet_count"`
	Sensitive    bool     `json:"possibly_sensitive" bson:"possibly_sensitive"`
	CreatedAt    RubyDate `json:"created_at"         bson:"created_at"`
}

type Client struct {
	consumerKey       string
	consumerSecret    string
	bearerAccessToken string
}

type FollowersIterator struct {
	client     *Client
	userID     uint64
	screenName string
	count      int
	cursor     int64
}

func (t *FollowersIterator) Next(data *[]uint64) error {
	if t.cursor == 0 {
		return errors.New("No more remaining pages")
	}
	url := fmt.Sprintf("%s/followers/ids.json?count=%d&cursor=%d",
		baseURL, t.count, t.cursor)
	if t.userID != 0 {
		url += fmt.Sprintf("&user_id=%d", t.userID)
	} else {
		url += "&screen_name=" + t.screenName
	}
	req, err := t.client.prepareRequest("GET", url)
	if err != nil {
		return err
	}
	var resp FollowersPage
	if err = exec(req, &resp); err != nil {
		return err
	}
	t.cursor = resp.NextCursor
	*data = resp.IDs
	return nil
}
