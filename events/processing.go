package events

import (
	"../bot"
	"../space"
	"fmt"
	"strings"
)

var projectsIndexById = make(map[uint64]*space.ProjectsEntry)
var projectsIndexByChat = make(map[string][]*space.ProjectsEntry)

func getProject(id uint64) *space.ProjectsEntry {
	project, ok := projectsIndexById[id]

	if !ok {
		project = space.Projects.Get(id)

		if !project.HasId() {
			return nil
		}

		projectsIndexById[id] = project
	}

	return project
}

func getProjectsByChat(id string) []*space.ProjectsEntry {
	_, ok := projectsIndexByChat[id]

	if !ok {
		allProjects, err := space.Projects.GetAll()

		if err == nil {
			for _, project := range allProjects {
				for _, id := range project.Chats {
					exists := false
					list, ok := projectsIndexByChat[id]

					if !ok {
						list = make([]*space.ProjectsEntry, 0, 1)
					}

					for _, p := range list {
						if p.Id == project.Id {
							exists = true
							break
						}
					}

					if !exists {
						list = append(list, &project)
						projectsIndexByChat[id] = list
					}
				}
			}
		} else {
			fmt.Println("getProjectsByChat fail:", err)
		}
	}

	return projectsIndexByChat[id]
}

func Processing(seqNum uint, events []bot.Event) uint {
	fmt.Println("Events processing:", len(events))

	for _, rawEvt := range events {
		seqNum = rawEvt.SeqNum + 1

		if rawEvt.Type == "im" {
			record := toRecordsEntry(rawEvt)

			if strings.Contains(record.Body, "timeline.bind") {
				binding(record, true)
			} else if strings.Contains(record.Body, "timeline.unbind") {
				binding(record, false)
			} else {
				trySaveRecord(record)
			}
		}
	}

	return seqNum
}

func toRecordsEntry(event bot.Event) *space.RecordsEntry {
	data := event.GetIMData()

	return &space.RecordsEntry{
		MsgId:  data.MsgID,
		SeqNum: event.SeqNum,
		TS:     data.Timestamp,
		Source: space.RecordsEntrySource{
			Id:   data.Source.AimID,
			Name: data.Source.Friendly,
		},
		Author: space.RecordsEntryAuthor{
			Login: data.MChatAttrs.Sender,
			Name:  data.MChatAttrs.SenderName,
		},
		Body: data.Message,
	}
}

func binding(record *space.RecordsEntry, bind bool) {
	detail := strings.Split(record.Body, " ")
	token, err := space.Tokens.GetByValue(detail[1])

	if err != nil {
		fmt.Println("token.err", err)
		return
	}

	fmt.Printf("token: %#v\n", token)
	project := getProject(token.ProjectId)

	if project == nil {
		fmt.Println("project not found")
		return
	}

	fmt.Printf("project: %#v\n", project)

	if bind {
		if project.HasChat(record.Source.Id) {
			fmt.Printf("project.chats: %#v\n", project.Chats)
			return
		}

		project.AddChat(record.Source.Id)
	} else {
		project.RemoveChat(record.Source.Id)
	}

	err = space.Projects.Save(project)

	if err != nil {
		fmt.Println("project save chat:", err)
		return
	}

	fmt.Printf("project.chats.added: %#v\n", project.Chats)
	delete(projectsIndexByChat, record.Source.Id)
	space.Projects.Delete(token)
}

func trySaveRecord(record *space.RecordsEntry) {
	projects := getProjectsByChat(record.Source.Id)

	fmt.Printf("Save message: %#v\n", record)

	for _, project := range projects {
		copy := *record
		copy.ProjectId = project.Id
		err := space.Records.Save(&copy)

		if err != nil {
			fmt.Printf("Failed save to %#v: %#v\n", project, err)
		}
	}
}
