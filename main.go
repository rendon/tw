// Package tw implements part of the Twitter API, aiming simplicity and
// flexibility.
//
// Information about the Twitter API can be found at
// https://dev.twitter.com/rest/public.
package tw

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL = "https://api.twitter.com/1.1"
	authURL = "https://api.twitter.com/oauth2/token"

	// MaxFollowersCount More info at
	// https://dev.twitter.com/rest/reference/get/followers/ids
	MaxFollowersCount = 5000

	// MaxFriendsCount More info at
	// https://dev.twitter.com/rest/reference/get/friends/ids
	MaxFriendsCount = 5000
)

var (
	// ErrMsgTooManyRequests rate limit message.
	ErrMsgTooManyRequests = "Too Many Requests"
	// ErrUnauthorized error for private profiles.
	ErrUnauthorized = errors.New("Authorization Required")
)

// RateLimitError indicates a "Too Many Requests" error and associated data.
type RateLimitError struct {
	ResetTime time.Time
}

func (t RateLimitError) Error() string {
	return ErrMsgTooManyRequests
}

// GetBearerAccessToken Authenticates  with Twitter using the  provided consumer
// key  and  consumer  secret.  Details  of the  algorithmn  are  available  at
// https://dev.twitter.com/oauth/application-only.
func (c *Client) GetBearerAccessToken() error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", authURL, nil)
	req.Header.Add("User-Agent", "My Twitter app")
	ck := url.QueryEscape(c.ConsumerKey)
	cs := url.QueryEscape(c.ConsumerSecret)
	bt := base64.StdEncoding.EncodeToString([]byte(ck + ":" + cs))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", bt))
	req.Header.Add("Content-Type",
		"application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", "29")
	req.Header.Add("Accept-Encoding", "gzip")

	body := []byte("grant_type=client_credentials")
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}
	c.BearerAccessToken = data["access_token"].(string)
	return nil
}

// NewClient Returns a new client with credentials set.
func NewClient(ck, cs string) *Client {
	return &Client{
		ConsumerKey:    ck,
		ConsumerSecret: cs,
	}
}

func (c *Client) prepareRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "My Twitter App")
	auth := fmt.Sprintf("Bearer %s", c.BearerAccessToken)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept-Encoding", "gzip")
	return req, err
}

func exec(req *http.Request, data interface{}) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	rb, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	// Too Many Requests
	if resp.StatusCode == 429 {
		var reset time.Time
		rateLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
		t, err := strconv.ParseInt(rateLimitReset, 10, 64)
		if err != nil {
			reset = time.Now().Add(16 * time.Minute)
		} else {
			reset = time.Unix(t, 0)
		}
		return RateLimitError{ResetTime: reset}
	}
	if resp.StatusCode == 401 {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	if err = json.Unmarshal(rb, data); err != nil {
		return err
	}
	return nil
}

// GetUsersShow Retrieves  user profile given  the user's screen name.  For more
// information see https://dev.twitter.com/rest/reference/get/users/show.
func (c *Client) GetUsersShow(screenName string) (*User, error) {
	screenName = url.QueryEscape(screenName)
	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", baseURL, screenName)
	req, err := c.prepareRequest("GET", url)
	if err != nil {
		return nil, err
	}
	var user User
	err = exec(req, &user)
	return &user, err
}

// GetUsersShowByID Retrieves user  profile given it's ID.  For more information
// see https://dev.twitter.com/rest/reference/get/users/show.
func (c *Client) GetUsersShowByID(id int64) (*User, error) {
	url := fmt.Sprintf("%s/users/show.json?user_id=%d", baseURL, id)
	req, err := c.prepareRequest("GET", url)
	if err != nil {
		return nil, err
	}
	var user User
	err = exec(req, &user)
	return &user, err
}

// GetTweets Retrieves latest tweets, limited by count. For more information see
// https://dev.twitter.com/rest/reference/get/statuses/user_timeline.
func (c *Client) GetTweets(screenName string, count uint) ([]Tweet, error) {
	screenName = url.QueryEscape(screenName)
	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d",
		baseURL, screenName, count)
	req, err := c.prepareRequest("GET", url)
	var tweets []Tweet
	if err != nil {
		return tweets, err
	}
	err = exec(req, &tweets)
	return tweets, err
}

// GetTweetsByID Retrieves  latest tweets  by user  ID, limited  by count.  See
// https://dev.twitter.com/rest/reference/get/statuses/user_timeline.
func (c *Client) GetTweetsByID(id int64, count uint) ([]Tweet, error) {
	url := fmt.Sprintf("%s/statuses/user_timeline.json?user_id=%d&count=%d",
		baseURL, id, count)
	var tweets []Tweet
	req, err := c.prepareRequest("GET", url)
	if err != nil {
		return tweets, err
	}
	err = exec(req, &tweets)
	return tweets, err
}

// GetFollowersIdsByID Returns an iterator which  you can call to retrieve pages
// of followers  IDs for the user  identified with ID. Count  specifies the page
// size. See https://dev.twitter.com/rest/reference/get/followers/ids.
func (c *Client) GetFollowersIdsByID(id int64, count int) *FollowersIterator {
	return &FollowersIterator{
		client: c,
		userID: id,
		count:  count,
		cursor: -1,
	}
}

// GetFriendsIdsByID Returns an iterator which you can call to retrieve pages of
// friends IDs for  the user identified with ID. Count  specifies the page size.
// See https://dev.twitter.com/rest/reference/get/friends/ids.
func (c *Client) GetFriendsIdsByID(id int64, count int) *FriendsIterator {
	return &FriendsIterator{
		client: c,
		userID: id,
		count:  count,
		cursor: -1,
	}
}
