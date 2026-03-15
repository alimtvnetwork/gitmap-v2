// Package cmd — scanprojects.go handles project detection during scan.
package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/detector"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

// detectAllProjects runs project detection across all scanned repos.
func detectAllProjects(records []model.ScanRecord) []detector.DetectionResult {
	var all []detector.DetectionResult
	repoCount := 0
	for _, rec := range records {
		results := detector.DetectProjects(rec.AbsolutePath, rec.ID, rec.RepoName)
		if len(results) > 0 {
			repoCount++
			all = append(all, results...)
		}
	}
	fmt.Printf(constants.MsgProjectDetectDone, len(all), repoCount)

	return all
}

// upsertProjectsToDB persists detected projects and metadata to SQLite.
func upsertProjectsToDB(results []detector.DetectionResult, records []model.ScanRecord, outputDir string) {
	if len(results) == 0 {
		return
	}
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

		return
	}
	upsertProjectRecords(db, results, records)
}

// upsertProjectRecords inserts projects, metadata, and cleans stale records.
func upsertProjectRecords(db *store.DB, results []detector.DetectionResult, records []model.ScanRecord) {
	count := 0
	repoIDs := collectRepoIDs(results)
	for i := range results {
		r := &results[i]
		err := db.UpsertDetectedProject(r.Project)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

			continue
		}
		if err := resolveDetectedProjectID(db, r); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

			continue
		}
		count++
		upsertProjectMetadata(db, *r)
	}
	cleanStaleProjects(db, repoIDs, results)
	fmt.Printf(constants.MsgProjectUpsertDone, count)
}

// resolveDetectedProjectID syncs the project ID with the persisted DB row.
func resolveDetectedProjectID(db *store.DB, r *detector.DetectionResult) error {
	id, err := db.SelectDetectedProjectID(
		r.Project.RepoID,
		r.Project.ProjectTypeID,
		r.Project.RelativePath,
	)
	if err != nil {
		return err
	}
	r.Project.ID = id

	return nil
}

// upsertProjectMetadata persists Go or C# metadata for a detection result.
func upsertProjectMetadata(db *store.DB, r detector.DetectionResult) {
	if r.GoMeta != nil {
		upsertGoProjectMeta(db, r)
	}
	if r.CSharp != nil {
		upsertCSharpProjectMeta(db, r)
	}
}

// upsertGoProjectMeta persists Go metadata and runnables.
func upsertGoProjectMeta(db *store.DB, r detector.DetectionResult) {
	r.GoMeta.DetectedProjectID = r.Project.ID
	if err := db.UpsertGoMetadata(*r.GoMeta); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoMetadataUpsert, err)

		return
	}
	saved, err := db.SelectGoMetadata(r.Project.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoMetadataUpsert, err)

		return
	}
	r.GoMeta.ID = saved.ID
	runnableIDs := upsertGoRunnables(db, r.GoMeta)
	_ = db.DeleteStaleGoRunnables(r.GoMeta.ID, runnableIDs)
}

// upsertGoRunnables persists all runnable files and returns their IDs.
func upsertGoRunnables(db *store.DB, meta *model.GoProjectMetadata) []string {
	var ids []string
	for _, run := range meta.Runnables {
		run.GoMetadataID = meta.ID
		if err := db.UpsertGoRunnable(run); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrGoRunnableUpsert, err)

			continue
		}
		ids = append(ids, run.ID)
	}

	return ids
}

// upsertCSharpProjectMeta persists C# metadata, project files, and key files.
func upsertCSharpProjectMeta(db *store.DB, r detector.DetectionResult) {
	r.CSharp.DetectedProjectID = r.Project.ID
	if err := db.UpsertCSharpMetadata(*r.CSharp); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCSharpMetaUpsert, err)

		return
	}
	saved, err := db.SelectCSharpMetadata(r.Project.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCSharpMetaUpsert, err)

		return
	}
	r.CSharp.ID = saved.ID
	fileIDs := upsertCSharpFiles(db, r.CSharp)
	keyIDs := upsertCSharpKeyFiles(db, r.CSharp)
	_ = db.DeleteStaleCSharpFiles(r.CSharp.ID, fileIDs)
	_ = db.DeleteStaleCSharpKeyFiles(r.CSharp.ID, keyIDs)
}

// upsertCSharpFiles persists C# project files and returns their IDs.
func upsertCSharpFiles(db *store.DB, meta *model.CSharpProjectMetadata) []string {
	var ids []string
	for _, f := range meta.ProjectFiles {
		f.CSharpMetadataID = meta.ID
		if err := db.UpsertCSharpProjectFile(f); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCSharpFileUpsert, err)

			continue
		}
		ids = append(ids, f.ID)
	}

	return ids
}

// upsertCSharpKeyFiles persists C# key files and returns their IDs.
func upsertCSharpKeyFiles(db *store.DB, meta *model.CSharpProjectMetadata) []string {
	var ids []string
	for _, f := range meta.KeyFiles {
		f.CSharpMetadataID = meta.ID
		if err := db.UpsertCSharpKeyFile(f); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCSharpKeyUpsert, err)

			continue
		}
		ids = append(ids, f.ID)
	}

	return ids
}

// collectRepoIDs extracts unique repo IDs from detection results.
func collectRepoIDs(results []detector.DetectionResult) map[string]bool {
	ids := make(map[string]bool)
	for _, r := range results {
		ids[r.Project.RepoID] = true
	}

	return ids
}

// cleanStaleProjects removes projects no longer detected for each repo.
func cleanStaleProjects(db *store.DB, repoIDs map[string]bool, results []detector.DetectionResult) {
	for repoID := range repoIDs {
		keepIDs := collectKeepIDs(repoID, results)
		cleaned, err := db.DeleteStaleProjects(repoID, keepIDs)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectCleanup, repoID, err)

			continue
		}
		if cleaned > 0 {
			fmt.Printf(constants.MsgProjectCleanedStale, cleaned)
		}
	}
}

// collectKeepIDs collects project IDs to keep for a given repo.
func collectKeepIDs(repoID string, results []detector.DetectionResult) []string {
	var ids []string
	for _, r := range results {
		if r.Project.RepoID == repoID {
			ids = append(ids, r.Project.ID)
		}
	}

	return ids
}
