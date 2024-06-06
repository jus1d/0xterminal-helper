package git

import (
	"bytes"
	"os/exec"
)

func LatestShortenedCommit() string {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%h")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return out.String()
}

func LatestCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return out.String()
}

func LatestCommitMessage() string {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return out.String()
}

func CurrentBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return out.String()
}
