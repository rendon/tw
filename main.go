package tw

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://api.twitter.com/1.1"
	authURL = "https://api.twitter.com/oauth2/token"
)

var (
	ErrTooManyRequests = errors.New("Too Many Requests")
)

type Client struct {
	consumerKey       string
	consumerSecret    string
	bearerAccessToken string
	httpClient        *http.Client
}

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
	return &Client{httpClient: &http.Client{}}
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

func (c *Client) GetUsersShow(user string) (map[string]interface{}, error) {
	user = url.QueryEscape(user)
	urlStr := fmt.Sprintf("%s/users/show.json?screen_name=%s", baseURL, user)
	req, err := c.prepareRequest("GET", urlStr)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	rb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Too Many Requests
	if resp.StatusCode == 429 {
		return nil, ErrTooManyRequests
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("RB: %s", rb)
		return nil, errors.New(resp.Status)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(rb, &data); err != nil {
		return nil, err
	}
	return data, nil
}
