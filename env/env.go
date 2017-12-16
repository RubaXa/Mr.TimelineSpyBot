package env

import (
	"encoding/json"
	"io/ioutil"
)

var env map[string]map[string]string

func Get(name string) map[string]string {
	return env[name]
}

func init() {
	body, err := ioutil.ReadFile("./.env.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &env)

	if err != nil {
		panic(err)
	}
}
