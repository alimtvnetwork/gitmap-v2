package model

// TempRelease represents a temporary release branch record.
type TempRelease struct {
	ID             string
	Branch         string
	VersionPrefix  string
	SequenceNumber int
	Commit         string
	CommitMessage  string
	CreatedAt      string
}
