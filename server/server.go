package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go_final_project/database"
	"go_final_project/task"
	"io"
	"net/http"
	"time"
)

type Server struct {
	DB *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		DB: db,
	}
}

func (s *Server) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		w.Write(Err{Error: "task id is not set"}.Bytes())
		return
	}

	task, err := database.GetTask(s.DB, taskID)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}
	w.Write(resp)
}

func (s *Server) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tsks, err := database.GetFutureTasks(s.DB)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	res := TasksResp{
		Tasks: tsks,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}
	w.Write(resp)
}

func taskFromReq(bytes []byte) (task.Task, error) {
	var t task.Task
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		return task.Task{}, err
	}

	if t.Title == "" {
		return task.Task{}, errors.New("title is not set")
	}

	y, m, d := time.Now().Date()
	nowDate := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	var tskDate time.Time
	if t.Date == "" || t.Date == "today" {
		tskDate = nowDate
	} else {
		tskDate, err = time.Parse(task.TimeFormat, t.Date)
		if err != nil {
			return task.Task{}, err
		}
		for tskDate.Before(nowDate) {
			tskDate, err = task.NextTime(nowDate, t.Repeat)
			if err != nil {
				return task.Task{}, err
			}
		}
	}
	t.Date = tskDate.Format(task.TimeFormat)

	return t, nil
}

func (s *Server) PostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	t, err := taskFromReq(bytes)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	idx, err := database.SaveTask(s.DB, t)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}
	res := ID{
		ID: idx,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}
	w.Write(resp)
}

func (s *Server) PutTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	t, err := taskFromReq(bytes)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	err = database.UpdateTask(s.DB, t)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Write([]byte("{}"))
}

func (s *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		w.Write(Err{Error: "task id is not set"}.Bytes())
		return
	}

	err := database.DeleteTask(s.DB, taskID)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Write([]byte("{}"))
}

func (s *Server) PostTaskDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		w.Write(Err{Error: "task id is not set"}.Bytes())
		return
	}

	tsk, err := database.GetTask(s.DB, taskID)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	if tsk.Repeat == "" {
		err := database.DeleteTask(s.DB, taskID)
		if err != nil {
			w.Write(Err{Error: err.Error()}.Bytes())
			return
		}
		w.Write([]byte("{}"))
		return
	}

	tsk.Date, err = tsk.NextDate()
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}
	err = database.UpdateTask(s.DB, tsk)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Write([]byte("{}"))
}

func (s *Server) GetNextDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	now := r.URL.Query().Get("now")
	if now == "" {
		w.Write([]byte{})
		return
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		w.Write([]byte{})
		return
	}
	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		w.Write([]byte{})
		return
	}

	tskDate, err := time.Parse(task.TimeFormat, date)
	if err != nil {
		w.Write([]byte{})
		return
	}
	nowDate, err := time.Parse(task.TimeFormat, now)
	if err != nil {
		w.Write([]byte{})
		return
	}

	resT, err := task.NextTime(tskDate, repeat)
	if err != nil {
		w.Write([]byte{})
		return
	}
	for tskDate.Before(nowDate) {
		tskDate, err = task.NextTime(tskDate, repeat)
		if err != nil {
			w.Write([]byte{})
			return
		}
		resT = tskDate
	}
	w.Write([]byte(resT.Format(task.TimeFormat)))
}
