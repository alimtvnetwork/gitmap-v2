// Package detector walks repository trees and classifies project types.
package detector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// DetectProjects scans a repo directory for all supported project types.
func DetectProjects(repoPath, repoID, repoName string) []DetectionResult {
	var results []DetectionResult
	slnDirs := map[string]bool{}

	collectSlnDirs(repoPath, slnDirs)
	walkRepo(repoPath, repoID, repoName, slnDirs, &results)

	return results
}

// DetectionResult holds a detected project and optional metadata.
type DetectionResult struct {
	Project model.DetectedProject
	GoMeta  *model.GoProjectMetadata
	CSharp  *model.CSharpProjectMetadata
}

// walkRepo walks the directory tree and detects projects.
func walkRepo(repoPath, repoID, repoName string, slnDirs map[string]bool, results *[]DetectionResult) {
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && shouldExcludeDir(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		detectFile(path, repoPath, repoID, repoName, slnDirs, results)

		return nil
	})
}

// detectFile checks a single file against all detection rules.
func detectFile(path, repoPath, repoID, repoName string, slnDirs map[string]bool, results *[]DetectionResult) {
	name := filepath.Base(path)
	dir := filepath.Dir(path)

	if name == constants.IndicatorGoMod {
		detectGo(dir, repoPath, repoID, repoName, results)
	}
	if name == constants.IndicatorPackageJSON {
		detectNodeOrReact(dir, path, repoPath, repoID, repoName, results)
	}
	detectCpp(name, dir, repoPath, repoID, repoName, results)
	detectCSharpFile(name, dir, repoPath, repoID, repoName, slnDirs, results)
}

// collectSlnDirs pre-scans for .sln files to enforce precedence.
func collectSlnDirs(repoPath string, slnDirs map[string]bool) {
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && shouldExcludeDir(info.Name()) {
			return filepath.SkipDir
		}
		if strings.HasSuffix(info.Name(), constants.ExtSln) {
			slnDirs[filepath.Dir(path)] = true
		}

		return nil
	})
}

// shouldExcludeDir checks if a directory name should be skipped.
func shouldExcludeDir(name string) bool {
	if strings.HasPrefix(name, constants.CMakeBuildPfx) {
		return true
	}
	for _, excluded := range constants.ProjectExcludeDirs {
		if name == excluded {
			return true
		}
	}

	return false
}

// buildRelativePath returns the relative path from repo root.
func buildRelativePath(dir, repoPath string) string {
	rel, err := filepath.Rel(repoPath, dir)
	if err != nil {
		return "."
	}

	return rel
}

// stableID returns a deterministic UUID derived from stable identity parts.
func stableID(parts ...string) string {
	seed := strings.Join(parts, "|")

	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()
}

// normalizePath creates a cross-platform canonical path for ID seeds.
func normalizePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

// projectID creates a stable ID for a detected project row.
func projectID(repoID, projectTypeID, relativePath string) string {
	return stableID("project", repoID, projectTypeID, normalizePath(relativePath))
}

// goMetadataID creates a stable ID for Go metadata linked to a project.
func goMetadataID(detectedProjectID string) string {
	return stableID("go_metadata", detectedProjectID)
}

// goRunnableID creates a stable ID for a Go runnable linked to metadata.
func goRunnableID(metadataID, relativePath string) string {
	return stableID("go_runnable", metadataID, normalizePath(relativePath))
}

// csharpMetadataID creates a stable ID for C# metadata linked to a project.
func csharpMetadataID(detectedProjectID string) string {
	return stableID("csharp_metadata", detectedProjectID)
}

// csharpProjectFileID creates a stable ID for a C# project file row.
func csharpProjectFileID(metadataID, relativePath string) string {
	return stableID("csharp_file", metadataID, normalizePath(relativePath))
}

// csharpKeyFileID creates a stable ID for a C# key file row.
func csharpKeyFileID(metadataID, relativePath string) string {
	return stableID("csharp_key_file", metadataID, normalizePath(relativePath))
}

