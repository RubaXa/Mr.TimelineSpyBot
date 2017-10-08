package main

import (
	"./flow"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Bot struct {
	UIN    string
	NICK   string
	AIMSID string
}

type Env struct {
	Bot Bot `json:"bot"`
}

func loadEnv() (env *Env, err error) {
	var body []byte

	err = flow.Go(
		func() (err error) {
			body, err = ioutil.ReadFile("./.env.json")
			return
		},

		func() (err error) {
			return json.Unmarshal(body, &env)
		},
	)

	return
}

func main() {
	var env *Env
	var response *http.Response
	var contents []byte

	err := flow.Go(
		func() (err error) {
			env, err = loadEnv()
			return
		},

		func() (err error) {
			endpoint := fmt.Sprintf(
				"https://botapi.icq.net/fetchEvents?aimsid=%s&timeout=%d",
				env.Bot.AIMSID,
				60000,
			)

			fmt.Println("NICK ", env.Bot.NICK)
			fmt.Println(string(endpoint))

			response, err = http.Get(endpoint)
			return
		},

		func() (err error) {
			contents, err = ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			return
		},
	)

	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

	fmt.Println(string(contents))
}
