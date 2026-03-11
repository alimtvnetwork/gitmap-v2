// Package model — gometadata.go defines Go-specific metadata structs.
package model

// GoProjectMetadata holds Go-specific metadata for a detected project.
type GoProjectMetadata struct {
	ID                string           `json:"id"`
	DetectedProjectID string           `json:"detectedProjectId"`
	GoModPath         string           `json:"goModPath"`
	GoSumPath         string           `json:"goSumPath"`
	ModuleName        string           `json:"moduleName"`
	GoVersion         string           `json:"goVersion"`
	Runnables         []GoRunnableFile `json:"runnables"`
}

// GoRunnableFile represents a main.go entry point inside a Go project.
type GoRunnableFile struct {
	ID           string `json:"id"`
	GoMetadataID string `json:"goMetadataId"`
	RunnableName string `json:"runnableName"`
	FilePath     string `json:"filePath"`
	RelativePath string `json:"relativePath"`
}

// GoProjectRecord combines a DetectedProject with its Go metadata for JSON output.
type GoProjectRecord struct {
	DetectedProject
	GoMetadata *GoProjectMetadata `json:"goMetadata,omitempty"`
}
