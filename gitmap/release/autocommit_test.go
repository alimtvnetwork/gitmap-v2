package release

import (
	"errors"
	"strings"
	"testing"
)

func TestIsNonFastForwardPushError(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{
			name: "fetch first rejection",
			output: "! [rejected] main -> main (fetch first)\n" +
				"error: failed to push some refs",
			want: true,
		},
		{
			name:   "explicit non fast forward",
			output: "error: non-fast-forward update rejected",
			want:   true,
		},
		{
			name:   "unrelated auth failure",
			output: "fatal: Authentication failed",
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isNonFastForwardPushError(tc.output)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFormatGitCommandErrorPrefersGitOutput(t *testing.T) {
	err := formatGitCommandError("error: failed to push some refs", errors.New("exit status 1"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to push some refs") {
		t.Fatalf("expected git output in error, got %q", err.Error())
	}
}

func TestTrimGitOutputFallback(t *testing.T) {
	got := trimGitOutput("   \n\t  ")
	if got != "unknown git error" {
		t.Fatalf("got %q, want %q", got, "unknown git error")
	}
}
