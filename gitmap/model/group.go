// Package model defines the core data structures for gitmap.
package model

// Group represents a named collection of repositories.
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	CreatedAt   string `json:"createdAt"`
}

// GroupRepo links a group to a repository.
type GroupRepo struct {
	GroupID string `json:"groupId"`
	RepoID  string `json:"repoId"`
}
