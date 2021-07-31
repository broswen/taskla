package task

import (
	"fmt"

	"github.com/broswen/taskla/pkg/group"
	"github.com/broswen/taskla/pkg/storage"
)

type Service struct {
	r storage.Repository
}

func NewService() (Service, error) {
	repo, err := storage.NewPostgres()
	if err != nil {
		return Service{}, fmt.Errorf("creating repository: %w", err)
	}
	return Service{
		r: repo,
	}, nil
}

type Task struct {
	TaskId      int64  `json:"id"`
	GroupId     int64  `json:"groupId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Username    string `json:"username"`
}

func (s Service) GetTasksByUser(username string) ([]Task, error) {
	rows, err := s.r.Db.Query("SELECT * FROM tasks WHERE username = $1", username)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.TaskId, &task.Username, &task.GroupId, &task.Name, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s Service) GetTask(groupId int64) (Task, error) {
	var task Task
	err := s.r.Db.QueryRow("SELECT * FROM tasks WHERE id = $1", groupId).Scan(&task.TaskId, &task.Username, &task.GroupId, &task.Name, &task.Description, &task.Status)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s Service) CreateTask(task Task) (Task, error) {
	var newTask Task
	err := s.r.Db.QueryRow("INSERT INTO tasks (group_id, username, name, description, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, group_id, username, name, description, status", task.GroupId, task.Username, task.Name, task.Description, task.Status).Scan(&newTask.TaskId, &newTask.GroupId, &newTask.Username, &newTask.Name, &newTask.Description, &newTask.Status)
	if err != nil {
		return Task{}, err
	}

	return newTask, nil
}

func (s Service) UpdateTask(task Task) (Task, error) {
	var newTask Task
	err := s.r.Db.QueryRow("UPDATE tasks SET name = $1, description = $2, status = $3 WHERE id = $4 AND username = $5 RETURNING id, group_id, username, name, description, status", task.Name, task.Description, task.Status, task.TaskId, task.Username).Scan(&task.TaskId, &newTask.GroupId, &newTask.Username, &newTask.Name, &newTask.Description, &task.Status)
	if err != nil {
		return Task{}, err
	}

	return newTask, nil
}

func (s Service) DeleteTask(task Task) (Task, error) {
	var deletedTask Task
	err := s.r.Db.QueryRow("DELETE FROM tasks WHERE id = $1 AND username = $2 RETURNING id, group_id, username, name, description, status", task.TaskId, task.Username).Scan(&deletedTask.TaskId, &deletedTask.GroupId, &deletedTask.Username, &deletedTask.Name, &deletedTask.Description, &deletedTask.Status)
	if err != nil {
		return Task{}, err
	}

	return deletedTask, nil
}

func (s Service) GetTasksByGroup(group group.Group) ([]Task, error) {
	rows, err := s.r.Db.Query("SELECT * FROM tasks WHERE username = $1 AND group_id = $2", group.Username, group.GroupId)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.TaskId, &task.Username, &task.GroupId, &task.Name, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
