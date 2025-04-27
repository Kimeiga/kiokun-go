package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/jmdict/testdata"
)

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Change to project root if needed
	if !filepath.IsAbs(cwd) || !filepath.HasPrefix(filepath.Base(cwd), "kiokun-go") {
		// Try to find project root
		for i := 0; i < 3; i++ { // Go up at most 3 levels
			if _, err := os.Stat(filepath.Join(cwd, "dictionaries")); err == nil {
				break
			}
			cwd = filepath.Dir(cwd)
		}
		if err := os.Chdir(cwd); err != nil {
			fmt.Printf("Error changing to root directory: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Running field finder tool...")
	if err := testdata.SimpleToolToFindExamples(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Field finder completed successfully!")
}
