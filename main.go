package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var config Config
	{
		confFile := flag.String("conf", "config.json", "Config File for OAuth")
		flag.Parse()
		data, err_file := ioutil.ReadFile(*confFile)
		check(err_file)
		err_json := json.Unmarshal(data, &config)
		check(err_json)
	}
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	v := url.Values{}
	v.Set("follow", "**********") //set his ID
	stream := api.PublicStreamFilter(v)
	for {
		select {
		case stream := <-stream.C:
			switch status := stream.(type) {
			case anaconda.Tweet:
				fmt.Printf("%s: %s\n", status.User.ScreenName, status.Text)
				encoded := base64.StdEncoding.EncodeToString([]byte(status.Text))
				if len(encoded) <= 140 {
					botTweet, err := api.PostTweet(encoded, nil)
					check(err)
					fmt.Printf("%s : tweet\n", botTweet.CreatedAt)
				}
			}
		}
	}
}
