package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Task struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Member    string `json:"member"`
	Deadline  string `json:"deadline"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./task.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Migration() error {
	sql := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		member TEXT,
		deadline TEXT,
		status TEXT,
		createdAt DATETIME,
		updatedAt DATETIME
	);`

	_, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// create task
func CreateTask(task Task) (int64, error) {
	result, err := DB.Exec("INSERT INTO tasks (name, member, deadline, status, createdAt, updatedAt) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))", task.Name, task.Member, task.Deadline, task.Status)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetTasks(page int, limit int) ([]Task, error) {
	rows, err := DB.Query("SELECT * FROM tasks LIMIT ? OFFSET ?", limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Member, &task.Deadline, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskById(id int) (Task, error) {
	stmt, err := DB.Prepare("SELECT id, name, member, deadline, status, createdAt, updatedAt from tasks WHERE id = ?")
	if err != nil {
		return Task{}, err
	}

	task := Task{}

	sqlErr := stmt.QueryRow(id).Scan(&task.Id, &task.Name, &task.Member, &task.Deadline, &task.Status, &task.CreatedAt, &task.UpdatedAt)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Task{}, nil
		}
		return Task{}, sqlErr
	}
	return task, nil
}

func UpdateTask(task Task) (int64, error) {
	result, err := DB.Exec("UPDATE tasks SET name = ?, member = ?, deadline = ?, status = ?, updatedAt = datetime('now') WHERE id = ?", task.Name, task.Member, task.Deadline, task.Status, task.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteTask(id int) (int64, error) {
	result, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func MarkDoneTask(id int) (int64, error) {
	result, err := DB.Exec("UPDATE tasks SET status = '1' WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
