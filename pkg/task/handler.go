package task

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/broswen/taskla/pkg/group"
	"github.com/broswen/taskla/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
)

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (gt *GetTasksResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Get(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		tasks, err := s.GetTasksByUser(r.Context().Value("subject").(string))
		if err != nil {
			oplog.Error().Err(err).Msgf("Couldn't get user tasks: %v", r.Context().Value("subject").(string))
			render.Render(w, r, models.ErrInternalServer(err))
			return
		}
		render.Render(w, r, &GetTasksResponse{tasks})
	}
}

type GetTaskResponse struct {
	Task Task `json:"task"`
}

func (gt *GetTaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetById(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var task Task
		if id := chi.URLParam(r, "id"); id != "" {
			taskId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse task id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			task, err = s.GetTask(taskId)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't get task: %v", taskId)
				render.Render(w, r, models.ErrNotFound(err))
				return
			}
			render.Render(w, r, &GetTaskResponse{task})
			return
		}
		oplog.Error().Msgf("Missing task id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("task id param not found")))
	}
}

type UpdateTaskRequest struct {
	*Task
}

func (ut UpdateTaskRequest) Bind(r *http.Request) error {
	return nil
}

type UpdateTaskResponse struct {
	Task Task `json:"task"`
}

func (ug *UpdateTaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Update(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		data := &UpdateTaskRequest{}
		if err := render.Bind(r, data); err != nil {
			oplog.Error().Err(err).Msgf("Couldn't bind UpdateTaskRequest")
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}

		if id := chi.URLParam(r, "id"); id != "" {
			taskId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse task id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			data.TaskId = taskId
			data.Username = r.Context().Value("subject").(string)
			task, err := s.UpdateTask(*data.Task)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't update task")
				render.Render(w, r, models.ErrInternalServer(err))
				return
			}

			render.Render(w, r, &UpdateTaskResponse{task})
			return
		}

		oplog.Error().Msgf("Missing task id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("task id param not found")))
	}
}

type CreateTaskRequest struct {
	*Task
}

func (ct CreateTaskRequest) Bind(r *http.Request) error {
	return nil
}

type CreateTaskResponse struct {
	Task Task `json:"Task"`
}

func (ct CreateTaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Create(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		data := &CreateTaskRequest{}
		if err := render.Bind(r, data); err != nil {
			oplog.Error().Err(err).Msgf("Couldn't bind CreateTaskRequest")
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}
		data.Username = r.Context().Value("subject").(string)

		task, err := s.CreateTask(*data.Task)
		if err != nil {
			oplog.Error().Err(err).Msgf("Couldn't create task")
			render.Render(w, r, models.ErrInternalServer(err))
			return
		}
		render.Render(w, r, &CreateTaskResponse{task})
	}
}

type DeleteTaskResponse struct {
	Task Task `json:"task"`
}

func (dg *DeleteTaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Delete(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		if id := chi.URLParam(r, "id"); id != "" {
			taskId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse task id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			var task Task
			task.TaskId = taskId
			task.Username = r.Context().Value("subject").(string)
			task, err = s.DeleteTask(task)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't delete task: %v", taskId)
				render.Render(w, r, models.ErrNotFound(err))
				return
			}

			render.Render(w, r, &DeleteTaskResponse{task})
			return
		}

		oplog.Error().Msgf("Missing task id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("task id param not found")))
	}
}

type GetTasksByGroupResponse struct {
	Tasks []Task `json:"tasks"`
}

func (gt *GetTasksByGroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetTasksByGroup(s Service, gs group.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var group group.Group
		if id := chi.URLParam(r, "id"); id != "" {
			groupId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse group id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			group, err = gs.GetGroup(groupId)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't get group id: %v", groupId)
				render.Render(w, r, models.ErrNotFound(err))
				return
			}
			tasks, err := s.GetTasksByGroup(group)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't get tasks for group: %v", groupId)
				render.Render(w, r, models.ErrInternalServer(err))
				return
			}

			render.Render(w, r, &GetTasksByGroupResponse{tasks})

		}
		oplog.Error().Msgf("Missing task id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("group id param not found")))
	}
}
