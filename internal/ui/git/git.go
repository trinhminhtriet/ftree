package git

import (
	"os/exec"
	"strings"
)

// GitStatus represents the status of a file in Git.
type GitStatus struct {
	Path     string
	Staged   bool
	Modified bool
	Untracked bool
}

// GetGitStatus retrieves the Git status for files in the current directory.
func GetGitStatus(dir string) (map[string]GitStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain", "-z")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
			return nil, nil // Not a Git repository, return empty map
		}
		return nil, err
	}

	statusMap := make(map[string]GitStatus)
	lines := strings.Split(string(output), "\x00")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status := GitStatus{
			Path:     line[3:],
			Staged:   strings.ContainsAny(line[:2], "M A D R C U"),
			Modified: strings.Contains(line[:2], "M"),
			Untracked: strings.HasPrefix(line, "??"),
		}
		statusMap[status.Path] = status
	}
	return statusMap, nil
}