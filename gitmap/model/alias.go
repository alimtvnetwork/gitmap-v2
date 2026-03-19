// Package model defines the core data structures for gitmap.
package model

// Alias links a short name to a repository for quick access.
type Alias struct {
	ID        string `json:"id"`
	Alias     string `json:"alias"`
	RepoID    string `json:"repoId"`
	CreatedAt string `json:"createdAt"`
}
