package group

import (
	"github.com/broswen/taskla/pkg/storage"
)

type Service struct {
	r storage.Repository
}

func NewService() (Service, error) {
	repo := storage.New()
	return Service{
		r: repo,
	}, nil
}

type Group struct {
	GroupId     int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Username    string `json:"username"`
}

func (s Service) GetGroupsByUser(username string, limit, offset int) ([]Group, error) {
	rows, err := s.r.DB().Query("SELECT * FROM groups WHERE username = $1 ORDER BY id LIMIT $2 OFFSET $3", username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]Group, 0)
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.GroupId, &group.Username, &group.Name, &group.Description)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (s Service) GetGroup(groupId int64, username string) (Group, error) {
	var group Group
	err := s.r.DB().QueryRow("SELECT * FROM groups WHERE id = $1 AND username = $2", groupId, username).Scan(&group.GroupId, &group.Username, &group.Name, &group.Description)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (s Service) CreateGroup(group Group) (Group, error) {
	var newGroup Group
	err := s.r.DB().QueryRow("INSERT INTO groups (username, name, description) VALUES ($1, $2, $3) RETURNING id, username, name, description", group.Username, group.Name, group.Description).Scan(&newGroup.GroupId, &newGroup.Username, &newGroup.Name, &newGroup.Description)
	if err != nil {
		return Group{}, err
	}

	return newGroup, nil
}

func (s Service) UpdateGroup(group Group) (Group, error) {
	var newGroup Group
	err := s.r.DB().QueryRow("UPDATE groups SET name = $1, description = $2 WHERE id = $3 AND username = $4 RETURNING id, username, name, description", group.Name, group.Description, group.GroupId, group.Username).Scan(&newGroup.GroupId, &newGroup.Username, &newGroup.Name, &newGroup.Description)
	if err != nil {
		return Group{}, err
	}

	return newGroup, nil
}

func (s Service) DeleteGroup(group Group) (Group, error) {
	var deletedGroup Group
	err := s.r.DB().QueryRow("DELETE FROM groups WHERE id = $1 AND username = $2 RETURNING id, username, name, description", group.GroupId, group.Username).Scan(&deletedGroup.GroupId, &deletedGroup.Username, &deletedGroup.Name, &deletedGroup.Description)
	if err != nil {
		return Group{}, err
	}

	return deletedGroup, nil
}
