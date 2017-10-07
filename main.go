package main

import (
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

type Response struct {
}

type Events struct {
}

func exit(err error) {
	panic(fmt.Sprintf("%#v", err))
}

func loadEnv() (env *Env, err error) {
	env = new(Env)
	body, err := ioutil.ReadFile("./.env.json")

	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(body), &env)

	return
}

func main() {
	env, err := loadEnv()

	if err != nil {
		exit(err)
	}

	fmt.Println("NICK ", env.Bot.NICK)

	endpoint := fmt.Sprintf(
		"https://botapi.icq.net/im/fetchEvents?aimsid=%s&seqNum=951681&timeout=%d",
		env.Bot.AIMSID,
		60000,
	)
	fmt.Print(string(endpoint))
	response, err := http.Get(endpoint)

	if err != nil {
		exit(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		exit(err)
	}

	fmt.Print(string(contents))
}
