package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var destinationPath string
var atomicDesignType string

func main() {
	var rootCmd = &cobra.Command{Use: "copy-files"}
	rootCmd.PersistentFlags().StringVarP(&destinationPath, "path", "p", "", "Destination path (required)")
	rootCmd.PersistentFlags().StringVarP(&atomicDesignType, "type", "t", "", "Type of atomic design (skeleton, ts, vue) (required)")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		err := validateInputs()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		err = copyFiles()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Files copied successfully.")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func validateInputs() error {
	if destinationPath == "" {
		return fmt.Errorf("please provide a destination path using --path option")
	}

	if atomicDesignType == "" {
		return fmt.Errorf("please provide a type of atomic design using --type option (skeleton)")
	}

	return nil
}

func copyFiles() error {
	//scriptDirectory, err := os.Executable();
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	scriptDirectory := filepath.Dir(executable)

	sourcePath := filepath.Join(scriptDirectory, "components", atomicDesignType)
	destinationPath := filepath.Join(destinationPath, "components")

	err = copyDir(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceFile := filepath.Join(src, entry.Name())
		destinationFile := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err := copyDir(sourceFile, destinationFile)
			if err != nil {
				return err
			}
		} else {
			err := copyFile(sourceFile, destinationFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
