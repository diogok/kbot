package kbot

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var telegram_api string

func getMe() (User, error) {
	var user UserResponse

	resp, err := http.Get(fmt.Sprintf("%s/getMe", telegram_api))

	if err != nil {
		return user.Result, err
	} else {
		contents, _ := ioutil.ReadAll(resp.Body)
		jerr := json.Unmarshal(contents, &user)

		if jerr != nil {
			return user.Result, jerr
		}

		if user.Ok != true {
			return user.Result, errors.New(user.Description)
		}

		resp.Body.Close()

		return user.Result, nil
	}
}

func getUpdates(since int) ([]Update, error) {
	var updates UpdateResponse

	resp, err := http.Get(fmt.Sprintf("%s/getUpdates?timeout=1&offset=%d", telegram_api, (since + 1)))

	if err != nil {
		return updates.Result, err
	} else {
		contents, _ := ioutil.ReadAll(resp.Body)
		jerr := json.Unmarshal(contents, &updates)

		if jerr != nil {
			return updates.Result, jerr
		}

		if updates.Ok != true {
			return updates.Result, errors.New(updates.Description)
		}

		resp.Body.Close()

		return updates.Result, nil
	}
}

func receiver(in chan<- Update) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var updt Update
		err := decoder.Decode(&updt)
		if err != nil {
			log.Println("Error parsing update", err)
		} else {
			in <- updt
		}
	}
}

func sendMessage(httpClient *http.Client, msg OutMessage) {
	content, _ := json.Marshal(msg)
	resp, err := httpClient.Post(fmt.Sprintf("%s/sendMessage", telegram_api), "application/json", bytes.NewBufferString(string(content)))
	if err != nil {
		log.Println("Error sending", msg, err)
	} else {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func sendInlineMessage(httpClient *http.Client, msg OutQuery) {
	content, _ := json.Marshal(msg)
	resp, err := httpClient.Post(fmt.Sprintf("%s/answerInlineQuery", telegram_api), "application/json", bytes.NewBufferString(string(content)))
	if err != nil {
		log.Println("Error sending", msg, err)
	} else {
		content, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
}

func Start(bot Bot) (chan bool, chan bool) {
	telegram_api = fmt.Sprintf("https://api.telegram.org/bot%s", bot.Token)

	user, err := getMe()
	if err != nil {
		panic(err)
	}
	log.Println("Your bot: ", user)

	in := make(chan Update, 256)
	done := make(chan bool, 1)
	stop := make(chan bool, 1)
	outMsg := make(chan OutMessage, 256)
	outQuery := make(chan OutQuery, 256)

	active := true

	for i := 0; i < 8; i = i + 1 {
		tr := &http.Transport{}
		httpClient := &http.Client{Transport: tr}

		go func(in <-chan Update) {
			for updt := range in {
				bot.Handler(updt, outMsg, outQuery)
			}
		}(in)

		go func(out <-chan OutMessage) {
			for msg := range out {
				sendMessage(httpClient, msg)
			}
		}(outMsg)

		go func(out <-chan OutQuery) {
			for msg := range out {
				sendInlineMessage(httpClient, msg)
			}
		}(outQuery)
	}

	if bot.Host == "" {
		go func(in chan<- Update) {
			since := 0
			for active {
				updates, err := getUpdates(since)
				if err != nil {
					log.Println(err)
				} else {
					for _, update := range updates {
						in <- update
						since = update.Update_id
					}
				}
			}
			done <- true
		}(in)
	} else {
		host := fmt.Sprintf("https://%s/bot%s", bot.Host, bot.Token)
		log.Printf("Setting webhook %s\n", host)

		lejson := fmt.Sprintf("{\"url\":\"%s\"}", host)
		resp, err := http.Post(fmt.Sprintf("%s/setWebhook", telegram_api), "application/json", bytes.NewBufferString(lejson))
		if err != nil {
			panic(err)
		}
		contents, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(contents))
		log.Println("Webhook set")

		go func(in chan<- Update) {
			http.HandleFunc(fmt.Sprintf("/bot%s", bot.Token), receiver(in))
			log.Fatal(http.ListenAndServe(":8080", nil))
		}(in)
	}

	go func(stop chan bool) {
		<-stop
		active = false
	}(stop)

	return stop, done
}

func RandId() string {
	n := 8
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
