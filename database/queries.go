package database

import (
	"database/sql"
	"errors"
	"go_final_project/task"
	"strconv"
)

func CreateTableSQL(db *sql.DB) error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT NOT NULL CHECK(length(repeat) <= 128)
	);
	`
	_, err := db.Exec(createTableSQL)
	return err
}

func CreateIndexSQL(db *sql.DB) error {
	const createIndexSQL = `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
	_, err := db.Exec(createIndexSQL)
	return err
}

func SaveTask(db *sql.DB, task task.Task) (string, error) {
	const InsertTask = "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(
		InsertTask,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return "", err
	}
	idx, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(idx)), nil
}

func DeleteTask(db *sql.DB, id string) error {
	const DeleteTask = "DELETE FROM scheduler WHERE id = ?"
	result, err := db.Exec(DeleteTask, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no such rows")
	}
	return nil
}

func GetFutureTasks(db *sql.DB) ([]task.Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE date >= strftime('%Y%m%d', 'now')
        ORDER BY date ASC
		LIMIT 50;
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []task.Task{}
	for rows.Next() {
		var tsk task.Task
		err := rows.Scan(&tsk.Id, &tsk.Date, &tsk.Title, &tsk.Comment, &tsk.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, tsk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(db *sql.DB, id string) (task.Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE id = ?;
    `
	row := db.QueryRow(query, id)
	var tsk task.Task
	err := row.Scan(&tsk.Id, &tsk.Date, &tsk.Title, &tsk.Comment, &tsk.Repeat)
	if err != nil {
		return task.Task{}, err
	}

	return tsk, nil
}

func UpdateTask(db *sql.DB, tsk task.Task) error {
	const Upd = `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?;
	`
	res, err := db.Exec(
		Upd,
		tsk.Date, tsk.Title, tsk.Comment, tsk.Repeat, tsk.Id,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no such rows")
	}
	return nil
}
