package main

import (
	"encoding/json"
	//"flag"
	"fmt"
	//"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/garyburd/go-oauth/oauth"
)

type client struct {
	client oauth.Client
	token  oauth.Credentials
}

func (c *client) get(urlStr string, params url.Values, v interface{}) error {
	resp, err := c.client.Get(nil, &c.token, urlStr, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("yelp status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func readCredentials(c *client) error {

	c.client.Credentials.Token = os.Getenv("YELP_CONS_ACCS_KEY")
	c.client.Credentials.Secret = os.Getenv("YELP_CONS_SCRT")
	c.token.Token = os.Getenv("YELP_ACCS_TOKEN")
	c.token.Secret = os.Getenv("YELP_ACCS_TOKEN_SCRT")

	return nil
}

func main() {
	log.Println("Testing go program")
	var c client
	if err := readCredentials(&c); err != nil {
		log.Fatal(err)
	}

	var data struct {
		Businesses []struct {
			Name     string
			Location struct {
				DisplayAddress []string `json:"display_address"`
			}
		}
	}
	form := url.Values{"term": {"food"}, "location": {"San Francisco"}}
	if err := c.get("http://api.yelp.com/v2/search", form, &data); err != nil {
		log.Fatal(err)
	}

	for _, b := range data.Businesses {
		addr := ""
		if len(b.Location.DisplayAddress) > 0 {
			addr = b.Location.DisplayAddress[0]
		}
		log.Println(b.Name, addr)
	}
}
