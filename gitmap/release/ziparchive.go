// Package release — ziparchive.go creates ZIP archives from zip groups
// with maximum compression (Deflate level 9) for release assets.
package release

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

// ZipGroupArchive holds the result of archiving a zip group.
type ZipGroupArchive struct {
	GroupName   string
	ArchivePath string
	ItemCount   int
}

// BuildZipGroupArchives resolves persistent zip groups from the DB and
// creates max-compression ZIP archives for each. Returns archive paths.
func BuildZipGroupArchives(db *store.DB, groupNames []string, stagingDir string) []string {
	var archives []string

	for _, name := range groupNames {
		archive, err := buildOneZipGroup(db, name, stagingDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrZGCompress, name, err)

			continue
		}

		archives = append(archives, archive)
	}

	return archives
}

// buildOneZipGroup loads a group's items and compresses them.
func buildOneZipGroup(db *store.DB, name, stagingDir string) (string, error) {
	group, err := db.FindZipGroupByName(name)
	if err != nil {
		return "", fmt.Errorf(constants.ErrZGGroupNotDB, name)
	}

	items, err := db.ListZipGroupItems(name)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		fmt.Printf(constants.MsgZGSkipEmpty, name)

		return "", fmt.Errorf("empty group")
	}

	archiveName := resolveArchiveName(group)
	archivePath := filepath.Join(stagingDir, archiveName)

	err = createMaxCompressZip(archivePath, items)
	if err != nil {
		return "", err
	}

	fmt.Printf(constants.MsgZGCompressed, name, archiveName)

	return archivePath, nil
}

// resolveArchiveName returns the archive filename for a group.
func resolveArchiveName(g model.ZipGroup) string {
	if len(g.ArchiveName) > 0 {
		return g.ArchiveName
	}

	return g.Name + ".zip"
}

// BuildAdHocArchive creates a ZIP from ad-hoc paths provided via -Z flags.
// If bundleName is set, all items go into one archive; otherwise each
// gets its own archive.
func BuildAdHocArchive(paths []string, bundleName, stagingDir string) []string {
	if len(bundleName) > 0 {
		return buildAdHocBundle(paths, bundleName, stagingDir)
	}

	return buildAdHocIndividual(paths, stagingDir)
}

// buildAdHocBundle bundles all ad-hoc paths into a single named archive.
func buildAdHocBundle(paths []string, bundleName, stagingDir string) []string {
	items := pathsToItems(paths)
	archivePath := filepath.Join(stagingDir, bundleName)

	err := createMaxCompressZip(archivePath, items)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGCompress, bundleName, err)

		return nil
	}

	fmt.Printf(constants.MsgZGCompressed, "ad-hoc", bundleName)

	return []string{archivePath}
}

// buildAdHocIndividual creates one archive per ad-hoc path.
func buildAdHocIndividual(paths []string, stagingDir string) []string {
	var archives []string

	for _, p := range paths {
		base := filepath.Base(p)
		archiveName := strings.TrimSuffix(base, filepath.Ext(base)) + ".zip"
		archivePath := filepath.Join(stagingDir, archiveName)

		items := pathsToItems([]string{p})

		err := createMaxCompressZip(archivePath, items)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrZGCompress, p, err)

			continue
		}

		fmt.Printf(constants.MsgZGCompressed, p, archiveName)
		archives = append(archives, archivePath)
	}

	return archives
}

// pathsToItems converts raw paths to ZipGroupItem entries.
func pathsToItems(paths []string) []model.ZipGroupItem {
	items := make([]model.ZipGroupItem, 0, len(paths))

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			fmt.Printf(constants.MsgZGSkipMissing, p)

			continue
		}

		items = append(items, model.ZipGroupItem{
			FullPath: p,
			Path:     p,
			IsFolder: info.IsDir(),
		})
	}

	return items
}

// createMaxCompressZip creates a ZIP archive with Deflate level 9.
func createMaxCompressZip(archivePath string, items []model.ZipGroupItem) error {
	outFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("create zip: %w", err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	// Register a custom compressor with max compression.
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return newMaxDeflateWriter(out), nil
	})

	for _, item := range items {
		itemPath := item.FullPath
		if len(itemPath) == 0 {
			itemPath = item.Path
		}

		if item.IsFolder {
			err = addFolderToZip(w, itemPath)
		} else {
			err = addSingleFileToZip(w, itemPath, filepath.Base(itemPath))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// addSingleFileToZip adds one file entry to a zip writer.
func addSingleFileToZip(w *zip.Writer, srcPath, entryName string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", srcPath, err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("stat %s: %w", srcPath, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("header %s: %w", srcPath, err)
	}

	header.Name = entryName
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create entry %s: %w", entryName, err)
	}

	_, err = io.Copy(writer, src)

	return err
}

// addFolderToZip recursively adds a directory's contents to the archive.
func addFolderToZip(w *zip.Writer, folderPath string) error {
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, relErr := filepath.Rel(filepath.Dir(folderPath), path)
		if relErr != nil {
			relPath = path
		}

		return addSingleFileToZip(w, path, filepath.ToSlash(relPath))
	})
}

// DryRunZipGroups prints what zip groups would produce without creating them.
func DryRunZipGroups(db *store.DB, groupNames []string) {
	if len(groupNames) == 0 {
		return
	}

	fmt.Printf(constants.MsgZGDryRunHeader, len(groupNames))

	for _, name := range groupNames {
		items, err := db.ListZipGroupItems(name)
		if err != nil {
			continue
		}

		paths := make([]string, len(items))
		for i, item := range items {
			paths[i] = item.FullPath
			if len(paths[i]) == 0 {
				paths[i] = item.Path
			}
		}

		group, _ := db.FindZipGroupByName(name)
		archiveName := resolveArchiveName(group)

		fmt.Printf(constants.MsgZGDryRunEntry, archiveName, len(items), strings.Join(paths, ", "))
	}
}

// DryRunAdHoc prints what ad-hoc zip items would produce without creating them.
func DryRunAdHoc(paths []string, bundleName string) {
	if len(paths) == 0 {
		return
	}

	if len(bundleName) > 0 {
		fmt.Printf(constants.MsgZGDryRunEntry, bundleName, len(paths), strings.Join(paths, ", "))

		return
	}

	for _, p := range paths {
		base := filepath.Base(p)
		archiveName := strings.TrimSuffix(base, filepath.Ext(base)) + ".zip"

		fmt.Printf(constants.MsgZGDryRunEntry, archiveName, 1, p)
	}
}
