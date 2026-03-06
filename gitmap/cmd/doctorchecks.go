package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// checkRepoPath reports whether RepoPath is embedded.
func checkRepoPath() int {
	if len(constants.RepoPath) == 0 {
		printIssue(constants.DoctorRepoPathMissing, constants.DoctorRepoPathDetail)
		printFix(constants.DoctorRepoPathFix)

		return 1
	}

	printOK(constants.DoctorRepoPathOKFmt, constants.RepoPath)

	return 0
}

// checkActiveBinary reports the gitmap binary on PATH.
func checkActiveBinary() int {
	path, err := exec.LookPath("gitmap")
	if err != nil {
		printIssue(constants.DoctorPathMissTitle, constants.DoctorPathMissDetail)
		printFix(constants.DoctorPathMissFix)

		return 1
	}

	absPath, _ := filepath.Abs(path)
	version := getBinaryVersion(absPath)
	printOK(constants.DoctorPathBinaryFmt, absPath, version)

	return 0
}

// checkDeployedBinary reports the deployed binary from powershell.json.
func checkDeployedBinary() int {
	if len(constants.RepoPath) == 0 {
		return 0
	}

	data, err := readPowershellJSON()
	if err != nil {
		printIssue(constants.DoctorDeployReadFail, constants.DoctorDeployReadDet)

		return 1
	}

	deployedBinary, issue := resolveDeployedFromData(data)
	if issue > 0 {
		return issue
	}

	version := getBinaryVersion(deployedBinary)
	printOK(constants.DoctorDeployOKFmt, deployedBinary, version)

	return 0
}

// readPowershellJSON reads the powershell.json config file.
func readPowershellJSON() ([]byte, error) {
	configPath := filepath.Join(constants.RepoPath, "gitmap", "powershell.json")

	return os.ReadFile(configPath)
}

// resolveDeployedFromData extracts and validates the deployed binary path.
func resolveDeployedFromData(data []byte) (string, int) {
	deployPath := extractJSONString(data, "deployPath")
	if len(deployPath) == 0 {
		printIssue(constants.DoctorNoDeployPath, constants.DoctorNoDeployDet)

		return "", 1
	}

	binaryName := extractJSONString(data, "binaryName")
	if len(binaryName) == 0 {
		binaryName = constants.DoctorDefaultBinary
	}

	deployedBinary := filepath.Join(deployPath, "gitmap", binaryName)
	if _, err := os.Stat(deployedBinary); err != nil {
		printIssue(constants.DoctorDeployNotFound, deployedBinary)
		printFix(constants.DoctorDeployRunFix)

		return "", 1
	}

	return deployedBinary, 0
}

// checkVersionMismatch compares PATH vs deployed vs source versions.
func checkVersionMismatch() int {
	sourceVersion := fmt.Sprintf(constants.MsgVersionFmt[:len(constants.MsgVersionFmt)-1], constants.Version)
	activeVersion, activePath := getActiveVersion()
	deployedVersion, deployedPath := getDeployedVersion()
	issues := 0

	issues += checkActiveVsSource(activeVersion, sourceVersion)
	issues += checkDeployedVsSource(deployedVersion, sourceVersion)
	issues += checkActiveVsDeployed(activeVersion, deployedVersion, activePath, deployedPath)

	if issues == 0 {
		printOK(constants.DoctorSourceOKFmt, sourceVersion)
	}

	return issues
}

// getActiveVersion returns version and path of the active PATH binary.
func getActiveVersion() (string, string) {
	path, err := exec.LookPath("gitmap")
	if err != nil {
		return "", ""
	}

	absPath, _ := filepath.Abs(path)

	return getBinaryVersion(absPath), absPath
}

// getDeployedVersion returns version and path of the deployed binary.
func getDeployedVersion() (string, string) {
	if len(constants.RepoPath) == 0 {
		return "", ""
	}

	data, err := readPowershellJSON()
	if err != nil {
		return "", ""
	}

	dp := extractJSONString(data, "deployPath")
	bn := extractJSONString(data, "binaryName")
	if len(bn) == 0 {
		bn = constants.DoctorDefaultBinary
	}

	if len(dp) == 0 {
		return "", ""
	}

	deployedPath := filepath.Join(dp, "gitmap", bn)

	return getBinaryVersion(deployedPath), deployedPath
}

// checkActiveVsSource reports if PATH binary differs from source.
func checkActiveVsSource(activeVersion, sourceVersion string) int {
	if len(activeVersion) > 0 && activeVersion != sourceVersion {
		printIssue(constants.DoctorVersionMismatch,
			fmt.Sprintf(constants.DoctorVMismatchFmt, activeVersion, sourceVersion))
		printFix(constants.DoctorVMismatchFix)

		return 1
	}

	return 0
}

// checkDeployedVsSource reports if deployed binary differs from source.
func checkDeployedVsSource(deployedVersion, sourceVersion string) int {
	if len(deployedVersion) > 0 && deployedVersion != sourceVersion {
		printIssue(constants.DoctorDeployMismatch,
			fmt.Sprintf(constants.DoctorDMismatchFmt, deployedVersion, sourceVersion))
		printFix(constants.DoctorDMismatchFix)

		return 1
	}

	return 0
}

// checkActiveVsDeployed reports if PATH and deployed binaries differ.
func checkActiveVsDeployed(activeVersion, deployedVersion, activePath, deployedPath string) int {
	if len(activeVersion) == 0 || len(deployedVersion) == 0 {
		return 0
	}

	if activeVersion == deployedVersion {
		return 0
	}

	absActive, _ := filepath.Abs(activePath)
	absDeployed, _ := filepath.Abs(deployedPath)
	if absActive == absDeployed {
		return 0
	}

	printIssue(constants.DoctorBinariesDiffer,
		fmt.Sprintf(constants.DoctorBDifferFmt, absActive, activeVersion, absDeployed, deployedVersion))
	printFix(constants.DoctorBDifferFix)

	return 1
}

// checkGit verifies git is available.
func checkGit() int {
	path, err := exec.LookPath(constants.GitBin)
	if err != nil {
		printIssue(constants.DoctorGitMissTitle, constants.DoctorGitMissDetail)

		return 1
	}

	version := getToolVersion(constants.GitBin, "--version")
	if len(version) == 0 {
		printOK(constants.DoctorGitOKPathFmt, path)

		return 0
	}

	printOK(constants.DoctorGitOKFmt, path, version)

	return 0
}

// checkGo verifies Go is available for building.
func checkGo() int {
	path, err := exec.LookPath("go")
	if err != nil {
		printWarn(constants.DoctorGoWarn)

		return 0
	}

	version := getToolVersion("go", "version")
	if len(version) == 0 {
		printOK(constants.DoctorGoOKPathFmt, path)

		return 0
	}

	printOK(constants.DoctorGoOKFmt, version)

	return 0
}

// getToolVersion runs a tool with an arg and returns trimmed output.
func getToolVersion(tool, arg string) string {
	cmd := exec.Command(tool, arg)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// checkChangelogFile verifies CHANGELOG.md exists.
func checkChangelogFile() int {
	if _, err := os.Stat(constants.ChangelogFile); err != nil {
		printWarn(constants.DoctorChangelogWarn)

		return 0
	}

	printOK(constants.DoctorChangelogOK)

	return 0
}
