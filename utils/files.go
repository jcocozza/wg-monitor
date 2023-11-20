package utils

import (
	"log/slog"
	"os"
	"os/exec"
)

// append to a file
func AppendTo(filePath string, data []byte) {
	slog.Debug("Appending to: " + filePath)
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the data to the end of the file
	if _, err = file.Write(data); err != nil {
		slog.Error("Error appending to file:", err)
		return
	}
}

// read file (with "cat" - for sudo purposes) of a given path
func ReadFile(path string) ([]byte, error) {
	slog.Debug("Reading file: " + path)
	cmd := exec.Command("cat", path)
	output, err := cmd.CombinedOutput()
	return output, err
}

//make a directory
func MkDir(dirPath string) {
	slog.Debug("Creating directory: " + dirPath)
	err := os.Mkdir(dirPath, 0755)

	if err != nil {
		slog.Error("Error creating directory:", err)
		return
	}
}

// Write a file
func WriteFile(filePath string, data []byte) {
	slog.Debug("Writing to file: " + filePath)
	err := os.WriteFile(filePath, data, 0644)

	if err != nil {
		slog.Error("Error writing file:", err)
		return
	}
}