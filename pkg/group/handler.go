package group

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/broswen/taskla/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
)

type GetGroupsResponse struct {
	Groups []Group `json:"groups"`
}

func (gg *GetGroupsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Get(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		groups, err := s.GetGroupsByUser(r.Context().Value("subject").(string))
		if err != nil {
			oplog.Error().Err(err).Msgf("Couldn't get groups for user: %v", r.Context().Value("subject").(string))
			render.Render(w, r, models.ErrInternalServer(err))
			return
		}
		render.Render(w, r, &GetGroupsResponse{groups})
	}
}

type GetGroupResponse struct {
	Group Group `json:"group"`
}

func (gg *GetGroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetById(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var group Group
		if id := chi.URLParam(r, "id"); id != "" {
			groupId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse group id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			group, err = s.GetGroup(groupId)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't get group: %v", groupId)
				render.Render(w, r, models.ErrNotFound(err))
				return
			}

			render.Render(w, r, &GetGroupResponse{group})
			return
		}

		oplog.Error().Msgf("Missing group id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("group id param not found")))

	}
}

type UpdateGroupRequest struct {
	*Group
}

func (ug UpdateGroupRequest) Bind(r *http.Request) error {
	return nil
}

type UpdateGroupResponse struct {
	Group Group `json:"group"`
}

func (ug *UpdateGroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Update(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		data := &CreateGroupRequest{}
		if err := render.Bind(r, data); err != nil {
			oplog.Error().Err(err).Msgf("Couldn't bind UpdateGroupRequest")
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}

		if id := chi.URLParam(r, "id"); id != "" {
			groupId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse group id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			data.GroupId = groupId
			data.Username = r.Context().Value("subject").(string)
			group, err := s.UpdateGroup(*data.Group)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't update group: %v", groupId)
				render.Render(w, r, models.ErrInternalServer(err))
				return
			}

			render.Render(w, r, &UpdateGroupResponse{group})
			return
		}

		oplog.Error().Msgf("Missing group id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("group id param not found")))
	}
}

type CreateGroupRequest struct {
	*Group
}

func (cg CreateGroupRequest) Bind(r *http.Request) error {
	return nil
}

type CreateGroupResponse struct {
	Group Group `json:"group"`
}

func (cg *CreateGroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Create(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		data := &CreateGroupRequest{}
		if err := render.Bind(r, data); err != nil {
			oplog.Error().Err(err).Msgf("Couldn't bind CreateGroupRequest")
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}

		data.Username = r.Context().Value("subject").(string)

		group, err := s.CreateGroup(*data.Group)
		if err != nil {
			oplog.Error().Err(err).Msgf("Couldn't create group")
			render.Render(w, r, models.ErrInternalServer(err))
			return
		}

		render.Render(w, r, &CreateGroupResponse{group})
	}
}

type DeleteGroupResponse struct {
	Group Group `json:"group"`
}

func (dg *DeleteGroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Delete(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		if id := chi.URLParam(r, "id"); id != "" {
			groupId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't parse group id")
				render.Render(w, r, models.ErrInvalidRequest(err))
				return
			}
			var group Group
			group.GroupId = groupId
			group.Username = r.Context().Value("subject").(string)
			group, err = s.DeleteGroup(group)
			if err != nil {
				oplog.Error().Err(err).Msgf("Couldn't delete group: %v", groupId)
				render.Render(w, r, models.ErrNotFound(err))
				return
			}

			render.Render(w, r, &DeleteGroupResponse{group})
			return
		}

		oplog.Error().Msgf("Missing group id")
		render.Render(w, r, models.ErrInvalidRequest(errors.New("group id param not found")))
	}
}
