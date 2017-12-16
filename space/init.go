package space

import (
	"../env"
	"../tnt"
)

var box *tnt.Box
var Projects *sProjects
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
	Records = &sRecords{}

	Projects.Init(box.GetSpace(spaceEnv["projects"]))
	Records.Init(box.GetSpace(spaceEnv["records"]))
}
