package utils

import (
	"fmt"
	"os"
)

// append to a file
func AppendTo(filePath string, data []byte) {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the data to the end of the file
	if _, err = file.Write(data); err != nil {
		fmt.Println("Error appending to file:", err)
		return
	}

	fmt.Println("Data appended to the file successfully.")
}
