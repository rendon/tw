#tw
A new and simple Twitter API client.

# Usage

    package main

    import (
        "fmt"
        "log"
        "os"

        "github.com/rendon/tw"
    )

    func main() {

        // Create new client
        tc := tw.NewClient()

        // Set keys
        ck := os.Getenv("TWITTER_CONSUMER_KEY")
        cs := os.Getenv("TWITTER_CONSUMER_SECRET")
        if err := tc.SetKeys(ck, cs); err != nil {
            log.Fatalf("Failed to setup client")
        }

        // GET users/show
        data, err := tc.UsersShow("twitterdev")
        if err != nil {
            log.Fatalf("Failed to obtain user: %s", err)
        }

        fmt.Printf("User ID: %s\n", data["id_str"].(string))
        fmt.Printf("User name: %s\n", data["screen_name"].(string))
    }

The client was designed to use multiple sets of keys, if you reach the rate limit, just call `SetKeys()` with a new set of keys and keep going.

# Detect "Too Many Requests"
If your client reaches the rate limit, `err` will be equal to `ErrTooManyRequests`, you can handle this case as follows:

    // GET users/show
    data, err := tc.UsersShow("twitterdev")
    if err == tw.ErrTooManyRequests {
        // Set a new set of keys and try again
    }

