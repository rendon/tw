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
)

const (
	baseURL = "https://api.twitter.com/1.1"
	authURL = "https://api.twitter.com/oauth2/token"

	MaxFollowersCount = 5000
	MaxFriendsCount   = 5000
)

var (
	ErrTooManyRequests = errors.New("Too Many Requests")
	ErrUnauthorized    = errors.New("Authorization Required")
)

func GetBearerAccessToken(consumerKey, consumerSecret string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", authURL, nil)
	req.Header.Add("User-Agent", "My Twitter app")
	ck := url.QueryEscape(consumerKey)
	cs := url.QueryEscape(consumerSecret)
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
		return "", err
	}
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return "", err
	}
	return data["access_token"].(string), nil
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetKeys(consumerKey, consumerSecret string) error {
	bat, err := GetBearerAccessToken(consumerKey, consumerSecret)
	if err != nil {
		return err
	}
	c.consumerSecret = consumerSecret
	c.consumerKey = consumerKey
	c.bearerAccessToken = bat
	return nil
}

func (c *Client) prepareRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "My Twitter App")
	auth := fmt.Sprintf("Bearer %s", c.bearerAccessToken)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept-Encoding", "gzip")
	return req, err
}

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
		return ErrTooManyRequests
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

func (c *Client) GetTweets(screenName string, count uint) ([]Tweet, error) {
	screenName = url.QueryEscape(screenName)
	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d",
		baseURL, screenName, count)
	req, err := c.prepareRequest("GET", url)
	tweets := make([]Tweet, 0)
	if err != nil {
		return tweets, err
	}
	err = exec(req, &tweets)
	return tweets, err
}

func (c *Client) GetTweetsByID(id int64, count uint) ([]Tweet, error) {
	url := fmt.Sprintf("%s/statuses/user_timeline.json?user_id=%d&count=%d",
		baseURL, id, count)
	tweets := make([]Tweet, 0)
	req, err := c.prepareRequest("GET", url)
	if err != nil {
		return tweets, err
	}
	err = exec(req, &tweets)
	return tweets, err
}

func (c *Client) GetFollowersIdsByID(id int64, count int) *FollowersIterator {
	return &FollowersIterator{
		client: c,
		userID: id,
		count:  count,
		cursor: -1,
	}
}

func (c *Client) GetFriendsIdsByID(id int64, count int) *FriendsIterator {
	return &FriendsIterator{
		client: c,
		userID: id,
		count:  count,
		cursor: -1,
	}
}
