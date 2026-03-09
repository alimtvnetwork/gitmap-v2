// Package cmd — seowriteloop.go handles the commit loop, rotation, and timing.
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/user/gitmap/constants"
)

// commitMessage holds a single title/description pair.
type commitMessage struct {
	title       string
	description string
}

// runCommitLoop executes the timed commit-and-push cycle.
func runCommitLoop(flags seoWriteFlags, messages []commitMessage, minSec, maxSec int) {
	pendingFiles := resolvePendingFiles(flags.files)
	printHeader(flags.maxCommits, minSec, maxSec)

	stop := setupSignalHandler()
	start := time.Now()
	count := 0

	for i, m := range messages {
		if shouldStop(stop, flags.maxCommits, count) {
			break
		}

		commitOne(flags, pendingFiles, m, i, count, len(messages))
		count++
		waitRandom(minSec, maxSec, stop)
	}

	runRotation(flags, messages, pendingFiles, stop, &count, minSec, maxSec)
	printDone(count, time.Since(start))
}

// commitOne stages, commits, and pushes a single file.
func commitOne(flags seoWriteFlags, files []string, m commitMessage, idx, count, total int) {
	file := pickFile(files, idx)
	gitStage(file)
	gitCommitWithAuthor(m.title, m.description, flags.authorName, flags.authorEmail)
	gitPush()
	printCommitLine(flags.maxCommits, count+1, total, m.title, file)
}

// runRotation handles rotation mode when pending files are exhausted.
func runRotation(flags seoWriteFlags, msgs []commitMessage, files []string, stop <-chan bool, count *int, minSec, maxSec int) {
	if flags.maxCommits > 0 && *count >= flags.maxCommits {
		return
	}
	if len(msgs) <= *count && flags.maxCommits == 0 {
		return
	}

	rotateFile := resolveRotateFile(flags.rotateFile)
	if rotateFile == "" {
		return
	}

	rotateLoop(flags, msgs, rotateFile, stop, count, minSec, maxSec)
}

// rotateLoop appends text, commits, reverts, commits in a cycle.
func rotateLoop(flags seoWriteFlags, msgs []commitMessage, file string, stop <-chan bool, count *int, minSec, maxSec int) {
	for flags.maxCommits == 0 || *count < flags.maxCommits {
		if shouldStop(stop, 0, 0) {
			break
		}

		m := msgs[*count%len(msgs)]
		appendToFile(file, m.description)
		gitStage(file)
		gitCommitWithAuthor(m.title, m.description, flags.authorName, flags.authorEmail)
		gitPush()
		printRotationLine(flags.maxCommits, *count+1, file)
		*count++

		revertFile(file, m.description)
		waitRandom(minSec, maxSec, stop)
	}
}

// resolvePendingFiles returns files to stage, filtered by glob if set.
func resolvePendingFiles(pattern string) []string {
	if pattern != "" {
		matches, _ := filepath.Glob(pattern)

		return matches
	}

	out, err := exec.Command("git", "ls-files", "--others", "--modified", "--exclude-standard").Output()
	if err != nil {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, l := range lines {
		if l != "" {
			result = append(result, l)
		}
	}

	return result
}

// pickFile selects a file from the list using round-robin.
func pickFile(files []string, idx int) string {
	if len(files) == 0 {
		return "."
	}

	return files[idx%len(files)]
}

// resolveRotateFile finds or validates the rotation target file.
func resolveRotateFile(explicit string) string {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrSEORotateNotFound, explicit)

			return ""
		}

		return explicit
	}

	return autoDetectRotateFile()
}

// autoDetectRotateFile finds the first .html or .txt file in the repo.
func autoDetectRotateFile() string {
	for _, ext := range []string{"*.html", "*.txt"} {
		matches, _ := filepath.Glob(ext)
		if len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// parseInterval parses "min-max" into two integers.
func parseInterval(s string) (int, int) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		fmt.Fprint(os.Stderr, constants.ErrSEOIntervalFmt)
		os.Exit(1)
	}

	min, err1 := strconv.Atoi(parts[0])
	max, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || min > max {
		fmt.Fprint(os.Stderr, constants.ErrSEOIntervalFmt)
		os.Exit(1)
	}

	return min, max
}

// waitRandom sleeps for a random duration between min and max seconds.
func waitRandom(minSec, maxSec int, stop <-chan bool) {
	delay := minSec + rand.Intn(maxSec-minSec+1)
	fmt.Printf(constants.MsgSEOWaiting, delay)

	timer := time.NewTimer(time.Duration(delay) * time.Second)
	select {
	case <-timer.C:
	case <-stop:
		timer.Stop()
	}
}

// setupSignalHandler returns a channel that closes on Ctrl+C.
func setupSignalHandler() <-chan bool {
	ch := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		fmt.Print(constants.MsgSEOGraceful)
		ch <- true
	}()

	return ch
}

// shouldStop checks if the loop should terminate.
func shouldStop(stop <-chan bool, maxCommits, count int) bool {
	select {
	case <-stop:
		return true
	default:
	}

	if maxCommits > 0 && count >= maxCommits {
		return true
	}

	return false
}

// gitStage runs git add for a file.
func gitStage(file string) {
	cmd := exec.Command("git", "add", file)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOGitStage, err)
	}
}

// gitCommit creates a commit with title and description.
func gitCommit(title, description string) {
	gitCommitWithAuthor(title, description, "", "")
}

// gitCommitWithAuthor creates a commit with optional author override.
func gitCommitWithAuthor(title, description, authorName, authorEmail string) {
	msg := title + "\n\n" + description

	if authorName != "" || authorEmail != "" {
		author := resolveAuthorFlag(authorName, authorEmail)
		cmd := exec.Command("git", "commit", "-m", msg, "--author", author)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrSEOGitCommit, err)
		}

		return
	}

	cmd := exec.Command("git", "commit", "-m", msg)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOGitCommit, err)
	}
}

// resolveAuthorFlag builds the --author "Name <email>" string.
func resolveAuthorFlag(name, email string) string {
	if name == "" {
		out, _ := exec.Command("git", "config", "user.name").Output()
		name = strings.TrimSpace(string(out))
	}

	if email == "" {
		out, _ := exec.Command("git", "config", "user.email").Output()
		email = strings.TrimSpace(string(out))
	}

	return name + " <" + email + ">"
}

// gitPush pushes to the remote.
func gitPush() {
	cmd := exec.Command("git", "push")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOGitPush, err)
	}
}

// appendToFile appends text to a file for rotation mode.
func appendToFile(path, text string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString("\n" + text)
}

// revertFile removes the appended text from the file.
func revertFile(path, text string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	cleaned := strings.Replace(string(data), "\n"+text, "", 1)
	_ = os.WriteFile(path, []byte(cleaned), 0o644)
}

// printHeader outputs the commit plan header.
func printHeader(max, minSec, maxSec int) {
	if max > 0 {
		fmt.Printf(constants.MsgSEOHeader, max, minSec, maxSec)

		return
	}

	fmt.Printf(constants.MsgSEOHeaderUnlimited, minSec, maxSec)
}

// printCommitLine outputs a single commit progress line.
func printCommitLine(max, current, total int, title, file string) {
	if max > 0 {
		fmt.Printf(constants.MsgSEOCommit, current, max, title, file)

		return
	}

	fmt.Printf(constants.MsgSEOCommitOpen, current, title, file)
}

// printRotationLine outputs a rotation progress line.
func printRotationLine(max, current int, file string) {
	if max > 0 {
		fmt.Printf(constants.MsgSEORotation, current, max, file)

		return
	}

	fmt.Printf(constants.MsgSEORotationOpen, current, file)
}

// printDone outputs the final summary line.
func printDone(count int, elapsed time.Duration) {
	fmt.Printf(constants.MsgSEODone, count, formatDuration(elapsed))
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}

	return fmt.Sprintf("%dm", m)
}
