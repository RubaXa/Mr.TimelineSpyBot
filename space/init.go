package space

import (
	"../env"
	"../tnt"
)

var box *tnt.Box
var Projects *sProjects
var Tokens *sTokens
var Records *sRecords

func init() {
	var err error

	boxEnv := env.Get("box")
	spaceEnv := env.Get("space")

	box, err = tnt.InitBox(boxEnv["host"], boxEnv["user"], boxEnv["pass"], 0)

	if err != nil {
		panic(err)
	}

	Projects = &sProjects{}
	Tokens = &sTokens{}
	Records = &sRecords{}

	Projects.Init(box.GetSpace(spaceEnv["projects"]))
	Tokens.Init(box.GetSpace(spaceEnv["tokens"]))
	Records.Init(box.GetSpace(spaceEnv["records"]))
}
