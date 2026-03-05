// Package formatter renders ScanRecords to terminal output.
package formatter

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// Terminal writes a professional colored output to the given writer.
func Terminal(w io.Writer, records []model.ScanRecord, outputDir string) error {
	printBanner(w, len(records))
	printRepoList(w, records)
	printFolderTree(w, records)
	printOutputFiles(w, outputDir)
	printCloneHelp(w)

	return nil
}

// printBanner writes the header section.
func printBanner(w io.Writer, count int) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorCyan+constants.TermBannerTop+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorCyan+constants.TermBannerTitle+constants.ColorReset+"\n", constants.Version)
	fmt.Fprintf(w, constants.ColorCyan+constants.TermBannerBottom+constants.ColorReset+"\n")
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorGreen+constants.TermFoundFmt+constants.ColorReset+"\n", count)
	fmt.Fprintln(w)
}

// printRepoList writes each repo with folder name and clone instruction.
func printRepoList(w io.Writer, records []model.ScanRecord) {
	fmt.Fprintf(w, constants.ColorYellow+constants.TermReposHeader+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	for i, r := range records {
		printOneRepo(w, r, i+1, len(records))
	}
	fmt.Fprintln(w)
}

// printOneRepo writes a single repo entry with index.
func printOneRepo(w io.Writer, r model.ScanRecord, idx, total int) {
	fmt.Fprintf(w, constants.ColorDim+"  %d/%d "+constants.ColorReset, idx, total)
	fmt.Fprintf(w, constants.ColorGreen+"📦 %s"+constants.ColorReset, r.RepoName)
	fmt.Fprintf(w, constants.ColorDim+" (%s)"+constants.ColorReset+"\n", r.Branch)
	fmt.Fprintf(w, constants.ColorDim+"       └─ "+constants.ColorReset)
	fmt.Fprintf(w, constants.ColorWhite+"%s"+constants.ColorReset+"\n", r.CloneInstruction)
}

// printFolderTree writes the folder structure to terminal.
func printFolderTree(w io.Writer, records []model.ScanRecord) {
	fmt.Fprintf(w, constants.ColorYellow+constants.TermTreeHeader+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	paths := collectTermPaths(records)
	tree := buildTermTree(paths)
	renderTermTree(w, tree, "  ")
	fmt.Fprintln(w)
}

// printOutputFiles shows the generated output files.
func printOutputFiles(w io.Writer, outputDir string) {
	fmt.Fprintf(w, constants.ColorYellow+"  ■ Output Files"+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+"  📁 %s/"+constants.ColorReset+"\n", outputDir)
	printOutputFile(w, outputDir, constants.DefaultCSVFile, "Repo data in CSV")
	printOutputFile(w, outputDir, constants.DefaultJSONFile, "Repo data in JSON")
	printOutputFile(w, outputDir, constants.DefaultStructureFile, "Folder tree")
	printOutputFile(w, outputDir, constants.DefaultCloneScript, "PowerShell clone script")
	printOutputFile(w, outputDir, constants.DefaultDesktopScript, "GitHub Desktop registration")
	fmt.Fprintln(w)
}

// printOutputFile shows one output file entry.
func printOutputFile(w io.Writer, dir, name, desc string) {
	fullPath := filepath.Join(dir, name)
	fmt.Fprintf(w, constants.ColorDim+"  ├── "+constants.ColorReset)
	fmt.Fprintf(w, constants.ColorCyan+"📄 %s"+constants.ColorReset, name)
	fmt.Fprintf(w, constants.ColorDim+"  %s"+constants.ColorReset+"\n", desc)
	_ = fullPath
}

// printCloneHelp writes instructions for cloning on another machine.
func printCloneHelp(w io.Writer) {
	fmt.Fprintf(w, constants.ColorYellow+constants.TermCloneHeader+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorWhite+constants.TermCloneStep1+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorCyan+constants.TermCloneCmd1+constants.ColorReset+"\n")
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorWhite+constants.TermCloneStep2+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorCyan+constants.TermCloneCmd2+constants.ColorReset+"\n")
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorWhite+constants.TermCloneStep3+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorCyan+constants.TermCloneCmd3+constants.ColorReset+"\n")
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorWhite+"  4. Or just run the PowerShell script:"+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorCyan+"     .\\clone.ps1 -TargetDir .\\projects"+constants.ColorReset+"\n")
	fmt.Fprintln(w)
}

// termPathEntry holds path info for tree rendering.
type termPathEntry struct {
	Path   string
	Branch string
}

// collectTermPaths extracts and sorts paths from records.
func collectTermPaths(records []model.ScanRecord) []termPathEntry {
	entries := make([]termPathEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, termPathEntry{
			Path: r.RelativePath, Branch: r.Branch,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	return entries
}

// termNode represents a folder or repo in the terminal tree.
type termNode struct {
	Name     string
	Children []*termNode
	IsRepo   bool
	Branch   string
}

// buildTermTree constructs a tree from path entries.
func buildTermTree(entries []termPathEntry) *termNode {
	root := &termNode{Name: "."}
	for _, e := range entries {
		insertTermNode(root, e)
	}

	return root
}

// insertTermNode adds a path into the tree.
func insertTermNode(root *termNode, entry termPathEntry) {
	normalized := strings.ReplaceAll(entry.Path, "\\", "/")
	parts := strings.Split(normalized, "/")
	current := root
	for i, part := range parts {
		child := findTermChild(current, part)
		if child == nil {
			child = &termNode{Name: part}
			current.Children = append(current.Children, child)
		}
		if i == len(parts)-1 {
			child.IsRepo = true
			child.Branch = entry.Branch
		}
		current = child
	}
}

// findTermChild looks for a child node by name.
func findTermChild(node *termNode, name string) *termNode {
	for _, c := range node.Children {
		if c.Name == name {
			return c
		}
	}

	return nil
}

// renderTermTree writes the colored tree to the writer.
func renderTermTree(w io.Writer, node *termNode, prefix string) {
	for i, child := range node.Children {
		connector := constants.TreeBranch
		nextPrefix := prefix + constants.TreePipe
		if i == len(node.Children)-1 {
			connector = constants.TreeCorner
			nextPrefix = prefix + constants.TreeSpace
		}
		renderTermNode(w, child, prefix, connector)
		if len(child.Children) > 0 {
			renderTermTree(w, child, nextPrefix)
		}
	}
}

// renderTermNode writes a single colored tree node.
func renderTermNode(w io.Writer, node *termNode, prefix, connector string) {
	if node.IsRepo {
		fmt.Fprintf(w, "%s%s%s 📦 %s%s%s %s(%s)%s\n",
			constants.ColorDim, prefix, connector,
			constants.ColorGreen, node.Name, constants.ColorReset,
			constants.ColorDim, node.Branch, constants.ColorReset)

		return
	}
	fmt.Fprintf(w, "%s%s%s 📁 %s%s%s\n",
		constants.ColorDim, prefix, connector,
		constants.ColorYellow, node.Name, constants.ColorReset)
}
