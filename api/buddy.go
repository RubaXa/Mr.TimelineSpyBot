package api

import (
	botAPI "../bot"
	"../env"
	"net/http"
)

type sListItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func BuddyList(w http.ResponseWriter, r *http.Request) {
	project, err := getProjectByQuery(r.URL.Query())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	botEnv := env.Get("bot")
	bot := botAPI.CreateBot(botEnv["uin"], botEnv["nick"], botEnv["token"])
	buddyList, err := bot.GetBuddyList()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list := make([]sListItem, len(project.Chats))

	for i := range list {
		list[i].Id = project.Chats[i]
		for _, buddy := range buddyList.Norm() {
			if buddy.AimId == list[i].Id {
				list[i].Name = buddy.Friendly
			}
		}
	}

	end(w, list)
}
