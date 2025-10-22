package action

import (
	"encoding/json"
	"os"
	"robotech/logger"
)

func init() {
	load()
}

func Add(a Action) int {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	a.ID = nextID
	nextID++
	actions[a.ID] = &a
	save()
	logger.Info("Added action:", a.Description)
	return a.ID
}

func List() []*Action {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	list := []*Action{}
	for _, a := range actions {
		list = append(list, a)
	}
	return list
}

func Bind(sessionIDs []string, actionIDs []int) {
	logger.Info("Bind sessions:", sessionIDs, "with actions:", actionIDs)
}

func save() {
	data, _ := json.MarshalIndent(actions, "", "  ")
	os.WriteFile(configFile, data, 0644)
}

func load() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &actions)
	for id := range actions {
		if id >= nextID {
			nextID = id + 1
		}
	}
}
