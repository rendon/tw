package tw

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrEndOfList Indicates  we've reached the  last page for those  queries
	// that fetch data in pages, e.g. get followers, get friends, etc.
	ErrEndOfList = errors.New("No more pages available")
)

// FollowersPage describes a page of follower IDs.
type FollowersPage struct {
	IDs            []int64 `json:"ids"`
	NextCursor     int64   `json:"next_cursor"`
	PreviousCursor int64   `json:"previous_cursor"`
}

// FriendsPage describes a page of friend IDs.
type FriendsPage struct {
	IDs            []int64 `json:"ids"`
	NextCursor     int64   `json:"next_cursor"`
	PreviousCursor int64   `json:"previous_cursor"`
}

// RubyDate is a custom type to handle Twitter dates.
type RubyDate struct {
	value time.Time
}

// MarshalJSON returns JSON representation of RubyDate type.
func (t *RubyDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.value.Format(time.RubyDate) + `"`), nil
}

// UnmarshalJSON gets RubyDate value from data.
func (t *RubyDate) UnmarshalJSON(data []byte) error {
	ts := strings.Trim(string(data), `"`)
	var err error
	t.value, err = time.Parse(time.RubyDate, ts)
	return err
}

// GetBSON returns value that can be represented in BSON format.
func (t *RubyDate) GetBSON() (interface{}, error) {
	return t.value, nil
}

// SetBSON sets RubyDate value from raw BSON value.
func (t *RubyDate) SetBSON(raw bson.Raw) error {
	return raw.Unmarshal(&t.value)
}

// User Represents a user with some important fields.
type User struct {
	ID              int64     `json:"id"                bson:"_id"`
	Name            string    `json:"name"              bson:"name"`
	ScreenName      string    `json:"screen_name"       bson:"screen_name"`
	Description     string    `json:"description"       bson:"description"`
	ProfileImageURL string    `json:"profile_image_url" bson:"profile_image_url"`
	Location        string    `json:"location"          bson:"location"`
	Lang            string    `json:"lang"              bson:"lang"`
	TimeZone        string    `json:"time_zone"         bson:"time_zone"`
	URL             string    `json:"url"               bson:"url"`
	Protected       bool      `json:"protected"         bson:"protected"`
	Verified        bool      `json:"verified"          bson:"verified"`
	FriendsCount    int       `json:"friends_count"     bson:"friends_count"`
	ListedCount     int       `json:"listed_count"      bson:"listed_count"`
	FavouritesCount int       `json:"favourites_count"  bson:"favourites_count"`
	FollowersCount  int       `json:"followers_count"   bson:"followers_count"`
	StatusesCount   int       `json:"statuses_count"    bson:"statuses_count"`
	CreatedAt       *RubyDate `json:"created_at"        bson:"created_at"`
	RetrievedAt     time.Time `json:"retrieved_at"      bson:"retrieved_at"`
}

// UserMention represents a mention in a tweet.
type UserMention struct {
	ID int64 `json:"id" bson:"id"`
}

// Entities represents a list of mentions in a tweet.
type Entities struct {
	UserMentions []UserMention `json:"user_mentions"   bson:"user_mentions"`
}

// Tweet Represents a tweet with some important fields.
type Tweet struct {
	ID            int64    `json:"id"                 bson:"id"`
	UserID        int64    `json:"user_id"            bson:"user_id"`
	Text          string   `json:"text"               bson:"text"`
	Retweeted     bool     `json:"retweeted"          bson:"retweeted"`
	IsRetweet     bool     `json:"is_retweet"         bson:"is_retweet"`
	RetweetCount  uint     `json:"retweet_count"      bson:"retweet_count"`
	FavoriteCount uint     `json:"favorite_count"     bson:"favorite_count"`
	Sensitive     bool     `json:"possibly_sensitive" bson:"possibly_sensitive"`
	Entities      Entities `json:"entities"           bson:"entities"`
	CreatedAt     RubyDate `json:"created_at"         bson:"created_at"`
}

// Client Represents a Twitter API client with the necessary access data and
// methods to query the API.
type Client struct {
	ConsumerKey       string
	ConsumerSecret    string
	BearerAccessToken string
}

// FollowersIterator Contains the necessary information to retrieve the next
// page of follower IDs.
type FollowersIterator struct {
	client     *Client
	userID     int64
	screenName string
	count      int
	cursor     int64
}

// FriendsIterator Contains the necessary information to retrieve the next
// page of friends IDs.
type FriendsIterator struct {
	client     *Client
	userID     int64
	screenName string
	count      int
	cursor     int64
}

// Next Returns the next page of follower IDs.
func (t *FollowersIterator) Next(data *[]int64) error {
	if t.cursor == 0 {
		return ErrEndOfList
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

// Next Returns the next page of friends IDs.
func (t *FriendsIterator) Next(data *[]int64) error {
	if t.cursor == 0 {
		return ErrEndOfList
	}
	url := fmt.Sprintf("%s/friends/ids.json?count=%d&cursor=%d",
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
	var resp FriendsPage
	if err = exec(req, &resp); err != nil {
		return err
	}
	t.cursor = resp.NextCursor
	*data = resp.IDs
	return nil
}
