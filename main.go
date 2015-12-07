package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	authURL = "https://api.twitter.com/oauth2/token"
)

var (
	consumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
)

func getAccessToken() (string, error) {
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

func main() {
	at, err := getAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Access token: %s\n", at)
}
