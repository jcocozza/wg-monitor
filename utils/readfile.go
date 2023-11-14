package utils

import "os/exec"

// read file (with "cat" - for sudo purposes) of a given path
func ReadFile(path string) ([]byte, error) {
	cmd := exec.Command("cat", path)
	output, err := cmd.CombinedOutput()
	return output, err
}