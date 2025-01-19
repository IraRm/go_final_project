package server

import (
	"encoding/json"
	"go_final_project/task"
)

type ID struct {
	ID string `json:"id"`
}

type Err struct {
	Error string `json:"error"`
}

func (err Err) Bytes() []byte {
	data, _ := json.Marshal(err)
	return data
}

type TasksResp struct {
	Tasks []task.Task `json:"tasks"`
}
