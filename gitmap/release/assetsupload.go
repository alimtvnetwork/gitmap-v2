// Package release — assetsupload.go uploads release assets via GitHub API.
package release

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// UploadAsset uploads a single file to a GitHub release.
func UploadAsset(owner, repo string, releaseID int, filePath, token string) error {
	filename := filepath.Base(filePath)
	url := fmt.Sprintf("https://uploads.github.com/repos/%s/%s/releases/%d/assets?name=%s",
		owner, repo, releaseID, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open asset: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodPost, url, file)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		return &uploadError{
			statusCode: resp.StatusCode,
			message:    string(respBody),
		}
	}

	return nil
}

// uploadError captures HTTP status for retry decisions.
type uploadError struct {
	statusCode int
	message    string
}

// Error implements the error interface.
func (e *uploadError) Error() string {
	return fmt.Sprintf("upload error %d: %s", e.statusCode, e.message)
}

// UploadAllAssets uploads all assets to a GitHub release with exponential backoff retry.
func UploadAllAssets(owner, repo string, releaseID int, assets []string, token string) {
	for _, asset := range assets {
		uploadSingleAsset(owner, repo, releaseID, asset, token)
	}
}

// uploadSingleAsset uploads one asset with retry logic.
func uploadSingleAsset(owner, repo string, releaseID int, asset, token string) {
	filename := filepath.Base(asset)

	err := withRetry(filename, constants.RetryMaxAttempts, func() error {
		return UploadAsset(owner, repo, releaseID, asset, token)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAssetUploadFinal, filename, err)

		return
	}

	fmt.Printf(constants.MsgAssetUploaded, filename)
}
