package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/release"
)

// loadReleasesFromRepo reads .gitmap/release/v*.json files and converts to records.
func loadReleasesFromRepo() []model.ReleaseRecord {
	metas, err := release.ListReleaseMetaFiles()
	if err != nil || len(metas) == 0 {
		return nil
	}

	records := convertMetasToRecords(metas)
	sortRecordsByDate(records)
	markLatestRecord(records)

	return records
}

// convertMetasToRecords converts ReleaseMeta slices to ReleaseRecord slices.
func convertMetasToRecords(metas []release.ReleaseMeta) []model.ReleaseRecord {
	records := make([]model.ReleaseRecord, 0, len(metas))

	for _, m := range metas {
		records = append(records, metaToRecord(m))
	}

	return records
}

// metaToRecord converts a single ReleaseMeta to a ReleaseRecord.
func metaToRecord(m release.ReleaseMeta) model.ReleaseRecord {
	return model.ReleaseRecord{
		Version:      m.Version,
		Tag:          m.Tag,
		Branch:       m.Branch,
		SourceBranch: m.SourceBranch,
		CommitSha:    m.Commit,
		Changelog:    strings.Join(m.Changelog, "\n"),
		Notes:        m.Notes,
		Draft:        m.Draft,
		PreRelease:   m.PreRelease,
		IsLatest:     m.IsLatest,
		Source:       model.SourceRepo,
		CreatedAt:    m.CreatedAt,
	}
}

// sortRecordsByDate sorts records by CreatedAt descending (newest first).
func sortRecordsByDate(records []model.ReleaseRecord) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt > records[j].CreatedAt
	})
}

// markLatestRecord sets IsLatest on the first record matching latest.json.
func markLatestRecord(records []model.ReleaseRecord) {
	latest, err := release.ReadLatest()
	if err != nil {
		return
	}

	for i := range records {
		if records[i].Tag == latest.Tag {
			records[i].IsLatest = true

			return
		}
	}
}

// loadReleasesFromDB opens the DB and fetches all releases.
func loadReleasesFromDB() []model.ReleaseRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrNoDatabase)
		os.Exit(1)
	}
	defer db.Close()

	releases, err := db.ListReleases()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrListReleasesFailed, err)
		os.Exit(1)
	}

	return releases
}
