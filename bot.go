package main

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	token, err := getenv("TOKEN_WIKIHOW_BOT")
	if err != nil {
		log.Fatal(err)
	}
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				text := ev.Text
				matched, query := concat(text)
				println(ev.Channel)
				if ev.User != info.User.ID && matched && ev.Channel == "C0G07RYJD" {
					url, err := get_wiki_url(query)
					if err != nil {
						url = "Data not found sozzzzz"
					}
					rtm.SendMessage(rtm.NewOutgoingMessage(url, ev.Channel))
					text = ""
					matched = false
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//later for more
			}
		}
	}
}

func concat(text string) (bool, string) {
	if len(text) < 6 {
		return false, ""
	}
	matched := strings.HasPrefix(text, "!howto ")
	if matched {
		query := strings.TrimPrefix(text, "!howto ")
		return true, query
	}
	return false, ""
}

func getenv(env_var string) (string, error) {
	token := os.Getenv(env_var)
	if token == "" {
		return "", errors.New("Token is empty")
	}
	return token, nil
}

func get_wiki_url(query string) (string, error) {

	query = replace_whitespace(query)
	url := "https://fr.wikihow.com/wikiHowTo?search=" + query

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(http.StatusText(resp.StatusCode))
	}
	body, _ := ioutil.ReadAll(resp.Body)

	ret := pars_body(body)
	return ret, nil
}

func replace_whitespace(str string) string {
	spaces := map[int]string{
		0: " ",
		1: "\t",
		2: "\n",
		3: "\v",
		4: "\f",
		5: "\r",
		6: string(0x85),
		7: string(0xA0),
	}
	for _, space := range spaces {
		str = strings.ReplaceAll(str, space, "+")
	}
	return str
}

func pars_body(body []byte) string {
	//GARBAGE COLLECTOR FUNC
	str := string(body)
	i := strings.Index(str, "searchresults_list")
	str = str[i:]
	y := strings.Index(str, "href")
	str2 := str[y:]
	j := strings.Index(str2, "\">")
	str = str[0 : j+y]
	i = strings.Index(str, "https")
	str = str[i:]
	j = strings.Index(str, "\" ")
	str = str[0:j]
	return str
}
